package utils

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// listenBoth поднимает слушателей на 127.0.0.1 и [::1] с одним и тем же портом
func listenBoth(t *testing.T) (v4, v6 net.Listener, port string) {
	l6, err := net.Listen("tcp6", "[::1]:0")
	if err != nil {
		t.Skipf("IPv6 loopback is not available: %s", err)
	}
	_, port, _ = net.SplitHostPort(l6.Addr().String())
	l4, err := net.Listen("tcp4", "127.0.0.1:"+port)
	if err != nil {
		l6.Close()
		t.Skipf("cannot bind the same port on IPv4: %s", err)
	}
	t.Cleanup(func() { l4.Close(); l6.Close() })
	return l4, l6, port
}

func accepted(l net.Listener) <-chan struct{} {
	ch := make(chan struct{}, 1)
	go func() {
		if c, err := l.Accept(); err == nil {
			c.Close()
			ch <- struct{}{}
		}
	}()
	return ch
}

func TestPreferIPv6DialContext_PrefersIPv6(t *testing.T) {
	l4, l6, port := listenBoth(t)
	got4, got6 := accepted(l4), accepted(l6)

	dial := preferIPv6DialContext(&net.Dialer{}, time.Second)
	// "localhost" резолвится в оба адреса: без предпочтения IPv6 выбор зависит от настроек системы
	conn, err := dial(context.Background(), "tcp", "localhost:"+port)
	require.NoError(t, err)
	defer conn.Close()

	assert.Equal(t, "::1", conn.RemoteAddr().(*net.TCPAddr).IP.String())
	select {
	case <-got6:
	case <-time.After(time.Second):
		t.Fatal("IPv6 listener did not get the connection")
	}
	select {
	case <-got4:
		t.Fatal("IPv4 listener must not get the connection")
	default:
	}
}

func TestPreferIPv6DialContext_FallsBackToIPv4(t *testing.T) {
	l4, err := net.Listen("tcp4", "127.0.0.1:0")
	require.NoError(t, err)
	defer l4.Close()
	got4 := accepted(l4)
	_, port, _ := net.SplitHostPort(l4.Addr().String())

	dial := preferIPv6DialContext(&net.Dialer{}, time.Second)

	// IPv4-литерал: IPv6-попытка невозможна, должен сработать обычный путь
	conn, err := dial(context.Background(), "tcp", "127.0.0.1:"+port)
	require.NoError(t, err)
	conn.Close()
	select {
	case <-got4:
	case <-time.After(time.Second):
		t.Fatal("IPv4 listener did not get the connection")
	}

	// имя, у которого слушает только IPv4: IPv6-попытка отклонена, переход на IPv4
	got4 = accepted(l4)
	conn, err = dial(context.Background(), "tcp", "localhost:"+port)
	require.NoError(t, err)
	conn.Close()
	select {
	case <-got4:
	case <-time.After(time.Second):
		t.Fatal("fallback to IPv4 did not happen")
	}
}

func TestPreferIPv6DialContext_NonTCPNetworkIsUntouched(t *testing.T) {
	l, err := net.Listen("tcp4", "127.0.0.1:0")
	require.NoError(t, err)
	defer l.Close()
	_, port, _ := net.SplitHostPort(l.Addr().String())
	_, err = strconv.Atoi(port)
	require.NoError(t, err)

	dial := preferIPv6DialContext(&net.Dialer{}, time.Second)
	conn, err := dial(context.Background(), "tcp4", "127.0.0.1:"+port)
	require.NoError(t, err)
	conn.Close()
}
