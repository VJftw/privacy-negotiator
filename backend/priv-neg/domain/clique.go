package domain

import (
	"errors"

	"github.com/satori/go.uuid"
)

// CacheClique - Representation of a user's cached clique memberships. Stored as `<userID>:cliques`: cliqueID: {}
type CacheClique struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Categories     []string `json:"categories"`
	UserCategories []string `json:"userCategories"`
}

// WebClique - Representation of a user's clique returned to web
type WebClique struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Categories []string `json:"categories"`
}

// WSClique - Representation of Clique submitted via Websocket
type WSClique struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	UserIDs    []string `json:"users"`
	Categories []string `json:"categories"`
}

// DBClique - Representation of a Clique stored in the database
type DBClique struct {
	ID            string         `gorm:"primary_key"`
	DBUserCliques []DBUserClique `gorm:"ForeignKey:CliqueID"`
}

// TableName - Returns the table name for the entity.
func (c DBClique) TableName() string {
	return "cliques"
}

// GetUserIDs - returns the user ides for a DBClique
func (c *DBClique) GetUserIDs() []string {
	userIDs := []string{}

	for _, userClique := range c.DBUserCliques {
		userIDs = append(userIDs, userClique.UserID)
	}

	return userIDs
}

// DBUserClique - Represents a user belonging to a Clique
type DBUserClique struct {
	CliqueID   string       `gorm:"primary_key"`
	UserID     string       `gorm:"primary_key"`
	Name       string       `gorm:"type:varchar(100)"`
	Categories []DBCategory `gorm:"many2many:user_clique_categories"`
}

// TableName - Returns the table name for the entity.
func (uC DBUserClique) TableName() string {
	return "user_cliques"
}

// GetUserCliqueForUserID - returns the dbUserClique for a given user id
func (c *DBClique) GetUserCliqueForUserID(uID string) (*DBUserClique, error) {
	for _, userClique := range c.DBUserCliques {
		if userClique.UserID == uID {
			return &userClique, nil
		}
	}

	return nil, errors.New("UserID not found")
}

// NewDBClique - Returns a new DBClique with UUID
func NewDBClique() *DBClique {
	return &DBClique{
		ID:             uuid.NewV4().String(),
	}
}

// CacheCliqueFromDBClique - Returns a CacheClique from a DBClique
func CacheCliqueFromDBClique(dbClique *DBClique) *CacheClique {
	return &CacheClique{
		ID: dbClique.ID,
	}
}

// CacheCliqueFromDBUserClique - Returns a cache clique and user id from a DBUserClique
func CacheCliqueFromDBUserClique(dbUserClique *DBUserClique) (*CacheClique, string) {
	categories := []string{}
	userCategories := []string{}
	for _, dbCat := range dbUserClique.Categories {
		if dbCat.UserID == "none" {
			categories = append(categories, dbCat.Name)
		} else {
			userCategories = append(userCategories, dbCat.Name)
		}
	}
	return &CacheClique{
		ID:             dbUserClique.CliqueID,
		Name:           dbUserClique.Name,
		Categories:     categories,
		UserCategories: userCategories,
	}, dbUserClique.UserID
}

// DBUserCliqueFromCacheCliqueAndUserID - Returns a DBUserClique from CacheClique and UserID
func DBUserCliqueFromCacheCliqueAndUserID(cacheClique *CacheClique, userID string) *DBUserClique {
	categories := []DBCategory{}
	for _, cat := range cacheClique.Categories {
		categories = append(categories, DBCategory{Name: cat, UserID: "none"})
	}
	for _, cat := range cacheClique.UserCategories {
		categories = append(categories, DBCategory{Name: cat, UserID: userID})
	}
	return &DBUserClique{
		CliqueID:   cacheClique.ID,
		Name:       cacheClique.Name,
		UserID:     userID,
		Categories: categories,
	}
}

// WebCliqueFromCacheClique - Returns a WebClique from a CacheClique
func WebCliqueFromCacheClique(clique CacheClique) WebClique {
	return WebClique{
		ID:         clique.ID,
		Name:       clique.Name,
		Categories: append(clique.Categories, clique.UserCategories...),
	}
}
