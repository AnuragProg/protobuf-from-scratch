package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"protobuf-from-scratch/decoders"
	"protobuf-from-scratch/encoders"
	"protobuf-from-scratch/types"
)

func main(){

	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		body, _ := io.ReadAll(r.Body)

		// parsing the project type
		var projectType types.ProjectType
		json.Unmarshal(body, &projectType)

		// mutating data 
		projectType.Description = "changed"

		// responding with mutated data
		jsonBody, _ := json.Marshal(projectType)
		w.Write(jsonBody)
	})

	http.HandleFunc("/proto", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		// parsing the project type
		projectType, _ := decoders.DecodeProjectType(r.Body)

		// mutating data
		projectType.Description = "changed"

		// responding with mutated data
		protoBody := encoders.EncodeProjectType(projectType)
		w.Write(protoBody)
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
