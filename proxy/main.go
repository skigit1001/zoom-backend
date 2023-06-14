package main

import (
	"log"
	"net/http"
	"strings"

	"go-fiber-auth/utilities"

	"github.com/elazarl/goproxy"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = false

	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			proxyAuthToken := r.Header.Get("Proxy-Authorization")

			if proxyAuthToken == "" {
				return r, goproxy.NewResponse(r,
					goproxy.ContentTypeText, http.StatusForbidden,
					"Don't waste your time!")
			}

			trimmedToken := strings.TrimSpace(proxyAuthToken)

			if trimmedToken == "" {
				return r, goproxy.NewResponse(r,
					goproxy.ContentTypeText, http.StatusForbidden,
					"Don't waste your time!")
			}

			claims, parsingError := utilities.ParseClaims(trimmedToken)
			if parsingError != nil {
				return r, goproxy.NewResponse(r,
					goproxy.ContentTypeText, http.StatusForbidden,
					"Not authorized!")
			}

			log.Println(claims.UserId)

			return r, nil
		})

	http.ListenAndServe(":8080", proxy)
}
