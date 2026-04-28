package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/alecthomas/kong"
	mdns "github.com/pion/mdns/v2"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

type Params struct {
	Domain   string `default:"local." help:"The domain to advertise the service under"`
	Hostname string `help:"The hostname of the advertised service"`
	Verbose  bool   `help:"Turn on verbose logging"`
}

var params Params

func main() {
	kong.Parse(&params, kong.Name("mdns-svc"), kong.Description("A tiny cli wrapper for Pions mDNS library"))
	slog.SetDefault(slog.Default())
	if params.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)
	signal.Notify(sigs, syscall.SIGINT)

	server, err := createServer(params)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer server.Close()
	<-sigs
}

func createServer(params Params) (*mdns.Conn, error) {
	var (
		err      error
		hostname string
	)
	if params.Hostname == "" {
		hostname, err = os.Hostname()
		if err != nil {
			return nil, err
		}
	} else {
		hostname = strings.TrimRight(params.Hostname, ".")
	}

	addr4, err := net.ResolveUDPAddr("udp4", mdns.DefaultAddressIPv4)
	if err != nil {
		return nil, err
	}
	addr6, err := net.ResolveUDPAddr("udp6", mdns.DefaultAddressIPv6)
	if err != nil {
		return nil, err
	}

	l4, err := net.ListenUDP("udp4", addr4)
	if err != nil {
		return nil, err
	}
	l6, err := net.ListenUDP("udp6", addr6)
	if err != nil {
		return nil, err
	}

	return mdns.Server(ipv4.NewPacketConn(l4), ipv6.NewPacketConn(l6), &mdns.Config{
		LocalNames: []string{fmt.Sprintf("%s.%s", hostname, strings.TrimRight(params.Domain, "."))},
	})
}
