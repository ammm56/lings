package protowire

import (
	"github.com/ammm56/lings/app/appmessage"
	"github.com/pkg/errors"
)

func (x *LingsMessage_InvRelayBlock) toAppMessage() (appmessage.Message, error) {
	if x == nil {
		return nil, errors.Wrapf(errorNil, "LingsMessage_InvRelayBlock is nil")
	}
	return x.InvRelayBlock.toAppMessage()
}

func (x *InvRelayBlockMessage) toAppMessage() (appmessage.Message, error) {
	if x == nil {
		return nil, errors.Wrapf(errorNil, "InvRelayBlockMessage is nil")
	}
	hash, err := x.Hash.toDomain()
	if err != nil {
		return nil, err
	}

	return &appmessage.MsgInvRelayBlock{Hash: hash}, nil

}

func (x *LingsMessage_InvRelayBlock) fromAppMessage(msgInvRelayBlock *appmessage.MsgInvRelayBlock) error {
	x.InvRelayBlock = &InvRelayBlockMessage{
		Hash: domainHashToProto(msgInvRelayBlock.Hash),
	}
	return nil
}
