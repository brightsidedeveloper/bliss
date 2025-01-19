package server

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"time"
)

func readEnvPort(filePath string) (string, error) {
	// Open the .env file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines and comments
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		// Check if the line starts with "PORT="
		if strings.HasPrefix(line, "PORT=") {
			// Extract the value after "PORT="
			return strings.TrimSpace(line[len("PORT="):]), nil
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	return "", fmt.Errorf("PORT not found in %s", filePath)
}

func IsServerRunning(serverPath string) bool {
	port, err := readEnvPort(path.Join(serverPath, ".env"))
	if err != nil {
		fmt.Println(err)
		return false
	}

	address := fmt.Sprintf("localhost:%s", port)

	timeout := 1 * time.Second
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()

	return true
}

func StartServer(serverPath string) error {
	if IsServerRunning(serverPath) {
		fmt.Println("Server is already running")
		return nil
	}

	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = serverPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Detach the process
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("error starting server: %w", err)
	}

	return nil
}

func StopServer(serverPath string) error {
	port, err := readEnvPort(path.Join(serverPath, ".env"))
	if err != nil {
		return fmt.Errorf("failed to read port: %w", err)
	}

	address := fmt.Sprintf("localhost:%s", port)
	timeout := 1 * time.Second
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return nil
	}
	defer conn.Close()

	cmd := exec.Command("lsof", "-t", fmt.Sprintf("-i:%s", port))
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("lsof failed: %w", err)
	}

	pids := strings.Fields(string(output))
	if len(pids) == 0 {
		return fmt.Errorf("no processes found on port %s", port)
	}

	for _, pid := range pids {
		killCmd := exec.Command("kill", "-9", pid)
		err := killCmd.Run()
		if err != nil {
			return fmt.Errorf("failed to kill process %s: %w", pid, err)
		}
	}

	return nil
}
