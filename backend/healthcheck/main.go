package main

import (
	"net/http"
	"os"
)

const URL = "http://127.0.0.1:8000/api/v1/health"

func main() {
	if _, err := http.Get(URL); err != nil {
		os.Exit(1)
	}
}
