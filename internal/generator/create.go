package generator

import (
	"master-gen/internal/clone"
	"os"
)

func createServerIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return clone.Repo(path, "https://github.com/brightsidedeveloper/bsd-solar-system.git")
	}
	return nil
}

func createWebIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return clone.Repo(path, "https://github.com/brightsidedeveloper/bsd-planet-web.git")
	}
	return nil
}
