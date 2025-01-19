package docker

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func CheckDockerDB() bool {
	cmd := exec.Command("docker", "ps")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error executing docker ps: %v\n", err)
		return false
	}

	// Read the output of `docker ps`
	output := out.String()
	lines := strings.Split(output, "\n")

	// Check if any line contains the container name "db"
	containerFound := false
	for _, line := range lines {
		if strings.Contains(line, "db") {
			containerFound = true
			break
		}
	}
	return containerFound
}
