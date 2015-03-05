package main

import (
	"crypto/sha256"
	"fmt"
	"os"
)

func FileExists(f string) bool {

	if _, err := os.Stat(f); os.IsNotExist(err) {

		return false
	}

	return true
}

func HashString(s string) string {

	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}
