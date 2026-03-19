//go:build testing

package mtx

import (
	"fmt"
	"testing"

	mtesting "github.com/zoobzio/mtx/testing"
)

func TestRunCreateSuccess(t *testing.T) {
	var createdName, createdTopic string
	creator := func(name, topic string) (*Room, error) {
		createdName = name
		createdTopic = topic
		return &Room{RoomID: "!new:localhost"}, nil
	}
	err := runCreate("general", "dev discussion", creator)
	mtesting.AssertNoError(t, err)
	mtesting.AssertEqual(t, createdName, "general")
	mtesting.AssertEqual(t, createdTopic, "dev discussion")
}

func TestRunCreateError(t *testing.T) {
	creator := func(_, _ string) (*Room, error) {
		return nil, fmt.Errorf("already exists")
	}
	err := runCreate("general", "", creator)
	mtesting.AssertError(t, err)
}

func TestRunCreateWithAliasSuccess(t *testing.T) {
	Homeserver = "http://localhost:8008"
	var createdAlias string
	creator := func(name, topic, alias string) (*Room, error) {
		createdAlias = alias
		return &Room{RoomID: "!new:localhost"}, nil
	}
	err := runCreateWithAlias("general", "dev chat", "general", creator)
	mtesting.AssertNoError(t, err)
	mtesting.AssertEqual(t, createdAlias, "general")
}

func TestRunCreateWithAliasError(t *testing.T) {
	creator := func(_, _, _ string) (*Room, error) {
		return nil, fmt.Errorf("alias taken")
	}
	err := runCreateWithAlias("general", "", "general", creator)
	mtesting.AssertError(t, err)
}
