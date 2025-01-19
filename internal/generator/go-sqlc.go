package generator

import (
	"os"
	"os/exec"
)

func genSqlc() error {
	cmd := exec.Command("sqlc", "generate")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
