package gen

import "math/rand"

const runes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var (
	buffer = make([]rune, 26)
)

func NewLicenseID() (id string) {
	for i := range buffer {
		buffer[i] = []rune(runes)[rand.Intn(len(runes))]
	}
	id = string(buffer)
	return
}
