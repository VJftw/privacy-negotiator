package domain

import "github.com/satori/go.uuid"

type DBClique struct {
	ID string
}

func NewDBClique() *DBClique {
	return &DBClique{
		ID: uuid.NewV4().String(),
	}
}
