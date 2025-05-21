package common

import (
	"encoding/json"

	"math/rand"
)

const (
	PASSWORD_LENGTH = 6
	letterBytes     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func AnyToAnyStructField(from, to any) any {
	inrec, _ := json.Marshal(from)
	json.Unmarshal(inrec, &to)
	return to
}

func GenerateRandomPassword() string {
	b := make([]byte, PASSWORD_LENGTH)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
