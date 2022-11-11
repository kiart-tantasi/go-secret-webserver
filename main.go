package main

import (
	"fmt"
	"net/http"
	"os"
)

func handlerHome(writer http.ResponseWriter, request *http.Request) {
	secret := os.Getenv("SECRET")
	fmt.Fprintf(writer, "SECRET: %v", secret)
}

func main() {
	http.HandleFunc("/", handlerHome)
	http.ListenAndServe(":8080", nil)
}
