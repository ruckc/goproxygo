package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	router := Router{}
	host := flag.String("host", "127.0.0.1", "which interface to bind to")
	port := flag.Int("port", 8080, "which port to bind to")

	flag.Parse()
	for _, arg := range flag.Args() {
		pieces := strings.SplitN(arg, ":", 2)
		router.Handle(pieces[0], pieces[1])
	}

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), &router)
	if err != nil {
		log.Fatal(err)
	}
}

type ProxyRoute struct {
	path  *string
	proxy *httputil.ReverseProxy
}

type Router struct {
	routes []*ProxyRoute
}

func (router *Router) Handle(path string, destination string) {
	destinationUrl, err := url.Parse(destination)
	if err != nil {
		log.Fatal(err)
	}
	handler := reverseProxyHandler(destinationUrl)

	log.Printf("Registering handler for %s to %s\n", path, destinationUrl)
	router.routes = append(router.routes, &ProxyRoute{&path, handler})
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range router.routes {
		matched := strings.HasPrefix(r.URL.Path, *route.path)
		//log.Printf("Checking if %s matches %s and the result is %t", r.URL.Path, route.path, matched)
		if matched {
			route.proxy.ServeHTTP(w, r)
			return
		}
	}
	log.Printf("Unhandled path: %s\n", r.URL.Path)
}

func reverseProxyHandler(destination *url.URL) *httputil.ReverseProxy {
	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(destination)
			r.SetXForwarded()
			r.Out.Host = r.In.Host
		},
	}

	return proxy
}
