package main

import (
	"net"
	"net/http"
	"os"

	"github.com/trojan-t/http/cmd/app"
	"github.com/trojan-t/http/pkg/banners"
)

func main() {
	host := "0.0.0.0"
	port := "1234"
	if err := execute(host, port); err != nil {
		os.Exit(1)
	}
}

func execute(host, port string) (err error) {
	mux := http.NewServeMux()
	bannersSvc := banners.NewService()
	serverHandler := app.NewServer(mux, bannersSvc)
	serverHandler.Init()

	srv := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: serverHandler,
	}
	return srv.ListenAndServe()
}
