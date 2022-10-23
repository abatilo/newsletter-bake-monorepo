package internal

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func ServeStaticAssets(staticAssetsPath string) (string, http.HandlerFunc) {
	var staticHandler http.Handler
	serveFromProxy := strings.HasPrefix(staticAssetsPath, "http://")

	if serveFromProxy {
		url, _ := url.Parse(staticAssetsPath)
		staticHandler = httputil.NewSingleHostReverseProxy(url)
	} else {
		staticHandler = http.FileServer(http.Dir(staticAssetsPath))
	}

	return "/", func(w http.ResponseWriter, r *http.Request) {
		if serveFromProxy {
			staticHandler.ServeHTTP(w, r)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/static") {
			w.Header().Set("Cache-Control", "public, max-age=31536000")
		} else {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		}

		filePath := filepath.Join(staticAssetsPath, r.URL.Path)
		// If the file doesn't exist, serve index.html so that react-router can
		// handle the 404.
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			indexPath := filepath.Join(staticAssetsPath, "index.html")
			index, _ := os.ReadFile(indexPath)
			w.Write(index)
			return
		}

		staticHandler.ServeHTTP(w, r)
	}
}
