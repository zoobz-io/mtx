//go:build testing

package mtx

import (
	"fmt"
	"testing"

	mtesting "github.com/zoobzio/mtx/testing"
)

func TestRunSendSuccess(t *testing.T) {
	var sentRoom, sentMsg string
	sender := func(roomID, message string) (string, error) {
		sentRoom = roomID
		sentMsg = message
		return "$evt1", nil
	}
	err := runSend("!room:localhost", "hello world", sender)
	mtesting.AssertNoError(t, err)
	mtesting.AssertEqual(t, sentRoom, "!room:localhost")
	mtesting.AssertEqual(t, sentMsg, "hello world")
}

func TestRunSendError(t *testing.T) {
	sender := func(_, _ string) (string, error) {
		return "", fmt.Errorf("send failed")
	}
	err := runSend("!room:localhost", "hello", sender)
	mtesting.AssertError(t, err)
}
