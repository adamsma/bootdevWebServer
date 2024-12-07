package main

import "net/http"

func handlerHealth(resp http.ResponseWriter, req *http.Request) {

	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resp.WriteHeader(http.StatusOK)

	resp.Write([]byte(http.StatusText(http.StatusOK)))

}
