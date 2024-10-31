package constants

import "math"

var (
	// BlockVersion represents the current block
	// 1 Pyrinhash
	// 2 HoohashV1
	// 3 HoohashV2
	BlockVersion uint16 = 1
)

const (
	DevFee        = 5
	DevFeeMin     = 1
	DevFeeAddress = "lings:qp4ad2eh72xc8dtjjyz4llxzq9utn6k26uyl644xxw70wskdfl85zsqj9k4vz"

	// MaxTransactionVersion is the current latest supported transaction version.
	MaxTransactionVersion uint16 = 0

	// MaxScriptPublicKeyVersion is the current latest supported public key script version.
	MaxScriptPublicKeyVersion uint16 = 0

	// SompiPerLings is the number of sompi in one lings (1 LSN).
	SompiPerLings = 100_000_000

	// MaxSompi is the maximum transaction amount allowed in sompi.
	MaxSompi = uint64(17_100_000_000 * SompiPerLings)

	// MaxTxInSequenceNum is the maximum sequence number the sequence field
	// of a transaction input can be.
	MaxTxInSequenceNum uint64 = math.MaxUint64

	// SequenceLockTimeDisabled is a flag that if set on a transaction
	// input's sequence number, the sequence number will not be interpreted
	// as a relative locktime.
	SequenceLockTimeDisabled uint64 = 1 << 63

	// SequenceLockTimeMask is a mask that extracts the relative locktime
	// when masked against the transaction input sequence number.
	SequenceLockTimeMask uint64 = 0x00000000ffffffff

	// LockTimeThreshold is the number below which a lock time is
	// interpreted to be a DAA score.
	LockTimeThreshold = 5e11 // Tue Nov 5 00:53:20 1985 UTC

	// UnacceptedDAAScore is used to for UTXOEntries that were created by transactions in the mempool, or otherwise
	// not-yet-accepted transactions.
	UnacceptedDAAScore = math.MaxUint64
)
