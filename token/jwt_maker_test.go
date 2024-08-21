package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/leedrum/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker("secret")
	require.Error(t, err)
	require.ErrorContains(t, err, "must be at least 32 characters")
	require.Equal(t, nil, maker)

	maker, err = NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute
	issuedAt := time.Now()
	role := util.DepositorRole

	token, payload, err := maker.CreateToken(username, role, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.Equal(t, username, payload.Username)
	require.Equal(t, role, payload.Role)
	require.NotZero(t, payload.ID)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, issuedAt.Add(duration), payload.ExpiredAt, time.Second)
}

func TestExpiredToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(util.RandomOwner(), util.DepositorRole, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.ErrorContains(t, err, jwt.ErrTokenExpired.Error())
	require.Nil(t, payload)
}

func TestInvalidToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err := maker.VerifyToken("invalid-token")
	require.Error(t, err)
	require.ErrorContains(t, err, jwt.ErrTokenMalformed.Error())
	require.Nil(t, payload)
}

func TestInvalidAglo(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(util.RandomOwner(), util.DepositorRole, time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	// Change the algorithm
	_, err = jwt.ParseWithClaims(token, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "signature is invalid")
}

func TestEmptyToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err := maker.VerifyToken("")
	require.Error(t, err)
	require.ErrorContains(t, err, "token contains an invalid number of segments")
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlgoNone(t *testing.T) {
	payload, err := NewPayload(util.RandomOwner(), util.DepositorRole, time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.ErrorContains(t, err, "token signature is invalid")
	require.Nil(t, payload)
}
