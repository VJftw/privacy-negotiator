package user

import (
	"time"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain/category"
)

// FacebookUser - The FacebookUser entity.
type FacebookUser struct {
	FacebookUserID  string              `json:"userID" gorm:"primary_key"`
	LongLivedToken  string              `json:"longLivedToken"`
	ShortLivedToken string              `json:"accessToken" gorm:"-"`
	TokenExpires    time.Time           `json:"tokenExpires"`
	Categories      []category.Category `json:"-"`
}
