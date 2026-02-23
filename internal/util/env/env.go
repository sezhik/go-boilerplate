package env

import (
	"fmt"
	"os"
)

func MustGet(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("env val '%s' is absent", key))
	}

	return v
}
