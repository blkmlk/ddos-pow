package env

import (
	"fmt"
	"os"
)

const (
	Host = "HOST"
)

func NewErrNotSet(env string) error {
	return fmt.Errorf("env %s isn't set", env)
}

func Get(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", NewErrNotSet(key)
	}
	return value, nil
}
