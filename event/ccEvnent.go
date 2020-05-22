package event

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

// ListenChainCode 监听链码事件
// 监听完毕后需要调用注销函数 eventHub.Unregister(reg)
// eventch 为事件监听通道
func (E *EventManager) ListenChainCode(user, ccID, eventSt, seekType string, num uint64) (*event.Client, fab.Registration, <-chan *fab.CCEvent, error) {
	evnetOpts, err := newEvnetOpts(seekType, num)
	if err != nil {
		return nil, nil, nil, err
	}

	eventHub, err := E.EventClient(user, evnetOpts...)
	if err != nil {
		return nil, nil, nil, err
	}

	breg, beventch, err := eventHub.RegisterChaincodeEvent(ccID, eventSt)
	return eventHub, breg, beventch, err

	// if err != nil {
	// 	return errors.WithMessage(err, "Error registering for block events")
	// }
	// defer eventHub.Unregister(breg)

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
