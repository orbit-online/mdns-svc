package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/hashicorp/mdns"
)

type Params struct {
	Instance    string   `default:"$HOST" help:"The name of the specifi instance of this service"`
	Service     string   `default:"_http._tcp" help:"The service type to advertise"`
	Domain      string   `default:"local." help:"The domain to advertise the service under"`
	Hostname    string   `default:"$HOST." help:"The hostname of the advertied service"`
	Port        int      `required:"" arg:"" help:"The port the advertised service is running on"`
	Description string   `help:"A text describing the advertised service"`
	Interface   string   `help:"Name of the interface to advertise on"`
	IPs         []string `name:"ip" help:"IP addresses to advertise"`
	Verbose     bool     `help:"Turn on verbose logging"`
}

var params Params

func main() {
	kong.Parse(&params, kong.Name("mdns-svc"), kong.Description("A tiny cli wrapper for hashicorps mDNS library"))
	slog.SetDefault(slog.Default())
	if params.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
	handleError := func(err error) {
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)
	signal.Notify(sigs, syscall.SIGINT)
	var err error
	var instance string
	var hostName string
	if params.Instance == "$HOST" {
		instance, err = os.Hostname()
		handleError(err)
	} else {
		instance = params.Instance
	}
	if params.Hostname == "$HOST." {
		hostName, err = os.Hostname()
		handleError(err)
		hostName = fmt.Sprintf("%s.", hostName)
	} else {
		hostName = params.Hostname
	}
	var ips []net.IP = nil
	if len(params.IPs) > 0 {
		ips = []net.IP{}
		for _, ip := range params.IPs {
			ips = append(ips, net.ParseIP(ip))
		}
	}
	var iface *net.Interface
	if params.Interface != "" {
		iface, err = net.InterfaceByName(params.Interface)
		handleError(err)
	}
	service, err := mdns.NewMDNSService(instance, params.Service, params.Domain, hostName, params.Port, ips, []string{params.Description})
	handleError(err)
	slog.Debug(fmt.Sprintf("Broadcasting %s port %d on %s%s", params.Service, params.Port, hostName, params.Domain))
	server, err := mdns.NewServer(&mdns.Config{Zone: service, Iface: iface})
	handleError(err)
	defer server.Shutdown()
	<-sigs
}
