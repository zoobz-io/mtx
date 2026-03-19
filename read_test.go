//go:build testing

package mtx

import (
	"fmt"
	"testing"

	mtesting "github.com/zoobzio/mtx/testing"
)

func TestRunReadSuccess(t *testing.T) {
	reader := func(roomID string, limit int) (*Messages, error) {
		mtesting.AssertEqual(t, roomID, "!room:localhost")
		mtesting.AssertEqual(t, limit, 10)
		return &Messages{
			Chunk: []Message{
				{Sender: "@bob:localhost", Type: "m.room.message", Content: map[string]interface{}{"body": "second"}},
				{Sender: "@alice:localhost", Type: "m.room.message", Content: map[string]interface{}{"body": "first"}},
			},
		}, nil
	}
	err := runRead("!room:localhost", 10, reader)
	mtesting.AssertNoError(t, err)
}

func TestRunReadError(t *testing.T) {
	reader := func(_ string, _ int) (*Messages, error) {
		return nil, fmt.Errorf("not found")
	}
	err := runRead("!room:localhost", 10, reader)
	mtesting.AssertError(t, err)
}

func TestRunReadFiltersNonMessages(t *testing.T) {
	reader := func(_ string, _ int) (*Messages, error) {
		return &Messages{
			Chunk: []Message{
				{Sender: "@bob:localhost", Type: "m.room.member", Content: map[string]interface{}{}},
				{Sender: "@alice:localhost", Type: "m.room.message", Content: map[string]interface{}{"body": "hello"}},
			},
		}, nil
	}
	// Should not error; non-message events are silently skipped.
	err := runRead("!room:localhost", 10, reader)
	mtesting.AssertNoError(t, err)
}

func TestRunReadJSONSuccess(t *testing.T) {
	reader := func(roomID string, limit int) (*Messages, error) {
		return &Messages{
			Chunk: []Message{
				{Sender: "@bob:localhost", Type: "m.room.message", Content: map[string]interface{}{"body": "hi"}, EventID: "$evt1"},
			},
		}, nil
	}
	err := runReadJSON("!room:localhost", 10, reader)
	mtesting.AssertNoError(t, err)
}

func TestRunReadJSONError(t *testing.T) {
	reader := func(_ string, _ int) (*Messages, error) {
		return nil, fmt.Errorf("not found")
	}
	err := runReadJSON("!room:localhost", 10, reader)
	mtesting.AssertError(t, err)
}

func TestRunReadSinceSuccess(t *testing.T) {
	getContext := func(roomID, eventID string) (string, error) {
		mtesting.AssertEqual(t, roomID, "!room:localhost")
		mtesting.AssertEqual(t, eventID, "$evt1")
		return "token_after_evt1", nil
	}
	readFrom := func(roomID, from string, limit int, dir string) (*Messages, error) {
		mtesting.AssertEqual(t, from, "token_after_evt1")
		mtesting.AssertEqual(t, dir, "f")
		return &Messages{
			Chunk: []Message{
				{Sender: "@alice:localhost", Type: "m.room.message", Content: map[string]interface{}{"body": "new msg"}, EventID: "$evt2"},
			},
		}, nil
	}
	err := runReadSince("!room:localhost", "$evt1", 20, false, getContext, readFrom)
	mtesting.AssertNoError(t, err)
}

func TestRunReadSinceJSON(t *testing.T) {
	getContext := func(_, _ string) (string, error) { return "tok", nil }
	readFrom := func(_, _ string, _ int, _ string) (*Messages, error) {
		return &Messages{
			Chunk: []Message{
				{Sender: "@alice:localhost", Type: "m.room.message", Content: map[string]interface{}{"body": "new"}, EventID: "$evt2"},
			},
		}, nil
	}
	err := runReadSince("!room:localhost", "$evt1", 20, true, getContext, readFrom)
	mtesting.AssertNoError(t, err)
}

func TestRunReadSinceContextError(t *testing.T) {
	getContext := func(_, _ string) (string, error) {
		return "", fmt.Errorf("event not found")
	}
	err := runReadSince("!room:localhost", "$bad", 20, false, getContext, nil)
	mtesting.AssertError(t, err)
}

func TestSenderMatches(t *testing.T) {
	mtesting.AssertEqual(t, senderMatches("@wintermute-argus:localhost", "wintermute-argus"), true)
	mtesting.AssertEqual(t, senderMatches("@wintermute-argus:localhost", "@wintermute-argus:localhost"), true)
	mtesting.AssertEqual(t, senderMatches("@wintermute-argus:localhost", "argus"), true)
	mtesting.AssertEqual(t, senderMatches("@wintermute-argus:localhost", "rockhopper"), false)
}

func TestRunReadFiltered(t *testing.T) {
	reader := func(_ string, _ int) (*Messages, error) {
		return &Messages{
			Chunk: []Message{
				{Sender: "@bob:localhost", Type: "m.room.message", Content: map[string]interface{}{"body": "from bob"}, EventID: "$evt1"},
				{Sender: "@alice:localhost", Type: "m.room.message", Content: map[string]interface{}{"body": "from alice"}, EventID: "$evt2"},
				{Sender: "@bob:localhost", Type: "m.room.message", Content: map[string]interface{}{"body": "also bob"}, EventID: "$evt3"},
			},
		}, nil
	}
	// Should not error — filters to only alice's messages.
	err := runReadFiltered("!room:localhost", 10, false, "alice", reader)
	mtesting.AssertNoError(t, err)
}

func TestRunReadFilteredJSON(t *testing.T) {
	reader := func(_ string, _ int) (*Messages, error) {
		return &Messages{
			Chunk: []Message{
				{Sender: "@bob:localhost", Type: "m.room.message", Content: map[string]interface{}{"body": "hi"}, EventID: "$evt1"},
				{Sender: "@alice:localhost", Type: "m.room.message", Content: map[string]interface{}{"body": "hey"}, EventID: "$evt2"},
			},
		}, nil
	}
	err := runReadFiltered("!room:localhost", 10, true, "bob", reader)
	mtesting.AssertNoError(t, err)
}
