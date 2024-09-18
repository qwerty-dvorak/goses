package main

import (
    "os"
    "log"
)

func real() string{
    data, err := os.ReadFile("pass.txt")
    if err != nil {
        log.Fatal(err)
    }
    content := string(data)
	return content
}