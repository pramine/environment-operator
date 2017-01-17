package main

import (
	"fmt"

	"github.com/pearsontechnology/environment-operator/version"
)

// This package adds environment-validator binary, which can be used to
// validate environments.bitesize file
func main() {
	fmt.Println(version.Version)
}
