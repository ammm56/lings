package protowire

import (
	"github.com/ammm56/lings/app/appmessage"
	"github.com/pkg/errors"
)

func (x *LingsMessage_Verack) toAppMessage() (appmessage.Message, error) {
	if x == nil {
		return nil, errors.Wrapf(errorNil, "LingsMessage_Verack is nil")
	}
	return &appmessage.MsgVerAck{}, nil
}

func (x *LingsMessage_Verack) fromAppMessage(_ *appmessage.MsgVerAck) error {
	return nil
}
