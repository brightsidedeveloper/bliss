package generator

import (
	"fmt"
	"master-gen/internal/clone"
	"os"
)

func createServerIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return clone.Repo(path, "https://github.com/brightsidedeveloper/bliss-server.git")
	}
	return nil
}

func createWebIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Cloning web")
		return clone.Repo(path, "https://github.com/brightsidedeveloper/bsd-planet-web.git")
	}
	return nil
}
