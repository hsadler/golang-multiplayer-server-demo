package main

import (
	"bytes"
	"encoding/json"
	"os"
)

func ConsoleLogJsonByteArray(message []byte) {
	var out bytes.Buffer
	json.Indent(&out, message, "", "  ")
	out.WriteTo(os.Stdout)
}
