package bip32

import "github.com/pkg/errors"

// BitcoinMainnetPrivate is the version that is used for
// bitcoin mainnet bip32 private extended keys.
// Ecnodes to xprv in base58.
var BitcoinMainnetPrivate = [4]byte{
	0x04,
	0x88,
	0xad,
	0xe4,
}

// BitcoinMainnetPublic is the version that is used for
// bitcoin mainnet bip32 public extended keys.
// Ecnodes to xpub in base58.
var BitcoinMainnetPublic = [4]byte{
	0x04,
	0x88,
	0xb2,
	0x1e,
}

// LingsMainnetPrivate is the version that is used for
// lings mainnet bip32 private extended keys.
// Ecnodes to xprv in base58.
var LingsMainnetPrivate = [4]byte{
	0x03,
	0x8f,
	0x2e,
	0xf4,
}

// LingsMainnetPublic is the version that is used for
// lings mainnet bip32 public extended keys.
// Ecnodes to kpub in base58.
var LingsMainnetPublic = [4]byte{
	0x03,
	0x8f,
	0x33,
	0x2e,
}

// LingsTestnetPrivate is the version that is used for
// lings testnet bip32 public extended keys.
// Ecnodes to ktrv in base58.
var LingsTestnetPrivate = [4]byte{
	0x03,
	0x90,
	0x9e,
	0x07,
}

// LingsTestnetPublic is the version that is used for
// lings testnet bip32 public extended keys.
// Ecnodes to ktub in base58.
var LingsTestnetPublic = [4]byte{
	0x03,
	0x90,
	0xa2,
	0x41,
}

// LingsevnetPrivate is the version that is used for
// lings devnet bip32 public extended keys.
// Ecnodes to kdrv in base58.
var LingsevnetPrivate = [4]byte{
	0x03,
	0x8b,
	0x3d,
	0x80,
}

// LingsevnetPublic is the version that is used for
// lings devnet bip32 public extended keys.
// Ecnodes to xdub in base58.
var LingsevnetPublic = [4]byte{
	0x03,
	0x8b,
	0x41,
	0xba,
}

// LingsSimnetPrivate is the version that is used for
// lings simnet bip32 public extended keys.
// Ecnodes to ksrv in base58.
var LingsSimnetPrivate = [4]byte{
	0x03,
	0x90,
	0x42,
	0x42,
}

// LingsSimnetPublic is the version that is used for
// lings simnet bip32 public extended keys.
// Ecnodes to xsub in base58.
var LingsSimnetPublic = [4]byte{
	0x03,
	0x90,
	0x46,
	0x7d,
}

func toPublicVersion(version [4]byte) ([4]byte, error) {
	switch version {
	case BitcoinMainnetPrivate:
		return BitcoinMainnetPublic, nil
	case LingsMainnetPrivate:
		return LingsMainnetPublic, nil
	case LingsTestnetPrivate:
		return LingsTestnetPublic, nil
	case LingsevnetPrivate:
		return LingsevnetPublic, nil
	case LingsSimnetPrivate:
		return LingsSimnetPublic, nil
	}

	return [4]byte{}, errors.Errorf("unknown version %x", version)
}

func isPrivateVersion(version [4]byte) bool {
	switch version {
	case BitcoinMainnetPrivate:
		return true
	case LingsMainnetPrivate:
		return true
	case LingsTestnetPrivate:
		return true
	case LingsevnetPrivate:
		return true
	case LingsSimnetPrivate:
		return true
	}

	return false
}
