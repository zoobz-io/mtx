//go:build testing

package mtx

import (
	"fmt"
	"testing"

	mtesting "github.com/zoobzio/mtx/testing"
)

func TestRunDMSendExistingRoom(t *testing.T) {
	Homeserver = "http://localhost:8008"
	whoami := func() (*WhoAmIResponse, error) {
		return &WhoAmIResponse{UserID: "@blue-vicky:localhost"}, nil
	}
	getDirect := func(_ string) (map[string][]string, error) {
		return map[string][]string{
			"@blue-flux:localhost": {"!dm:localhost"},
		}, nil
	}
	setDirect := func(_ string, _ map[string][]string) error { return nil }

	var sentRoom, sentMsg string
	sender := func(roomID, message string) (string, error) {
		sentRoom = roomID
		sentMsg = message
		return "$evt1", nil
	}
	creator := func(_ string) (*Room, error) { return nil, fmt.Errorf("should not create") }
	checker := func(_ string) error { return nil }
	aliaser := func(_, _ string) error { return nil }

	resolver := func(_ string) (*AliasResponse, error) { return nil, fmt.Errorf("not found") }
	joiner := func(_ string) (string, error) { return "", fmt.Errorf("not found") }

	err := runDMSend("blue-flux", "hello", whoami, getDirect, setDirect, sender, creator, checker, aliaser, resolver, joiner)
	mtesting.AssertNoError(t, err)
	mtesting.AssertEqual(t, sentRoom, "!dm:localhost")
	mtesting.AssertEqual(t, sentMsg, "hello")
}

func TestRunDMSendCreatesRoom(t *testing.T) {
	Homeserver = "http://localhost:8008"
	whoami := func() (*WhoAmIResponse, error) {
		return &WhoAmIResponse{UserID: "@blue-vicky:localhost"}, nil
	}
	getDirect := func(_ string) (map[string][]string, error) {
		return map[string][]string{}, nil
	}
	var directsSet bool
	setDirect := func(_ string, rooms map[string][]string) error {
		directsSet = true
		return nil
	}
	sender := func(_, _ string) (string, error) { return "$evt1", nil }
	var createdInvite string
	creator := func(invite string) (*Room, error) {
		createdInvite = invite
		return &Room{RoomID: "!newdm:localhost"}, nil
	}
	checker := func(_ string) error { return nil }
	var aliasSet string
	aliaser := func(alias, _ string) error {
		aliasSet = alias
		return nil
	}

	resolver := func(_ string) (*AliasResponse, error) { return nil, fmt.Errorf("not found") }
	joiner := func(_ string) (string, error) { return "", fmt.Errorf("not found") }

	err := runDMSend("blue-flux", "hey", whoami, getDirect, setDirect, sender, creator, checker, aliaser, resolver, joiner)
	mtesting.AssertNoError(t, err)
	mtesting.AssertEqual(t, createdInvite, "@blue-flux:localhost")
	mtesting.AssertEqual(t, directsSet, true)
	mtesting.AssertEqual(t, aliasSet, "#dm-blue-flux-blue-vicky:localhost")
}

func TestRunDMSendWithFullUserID(t *testing.T) {
	Homeserver = "http://localhost:8008"
	whoami := func() (*WhoAmIResponse, error) {
		return &WhoAmIResponse{UserID: "@blue-vicky:localhost"}, nil
	}
	getDirect := func(_ string) (map[string][]string, error) {
		return map[string][]string{
			"@red-flux:example.com": {"!dm:example.com"},
		}, nil
	}
	setDirect := func(_ string, _ map[string][]string) error { return nil }
	var sentRoom string
	sender := func(roomID, _ string) (string, error) {
		sentRoom = roomID
		return "$evt1", nil
	}
	creator := func(_ string) (*Room, error) { return nil, fmt.Errorf("should not create") }
	checker := func(_ string) error { return nil }
	aliaser := func(_, _ string) error { return nil }

	resolver := func(_ string) (*AliasResponse, error) { return nil, fmt.Errorf("not found") }
	joiner := func(_ string) (string, error) { return "", fmt.Errorf("not found") }

	err := runDMSend("@red-flux:example.com", "hi", whoami, getDirect, setDirect, sender, creator, checker, aliaser, resolver, joiner)
	mtesting.AssertNoError(t, err)
	mtesting.AssertEqual(t, sentRoom, "!dm:example.com")
}

func TestRunDMSendWhoAmIError(t *testing.T) {
	whoami := func() (*WhoAmIResponse, error) {
		return nil, fmt.Errorf("unauthorized")
	}
	err := runDMSend("someone", "hi", whoami, nil, nil, nil, nil, nil, nil, nil, nil)
	mtesting.AssertError(t, err)
}

func TestRunDMSendNonexistentUser(t *testing.T) {
	Homeserver = "http://localhost:8008"
	whoami := func() (*WhoAmIResponse, error) {
		return &WhoAmIResponse{UserID: "@blue-vicky:localhost"}, nil
	}
	getDirect := func(_ string) (map[string][]string, error) {
		return map[string][]string{}, nil
	}
	setDirect := func(_ string, _ map[string][]string) error { return nil }
	sender := func(_, _ string) (string, error) { return "$evt1", nil }
	creator := func(_ string) (*Room, error) { return nil, fmt.Errorf("should not create") }
	checker := func(_ string) error { return fmt.Errorf("not found") }
	aliaser := func(_, _ string) error { return nil }

	resolver := func(_ string) (*AliasResponse, error) { return nil, fmt.Errorf("not found") }
	joiner := func(_ string) (string, error) { return "", fmt.Errorf("not found") }

	err := runDMSend("ghost", "hello", whoami, getDirect, setDirect, sender, creator, checker, aliaser, resolver, joiner)
	mtesting.AssertError(t, err)
}

func TestRunDMSendAliasFallback(t *testing.T) {
	Homeserver = "http://localhost:8008"
	whoami := func() (*WhoAmIResponse, error) {
		return &WhoAmIResponse{UserID: "@blue-vicky:localhost"}, nil
	}
	getDirect := func(_ string) (map[string][]string, error) {
		return map[string][]string{}, nil // no m.direct entry
	}
	var directsUpdated map[string][]string
	setDirect := func(_ string, rooms map[string][]string) error {
		directsUpdated = rooms
		return nil
	}
	var sentRoom string
	sender := func(roomID, _ string) (string, error) {
		sentRoom = roomID
		return "$evt1", nil
	}
	creator := func(_ string) (*Room, error) { return nil, fmt.Errorf("should not create") }
	checker := func(_ string) error { return nil }
	aliaser := func(_, _ string) error { return nil }
	resolver := func(alias string) (*AliasResponse, error) {
		return &AliasResponse{RoomID: "!existing-dm:localhost"}, nil
	}
	joiner := func(roomID string) (string, error) {
		return roomID, nil
	}

	err := runDMSend("blue-flux", "hello", whoami, getDirect, setDirect, sender, creator, checker, aliaser, resolver, joiner)
	mtesting.AssertNoError(t, err)
	mtesting.AssertEqual(t, sentRoom, "!existing-dm:localhost")
	mtesting.AssertEqual(t, len(directsUpdated["@blue-flux:localhost"]), 1)
}

func TestRunDMReadSuccess(t *testing.T) {
	Homeserver = "http://localhost:8008"
	whoami := func() (*WhoAmIResponse, error) {
		return &WhoAmIResponse{UserID: "@blue-vicky:localhost"}, nil
	}
	getDirect := func(_ string) (map[string][]string, error) {
		return map[string][]string{
			"@blue-flux:localhost": {"!dm:localhost"},
		}, nil
	}
	reader := func(roomID string, limit int) (*Messages, error) {
		mtesting.AssertEqual(t, roomID, "!dm:localhost")
		return &Messages{
			Chunk: []Message{
				{Sender: "@blue-flux:localhost", Type: "m.room.message", Content: map[string]interface{}{"body": "hey"}},
			},
		}, nil
	}
	resolver := func(_ string) (*AliasResponse, error) {
		t.Fatal("resolver should not be called when m.direct has the room")
		return nil, nil
	}
	joiner := func(_ string) (string, error) {
		t.Fatal("joiner should not be called when m.direct has the room")
		return "", nil
	}

	err := runDMRead("blue-flux", 20, false, whoami, getDirect, reader, resolver, joiner)
	mtesting.AssertNoError(t, err)
}

func TestRunDMReadAliasFallback(t *testing.T) {
	Homeserver = "http://localhost:8008"
	whoami := func() (*WhoAmIResponse, error) {
		return &WhoAmIResponse{UserID: "@blue-vicky:localhost"}, nil
	}
	getDirect := func(_ string) (map[string][]string, error) {
		return map[string][]string{}, nil // no m.direct entry
	}
	reader := func(roomID string, limit int) (*Messages, error) {
		mtesting.AssertEqual(t, roomID, "!dm:localhost")
		return &Messages{
			Chunk: []Message{
				{Sender: "@blue-flux:localhost", Type: "m.room.message", Content: map[string]interface{}{"body": "hey"}},
			},
		}, nil
	}
	var resolvedAlias string
	resolver := func(alias string) (*AliasResponse, error) {
		resolvedAlias = alias
		return &AliasResponse{RoomID: "!dm:localhost"}, nil
	}
	var joinedRoom string
	joiner := func(roomID string) (string, error) {
		joinedRoom = roomID
		return roomID, nil
	}

	err := runDMRead("blue-flux", 20, false, whoami, getDirect, reader, resolver, joiner)
	mtesting.AssertNoError(t, err)
	mtesting.AssertEqual(t, resolvedAlias, "#dm-blue-flux-blue-vicky:localhost")
	mtesting.AssertEqual(t, joinedRoom, "!dm:localhost")
}

func TestRunDMReadNoDMRoom(t *testing.T) {
	Homeserver = "http://localhost:8008"
	whoami := func() (*WhoAmIResponse, error) {
		return &WhoAmIResponse{UserID: "@blue-vicky:localhost"}, nil
	}
	getDirect := func(_ string) (map[string][]string, error) {
		return map[string][]string{}, nil
	}
	resolver := func(_ string) (*AliasResponse, error) {
		return nil, fmt.Errorf("not found")
	}
	joiner := func(_ string) (string, error) {
		return "", fmt.Errorf("not found")
	}

	err := runDMRead("ghost", 20, false, whoami, getDirect, nil, resolver, joiner)
	mtesting.AssertError(t, err)
}

func TestLocalpart(t *testing.T) {
	mtesting.AssertEqual(t, localpart("@alice:localhost"), "alice")
	mtesting.AssertEqual(t, localpart("@blue-vicky:example.com"), "blue-vicky")
	mtesting.AssertEqual(t, localpart("alice"), "alice")
}

func TestDMAliasName(t *testing.T) {
	// Should sort alphabetically regardless of order.
	mtesting.AssertEqual(t, dmAliasName("@blue-vicky:localhost", "@blue-flux:localhost"), "dm-blue-flux-blue-vicky")
	mtesting.AssertEqual(t, dmAliasName("@blue-flux:localhost", "@blue-vicky:localhost"), "dm-blue-flux-blue-vicky")
}
