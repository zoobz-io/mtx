//go:build testing

package mtx

import (
	"fmt"
	"testing"

	mtesting "github.com/zoobzio/mtx/testing"
)

func TestRunRegisterSuccess(t *testing.T) {
	registerer := func(user, pass, token string) (*Registration, error) {
		return &Registration{
			UserID:      "@" + user + ":localhost",
			AccessToken: "tok_abc",
		}, nil
	}

	err := runRegister("agent", "pass", "mtx", registerer)
	mtesting.AssertNoError(t, err)
}

func TestRunRegisterError(t *testing.T) {
	registerer := func(_, _, _ string) (*Registration, error) {
		return nil, fmt.Errorf("registration failed")
	}
	err := runRegister("agent", "pass", "mtx", registerer)
	mtesting.AssertError(t, err)
}
