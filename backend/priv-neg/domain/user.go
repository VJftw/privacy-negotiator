package domain

import "time"

// User - The FacebookUser entity.
type User interface{}

// CacheUser - a representation of a User stored in the cache.
type CacheUser struct {
	ID           string    `json:"id"`
	TokenExpires time.Time `json:"tokenExpires"`
}

// AuthUser - a representation of a User communicated from the web_app for authentication.
type AuthUser struct {
	ID              string `json:"id"`
	ShortLivedToken string `json:"accessToken"`
}

// DBUser - The representation of a User stored in the database.
type DBUser struct {
	ID             string `gorm:"primary_key"`
	LongLivedToken string
	TokenExpires   time.Time

	Categories []DBCategory

	DBUserCliques []DBUserClique `gorm:"ForeignKey:UserID;AssociationForeignKey:ID"`
}

// CacheUserFromAuthUser - Translates a AuthUser to a CacheUser.
func CacheUserFromAuthUser(authUser *AuthUser) *CacheUser {
	return &CacheUser{
		ID: authUser.ID,
	}
}

// CacheUserFromDatabaseUser - Translates a DatabaseUser to a CacheUser.
func CacheUserFromDatabaseUser(user *DBUser) *CacheUser {
	cU := CacheUser{
		ID:           user.ID,
		TokenExpires: user.TokenExpires,
	}

	return &cU
}

// DBUserFromAuthUser - Translates an AuthUser to a DBUser.
func DBUserFromAuthUser(authUser *AuthUser) *DBUser {
	return &DBUser{
		ID: authUser.ID,
	}
}
