package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	// initialize a new servemux
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/todo/view", app.todoView) // fixed path

	return mux
}
