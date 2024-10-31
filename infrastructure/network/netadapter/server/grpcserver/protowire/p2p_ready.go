package protowire

import (
	"github.com/ammm56/lings/app/appmessage"
	"github.com/pkg/errors"
)

func (x *LingsMessage_Ready) toAppMessage() (appmessage.Message, error) {
	if x == nil {
		return nil, errors.Wrapf(errorNil, "LingsMessage_Ready is nil")
	}
	return &appmessage.MsgReady{}, nil
}

func (x *LingsMessage_Ready) fromAppMessage(_ *appmessage.MsgReady) error {
	return nil
}
