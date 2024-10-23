package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

type ProxyConfig struct {
	Port    int64   `json:"port"`
	Proxies []Proxy `json:"proxies"`
}

type Proxy struct {
	Path      string            `json:"path"`
	RemoteURL string            `json:"remote_url"`
	Headers   map[string]string `json:"headers"`
}

func main() {
	file, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	config := new(ProxyConfig)
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	for _, proxy := range config.Proxies {
		handler := func(c *gin.Context) {
			remote, err := url.Parse(proxy.RemoteURL)
			if err != nil {
				panic(err)
			}

			reverseProxy := httputil.NewSingleHostReverseProxy(remote)
			reverseProxy.Director = func(req *http.Request) {
				req.Header = c.Request.Header
				req.Host = remote.Host
				req.URL.Scheme = remote.Scheme
				req.URL.Host = remote.Host
				for key, value := range proxy.Headers {
					req.Header.Set(key, value)
				}

				log.Printf("Forwarding request from %s to %s\n", req.RemoteAddr, req.URL)
			}

			reverseProxy.ServeHTTP(c.Writer, c.Request)
		}

		r.Any(fmt.Sprintf("%s/*proxyPath", proxy.Path), handler)
	}

	err = r.Run(fmt.Sprintf(":%d", config.Port))
	if err != nil {
		panic(err)
	}
}
