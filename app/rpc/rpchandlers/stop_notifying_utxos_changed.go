package rpchandlers

import (
	"github.com/ammm56/lings/app/appmessage"
	"github.com/ammm56/lings/app/rpc/rpccontext"
	"github.com/ammm56/lings/infrastructure/network/netadapter/router"
)

// HandleStopNotifyingUTXOsChanged handles the respectively named RPC command
func HandleStopNotifyingUTXOsChanged(context *rpccontext.Context, router *router.Router, request appmessage.Message) (appmessage.Message, error) {
	if !context.Config.UTXOIndex {
		errorMessage := appmessage.NewStopNotifyingUTXOsChangedResponseMessage()
		errorMessage.Error = appmessage.RPCErrorf("Method unavailable when lings is run without --utxoindex")
		return errorMessage, nil
	}

	stopNotifyingUTXOsChangedRequest := request.(*appmessage.StopNotifyingUTXOsChangedRequestMessage)
	addresses, err := context.ConvertAddressStringsToUTXOsChangedNotificationAddresses(stopNotifyingUTXOsChangedRequest.Addresses)
	if err != nil {
		errorMessage := appmessage.NewNotifyUTXOsChangedResponseMessage()
		errorMessage.Error = appmessage.RPCErrorf("Parsing error: %s", err)
		return errorMessage, nil
	}

	listener, err := context.NotificationManager.Listener(router)
	if err != nil {
		return nil, err
	}
	context.NotificationManager.StopPropagatingUTXOsChangedNotifications(listener, addresses)

	response := appmessage.NewStopNotifyingUTXOsChangedResponseMessage()
	return response, nil
}
