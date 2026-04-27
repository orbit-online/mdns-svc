# mdns-svc

A tiny cli wrapper for [hashicorps mDNS library](https://github.com/hashicorp/mdns).

## Usage

```
Usage: mdns-svc <port> [flags]

A tiny cli wrapper for hashicorps mDNS library

Arguments:
  <port>    The port the advertised service is running on

Flags:
  -h, --help                    Show context-sensitive help.
      --instance="$HOST"        The name of the specifi instance of this service
      --service="_http._tcp"    The service type to advertise
      --domain="local."         The domain to advertise the service under
      --host-name="$HOST."      The hostname of the advertied service
      --description=STRING      A text describing the advertised service
      --interface=STRING        Name of the interface to advertise on
      --verbose                 Turn on verbose logging
```
