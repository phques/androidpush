package main

import (
	"fmt"
	"math/rand"
	"time"
)

const base62Digits = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// Convert int64 to a string of 'base62' digits
func ToBase62(val int64) string {
	var str []byte
	lenDigits := int64(len(base62Digits))
	atleastonce := true

	for val > 0 || atleastonce {
		atleastonce = false
		digitVal := val % lenDigits
		val = val / lenDigits
		str = append(str, byte(base62Digits[digitVal]))
	}

	return string(str)
}

// Generate a random 'base62' Id string
func GenBase62Id() string {
	rnd := rand.Int63()
	fmt.Println(rnd)
	return ToBase62(rnd)
}

func main() {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 100; i++ {
		fmt.Println(i, ToBase62(int64(i)))
	}

	for i := 0; i < 5; i++ {
		id := GenBase62Id()
		fmt.Println(id, "\n")
	}
}
