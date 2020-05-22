package event

import (
	"errors"

	"fabric-tools/manager"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/events/deliverclient/seek"
)

// EventManager 提供事件管理模块
type EventManager struct {
	manager.Manager
}

// NewEventManager 返回事件管理模块
func NewEventManager(m *manager.Manager) (*EventManager, error) {
	return &EventManager{*m}, nil
}

func newEvnetOpts(seekType string, blockNum uint64) ([]event.ClientOption, error) {
	if seekType != "oldest" && seekType != "newest" && seekType != "from" {
		return nil, errors.New("seek type error,must be one of [oldest,newest,from]")
	}

	evnetOpts := []event.ClientOption{}
	evnetOpts = append(evnetOpts, event.WithBlockEvents())
	evnetOpts = append(evnetOpts, event.WithSeekType(seek.Type(seekType)))

	if seekType == "from" {
		evnetOpts = append(evnetOpts, event.WithBlockNum(blockNum))
	}
	return evnetOpts, nil
}
