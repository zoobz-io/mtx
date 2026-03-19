//go:build testing

package mtx

import (
	"fmt"
	"testing"

	mtesting "github.com/zoobzio/mtx/testing"
)

func TestRunLoginSuccess(t *testing.T) {
	authenticator := func(user, pass string) (*Registration, error) {
		return &Registration{
			UserID:      "@" + user + ":localhost",
			AccessToken: "tok_login",
		}, nil
	}

	err := runLogin("operator", "pass", authenticator)
	mtesting.AssertNoError(t, err)
}

func TestRunLoginError(t *testing.T) {
	authenticator := func(_, _ string) (*Registration, error) {
		return nil, fmt.Errorf("bad credentials")
	}
	err := runLogin("x", "y", authenticator)
	mtesting.AssertError(t, err)
}
