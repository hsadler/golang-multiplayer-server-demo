package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
)

// logging helpers

func LogData(logHeader string, logData interface{}) {
	if DO_LOGGING {
		fmt.Println("======================")
		fmt.Println(logHeader, logData)
		fmt.Println("======================")
	}
}

func LogDataForce(logHeader string, logData interface{}) {
	fmt.Println("======================")
	fmt.Println(logHeader, logData)
	fmt.Println("======================")
}

func LogJson(logHeader string, messageJson []byte) {
	if DO_LOGGING {
		fmt.Println("======================")
		fmt.Println(logHeader)
		var out bytes.Buffer
		json.Indent(&out, messageJson, "", "  ")
		out.WriteTo(os.Stdout)
		fmt.Println("======================")
	}
}

func LogJsonForce(logHeader string, messageJson []byte) {
	fmt.Println("======================")
	fmt.Println(logHeader)
	var out bytes.Buffer
	json.Indent(&out, messageJson, "", "  ")
	out.WriteTo(os.Stdout)
	fmt.Println("======================")
}

// utils

func GenUUID() string {
	return uuid.New().String()
}
