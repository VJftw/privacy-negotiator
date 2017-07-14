package domain

import "time"

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

	DBUserCliques []DBUserClique `gorm:"ForeignKey:UserID;AssociationForeignKey:ID"`
}

// TableName - Returns the table name for the entity.
func (u DBUser) TableName() string {
	return "users"
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

// CacheUserProfile - Cache representation of a User Profile
type CacheUserProfile struct {
	Gender              string   `json:"gender"`
	AgeRange            string   `json:"ageRange"`
	Hometown            string   `json:"hometown"`
	Location            string   `json:"location"`
	Education           []string `json:"education"`
	FavouriteTeams      []string `json:"favouriteTeams"`
	InspirationalPeople []string `json:"inspirationalPeople"`
	Languages           []string `json:"languages"`
	Sports              []string `json:"sports"`
	Work                []string `json:"work"`
	Family              []string `json:"family"`
	Music               []string `json:"music"`
	Movies              []string `json:"movies"`
	Likes               []string `json:"likes"`
	Groups              []string `json:"groups"`
	Events              []string `json:"events"`
	Political           string   `json:"political"`
	Religion            string   `json:"religion"`
}
