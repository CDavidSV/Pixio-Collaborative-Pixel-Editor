package services

import "github.com/CDavidSV/Pixio/data"

type Services struct {
	AuthService *AuthService
}

func NewServices(queries *data.Queries) *Services {
	return &Services{
		AuthService: &AuthService{queries},
	}

}
