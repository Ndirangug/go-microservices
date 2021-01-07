package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		log.Println("Hello world /")
		d, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(rw, "OOps!", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(rw, "Hello %s", d)

	})

	http.HandleFunc("/george", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello george")
	})

	http.ListenAndServe("localhost:9090", nil)
}
