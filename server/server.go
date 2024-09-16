package main

import (
	"fehdle/routes/classic"
	"fehdle/routes/common"
	"fehdle/routes/shadows"
	"net/http"
)

func corsMiddleware(next http.Handler) http.Handler {
	var corsMiddleware = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// writer.Header().Add("Access-Control-Allow-Origin", os.Getenv("CORS_DOMAIN"))
		writer.Header().Add("Access-Control-Allow-Origin", "*")
		writer.Header().Add("Access-Control-Allow-Methods", "POST, GET")
		next.ServeHTTP(writer, request)
	})

	return corsMiddleware
}

func main() {
	mux := http.NewServeMux()

	go classic.UpdateGoroutine()
	go shadows.UpdateGoroutine()

	var guessClassic = http.HandlerFunc(classic.Guess)
	var getCurrentShadow = http.HandlerFunc(shadows.Route)
	var guessShadow = http.HandlerFunc(shadows.GuessShadow)
	var searchHero = http.HandlerFunc(common.SearchRoute)

	mux.Handle("POST /classic", corsMiddleware(guessClassic))
	mux.Handle("GET /shadow", corsMiddleware(getCurrentShadow))
	mux.Handle("POST /shadow", corsMiddleware(guessShadow))
	mux.Handle("GET /search", corsMiddleware(searchHero))

	http.ListenAndServe(":4444", mux)
}
