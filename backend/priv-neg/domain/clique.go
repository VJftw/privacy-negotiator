package domain

import "github.com/satori/go.uuid"

// CacheClique - Representation of a user's cached clique memberships. Stored as `<userID>:cliques`: cliqueID: {}
type CacheClique struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type WebClique struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	UserIDs []string `json:"users"`
}

// DBClique - Representation of a Clique stored in the database
type DBClique struct {
	ID            string         `gorm:"primary_key"`
	DBUserCliques []DBUserClique `gorm:"ForeignKey:CliqueID"`
}

type DBUserClique struct {
	CliqueID string `gorm:"primary_key"`
	UserID   string `gorm:"primary_key"`
	Name     string
}

func (c *DBClique) GetUserIDs() []string {
	userIDs := []string{}

	for _, userClique := range c.DBUserCliques {
		userIDs = append(userIDs, userClique.UserID)
	}

	return userIDs
}

// NewCacheClique - Returns a new CacheClique with UUID
func NewCacheClique() *CacheClique {
	return &CacheClique{
		ID: uuid.NewV4().String(),
	}
}

func DBCliqueFromCacheClique(cacheClique *CacheClique) *DBClique {
	return &DBClique{
		ID: cacheClique.ID,
	}
}

func DBUserCliqueFromCacheCliqueAndUserID(cacheClique *CacheClique, userID string) *DBUserClique {
	return &DBUserClique{
		CliqueID: cacheClique.ID,
		Name:     cacheClique.Name,
		UserID:   userID,
	}
}
