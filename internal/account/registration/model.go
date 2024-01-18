package registration

// type RegistrationRequest struct {
// 	Name     string `json:"name" validate:"required"`
// 	Email    string `json:"email" validate:"required"`
// 	Password string `json:"password" validate:"required"`
// }
// type RegistrationResult struct {
// 	Token     string               `json:"token"`
// 	Succeeded bool                 `json:"succeeded"`
// 	Claims    *account.TokenClaims `json:"claims"`
// }

// func Failed() *RegistrationResult {
// 	return &RegistrationResult{
// 		Token:     "",
// 		Succeeded: false,
// 		Claims:    nil,
// 	}
// }
// func Succeeded(token string, claims *account.TokenClaims) *RegistrationResult {
// 	return &RegistrationResult{
// 		Token:     token,
// 		Succeeded: true,
// 		Claims:    claims,
// 	}
// }
