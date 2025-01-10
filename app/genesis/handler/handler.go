package handler

import (
	"app/genesis/util"
	"database/sql"
)

type Handler struct {
	DB   *sql.DB
	JSON *util.JSON
}
