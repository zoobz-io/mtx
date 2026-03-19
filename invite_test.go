//go:build testing

package mtx

import (
	"fmt"
	"testing"

	mtesting "github.com/zoobzio/mtx/testing"
)

func TestRunInviteSuccess(t *testing.T) {
	var invitedRoom, invitedUser string
	inviter := func(roomID, userID string) error {
		invitedRoom = roomID
		invitedUser = userID
		return nil
	}
	err := runInvite("!room:localhost", "@agent:localhost", inviter)
	mtesting.AssertNoError(t, err)
	mtesting.AssertEqual(t, invitedRoom, "!room:localhost")
	mtesting.AssertEqual(t, invitedUser, "@agent:localhost")
}

func TestRunInviteError(t *testing.T) {
	inviter := func(_, _ string) error {
		return fmt.Errorf("forbidden")
	}
	err := runInvite("!room:localhost", "@agent:localhost", inviter)
	mtesting.AssertError(t, err)
}
