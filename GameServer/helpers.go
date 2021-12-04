package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
)

func ConsoleLogJsonByteArray(logHeader string, message []byte) {
	if DO_LOGGING {
		fmt.Println(logHeader)
		var out bytes.Buffer
		json.Indent(&out, message, "", "  ")
		out.WriteTo(os.Stdout)
	}
}

func NewUUID() string {
	return uuid.New().String()
}
