package gist6197035

import (
	"os"
)

func GetenvOr(key, fallback string) string {
	s := os.Getenv(key)
	if s == "" {
		return fallback
	}
	return s
}

func main() {
	println(GetenvOr("PATH", "PATH was not set"))
	println(GetenvOr("Non-existant Env Var", "falling back to default"))
}
