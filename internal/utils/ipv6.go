package utils

import (
	"context"
	"net"
	"net/http"
	"time"
)

// ipv6AttemptTimeout - сколько ждем подключения по IPv6, прежде чем перейти на обычный выбор адреса
const ipv6AttemptTimeout = 5 * time.Second

// PreferIPv6 заставляет все HTTP-клиенты по умолчанию (http.DefaultTransport: colly, requester и др.)
// сначала подключаться по IPv6, а при неудаче (нет AAAA, нет маршрута, таймаут) - как обычно.
//
// Зачем: в Docker-сети с IPv6 контейнер получает адрес из ULA (fd00::/8). По RFC 6724 (правило 5, совпадение меток)
// резолвер Go в таком случае выбирает IPv4, хотя наружу контейнер ходит и по IPv6 (через NAT66).
// А rutracker привязывает cf_clearance к адресу, с которого он получен, и с IPv4 сервера cookies отклоняются.
func PreferIPv6() {
	t, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return
	}
	base := &net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}
	t.DialContext = preferIPv6DialContext(base, ipv6AttemptTimeout)
}

func preferIPv6DialContext(base *net.Dialer, v6Timeout time.Duration) func(ctx context.Context, network, addr string) (net.Conn, error) {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		if network != "tcp" {
			return base.DialContext(ctx, network, addr)
		}

		v6Ctx, cancel := context.WithTimeout(ctx, v6Timeout)
		defer cancel()
		conn, err := base.DialContext(v6Ctx, "tcp6", addr)
		if err == nil {
			return conn, nil
		}
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		return base.DialContext(ctx, network, addr)
	}
}
