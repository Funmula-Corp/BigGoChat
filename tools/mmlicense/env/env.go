package env

import (
	"log"
	"os"
	"strconv"
)

func Bool(key string) (val *bool) {
	val = new(bool)
	if envVal := os.Getenv(key); envVal != "" {
		if value, err := strconv.ParseBool(envVal); err != nil {
			log.Println("[ERROR]:", key, "->", err)
		} else {
			*val = value
		}
	}
	return
}

func Int(key string) (val *int) {
	val = new(int)
	if envVal := os.Getenv(key); envVal != "" {
		if value, err := strconv.ParseInt(envVal, 10, 32); err != nil {
			log.Println("[ERROR]:", key, "->", err)
		} else {
			*val = int(value)
		}
	}
	return
}

func String(key string) (val *string) {
	val = new(string)
	if envVal := os.Getenv(key); envVal != "" {
		*val = envVal
	}
	return
}
