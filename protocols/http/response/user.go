package response

import "kc-ewallet/domains/repository/postgres"

type GetUserByIDResponse struct {
	ID       int32   `json:"id"`
	Username string  `json:"username"`
	Balance  float64 `json:"balance,omitempty"`
}

type LoginResponse struct {
	AccessToken string              `json:"access_token"`
	User        GetUserByIDResponse `json:"user"`
}

func NewGetUserByIDResponse(user postgres.User) GetUserByIDResponse {
	return GetUserByIDResponse{
		ID:       user.ID,
		Username: user.Username,
		Balance:  user.Balance,
	}
}

func NewLoginResponse(token string, user postgres.User) LoginResponse {
	return LoginResponse{
		AccessToken: token,
		User:        NewGetUserByIDResponse(user),
	}
}
