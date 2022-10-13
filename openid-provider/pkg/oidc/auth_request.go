package oidc

// auth request symbols the authorizate request
// this auth request is switched to auth code after validation of User
// In general this is the same type as auth code,
// but it has different storage for security and simplicity reasons
type AuthRequest AuthCode

func (r AuthRequest) GetID() ID {
	return r.ID
}
