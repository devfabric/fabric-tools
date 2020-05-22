package event

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/events/deliverclient/seek"
)

// ListenTx 监听交易事件
// 监听完毕后需要调用注销函数 eventHub.Unregister(reg)
// eventch 为事件监听通道
func (E *EventManager) ListenTx(user, txID string) (*event.Client, fab.Registration, <-chan *fab.TxStatusEvent, error) {
	evnetOpts, err := newEvnetOpts(seek.Oldest, 0)
	if err != nil {
		return nil, nil, nil, err
	}

	eventHub, err := E.EventClient(user, evnetOpts...)
	if err != nil {
		return nil, nil, nil, err
	}

	E.Log().Infof("Registering TX event for TxID [%s]\n", txID)

	reg, eventch, err := eventHub.RegisterTxStatusEvent(txID)
	return eventHub, reg, eventch, err
	// if err != nil {
	// 	return errors.WithMessage(err, "Error registering for block events")
	// }
	// defer eventHub.Unregister(reg)

	// select {
	// case _, _ = <-done:
	// 	return nil
	// case event, ok := <-eventch:
	// 	if !ok {
	// 		return errors.WithMessage(err, "unexpected closed channel while waiting for tx status event")
	// 	}
	// 	E.Log().Infof("Received TX event. TxID: %s, Code: %s, Error: %s\n", event.TxID, event.TxValidationCode, err)
	// }

	// return nil
}
