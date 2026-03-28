package util

import "math/rand"

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

func GetRandomString(l int) string {
	b := make([]byte, l)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
