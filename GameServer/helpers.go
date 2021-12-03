package main

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/google/uuid"
)

func ConsoleLogJsonByteArray(message []byte) {
	// TODO: enable logging only for dev and an override
	var out bytes.Buffer
	json.Indent(&out, message, "", "  ")
	out.WriteTo(os.Stdout)
}

func NewUUID() string {
	return uuid.New().String()
}
