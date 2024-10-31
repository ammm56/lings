package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ammm56/lings/infrastructure/network/netadapter/server/grpcserver/protowire"
)

var commandTypes = []reflect.Type{
	reflect.TypeOf(protowire.LingsMessage_AddPeerRequest{}),
	reflect.TypeOf(protowire.LingsMessage_GetConnectedPeerInfoRequest{}),
	reflect.TypeOf(protowire.LingsMessage_GetPeerAddressesRequest{}),
	reflect.TypeOf(protowire.LingsMessage_GetCurrentNetworkRequest{}),
	reflect.TypeOf(protowire.LingsMessage_GetInfoRequest{}),

	reflect.TypeOf(protowire.LingsMessage_GetBlockRequest{}),
	reflect.TypeOf(protowire.LingsMessage_GetBlocksRequest{}),
	reflect.TypeOf(protowire.LingsMessage_GetHeadersRequest{}),
	reflect.TypeOf(protowire.LingsMessage_GetBlockCountRequest{}),
	reflect.TypeOf(protowire.LingsMessage_GetBlockDagInfoRequest{}),
	reflect.TypeOf(protowire.LingsMessage_GetSelectedTipHashRequest{}),
	reflect.TypeOf(protowire.LingsMessage_GetVirtualSelectedParentBlueScoreRequest{}),
	reflect.TypeOf(protowire.LingsMessage_GetVirtualSelectedParentChainFromBlockRequest{}),
	reflect.TypeOf(protowire.LingsMessage_ResolveFinalityConflictRequest{}),
	reflect.TypeOf(protowire.LingsMessage_EstimateNetworkHashesPerSecondRequest{}),

	reflect.TypeOf(protowire.LingsMessage_GetBlockTemplateRequest{}),
	reflect.TypeOf(protowire.LingsMessage_SubmitBlockRequest{}),

	reflect.TypeOf(protowire.LingsMessage_GetMempoolEntryRequest{}),
	reflect.TypeOf(protowire.LingsMessage_GetMempoolEntriesRequest{}),
	reflect.TypeOf(protowire.LingsMessage_GetMempoolEntriesByAddressesRequest{}),

	reflect.TypeOf(protowire.LingsMessage_SubmitTransactionRequest{}),

	reflect.TypeOf(protowire.LingsMessage_GetUtxosByAddressesRequest{}),
	reflect.TypeOf(protowire.LingsMessage_GetBalanceByAddressRequest{}),
	reflect.TypeOf(protowire.LingsMessage_GetCoinSupplyRequest{}),

	reflect.TypeOf(protowire.LingsMessage_BanRequest{}),
	reflect.TypeOf(protowire.LingsMessage_UnbanRequest{}),
}

type commandDescription struct {
	name       string
	parameters []*parameterDescription
	typeof     reflect.Type
}

type parameterDescription struct {
	name   string
	typeof reflect.Type
}

func commandDescriptions() []*commandDescription {
	commandDescriptions := make([]*commandDescription, len(commandTypes))

	for i, commandTypeWrapped := range commandTypes {
		commandType := unwrapCommandType(commandTypeWrapped)

		name := strings.TrimSuffix(commandType.Name(), "RequestMessage")
		numFields := commandType.NumField()

		var parameters []*parameterDescription
		for i := 0; i < numFields; i++ {
			field := commandType.Field(i)

			if !isFieldExported(field) {
				continue
			}

			parameters = append(parameters, &parameterDescription{
				name:   field.Name,
				typeof: field.Type,
			})
		}
		commandDescriptions[i] = &commandDescription{
			name:       name,
			parameters: parameters,
			typeof:     commandTypeWrapped,
		}
	}

	return commandDescriptions
}

func (cd *commandDescription) help() string {
	sb := &strings.Builder{}
	sb.WriteString(cd.name)
	for _, parameter := range cd.parameters {
		_, _ = fmt.Fprintf(sb, " [%s]", parameter.name)
	}
	return sb.String()
}
