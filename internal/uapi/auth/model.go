package auth

type authorizeAndCompleteRegistrationApiResponse struct {
	TimeWait                     int    `json:"timeWait"`
	AuthenticationRequestId      string `json:"authenticationRequestId"`
	AuthenticationRequestExpired bool   `json:"authenticationRequestExpired"`
	AccessToken                  string `json:"accessToken"`
	RefreshToken                 string `json:"refreshToken"`
	ExpiresIn                    int64  `json:"expiresIn"`
}

type AuthorizeAndCompleteRegistrationResponse struct {
	AuthenticationRequestId string `json:"authenticationRequestId"`
	AccessToken             string `json:"accessToken"`
	RefreshToken            string `json:"refreshToken"`
	ExpiresIn               int64  `json:"expiresIn"`
}

type UserInfo struct {
	Email string `json:"email"`
	Id    int    `json:"id"`
	Phone string `json:"phone_number"`
}

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type ConfirmEmailRequest struct {
	Code string `json:"code" validate:"required"`
}
