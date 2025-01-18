package generator

import "master-gen/internal/clone"

func createServerIfNotExists(path string) error {
	return clone.Repo(path, "https://github.com/brightsidedeveloper/bsd-solar-system.git")
}
