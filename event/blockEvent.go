package event

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

// ListenBlock 监听区块事件
func (E *EventManager) ListenBlock(user, seekType string, blockNum uint64) (*event.Client, fab.Registration, <-chan *fab.BlockEvent, error) {
	evnetOpts, err := newEvnetOpts(seekType, blockNum)
	if err != nil {
		return nil, nil, nil, err
	}

	eventClient, err := E.EventClient(user, evnetOpts...)
	if err != nil {
		return nil, nil, nil, err
	}

	breg, beventch, err := eventClient.RegisterBlockEvent()
	return eventClient, breg, beventch, err
	// if err != nil {
	// 	return errors.WithMessage(err, "Error registering for block events")
	// }
	// defer eventClient.Unregister(breg)

	// for {
	// 	select {
	// 	case _, _ = <-done:
	// 		return nil
	// 	case event, ok := <-beventch:
	// 		if !ok {
	// 			return errors.WithMessage(err, "unexpected closed channel while waiting for block event")
	// 		}
	// 		E.Log().Infof("%+v", event)
	// 	}
	// }

}
