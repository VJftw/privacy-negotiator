package user

import "time"

// FacebookUser - The FacebookUser entity.
type FacebookUser struct {
	FacebookUserID  string    `json:"userID" gorm:"primary_key"`
	LongLivedToken  string    `json:"longLivedToken"`
	ShortLivedToken string    `json:"accessToken" gorm:"-"`
	TokenExpires    time.Time `json:"tokenExpires"`
}
