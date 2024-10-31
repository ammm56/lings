package consensushashing_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/ammm56/lings/domain/consensus/utils/subnetworks"

	"github.com/kaspanet/go-secp256k1"

	"github.com/ammm56/lings/domain/consensus/utils/consensushashing"
	"github.com/ammm56/lings/domain/consensus/utils/txscript"
	"github.com/ammm56/lings/domain/consensus/utils/utxo"
	"github.com/ammm56/lings/domain/dagconfig"
	"github.com/ammm56/lings/util"

	"github.com/ammm56/lings/domain/consensus/model/externalapi"
)

// shortened versions of SigHash types to fit in single line of test case
const (
	all                = consensushashing.SigHashAll
	none               = consensushashing.SigHashNone
	single             = consensushashing.SigHashSingle
	allAnyoneCanPay    = consensushashing.SigHashAll | consensushashing.SigHashAnyOneCanPay
	noneAnyoneCanPay   = consensushashing.SigHashNone | consensushashing.SigHashAnyOneCanPay
	singleAnyoneCanPay = consensushashing.SigHashSingle | consensushashing.SigHashAnyOneCanPay
)

func modifyOutput(outputIndex int) func(tx *externalapi.DomainTransaction) *externalapi.DomainTransaction {
	return func(tx *externalapi.DomainTransaction) *externalapi.DomainTransaction {
		clone := tx.Clone()
		clone.Outputs[outputIndex].Value = 100
		return clone
	}
}

func modifyInput(inputIndex int) func(tx *externalapi.DomainTransaction) *externalapi.DomainTransaction {
	return func(tx *externalapi.DomainTransaction) *externalapi.DomainTransaction {
		clone := tx.Clone()
		clone.Inputs[inputIndex].PreviousOutpoint.Index = 2
		return clone
	}
}

func modifyAmountSpent(inputIndex int) func(tx *externalapi.DomainTransaction) *externalapi.DomainTransaction {
	return func(tx *externalapi.DomainTransaction) *externalapi.DomainTransaction {
		clone := tx.Clone()
		utxoEntry := clone.Inputs[inputIndex].UTXOEntry
		clone.Inputs[inputIndex].UTXOEntry = utxo.NewUTXOEntry(666, utxoEntry.ScriptPublicKey(), false, 100)
		return clone
	}
}

func modifyScriptPublicKey(inputIndex int) func(tx *externalapi.DomainTransaction) *externalapi.DomainTransaction {
	return func(tx *externalapi.DomainTransaction) *externalapi.DomainTransaction {
		clone := tx.Clone()
		utxoEntry := clone.Inputs[inputIndex].UTXOEntry
		scriptPublicKey := utxoEntry.ScriptPublicKey()
		scriptPublicKey.Script = append(scriptPublicKey.Script, 1, 2, 3)
		clone.Inputs[inputIndex].UTXOEntry = utxo.NewUTXOEntry(utxoEntry.Amount(), scriptPublicKey, false, 100)
		return clone
	}
}

func modifySequence(inputIndex int) func(tx *externalapi.DomainTransaction) *externalapi.DomainTransaction {
	return func(tx *externalapi.DomainTransaction) *externalapi.DomainTransaction {
		clone := tx.Clone()
		clone.Inputs[inputIndex].Sequence = 12345
		return clone
	}
}

func modifyPayload(tx *externalapi.DomainTransaction) *externalapi.DomainTransaction {
	clone := tx.Clone()
	clone.Payload = []byte{6, 6, 6, 4, 2, 0, 1, 3, 3, 7}
	return clone
}

func modifyGas(tx *externalapi.DomainTransaction) *externalapi.DomainTransaction {
	clone := tx.Clone()
	clone.Gas = 1234
	return clone
}

func modifySubnetworkID(tx *externalapi.DomainTransaction) *externalapi.DomainTransaction {
	clone := tx.Clone()
	clone.SubnetworkID = externalapi.DomainSubnetworkID{6, 6, 6, 4, 2, 0, 1, 3, 3, 7}
	return clone
}

func TestCalculateSignatureHashSchnorr(t *testing.T) {
	nativeTx, subnetworkTx, err := generateTxs()
	if err != nil {
		t.Fatalf("Error from generateTxs: %+v", err)
	}

	// Note: Expected values were generated by the same code that they test,
	// As long as those were not verified using 3rd-party code they only check for regression, not correctness
	tests := []struct {
		name                  string
		tx                    *externalapi.DomainTransaction
		hashType              consensushashing.SigHashType
		inputIndex            int
		modificationFunction  func(*externalapi.DomainTransaction) *externalapi.DomainTransaction
		expectedSignatureHash string
	}{
		// native transactions

		// sigHashAll
		{name: "native-all-0", tx: nativeTx, hashType: all, inputIndex: 0,
			expectedSignatureHash: "b3b98e28bd332255f113ca542c86658d7891124be669a3030d9fbe98a9802750"},
		{name: "native-all-0-modify-input-1", tx: nativeTx, hashType: all, inputIndex: 0,
			modificationFunction:  modifyInput(1), // should change the hash
			expectedSignatureHash: "d0d8ca04a548abf7dde5e101269bbf9414d5db85ba68523e52d1209af1657963"},
		{name: "native-all-0-modify-output-1", tx: nativeTx, hashType: all, inputIndex: 0,
			modificationFunction:  modifyOutput(1), // should change the hash
			expectedSignatureHash: "f01acf277983d78688fea38056dd7c438aab35c3869776324c84ca6f17f39dd2"},
		{name: "native-all-0-modify-sequence-1", tx: nativeTx, hashType: all, inputIndex: 0,
			modificationFunction:  modifySequence(1), // should change the hash
			expectedSignatureHash: "fa7b02f5df77106b08908b00f994333db43c55480f414e7c28b414e79e37a253"},
		{name: "native-all-anyonecanpay-0", tx: nativeTx, hashType: allAnyoneCanPay, inputIndex: 0,
			expectedSignatureHash: "502d348132a32716d5ee203d4530376bfaa7818fb61bf455b53451af472ad93a"},
		{name: "native-all-anyonecanpay-0-modify-input-0", tx: nativeTx, hashType: allAnyoneCanPay, inputIndex: 0,
			modificationFunction:  modifyInput(0), // should change the hash
			expectedSignatureHash: "b34bb58b8e12b1ec9ab99201a09baf01f47e24f822b3f111927ffe64e6d65ba0"},
		{name: "native-all-anyonecanpay-0-modify-input-1", tx: nativeTx, hashType: allAnyoneCanPay, inputIndex: 0,
			modificationFunction:  modifyInput(1), // shouldn't change the hash
			expectedSignatureHash: "502d348132a32716d5ee203d4530376bfaa7818fb61bf455b53451af472ad93a"},
		{name: "native-all-anyonecanpay-0-modify-sequence", tx: nativeTx, hashType: allAnyoneCanPay, inputIndex: 0,
			modificationFunction:  modifySequence(1), // shouldn't change the hash
			expectedSignatureHash: "502d348132a32716d5ee203d4530376bfaa7818fb61bf455b53451af472ad93a"},

		// sigHashNone
		{name: "native-none-0", tx: nativeTx, hashType: none, inputIndex: 0,
			expectedSignatureHash: "f87e01d8ddb149c946f9531971dcf5415abfa45bd8562e0adc94f4c77e63c662"},
		{name: "native-none-0-modify-output-1", tx: nativeTx, hashType: none, inputIndex: 0,
			modificationFunction:  modifyOutput(1), // shouldn't change the hash
			expectedSignatureHash: "f87e01d8ddb149c946f9531971dcf5415abfa45bd8562e0adc94f4c77e63c662"},
		{name: "native-none-0-modify-sequence-0", tx: nativeTx, hashType: none, inputIndex: 0,
			modificationFunction:  modifySequence(0), // should change the hash
			expectedSignatureHash: "7e3083dd0856c0eda0fb0efa3062e110f8b9339df68909f28001cd985f91acaa"},
		{name: "native-none-0-modify-sequence-1", tx: nativeTx, hashType: none, inputIndex: 0,
			modificationFunction:  modifySequence(1), // shouldn't change the hash
			expectedSignatureHash: "f87e01d8ddb149c946f9531971dcf5415abfa45bd8562e0adc94f4c77e63c662"},
		{name: "native-none-anyonecanpay-0", tx: nativeTx, hashType: noneAnyoneCanPay, inputIndex: 0,
			expectedSignatureHash: "b0830406f267a7d0cda4d4296054c2e93c6d40e81dfedbd83c14bec375b2ec01"},
		{name: "native-none-anyonecanpay-0-modify-amount-spent", tx: nativeTx, hashType: noneAnyoneCanPay, inputIndex: 0,
			modificationFunction:  modifyAmountSpent(0), // should change the hash
			expectedSignatureHash: "64eb64a015d9ca910e7fdeab3e0fe518c716c42646f7a89181723478df904dd3"},
		{name: "native-none-anyonecanpay-0-modify-script-public-key", tx: nativeTx, hashType: noneAnyoneCanPay, inputIndex: 0,
			modificationFunction:  modifyScriptPublicKey(0), // should change the hash
			expectedSignatureHash: "6a1e33e82a2d699674789d19c2e96e0faf27670dc04b5910c67232039da5470d"},

		// sigHashSingle
		{name: "native-single-0", tx: nativeTx, hashType: single, inputIndex: 0,
			expectedSignatureHash: "5b10da8caf353ca63d8fff78ef738a9d2bdb39512a1905781300576a4da0d8db"},
		{name: "native-single-0-modify-output-0", tx: nativeTx, hashType: single, inputIndex: 0,
			modificationFunction:  modifyOutput(0), // should change the hash
			expectedSignatureHash: "26aed36d0bb609beca93b2b29ad22b6f62ca7b6e9ccc6cae28e3dad11bfe951f"},
		{name: "native-single-0-modify-output-1", tx: nativeTx, hashType: single, inputIndex: 0,
			modificationFunction:  modifyOutput(1), // shouldn't change the hash
			expectedSignatureHash: "5b10da8caf353ca63d8fff78ef738a9d2bdb39512a1905781300576a4da0d8db"},
		{name: "native-single-0-modify-sequence-0", tx: nativeTx, hashType: single, inputIndex: 0,
			modificationFunction:  modifySequence(0), // should change the hash
			expectedSignatureHash: "2743f74f58ecc577c853049cb42075a6e6b44251515a6a8daa1f4f679692c81c"},
		{name: "native-single-0-modify-sequence-1", tx: nativeTx, hashType: single, inputIndex: 0,
			modificationFunction:  modifySequence(1), // shouldn't change the hash
			expectedSignatureHash: "5b10da8caf353ca63d8fff78ef738a9d2bdb39512a1905781300576a4da0d8db"},
		{name: "native-single-2-no-corresponding-output", tx: nativeTx, hashType: single, inputIndex: 2,
			expectedSignatureHash: "a6aaa4b8012fdc46b8f64601d25fe47f0ef511fc4230a446603ee2488f9da428"},
		{name: "native-single-2-no-corresponding-output-modify-output-1", tx: nativeTx, hashType: single, inputIndex: 2,
			modificationFunction:  modifyOutput(1), // shouldn't change the hash
			expectedSignatureHash: "a6aaa4b8012fdc46b8f64601d25fe47f0ef511fc4230a446603ee2488f9da428"},
		{name: "native-single-anyonecanpay-0", tx: nativeTx, hashType: singleAnyoneCanPay, inputIndex: 0,
			expectedSignatureHash: "0c2b38530a929e3cd6cf94b6ed7cdf3997adf8370e334212df49777e29796d7d"},
		{name: "native-single-anyonecanpay-2-no-corresponding-output", tx: nativeTx, hashType: singleAnyoneCanPay, inputIndex: 2,
			expectedSignatureHash: "f81f68f88381dd703c1d9c2c7a9f1b366f665db0b73a62c1e1f52f5f06fa19b5"},

		// subnetwork transaction
		{name: "subnetwork-all-0", tx: subnetworkTx, hashType: all, inputIndex: 0,
			expectedSignatureHash: "e4550587d8888be32d23a239d0bc9e4bbf099e5fa972f58b20d0a5ad67615a03"},
		{name: "subnetwork-all-modify-payload", tx: subnetworkTx, hashType: all, inputIndex: 0,
			modificationFunction:  modifyPayload, // should change the hash
			expectedSignatureHash: "167196adb7df0343d39207cdc7505381a3d324755bebb9814853dc3ef03583af"},
		{name: "subnetwork-all-modify-gas", tx: subnetworkTx, hashType: all, inputIndex: 0,
			modificationFunction:  modifyGas, // should change the hash
			expectedSignatureHash: "7b307cad2a13289dc0cdbf0e38247fb05a54c670ea0c950fe13cc51b61d167e6"},
		{name: "subnetwork-all-subnetwork-id", tx: subnetworkTx, hashType: all, inputIndex: 0,
			modificationFunction:  modifySubnetworkID, // should change the hash
			expectedSignatureHash: "c35c5f9ddb64325580703b6582bce6a8f354381af7547e971945d8ba26ad3bb0"},
	}

	for _, test := range tests {
		tx := test.tx
		if test.modificationFunction != nil {
			tx = test.modificationFunction(tx)
		}

		actualSignatureHash, err := consensushashing.CalculateSignatureHashSchnorr(
			tx, test.inputIndex, test.hashType, &consensushashing.SighashReusedValues{})
		if err != nil {
			t.Errorf("%s: Error from CalculateSignatureHashSchnorr: %+v", test.name, err)
			continue
		}

		if actualSignatureHash.String() != test.expectedSignatureHash {
			t.Errorf("%s: expected signature hash: '%s'; but got: '%s'",
				test.name, test.expectedSignatureHash, actualSignatureHash)
		}
	}
}

func TestCalculateSignatureHashECDSA(t *testing.T) {
	nativeTx, subnetworkTx, err := generateTxs()
	if err != nil {
		t.Fatalf("Error from generateTxs: %+v", err)
	}

	// Note: Expected values were generated by the same code that they test,
	// As long as those were not verified using 3rd-party code they only check for regression, not correctness
	tests := []struct {
		name                  string
		tx                    *externalapi.DomainTransaction
		hashType              consensushashing.SigHashType
		inputIndex            int
		modificationFunction  func(*externalapi.DomainTransaction) *externalapi.DomainTransaction
		expectedSignatureHash string
	}{
		// native transactions

		// sigHashAll
		{name: "native-all-0", tx: nativeTx, hashType: all, inputIndex: 0,
			expectedSignatureHash: "2a2f782c03f4d1f50f4c7775940a58b3874d42fa294b2cd8c089af8626716a6a"},
		{name: "native-all-0-modify-input-1", tx: nativeTx, hashType: all, inputIndex: 0,
			modificationFunction:  modifyInput(1), // should change the hash
			expectedSignatureHash: "ea0858d0fbf2f0ed939b578b8ecbacf7f5e60eeaf807e72f0920d9b5437f1214"},
		{name: "native-all-0-modify-output-1", tx: nativeTx, hashType: all, inputIndex: 0,
			modificationFunction:  modifyOutput(1), // should change the hash
			expectedSignatureHash: "eb0bc3c0f2769589019903c60ddda64f84be6f053cfaf0cb67e7d8fd1a43cf56"},
		{name: "native-all-0-modify-sequence-1", tx: nativeTx, hashType: all, inputIndex: 0,
			modificationFunction:  modifySequence(1), // should change the hash
			expectedSignatureHash: "d20a523aea828e89369bddcc348f88757f0847325a8ecd4cb2768768c005b2ef"},
		{name: "native-all-anyonecanpay-0", tx: nativeTx, hashType: allAnyoneCanPay, inputIndex: 0,
			expectedSignatureHash: "002b08d32620f420114c77822f988e5ac5bb85c41fb3bbcd9b22031a58af87b6"},
		{name: "native-all-anyonecanpay-0-modify-input-0", tx: nativeTx, hashType: allAnyoneCanPay, inputIndex: 0,
			modificationFunction:  modifyInput(0), // should change the hash
			expectedSignatureHash: "393f527931bb533dff52c3da2df23c42e267da5f056d82e8781fc95f08ef906c"},
		{name: "native-all-anyonecanpay-0-modify-input-1", tx: nativeTx, hashType: allAnyoneCanPay, inputIndex: 0,
			modificationFunction:  modifyInput(1), // shouldn't change the hash
			expectedSignatureHash: "002b08d32620f420114c77822f988e5ac5bb85c41fb3bbcd9b22031a58af87b6"},
		{name: "native-all-anyonecanpay-0-modify-sequence", tx: nativeTx, hashType: allAnyoneCanPay, inputIndex: 0,
			modificationFunction:  modifySequence(1), // shouldn't change the hash
			expectedSignatureHash: "002b08d32620f420114c77822f988e5ac5bb85c41fb3bbcd9b22031a58af87b6"},

		// sigHashNone
		{name: "native-none-0", tx: nativeTx, hashType: none, inputIndex: 0,
			expectedSignatureHash: "7d1f83ebf610f59826a4f5d0be3abd6a6b9fa91cfcde11ecde785a9a62717167"},
		{name: "native-none-0-modify-output-1", tx: nativeTx, hashType: none, inputIndex: 0,
			modificationFunction:  modifyOutput(1), // shouldn't change the hash
			expectedSignatureHash: "7d1f83ebf610f59826a4f5d0be3abd6a6b9fa91cfcde11ecde785a9a62717167"},
		{name: "native-none-0-modify-sequence-0", tx: nativeTx, hashType: none, inputIndex: 0,
			modificationFunction:  modifySequence(0), // should change the hash
			expectedSignatureHash: "7ebb530e8b46fc4312d7ee63eb4373e5a87ea496fd671bd2cf5f51c3d41951b6"},
		{name: "native-none-0-modify-sequence-1", tx: nativeTx, hashType: none, inputIndex: 0,
			modificationFunction:  modifySequence(1), // shouldn't change the hash
			expectedSignatureHash: "7d1f83ebf610f59826a4f5d0be3abd6a6b9fa91cfcde11ecde785a9a62717167"},
		{name: "native-none-anyonecanpay-0", tx: nativeTx, hashType: noneAnyoneCanPay, inputIndex: 0,
			expectedSignatureHash: "993ec593e006a7dc2718de346e570224b5489d79c3ab23005e5e19880ad2549e"},
		{name: "native-none-anyonecanpay-0-modify-amount-spent", tx: nativeTx, hashType: noneAnyoneCanPay, inputIndex: 0,
			modificationFunction:  modifyAmountSpent(0), // should change the hash
			expectedSignatureHash: "d81416c4021b3f5323c4c9f8e85cdf35a0e14d506fd2327e68a2304d1748a05b"},
		{name: "native-none-anyonecanpay-0-modify-script-public-key", tx: nativeTx, hashType: noneAnyoneCanPay, inputIndex: 0,
			modificationFunction:  modifyScriptPublicKey(0), // should change the hash
			expectedSignatureHash: "d94de4924bd7859527c8d87f11e954ad93601d54d7921b525f15cf3005eaaeb2"},

		// sigHashSingle
		{name: "native-single-0", tx: nativeTx, hashType: single, inputIndex: 0,
			expectedSignatureHash: "4ba903d849fa29e31960a9300a4e51ba2edf23fc1cb76e13c3324a969931644e"},
		{name: "native-single-0-modify-output-0", tx: nativeTx, hashType: single, inputIndex: 0,
			modificationFunction:  modifyOutput(0), // should change the hash
			expectedSignatureHash: "24becf355af5949ea5bb56369deb8d204345c8145beae624c9769bd141d8a4a3"},
		{name: "native-single-0-modify-output-1", tx: nativeTx, hashType: single, inputIndex: 0,
			modificationFunction:  modifyOutput(1), // shouldn't change the hash
			expectedSignatureHash: "4ba903d849fa29e31960a9300a4e51ba2edf23fc1cb76e13c3324a969931644e"},
		{name: "native-single-0-modify-sequence-0", tx: nativeTx, hashType: single, inputIndex: 0,
			modificationFunction:  modifySequence(0), // should change the hash
			expectedSignatureHash: "20b9231ca5f4b86fc36f66782ce7abe137a1b0f23958a2c1eb5dd91e4b2f3fec"},
		{name: "native-single-0-modify-sequence-1", tx: nativeTx, hashType: single, inputIndex: 0,
			modificationFunction:  modifySequence(1), // shouldn't change the hash
			expectedSignatureHash: "4ba903d849fa29e31960a9300a4e51ba2edf23fc1cb76e13c3324a969931644e"},
		{name: "native-single-2-no-corresponding-output", tx: nativeTx, hashType: single, inputIndex: 2,
			expectedSignatureHash: "ff91fd8342d8fdf15dd0fa06194d35492a9d2a45cf14fee694012f1a379dadd1"},
		{name: "native-single-2-no-corresponding-output-modify-output-1", tx: nativeTx, hashType: single, inputIndex: 2,
			modificationFunction:  modifyOutput(1), // shouldn't change the hash
			expectedSignatureHash: "ff91fd8342d8fdf15dd0fa06194d35492a9d2a45cf14fee694012f1a379dadd1"},
		{name: "native-single-anyonecanpay-0", tx: nativeTx, hashType: singleAnyoneCanPay, inputIndex: 0,
			expectedSignatureHash: "0e6dc0fea8cfce7ffd870cb6e6479d6001d75f46881c4efc9e830e84f618d731"},
		{name: "native-single-anyonecanpay-2-no-corresponding-output", tx: nativeTx, hashType: singleAnyoneCanPay, inputIndex: 2,
			expectedSignatureHash: "87fe06aa0fde636a9c6ef032c9ce63e6c11ff8fe9bb229c298ddefe7e15dcb57"},

		// subnetwork transaction
		{name: "subnetwork-all-0", tx: subnetworkTx, hashType: all, inputIndex: 0,
			expectedSignatureHash: "14355ef873b99c90c9674107446f9753d72b2c5716cef7e196289b7d50887e57"},
		{name: "subnetwork-all-modify-payload", tx: subnetworkTx, hashType: all, inputIndex: 0,
			modificationFunction:  modifyPayload, // should change the hash
			expectedSignatureHash: "ef0e995b07743a00e415eef79f927703f007f3ca7fd12dc1c7824b289c91cfaf"},
		{name: "subnetwork-all-modify-gas", tx: subnetworkTx, hashType: all, inputIndex: 0,
			modificationFunction:  modifyGas, // should change the hash
			expectedSignatureHash: "5a56d5ef318aca0183c026bfef20142b3e51cc8370b8a5ad72e2b3bfec1229d5"},
		{name: "subnetwork-all-subnetwork-id", tx: subnetworkTx, hashType: all, inputIndex: 0,
			modificationFunction:  modifySubnetworkID, // should change the hash
			expectedSignatureHash: "a7c43433eb809d95f68398a476bb9498687db6c777a863c57bfde26f9c66d493"},
	}

	for _, test := range tests {
		tx := test.tx
		if test.modificationFunction != nil {
			tx = test.modificationFunction(tx)
		}

		actualSignatureHash, err := consensushashing.CalculateSignatureHashECDSA(
			tx, test.inputIndex, test.hashType, &consensushashing.SighashReusedValues{})
		if err != nil {
			t.Errorf("%s: Error from CalculateSignatureHashECDSA: %+v", test.name, err)
			continue
		}

		if actualSignatureHash.String() != test.expectedSignatureHash {
			t.Errorf("%s: expected signature hash: '%s'; but got: '%s'",
				test.name, test.expectedSignatureHash, actualSignatureHash)
		}
	}
}

func TestSignatureHash(t *testing.T) {
	genesisCoinbase := dagconfig.MainnetParams.GenesisBlock.Transactions[0]
	genesisCoinbaseTransactionID := consensushashing.TransactionID(genesisCoinbase)
	fmt.Printf("%s\r\n", genesisCoinbaseTransactionID)

	address1Str := "lings:qzdz9t3j5pmrmqwzj3htquk82kxj77uqdh47ucq23wv00evkxnff6sefjdrg4"
	address1, err := util.DecodeAddress(address1Str, util.Bech32PrefixLings)
	if err != nil {
		t.Errorf("error decoding address1: %+v", err)
	}
	address1ToScript, err := txscript.PayToAddrScript(address1)
	if err != nil {
		t.Errorf("error generating script: %+v", err)
	}

	address2Str := "lings:qrt7lyuplddukhjnm8wlu9dy3npz2j58a6qhpm9cmdawfw5ymmhxgpvmwfnsp"
	address2, err := util.DecodeAddress(address2Str, util.Bech32PrefixLings)
	if err != nil {
		t.Errorf("error decoding address2: %+v", err)
	}
	address2ToScript, err := txscript.PayToAddrScript(address2)
	if err != nil {
		t.Errorf("error generating script: %+v", err)
	}

	txIns := []*externalapi.DomainTransactionInput{
		{
			PreviousOutpoint: *externalapi.NewDomainOutpoint(genesisCoinbaseTransactionID, 1),
			Sequence:         1,
			UTXOEntry:        utxo.NewUTXOEntry(200, address1ToScript, false, 0),
		},
	}
	txOuts := []*externalapi.DomainTransactionOutput{
		{
			Value:           300,
			ScriptPublicKey: address2ToScript,
		},
		{
			Value:           300,
			ScriptPublicKey: address1ToScript,
		},
	}
	nativeTx := &externalapi.DomainTransaction{
		Version:      0,
		Inputs:       txIns,
		Outputs:      txOuts,
		LockTime:     1615462089000,
		SubnetworkID: subnetworks.SubnetworkIDNative,
	}
	t.Errorf("er %s", nativeTx)
}

func generateTxs() (nativeTx, subnetworkTx *externalapi.DomainTransaction, err error) {
	genesisCoinbase := dagconfig.MainnetParams.GenesisBlock.Transactions[0]
	genesisCoinbaseTransactionID := consensushashing.TransactionID(genesisCoinbase)
	fmt.Printf("%s\r\n", genesisCoinbaseTransactionID)

	address1Str := "lings:qrt7lyuplddukhjnm8wlu9dy3npz2j58a6qhpm9cmdawfw5ymmhxgpvmwfnsp"
	address1, err := util.DecodeAddress(address1Str, util.Bech32PrefixLings)
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding address1: %+v", err)
	}
	address1ToScript, err := txscript.PayToAddrScript(address1)
	if err != nil {
		return nil, nil, fmt.Errorf("error generating script: %+v", err)
	}

	address2Str := "lings:qrt7lyuplddukhjnm8wlu9dy3npz2j58a6qhpm9cmdawfw5ymmhxgpvmwfnsp"
	address2, err := util.DecodeAddress(address2Str, util.Bech32PrefixLings)
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding address2: %+v", err)
	}
	address2ToScript, err := txscript.PayToAddrScript(address2)
	if err != nil {
		return nil, nil, fmt.Errorf("error generating script: %+v", err)
	}

	txIns := []*externalapi.DomainTransactionInput{
		{
			PreviousOutpoint: *externalapi.NewDomainOutpoint(genesisCoinbaseTransactionID, 0),
			Sequence:         0,
			UTXOEntry:        utxo.NewUTXOEntry(100, address1ToScript, false, 0),
		},
		{
			PreviousOutpoint: *externalapi.NewDomainOutpoint(genesisCoinbaseTransactionID, 1),
			Sequence:         1,
			UTXOEntry:        utxo.NewUTXOEntry(200, address2ToScript, false, 0),
		},
		{
			PreviousOutpoint: *externalapi.NewDomainOutpoint(genesisCoinbaseTransactionID, 2),
			Sequence:         2,
			UTXOEntry:        utxo.NewUTXOEntry(300, address2ToScript, false, 0),
		},
	}

	txOuts := []*externalapi.DomainTransactionOutput{
		{
			Value:           300,
			ScriptPublicKey: address2ToScript,
		},
		{
			Value:           300,
			ScriptPublicKey: address1ToScript,
		},
	}

	nativeTx = &externalapi.DomainTransaction{
		Version:      0,
		Inputs:       txIns,
		Outputs:      txOuts,
		LockTime:     1615462089000,
		SubnetworkID: subnetworks.SubnetworkIDNative,
	}
	subnetworkTx = &externalapi.DomainTransaction{
		Version:      0,
		Inputs:       txIns,
		Outputs:      txOuts,
		LockTime:     1615462089000,
		SubnetworkID: externalapi.DomainSubnetworkID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		Gas:          250,
		Payload:      []byte{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
	}

	return nativeTx, subnetworkTx, nil
}

func BenchmarkCalculateSignatureHashSchnorr(b *testing.B) {
	sigHashTypes := []consensushashing.SigHashType{
		consensushashing.SigHashAll,
		consensushashing.SigHashNone,
		consensushashing.SigHashSingle,
		consensushashing.SigHashAll | consensushashing.SigHashAnyOneCanPay,
		consensushashing.SigHashNone | consensushashing.SigHashAnyOneCanPay,
		consensushashing.SigHashSingle | consensushashing.SigHashAnyOneCanPay}

	for _, size := range []int{10, 100, 1000} {
		tx := generateTransaction(b, sigHashTypes, size)

		b.Run(fmt.Sprintf("%d-inputs-and-outputs", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				reusedValues := &consensushashing.SighashReusedValues{}
				for inputIndex := range tx.Inputs {
					sigHashType := sigHashTypes[inputIndex%len(sigHashTypes)]
					_, err := consensushashing.CalculateSignatureHashSchnorr(tx, inputIndex, sigHashType, reusedValues)
					if err != nil {
						b.Fatalf("Error from CalculateSignatureHashSchnorr: %+v", err)
					}
				}
			}
		})
	}
}

func generateTransaction(b *testing.B, sigHashTypes []consensushashing.SigHashType, inputAndOutputSizes int) *externalapi.DomainTransaction {
	sourceScript := getSourceScript(b)
	tx := &externalapi.DomainTransaction{
		Version:      0,
		Inputs:       generateInputs(inputAndOutputSizes, sourceScript),
		Outputs:      generateOutputs(inputAndOutputSizes, sourceScript),
		LockTime:     123456789,
		SubnetworkID: externalapi.DomainSubnetworkID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		Gas:          125,
		Payload:      []byte{9, 8, 7, 6, 5, 4, 3, 2, 1},
		Fee:          0,
		Mass:         0,
		ID:           nil,
	}
	signTx(b, tx, sigHashTypes)
	return tx
}

func signTx(b *testing.B, tx *externalapi.DomainTransaction, sigHashTypes []consensushashing.SigHashType) {
	sourceAddressPKStr := "a4d85b7532123e3dd34e58d7ce20895f7ca32349e29b01700bb5a3e72d2570eb"
	privateKeyBytes, err := hex.DecodeString(sourceAddressPKStr)
	if err != nil {
		b.Fatalf("Error parsing private key hex: %+v", err)
	}
	keyPair, err := secp256k1.DeserializeSchnorrPrivateKeyFromSlice(privateKeyBytes)
	if err != nil {
		b.Fatalf("Error deserializing private key: %+v", err)
	}
	for i, txIn := range tx.Inputs {
		signatureScript, err := txscript.SignatureScript(
			tx, i, sigHashTypes[i%len(sigHashTypes)], keyPair, &consensushashing.SighashReusedValues{})
		if err != nil {
			b.Fatalf("Error from SignatureScript: %+v", err)
		}
		txIn.SignatureScript = signatureScript
	}

}

func generateInputs(size int, sourceScript *externalapi.ScriptPublicKey) []*externalapi.DomainTransactionInput {
	inputs := make([]*externalapi.DomainTransactionInput, size)

	for i := 0; i < size; i++ {
		inputs[i] = &externalapi.DomainTransactionInput{
			PreviousOutpoint: *externalapi.NewDomainOutpoint(
				externalapi.NewDomainTransactionIDFromByteArray(&[32]byte{12, 3, 4, 5}), 1),
			SignatureScript: nil,
			Sequence:        uint64(i),
			UTXOEntry:       utxo.NewUTXOEntry(uint64(i), sourceScript, false, 12),
		}
	}

	return inputs
}

func getSourceScript(b *testing.B) *externalapi.ScriptPublicKey {
	sourceAddressStr := "lingssim:qz6f9z6l3x4v3lf9mgf0t934th4nx5kgzu663x9yjh"

	sourceAddress, err := util.DecodeAddress(sourceAddressStr, util.Bech32PrefixLingsSim)
	if err != nil {
		b.Fatalf("Error from DecodeAddress: %+v", err)
	}

	sourceScript, err := txscript.PayToAddrScript(sourceAddress)
	if err != nil {
		b.Fatalf("Error from PayToAddrScript: %+v", err)
	}
	return sourceScript
}

func generateOutputs(size int, script *externalapi.ScriptPublicKey) []*externalapi.DomainTransactionOutput {
	outputs := make([]*externalapi.DomainTransactionOutput, size)

	for i := 0; i < size; i++ {
		outputs[i] = &externalapi.DomainTransactionOutput{
			Value:           uint64(i),
			ScriptPublicKey: script,
		}
	}

	return outputs
}
