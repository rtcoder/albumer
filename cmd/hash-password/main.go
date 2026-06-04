package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("read password: %v", err)
	}

	value := strings.TrimSuffix(string(password), "\n")
	if value == "" {
		log.Fatal("password cannot be empty")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("hash password: %v", err)
	}
	fmt.Println(string(hash))
}
