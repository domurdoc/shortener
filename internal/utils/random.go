package utils

import "crypto/rand"

const ALPHA = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateRandomString(charSet string, length int) (string, error) {
	// https://stackoverflow.com/a/67035900
	charSetLength := len(charSet)
	buf := make([]byte, length)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	for i := range length {
		buf[i] = charSet[int(buf[i])%charSetLength]
	}
	return string(buf), nil
}

func MustGenerateRandomString(charSet string, length int) string {
	s, err := GenerateRandomString(charSet, length)
	if err != nil {
		panic(err)
	}
	return s
}
