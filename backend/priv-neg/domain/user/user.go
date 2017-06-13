package user

import (
	"time"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/category"
)

// User - The FacebookUser entity.
type User interface{}

// CacheUser - a representation of a User stored in the cache.
type CacheUser struct {
	ID           string    `json:"id"`
	TokenExpires time.Time `json:"tokenExpires"`
	Categories   []string  `json:"categories"`
}

// WebUser - a representation of a User communicated with the web_app.
type WebUser struct {
	ID              string `json:"id"`
	ShortLivedToken string `json:"accessToken"`
}

// DatabaseUser - The representation of a User stored in the database.
type DatabaseUser struct {
	ID             string `gorm:"primary_key"`
	LongLivedToken string
	TokenExpires   time.Time

	Categories []category.Category
}

func CacheUserFromWebUser(webUser *WebUser) *CacheUser {
	return &CacheUser{
		ID: webUser.ID,
	}
}

func CacheUserFromDatabaseUser(user *DatabaseUser) *CacheUser {
	cU := CacheUser{
		ID:           user.ID,
		TokenExpires: user.TokenExpires,
	}
	for _, category := range user.Categories {
		cU.Categories = append(cU.Categories, category.Name)
	}

	return &cU
}
