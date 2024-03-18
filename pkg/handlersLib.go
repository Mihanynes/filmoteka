package pkg

import (
	"log"
	"net/http"
)

func CloseBody(r *http.Request) {
	err := r.Body.Close()
	if err != nil {
		log.Println("error closing request body:", err)
	}
}
