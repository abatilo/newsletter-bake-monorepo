package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/abatilo/newsletter-bake-monorepo/internal"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	pflag.String(
		"static-assets-path",
		"http://localhost:3000",
		"Where to find the static assets to serve",
	)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	staticAssetsPath := viper.GetString("static-assets-path")

	mux := http.NewServeMux()
	mux.HandleFunc(internal.ServeStaticAssets(staticAssetsPath))
	mux.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Listening on :8080")
	httpServer.ListenAndServe()
}
