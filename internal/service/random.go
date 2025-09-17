package service

import "crypto/rand"

const (
	charSet         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	shortCodeLength = 6
	charSetLength   = len(charSet)
)

func generateShortCode() string {
	// https://stackoverflow.com/a/67035900
	buf := make([]byte, shortCodeLength)
	rand.Read(buf)
	for i := range shortCodeLength {
		buf[i] = charSet[int(buf[i])%charSetLength]
	}
	return string(buf)
}
