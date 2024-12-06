package main

import "net/http"

func main() {

	sMux := http.NewServeMux()

	server := http.Server{
		Handler: sMux,
		Addr:    ":8080",
	}

	server.ListenAndServe()

}
