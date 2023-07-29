package auth_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/giornetta/microshop/auth"
)

const testPrivateKey string = `
	{
		"p": "5IyvyKLZG40Pv5ZoxVhGSXas9dVa6OLIVuLHbx9qoUU",
		"kty": "RSA",
		"q": "4Pih6MSQm9SYJPfDT98PC9HoF4PskeLcW9-wx0UVh7s",
		"d": "vZ2yLYZrOxztmsIxPMSWY8t4wmhAyNj1cM6X6HuljeqsMlP4F5oELXo21PUjAKUeufVr1Vr89xOG2IZSn_z_gQ",
		"e": "AQAB",
		"use": "sig",
		"kid": "q45HdfPW-oJ5xxfenWJtrKf9lvobt_oQzJ5h5zcYHMg",
		"qi": "GY6akUQr8D0D-mTLG0QSLBNBYWAJrSf0-yKnuNJ3mw4",
		"dp": "1bCk-tcoX5Y4z012kG3E6hNIDGJ8KZtA7dwD1GZvcHE",
		"alg": "RS256",
		"dq": "gT72Lbb36T665codwD5C86R1NUVKXQm7VWDuu5y54M8",
		"n": "yNkSpiJlB4O10f1E6UA6bJPhE0TL1og8Jq2dlo5kAN5UXtteTZpsRiz5mpnxjsCkhcTwh5UCwgEbRmWXAZgwZw"
	}
`

const testPublicKey string = `
	{
		"kty": "RSA",
		"e": "AQAB",
		"use": "sig",
		"kid": "q45HdfPW-oJ5xxfenWJtrKf9lvobt_oQzJ5h5zcYHMg",
		"alg": "RS256",
		"n": "yNkSpiJlB4O10f1E6UA6bJPhE0TL1og8Jq2dlo5kAN5UXtteTZpsRiz5mpnxjsCkhcTwh5UCwgEbRmWXAZgwZw"
	}
`

const wrongPublicKey string = `
	{
		"kty": "RSA",
		"e": "AQAB",
		"use": "sig",
		"kid": "q45HdfPW-oJ5xxfenWJtrKf9lvobt_oQzJ5h5zcYHMg",
		"alg": "RS256",
		"n": "yNkSpiJlB4O10f1E6UA6bJPhE0TL1og8Jq2dlo53AN5UXtteTZpsRiz5mpnxjsCkhcTwh5UCwgEbRmWXAZgwZw"
	}
`

func TestJWT(t *testing.T) {
	issuer, err := auth.NewJWTIssuer([]byte(testPrivateKey), auth.WithName("test_issuer"), auth.WithTokenDuration(time.Minute))
	if err != nil {
		t.Fatalf("could not create JWT issuer: %v", err)
	}

	signedToken, err := issuer.Issue("test_user", []auth.Role{auth.CustomerRole})
	if err != nil {
		t.Fatalf("could not issue JWT: %v", err)
	}

	t.Logf("signed token: %s", signedToken)

	verifier, err := auth.NewJWTVerifier([]byte(testPublicKey), auth.WithAllowedIssuer("test_issuer"))
	if err != nil {
		t.Fatalf("could not create JWT verifier: %v", err)
	}

	token, err := verifier.Verify(signedToken)
	if err != nil {
		t.Fatalf("could not verify token: %v", err)
	}

	if token.Roles()[0] != auth.CustomerRole {
		t.FailNow()
	}

	if token.Subject() != "test_user" {
		t.FailNow()
	}

	ctx := context.WithValue(context.Background(), auth.ContextKey, token)

	tokenFromCtx, err := auth.FromContext(ctx)
	if err != nil {
		t.FailNow()
	}

	if !reflect.DeepEqual(token, tokenFromCtx) {
		t.FailNow()
	}
}

func TestInvalidIssuer(t *testing.T) {
	issuer, err := auth.NewJWTIssuer([]byte(testPrivateKey), auth.WithName("test_issuer"), auth.WithTokenDuration(time.Minute))
	if err != nil {
		t.Fatalf("could not create JWT issuer: %v", err)
	}

	signedToken, err := issuer.Issue("test_user", []auth.Role{auth.CustomerRole})
	if err != nil {
		t.Fatalf("could not issue JWT: %v", err)
	}

	verifier, err := auth.NewJWTVerifier([]byte(testPublicKey))
	if err != nil {
		t.Fatalf("could not create JWT verifier: %v", err)
	}

	_, err = verifier.Verify(signedToken)
	if err == nil {
		t.Fatalf("token issue verification didn't fail: %v", err)
	}
}

func TestInvalidPublicKey(t *testing.T) {
	issuer, err := auth.NewJWTIssuer([]byte(testPrivateKey), auth.WithName("test_issuer"), auth.WithTokenDuration(time.Minute))
	if err != nil {
		t.Fatalf("could not create JWT issuer: %v", err)
	}

	signedToken, err := issuer.Issue("test_user", []auth.Role{auth.CustomerRole})
	if err != nil {
		t.Fatalf("could not issue JWT: %v", err)
	}

	verifier, err := auth.NewJWTVerifier([]byte(wrongPublicKey), auth.WithAllowedIssuer("test_issuer"))
	if err != nil {
		t.Fatalf("could not create JWT verifier: %v", err)
	}

	token, err := verifier.Verify(signedToken)
	if err == nil {
		t.Fatalf("token verification didn't fail, got: %v", token.Subject())
	}
}
