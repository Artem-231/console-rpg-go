package server

import (
	"log"
	"net/http"
	"time"
)

// LoggerMiddleware логирует каждое действие клиента
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("Адрес: %s Метод: %s Путь: %s", r.RemoteAddr, r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		t := time.Since(start)

		log.Printf("Время: %v\n", t)
	})
}
