package authservice

import "github.com/orekhovskiy/shrtn/config"

type AuthService struct {
	options config.Config
}

func NewService(opts config.Config) *AuthService {
	return &AuthService{
		options: opts,
	}
}
