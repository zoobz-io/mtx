//go:build testing

package mtx

import (
	"fmt"
	"testing"

	mtesting "github.com/zoobzio/mtx/testing"
)

func TestRunMembersSuccess(t *testing.T) {
	lister := func(roomID string) ([]Member, error) {
		return []Member{
			{UserID: "@alice:localhost", DisplayName: "Alice"},
			{UserID: "@bob:localhost", DisplayName: ""},
		}, nil
	}
	err := runMembers("!room:localhost", lister)
	mtesting.AssertNoError(t, err)
}

func TestRunMembersError(t *testing.T) {
	lister := func(_ string) ([]Member, error) {
		return nil, fmt.Errorf("forbidden")
	}
	err := runMembers("!room:localhost", lister)
	mtesting.AssertError(t, err)
}
