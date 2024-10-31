package main

import (
	"context"
	"fmt"

	"github.com/ammm56/lings/cmd/lingswallet/daemon/client"
	"github.com/ammm56/lings/cmd/lingswallet/daemon/pb"
	"github.com/ammm56/lings/cmd/lingswallet/utils"
)

func balance(conf *balanceConfig) error {
	daemonClient, tearDown, err := client.Connect(conf.DaemonAddress)
	if err != nil {
		return err
	}
	defer tearDown()

	ctx, cancel := context.WithTimeout(context.Background(), daemonTimeout)
	defer cancel()
	response, err := daemonClient.GetBalance(ctx, &pb.GetBalanceRequest{})
	if err != nil {
		return err
	}

	pendingSuffix := ""
	if response.Pending > 0 {
		pendingSuffix = " (pending)"
	}
	if conf.Verbose {
		pendingSuffix = ""
		println("Address                                                                       Available             Pending")
		println("-----------------------------------------------------------------------------------------------------------")
		for _, addressBalance := range response.AddressBalances {
			fmt.Printf("%s %s %s\n", addressBalance.Address, utils.FomatLSN(addressBalance.Available), utils.FomatLSN(addressBalance.Pending))
		}
		println("-----------------------------------------------------------------------------------------------------------")
		print("                                                 ")
	}
	fmt.Printf("Total balance, LSN %s %s%s\n", utils.FomatLSN(response.Available), utils.FomatLSN(response.Pending), pendingSuffix)

	return nil
}
