package gen

import (
	"math/rand"
)

const runes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func NewLicenseID() (id string) {
	buffer := make([]rune, 26)
	for i := range buffer {
		buffer[i] = []rune(runes)[rand.Intn(len(runes))]
	}
	id = string(buffer)
	return
}
