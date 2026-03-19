//go:build testing

package mtx

import (
	"fmt"
	"testing"

	mtesting "github.com/zoobzio/mtx/testing"
)

func noopInfoGetter(_ string) (*RoomInfo, error) {
	return &RoomInfo{}, nil
}

func TestRunRoomsSuccess(t *testing.T) {
	lister := func() (*JoinedRooms, error) {
		return &JoinedRooms{Rooms: []string{"!a:localhost", "!b:localhost"}}, nil
	}
	infoGetter := func(roomID string) (*RoomInfo, error) {
		return &RoomInfo{Name: "test-room", Topic: "a topic"}, nil
	}
	err := runRooms(lister, infoGetter)
	mtesting.AssertNoError(t, err)
}

func TestRunRoomsError(t *testing.T) {
	lister := func() (*JoinedRooms, error) {
		return nil, fmt.Errorf("unauthorized")
	}
	err := runRooms(lister, noopInfoGetter)
	mtesting.AssertError(t, err)
}

func TestRunPublicRoomsSuccess(t *testing.T) {
	lister := func() (*PublicRoomsResponse, error) {
		return &PublicRoomsResponse{
			Chunk: []PublicRoom{
				{RoomID: "!a:localhost", Name: "general", Topic: "general chat", NumJoined: 5},
				{RoomID: "!b:localhost", Name: "dev", NumJoined: 3},
			},
		}, nil
	}
	err := runPublicRooms(lister)
	mtesting.AssertNoError(t, err)
}

func TestRunPublicRoomsError(t *testing.T) {
	lister := func() (*PublicRoomsResponse, error) {
		return nil, fmt.Errorf("server error")
	}
	err := runPublicRooms(lister)
	mtesting.AssertError(t, err)
}
