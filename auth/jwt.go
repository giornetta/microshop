package auth

import (
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"golang.org/x/exp/slices"
)

type jwtToken struct {
	token jwt.Token
}

func (t *jwtToken) Subject() string {
	return t.token.Subject()
}

func (t *jwtToken) Roles() []Role {
	claims := t.token.PrivateClaims()["roles"].([]any)

	roles := make([]Role, len(claims))
	for i, c := range claims {
		roles[i] = Role(c.(string))
	}

	return roles
}

func (t *jwtToken) IsAdmin() bool {
	return slices.Contains(t.Roles(), AdminRole)
}

type jwtIssuer struct {
	pemKey        bool
	name          string
	tokenDuration time.Duration
	privateKey    jwk.Key
}

func NewJWTIssuer(privateKey []byte, opts ...JWTIssuerOpt) (Issuer, error) {
	iss := &jwtIssuer{
		tokenDuration: time.Hour * 24,
		pemKey:        false,
	}

	for _, opt := range opts {
		opt(iss)
	}

	k, err := jwk.ParseKey(privateKey, jwk.WithPEM(iss.pemKey))
	if err != nil {
		return nil, fmt.Errorf("could not parse JWK: %v", err)
	}

	iss.privateKey = k

	return iss, nil
}

type JWTIssuerOpt func(*jwtIssuer)

func WithPEM(pemKey bool) JWTIssuerOpt {
	return func(i *jwtIssuer) {
		i.pemKey = pemKey
	}
}

func WithName(name string) JWTIssuerOpt {
	return func(i *jwtIssuer) {
		i.name = name
	}
}

func WithTokenDuration(duration time.Duration) JWTIssuerOpt {
	return func(i *jwtIssuer) {
		i.tokenDuration = duration
	}
}

func (i *jwtIssuer) Issue(subject string, roles []Role) (string, error) {
	timeNow := time.Now()

	// TODO Verify Roles?

	tok, err := jwt.NewBuilder().
		Issuer(i.name).
		IssuedAt(timeNow).
		Expiration(timeNow.Add(i.tokenDuration)).
		Subject(subject).
		Claim("roles", roles).
		Build()
	if err != nil {
		return "", fmt.Errorf("failed to build token: %v", err)
	}

	signed, err := jwt.Sign(tok, jwt.WithKey(i.privateKey.Algorithm(), i.privateKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return string(signed), nil
}

type jwtVerifier struct {
	publicKey     jwk.Key
	allowedIssuer string
}

func NewJWTVerifier(publicKey []byte, opts ...JWTVerifierOpt) (Verifier, error) {
	k, err := jwk.ParseKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("could not parse JWK: %v", err)
	}

	v := &jwtVerifier{
		publicKey: k,
	}

	for _, opt := range opts {
		opt(v)
	}

	return v, nil
}

type JWTVerifierOpt func(*jwtVerifier)

func WithAllowedIssuer(issuer string) JWTVerifierOpt {
	return func(v *jwtVerifier) {
		v.allowedIssuer = issuer
	}
}

func (v *jwtVerifier) Verify(signed string) (Token, error) {
	tok, err := jwt.Parse(
		[]byte(signed),
		jwt.WithKey(v.publicKey.Algorithm(), v.publicKey),
		jwt.WithValidate(true),
		jwt.WithIssuer(v.allowedIssuer),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to verify JWT: %v", err)
	}

	return &jwtToken{token: tok}, nil
}
