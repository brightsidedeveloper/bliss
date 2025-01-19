package generator

import "os/exec"

func genSqlc() error {
	return exec.Command("sqlc", "generate").Run()
}
