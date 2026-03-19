//go:build testing

package mtx

import (
	"fmt"
	"testing"

	mtesting "github.com/zoobzio/mtx/testing"
)

func TestRunWhoAmISuccess(t *testing.T) {
	getter := func() (*WhoAmIResponse, error) {
		return &WhoAmIResponse{UserID: "@blue-vicky:localhost"}, nil
	}
	err := runWhoAmI(getter)
	mtesting.AssertNoError(t, err)
}

func TestRunWhoAmIError(t *testing.T) {
	getter := func() (*WhoAmIResponse, error) {
		return nil, fmt.Errorf("unauthorized")
	}
	err := runWhoAmI(getter)
	mtesting.AssertError(t, err)
}
