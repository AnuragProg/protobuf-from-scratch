package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"protobuf-from-scratch/encoders"
	"protobuf-from-scratch/types"
	"time"
)

func main() {
	url := "http://localhost:3000"

	data := types.ProjectType{
		Name:        "John Doe",
		Description: "Lorum ipsum",
		Timestamp:   uint64(time.Now().Unix()),
		Tags:        []string{"tag1", "tag2", "tag3"},
	}


	fmt.Println("Requests started...")
	jsonBytes, _ := json.Marshal(data)
	jsonStream := bytes.NewReader(jsonBytes)
	protoStream := encoders.EncodeProjectType(data)
	http.Post(url+"/json", "application/json", jsonStream)
	http.Post(url+"/proto", "application/octet-stream", protoStream)

	fmt.Println("Requests completed")
}
