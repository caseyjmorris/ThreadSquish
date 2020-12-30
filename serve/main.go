package serve

import (
	"log"
	"net/http"
)

func main() {
	//http.HandleFunc("/", sayhelloName) // setting router rule
	//http.HandleFunc("/login", login)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
