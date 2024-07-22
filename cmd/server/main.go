package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"protobuf-from-scratch/decoders"
	"protobuf-from-scratch/types"
)

func main(){

	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var projectType types.ProjectType
		json.Unmarshal(body, &projectType)

		fmt.Printf("json size: %v bytes\n", len(body))
		fmt.Printf("json data: %+v\n", projectType)
	})

	http.HandleFunc("/proto", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		projectType, _ := decoders.DecodeProjectType(bytes.NewReader(body))

		fmt.Printf("proto size: %v bytes\n", len(body))
		fmt.Printf("proto data: %+v\n", projectType)
	})

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, os.Kill)
	go func(){
		if err := http.ListenAndServe(":3000", nil); err != nil {
			done<- nil
		}
	}()

	fmt.Println("Listening on port 3000")
	<-done
	fmt.Println("Shutting down the server")
}
