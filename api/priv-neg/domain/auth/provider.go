package auth

// type Provider interface {
// 	NewFromUser(*user.User) *Auth
// }
//
// type AuthProvider struct {
// }
//
// func NewProvider() Provider {
// 	return &AuthProvider{}
// }
//
// func (p AuthProvider) NewFromUser(u *user.User) *Auth {
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"fb_user_id": u.UUID,
// 		"nbf":        time.Now().Unix(),
// 	})
// 	tokenString, _ := token.SignedString([]byte("hmacSecret"))
//
// 	return &Auth{
// 		Token: tokenString,
// 	}
//
// }
