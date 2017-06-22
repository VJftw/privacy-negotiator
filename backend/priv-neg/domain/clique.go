package domain

import "github.com/satori/go.uuid"

// CacheClique - Representation of a user's cached clique memberships. Stored as `<userID>:cliques`: cliqueID: {}
type CacheClique struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// DBClique - Representation of a Clique stored in the database
type DBClique struct {
	ID string
}

// NewCacheClique - Returns a new CacheClique with UUID
func NewCacheClique() *CacheClique {
	return &CacheClique{
		ID: uuid.NewV4().String(),
	}
}

// NewDBClique - Returns a new DBClique with UUID
func NewDBClique() *DBClique {
	return &DBClique{
		ID: uuid.NewV4().String(),
	}
}
