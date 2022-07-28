package datadog

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func DDTagsFromHTTPRequest(request *http.Request) []ddtrace.StartSpanOption {
	opts := []ddtrace.StartSpanOption{}
	opts = append(opts, httpClientTagsFromHTTPRequest(request)...)
	opts = append(opts, basicHTTPTagsFromHTTPRequest(request)...)
	opts = append(opts, netTagsFromHTTPRequest(request)...)
	return opts
}

func httpClientTagsFromHTTPRequest(request *http.Request) []ddtrace.StartSpanOption {
	opts := []ddtrace.StartSpanOption{
		tracer.Tag(ext.HTTPMethod, request.Method),
		tracer.Tag(ext.HTTPURL, request.URL.Path),
		tracer.Tag("http.url.qs", request.URL.RawQuery),
	}

	if ua := request.UserAgent(); ua != "" {
		opts = append(opts, tracer.Tag("http.user_agent", ua))
	}
	if request.ContentLength > 0 {
		opts = append(opts, tracer.Tag("http.request_content_length", request.ContentLength))
	}

	return opts
}

func basicHTTPTagsFromHTTPRequest(request *http.Request) []ddtrace.StartSpanOption {
	opts := []ddtrace.StartSpanOption{}

	if request.TLS != nil {
		opts = append(opts, tracer.Tag("http.scheme", "https"))
	} else {
		opts = append(opts, tracer.Tag("http.scheme", "http"))
	}

	if request.Host != "" {
		opts = append(opts, tracer.Tag("http.host", request.Host))
	} else if request.URL != nil && request.URL.Host != "" {
		opts = append(opts, tracer.Tag("http.host", request.URL.Host))
	}

	flavor := ""
	if request.ProtoMajor == 1 {
		flavor = fmt.Sprintf("1.%d", request.ProtoMinor)
	} else if request.ProtoMajor == 2 {
		flavor = "2"
	}
	if flavor != "" {
		opts = append(opts, tracer.Tag("http.flavor", flavor))
	}

	if values, ok := request.Header["X-Forwarded-For"]; ok && len(values) > 0 {
		if addresses := strings.SplitN(values[0], ",", 2); len(addresses) > 0 {
			opts = append(opts, tracer.Tag("http.client_ip", addresses[0]))
		}
	}
	return opts
}

func netTagsFromHTTPRequest(request *http.Request) []ddtrace.StartSpanOption {
	opts := []ddtrace.StartSpanOption{}

	peerIP, peerName, peerPort := hostIPNamePort(request.RemoteAddr)
	if peerIP != "" {
		opts = append(opts, tracer.Tag("net.peer.ip", peerIP))
	}
	if peerName != "" {
		opts = append(opts, tracer.Tag("net.peer.name", peerName))
	}
	if peerPort != 0 {
		opts = append(opts, tracer.Tag("net.peer.port", peerPort))
	}

	hostIP, hostName, hostPort := "", "", 0
	for _, someHost := range []string{request.Host, request.Header.Get("Host"), request.URL.Host} {
		hostIP, hostName, hostPort = hostIPNamePort(someHost)
		if hostIP != "" || hostName != "" || hostPort != 0 {
			break
		}
	}
	if hostIP != "" {
		opts = append(opts, tracer.Tag("net.host.ip", hostIP))
	}
	if hostName != "" {
		opts = append(opts, tracer.Tag("net.host.name", hostName))
	}
	if hostPort != 0 {
		opts = append(opts, tracer.Tag("net.host.port", hostPort))
	}

	return opts
}

func hostIPNamePort(hostWithPort string) (ip string, name string, port int) {
	var (
		hostPart, portPart string
		parsedPort         uint64
		err                error
	)
	if hostPart, portPart, err = net.SplitHostPort(hostWithPort); err != nil {
		hostPart, portPart = hostWithPort, ""
	}
	if parsedIP := net.ParseIP(hostPart); parsedIP != nil {
		ip = parsedIP.String()
	} else {
		name = hostPart
	}
	if parsedPort, err = strconv.ParseUint(portPart, 10, 16); err == nil {
		port = int(parsedPort)
	}
	return
}
