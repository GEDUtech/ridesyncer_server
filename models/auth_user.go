package models

type AuthUser struct {
	User

	authenticated bool
}

func NewAuthUser(user User, authenticated bool) AuthUser {
	return AuthUser{user, authenticated}
}

func (authUser *AuthUser) SetAuthenticated(authenticated bool) {
	authUser.authenticated = authenticated
}

func (authUser *AuthUser) IsAuthenticated() bool {
	return authUser.authenticated
}
