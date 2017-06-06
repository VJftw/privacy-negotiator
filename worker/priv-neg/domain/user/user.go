package user

import "time"

// FacebookUser - The FacebookUser entity.
type FacebookUser struct {
	FacebookUserID  string    `json:"userId" gorm:"primary_key"`
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`
	LongLivedToken  string    `json:"longLivedToken"`
	ShortLivedToken string    `json:"accessToken" gorm:"-"`
	TokenExpires    time.Time `json:"tokenExpires"`
}
