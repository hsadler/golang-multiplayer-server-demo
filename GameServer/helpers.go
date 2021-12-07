package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

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

func RandFloat(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Float64()*(max-min)
}

func GenRandPosition(gs *GameState) Position {
	return Position{
		X: RandFloat(-float64(gs.MapWidth/2), float64(gs.MapWidth/2)),
		Y: RandFloat(-float64(gs.MapHeight/2), float64(gs.MapHeight/2)),
	}
}
