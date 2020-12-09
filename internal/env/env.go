package env

import (
	"log"
	"os"
)

func Init(k string) string {
	// k: key v: value
	v, ok := os.LookupEnv(k)
	// if unset env key:value, Fatal
	if !ok {
		log.Fatal("Environment value is not set, key: ", k)
	}
	return v
}
