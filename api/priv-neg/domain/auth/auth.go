package auth

// ApiAuth - Used to authenticate with this API.
type ApiAuth struct {
	Token string `json:"authToken"`
}

// FacebookAuth - Used to authenticate with the Facebook API.
type FacebookAuth struct {
	AccessToken string `json:"accessToken"`
	UserID      string `json:"userID"`
}
