package appmessage

import (
	"github.com/ammm56/lings/domain/consensus/model/externalapi"
)

// MsgIBDBlockLocator represents a lings ibdBlockLocator message
type MsgIBDBlockLocator struct {
	baseMessage
	TargetHash         *externalapi.DomainHash
	BlockLocatorHashes []*externalapi.DomainHash
}

// Command returns the protocol command string for the message
func (msg *MsgIBDBlockLocator) Command() MessageCommand {
	return CmdIBDBlockLocator
}

// NewMsgIBDBlockLocator returns a new lings ibdBlockLocator message
func NewMsgIBDBlockLocator(targetHash *externalapi.DomainHash,
	blockLocatorHashes []*externalapi.DomainHash) *MsgIBDBlockLocator {

	return &MsgIBDBlockLocator{
		TargetHash:         targetHash,
		BlockLocatorHashes: blockLocatorHashes,
	}
}
