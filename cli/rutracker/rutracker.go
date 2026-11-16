// Утилита для отдельной проверки провайдера rutracker (авторизация + поиск).
//
// Два режима (секреты лучше передавать через переменные окружения, а не флагами, чтобы не светить их в истории):
//
//  1. Готовая сессия из обычного браузера (rutracker закрыт Cloudflare, логин через headless-браузер не проходит):
//     RUTRACKER_COOKIES='bb_session=...; cf_clearance=...' RUTRACKER_UA='Mozilla/5.0 ...' go run ./cli/rutracker -query "Матрица" -download
//     Тот же код, что использует сервис для аккаунтов с полями cookies / userAgent.
//
//  2. Логин и пароль (авторизация через headless-браузер):
//     RUTRACKER_LOGIN=user RUTRACKER_PASSWORD=pass go run ./cli/rutracker -query "Матрица"
package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/http/httptrace"
	"os"
	"strings"
	"time"

	"github.com/RacoonMediaServer/rms-media-discovery/internal/utils"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/media"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/navigator"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider"
	"github.com/RacoonMediaServer/rms-media-discovery/pkg/provider/rutracker"
	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
)

// accessor отдает один фиксированный аккаунт; после MarkUnaccesible аккаунт перестает выдаваться (как в настоящем сервисе)
type accessor struct {
	cred   model.Credentials
	marked bool
}

func (a *accessor) GetCredentials(serviceId string) (model.Credentials, error) {
	if a.marked {
		return model.Credentials{}, errors.New("account not found")
	}
	return a.cred, nil
}

func (a *accessor) GetApiKey(serviceId string) (model.ApiKey, error) {
	return model.ApiKey{}, fmt.Errorf("not supported")
}

func (a *accessor) MarkUnaccesible(accountId string) {
	fmt.Fprintf(os.Stderr, "!!! account %q marked as unaccessible\n", accountId)
	a.marked = true
}

// stdinSolver просит ввести капчу руками
type stdinSolver struct{}

func (stdinSolver) ID() string { return "stdin" }

func (stdinSolver) Solve(ctx context.Context, c provider.Captcha) (string, error) {
	fmt.Fprintf(os.Stderr, "Captcha required: %s\nEnter code: ", c.Url)
	s, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(s), nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// tracingTransport логирует каждый запрос: реальное соединение (адреса, TLS, ALPN), отправленные заголовки и ответ.
// Нужен, чтобы сравнивать запросы из разных окружений (хост / контейнер), когда сайт отвечает по-разному
type tracingTransport struct {
	base http.RoundTripper
}

func (t tracingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	trace := &httptrace.ClientTrace{
		GotConn: func(i httptrace.GotConnInfo) {
			log.Infof("conn: %s -> %s (reused=%v)", i.Conn.LocalAddr(), i.Conn.RemoteAddr(), i.Reused)
		},
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			log.Infof("tls: version=0x%04x alpn=%q server_name=%q err=%v", cs.Version, cs.NegotiatedProtocol, cs.ServerName, err)
		},
		WroteHeaderField: func(key string, value []string) {
			if strings.EqualFold(key, "cookie") {
				value = []string{"<redacted>"}
			}
			log.Infof("> %s: %s", key, strings.Join(value, ", "))
		},
	}
	resp, err := t.base.RoundTrip(req.WithContext(httptrace.WithClientTrace(req.Context(), trace)))
	if err != nil {
		log.Warnf("< round trip failed: %s", err)
		return resp, err
	}
	log.Infof("< %s %s server=%q cf-ray=%q cf-mitigated=%q", resp.Proto, resp.Status, resp.Header.Get("server"), resp.Header.Get("cf-ray"), resp.Header.Get("cf-mitigated"))
	return resp, nil
}

// reportRoute показывает, каким адресом (IPv4/IPv6) пойдут запросы к сайту: эта проверка использует тот же транспорт,
// что и провайдер, а cf_clearance может приниматься только для того семейства адресов, для которого он выдан
func reportRoute() {
	t, ok := http.DefaultTransport.(*http.Transport)
	if !ok || t.DialContext == nil {
		t = &http.Transport{}
	}
	dial := t.DialContext
	if dial == nil {
		dial = (&net.Dialer{Timeout: 15 * time.Second}).DialContext
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	conn, err := dial(ctx, "tcp", "rutracker.org:443")
	if err != nil {
		log.Warnf("Route check failed: %s", err)
		return
	}
	defer conn.Close()
	family := "IPv4"
	if addr, ok := conn.RemoteAddr().(*net.TCPAddr); ok && addr.IP.To4() == nil {
		family = "IPv6"
	}
	log.Infof("Route to rutracker.org: %s -> %s (%s)", conn.LocalAddr(), conn.RemoteAddr(), family)
}

func main() {
	query := flag.String("query", "", "search query")
	login := flag.String("login", getenv("RUTRACKER_LOGIN", ""), "rutracker login (or RUTRACKER_LOGIN)")
	password := flag.String("password", getenv("RUTRACKER_PASSWORD", ""), "rutracker password (or RUTRACKER_PASSWORD)")
	dumps := flag.String("dumps", "/tmp/rutracker-dumps", "directory for screenshots/html dumps on error")
	timeout := flag.Duration("timeout", 3*time.Minute, "total timeout")
	limit := flag.Uint("limit", 10, "max results to print")
	download := flag.Bool("download", false, "also try to download .torrent of the first result")
	cookies := flag.String("cookies", getenv("RUTRACKER_COOKIES", ""), "cookie header string 'a=b; c=d' (or RUTRACKER_COOKIES); login via browser is skipped")
	ua := flag.String("ua", getenv("RUTRACKER_UA", ""), "User-Agent used to obtain cf_clearance (or RUTRACKER_UA)")
	ipv6 := flag.Bool("ipv6", false, "prefer IPv6 for outgoing connections (same as RMS_PREFERIPV6 of the service)")
	flag.Parse()

	if *query == "" {
		fmt.Fprintln(os.Stderr, "-query must be set")
		os.Exit(1)
	}

	log.SetHandler(text.New(os.Stderr))
	log.SetLevel(log.DebugLevel)

	if *ipv6 {
		utils.PreferIPv6()
		log.Info("Outgoing connections prefer IPv6")
	}
	reportRoute()
	http.DefaultTransport = tracingTransport{base: http.DefaultTransport}

	acc := &accessor{cred: model.Credentials{AccountId: "cli"}}
	if *cookies != "" {
		if *ua == "" {
			fmt.Fprintln(os.Stderr, "-ua (RUTRACKER_UA) must be set together with cookies: cf_clearance is bound to User-Agent")
			os.Exit(1)
		}
		acc.cred.Cookies = *cookies
		acc.cred.UserAgent = *ua
		log.Infof("Using ready cookies (browser login skipped), UA: %s", *ua)
	} else {
		if *login == "" || *password == "" {
			fmt.Fprintln(os.Stderr, "set either cookies (RUTRACKER_COOKIES + RUTRACKER_UA) or credentials (RUTRACKER_LOGIN / RUTRACKER_PASSWORD)")
			os.Exit(1)
		}
		acc.cred.Login = *login
		acc.cred.Password = *password

		log.Info("Headless browser engine initializing...")
		if err := navigator.Initialize(); err != nil {
			log.Fatalf("Failed: %s", err)
		}
		navigator.SetSettings(navigator.Settings{StoreDumpOnError: true, DefaultDumpLocation: *dumps})
	}
	p := rutracker.NewProvider(acc, stdinSolver{})

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	start := time.Now()
	results, err := p.SearchTorrents(ctx, model.SearchQuery{
		Query: *query,
		Type:  media.Movies,
		Limit: *limit,
	})
	if err != nil {
		log.Errorf("SearchTorrents failed after %s: %s (dumps: %s)", time.Since(start), err, *dumps)
		os.Exit(1)
	}

	fmt.Printf("Found %d results in %s\n", len(results), time.Since(start))
	for i, t := range results {
		if uint(i) >= *limit {
			break
		}
		fmt.Printf("%2d. %s | %d MB | seeders=%d | link=%q detail=%q\n", i+1, t.Title, t.SizeMB, t.Seeders, t.Link, t.DetailLink)
	}

	if *download && len(results) > 0 {
		data, err := results[0].Downloader(ctx)
		if err != nil {
			log.Errorf("download failed: %s", err)
			os.Exit(1)
		}
		fmt.Printf("Downloaded %d bytes of .torrent\n", len(data))
	}
}
