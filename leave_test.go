//go:build testing

package mtx

import (
	"fmt"
	"testing"

	mtesting "github.com/zoobzio/mtx/testing"
)

func TestRunLeaveSuccess(t *testing.T) {
	var leftRoom string
	leaver := func(roomID string) error {
		leftRoom = roomID
		return nil
	}
	err := runLeave("!room:localhost", leaver)
	mtesting.AssertNoError(t, err)
	mtesting.AssertEqual(t, leftRoom, "!room:localhost")
}

func TestRunLeaveError(t *testing.T) {
	leaver := func(_ string) error {
		return fmt.Errorf("forbidden")
	}
	err := runLeave("!room:localhost", leaver)
	mtesting.AssertError(t, err)
}
