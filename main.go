package main

import (
	"log"
	"net/http"

	"github.com/elazarl/goproxy"

	"github.com/peti2001/csrf_changer/handler"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	handler := handler.NewCookieSwitcherHandler([]string{"laravel_session", "XSRF-TOKEN"})

	proxy.OnRequest().DoFunc(handler.RequestHandler)
	proxy.OnResponse().DoFunc(handler.ResponseHandler)

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", proxy))
}
