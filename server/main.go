package main

import (
	"fmt"
	"net/http"
	"strconv"
)

var request_counter = 0

func helloHandler(w http.ResponseWriter, r *http.Request) {
	size_param := r.URL.Query()["size"]
	var size int64 = 1
	if len(size_param) > 0 {
		size_parsed, err := strconv.ParseInt(size_param[0], 10, 64)
		if err == nil {
			size = size_parsed
		}
	}

	request_counter += 1
	fmt.Println("Received request " + strconv.Itoa(request_counter) + " from " + r.Host + " with " + fmt.Sprintf("%d", size))

	var i int64 = 0
	for ; i < size; i++ {
		fmt.Fprint(w, "Hello, World!")
	}
}

func main() {
	http.HandleFunc("/", helloHandler)
	http.ListenAndServe(":8080", nil)
}
