package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"go.m3o.com/helloworld"
)

type Request struct {
	Name string `json:"name"`
}

type Response struct {
	Message string `json:"message"`
}

func main() {
	// token := os.Getenv("M3O_API_TOKEN")
	// if len(token) == 0 {
	// 	fmt.Println("Missing M3O_API_TOKEN")
	// 	return
	// }

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		message := "Hello world!"

		ct := r.Header.Get("Content-Type")

		if ct == "application/json" {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(500)
				return
			}

			var req Request

			if err := json.Unmarshal(b, &req); err != nil {
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))
				return
			}

			// fmt.Println(rsp, err)
			// if len(req.Name) > 0 {

			// 	message = "Hello " + req.Name + "!"
			// }

			message = call(req.Name)

			if err := json.NewEncoder(w).Encode(Response{Message: message}); err != nil {
				w.WriteHeader(500)
				return
			}

			return
		}

		r.ParseForm()
		name := r.Form.Get("name")

		if len(name) <= 0 {
			name = "Default"
		}

		message = call(name)

		w.Write([]byte(message))
	})

	http.ListenAndServe(":8080", nil)
}

func call(name string) string {
	helloworldService := helloworld.NewHelloworldService(os.Getenv("M3O_API_TOKEN"))
	rsp, err := helloworldService.Call(&helloworld.CallRequest{
		Name: name,
	})

	if err != nil {
		return err.Error()
	}

	return rsp.Message
}
