package integration

import (
	"testing"

	"github.com/ammm56/lings/app/appmessage"
)

func TestVirtualSelectedParentBlueScoreAndVirtualDAAScore(t *testing.T) {
	// Setup a single lings instance
	harnessParams := &harnessParams{
		p2pAddress:              p2pAddress1,
		rpcAddress:              rpcAddress1,
		miningAddress:           miningAddress1,
		miningAddressPrivateKey: miningAddress1PrivateKey,
		utxoIndex:               true,
	}
	lings, teardown := setupHarness(t, harnessParams)
	defer teardown()

	// Make sure that the initial selected parent blue score is 0
	response, err := lings.rpcClient.GetVirtualSelectedParentBlueScore()
	if err != nil {
		t.Fatalf("Error getting virtual selected parent blue score: %s", err)
	}
	if response.BlueScore != 0 {
		t.Fatalf("Unexpected virtual selected parent blue score. Want: %d, got: %d",
			0, response.BlueScore)
	}

	// Register to virtual selected parent blue score changes
	onVirtualSelectedParentBlueScoreChangedChan := make(chan *appmessage.VirtualSelectedParentBlueScoreChangedNotificationMessage)
	err = lings.rpcClient.RegisterForVirtualSelectedParentBlueScoreChangedNotifications(
		func(notification *appmessage.VirtualSelectedParentBlueScoreChangedNotificationMessage) {
			onVirtualSelectedParentBlueScoreChangedChan <- notification
		})
	if err != nil {
		t.Fatalf("Failed to register for virtual selected parent "+
			"blue score change notifications: %s", err)
	}

	// Register to virtual DAA score changes
	onVirtualDaaScoreChangedChan := make(chan *appmessage.VirtualDaaScoreChangedNotificationMessage)
	err = lings.rpcClient.RegisterForVirtualDaaScoreChangedNotifications(
		func(notification *appmessage.VirtualDaaScoreChangedNotificationMessage) {
			onVirtualDaaScoreChangedChan <- notification
		})
	if err != nil {
		t.Fatalf("Failed to register for virtual DAA score change notifications: %s", err)
	}

	// Mine some blocks and make sure that the notifications
	// report correct values
	const blockAmountToMine = 100
	for i := 0; i < blockAmountToMine; i++ {
		mineNextBlock(t, lings)
		blueScoreChangedNotification := <-onVirtualSelectedParentBlueScoreChangedChan
		if blueScoreChangedNotification.VirtualSelectedParentBlueScore != 1+uint64(i) {
			t.Fatalf("Unexpected virtual selected parent blue score. Want: %d, got: %d",
				1+uint64(i), blueScoreChangedNotification.VirtualSelectedParentBlueScore)
		}
		daaScoreChangedNotification := <-onVirtualDaaScoreChangedChan
		if daaScoreChangedNotification.VirtualDaaScore > 1+uint64(i) {
			t.Fatalf("Unexpected virtual DAA score. Want: %d, got: %d",
				1+uint64(i), daaScoreChangedNotification.VirtualDaaScore)
		}
	}

	// Make sure that the blue score after all that mining is as expected
	response, err = lings.rpcClient.GetVirtualSelectedParentBlueScore()
	if err != nil {
		t.Fatalf("Error getting virtual selected parent blue score: %s", err)
	}
	if response.BlueScore != blockAmountToMine {
		t.Fatalf("Unexpected virtual selected parent blue score. Want: %d, got: %d",
			blockAmountToMine, response.BlueScore)
	}
}
