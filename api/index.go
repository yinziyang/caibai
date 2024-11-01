package handler

import (
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/json":
		JsonHandler(w, r)
	case "/today":
		TodayHandler(w, r)
	case "/history":
		HistoryHandler(w, r)
	default:
		JsonHandler(w, r)
	}
}
