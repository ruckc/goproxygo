package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
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
	path  *regexp.Regexp
	proxy *httputil.ReverseProxy
}

type Router struct {
	routes []*ProxyRoute
}

func (router *Router) Handle(path string, destination string) {
	regexpPath := regexp.MustCompile(path)
	destinationUrl, err := url.Parse(destination)
	if err != nil {
		log.Fatal(err)
	}
	handler := reverseProxyHandler(destinationUrl)

	log.Printf("Registering handler for %s to %s\n", regexpPath, destinationUrl)
	router.routes = append(router.routes, &ProxyRoute{regexpPath, handler})
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range router.routes {
		matched := route.path.MatchString(r.URL.Path)
		//log.Printf("Checking if %s matches %s and the result is %t", r.URL.Path, route.path, matched)
		if matched {
			route.proxy.ServeHTTP(w, r)
			return
		}
	}
	log.Printf("Unhandled path: %s\n", r.URL.Path)
}

func reverseProxyHandler(destination *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(destination)

	proxy.Director = func(req *http.Request) {
		fmt.Println(req.Header)
		fmt.Println(req.Host)

		req.Header.Set("Origin", fmt.Sprintf("%s://%s", destination.Scheme, destination.Host))
		req.Host = destination.Host
		req.URL.Scheme = destination.Scheme
		req.URL.Host = destination.Host

		fmt.Println(req)
	}

	return proxy
}
