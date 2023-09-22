package main

import (
	"encoding/json"
	"os"
)

func main() {
	switch os.Args[1] {
	case "info":
		json.NewEncoder(os.Stdout).Encode(map[string]string{
			"short_description": "sample plugin short desc",
		})
	case "version":
		json.NewEncoder(os.Stdout).Encode(map[string]string{
			"sem_ver": "v0.0.1",
		})
	}
}
