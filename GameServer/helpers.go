package main

import (
	"bytes"
	"encoding/json"
	"os"
)

func ConsoleLogJsonByteArray(message []byte) {
	// TODO: enable logging only for dev and an override
	var out bytes.Buffer
	json.Indent(&out, message, "", "  ")
	out.WriteTo(os.Stdout)
}
