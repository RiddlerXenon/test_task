package main

import (
	"net/http"

	"github.com/RiddlerXenon/test_task/internal/cache"
	"github.com/RiddlerXenon/test_task/internal/handler"
	"go.uber.org/zap"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
}

func main() {
	c, err := cache.New()
	if err != nil {
		zap.S().Fatal(err)
	}

	http.HandleFunc("/api/add", handler.AddSetHandler(c.Add))
	http.HandleFunc("/api/set", handler.AddSetHandler(c.Set))
	http.HandleFunc("/api/get", handler.GetHandler(c.Get))
	http.HandleFunc("/api/del", handler.DelHandler(c.Del))

	zap.S().Info("Server starting at http://127.0.0.1:8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		zap.S().Fatal(err)
	}
}
