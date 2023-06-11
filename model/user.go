package model

import "database/sql"

type User struct {
	Id           sql.NullInt64  `json:"id"`
	Password     sql.NullString `json:"password"`
	Phoneoremail sql.NullString `json:"phoneoremail"`
}
