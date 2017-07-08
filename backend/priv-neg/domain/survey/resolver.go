package survey

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/domain"
)

// FromRequest - Returns a Photo from a http request.
func FromRequest(r *http.Request, user *domain.CacheUser) (*domain.DBSurvey, error) {
	s := &domain.DBSurvey{}

	requestSurvey := &requestSurvey{}
	err := json.NewDecoder(r.Body).Decode(requestSurvey)
	if err != nil {
		return nil, err
	}

	if requestSurvey.Type == "" {
		return nil, errors.New("Missing type")
	}

	if requestSurvey.Data == nil {
		return nil, errors.New("Missing data")
	}

	s.UserID = user.ID
	if requestSurvey.Type == "photo" {
		s.PhotoID = requestSurvey.PhotoID
	} else if requestSurvey.Type == "general" {
		s.PhotoID = "general"
	}

	b, err := json.Marshal(requestSurvey.Data)
	if err != nil {
		return nil, err
	}
	s.RawJSON = string(b)

	return s, nil
}

type requestSurvey struct {
	Type    string                 `json:"type"`
	PhotoID string                 `json:"photoID"`
	Data    map[string]interface{} `json:"data"`
}
