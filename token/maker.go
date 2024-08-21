package token

import "time"

type Maker interface {
	// Create a new token for a specific username and duration
	CreateToken(username string, role string, duration time.Duration) (string, *Payload, error)

	// VerifyToken check if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
