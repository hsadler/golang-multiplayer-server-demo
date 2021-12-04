package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
)

func LogJson(logHeader string, messageJson []byte) {
	if DO_LOGGING {
		fmt.Println(logHeader)
		var out bytes.Buffer
		json.Indent(&out, messageJson, "", "  ")
		out.WriteTo(os.Stdout)
		fmt.Println("======================")
	}
}

func LogJsonForce(logHeader string, messageJson []byte) {
	fmt.Println(logHeader)
	var out bytes.Buffer
	json.Indent(&out, messageJson, "", "  ")
	out.WriteTo(os.Stdout)
	fmt.Println("======================")
}

func GenUUID() string {
	return uuid.New().String()
}
