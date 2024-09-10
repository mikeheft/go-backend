package token

import "time"

// Maker is an interface for managing tokens
type Maker interface {
	// CreateToken create a new token for a specific duration and username
	CreateToken(username string, duration time.Duration) (string, error)
	// VerifyToken checker if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
