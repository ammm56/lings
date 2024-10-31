wire
====

[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](https://choosealicense.com/licenses/isc/)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/ammm56/lings/wire)
=======

Package wire implements the lings wire protocol.

## Hoosat Message Overview

The lings protocol consists of exchanging messages between peers. Each message
is preceded by a header which identifies information about it such as which
lings network it is a part of, its type, how big it is, and a checksum to
verify validity. All encoding and decoding of message headers is handled by this
package.

To accomplish this, there is a generic interface for lings messages named
`Message` which allows messages of any type to be read, written, or passed
around through channels, functions, etc. In addition, concrete implementations
of most all lings messages are provided. All of the details of marshalling and 
unmarshalling to and from the wire using lings encoding are handled so the 
caller doesn't have to concern themselves with the specifics.

## Reading Messages Example

In order to unmarshal lings messages from the wire, use the `ReadMessage`
function. It accepts any `io.Reader`, but typically this will be a `net.Conn`
to a remote node running a lings peer. Example syntax is:

```Go
	// Use the most recent protocol version supported by the package and the
	// main lings network.
	pver := wire.ProtocolVersion
	lingsnet := wire.Mainnet

	// Reads and validates the next lings message from conn using the
	// protocol version pver and the lings network lingsnet. The returns
	// are a appmessage.Message, a []byte which contains the unmarshalled
	// raw payload, and a possible error.
	msg, rawPayload, err := wire.ReadMessage(conn, pver, lingsnet)
	if err != nil {
		// Log and handle the error
	}
```

See the package documentation for details on determining the message type.

## Writing Messages Example

In order to marshal lings messages to the wire, use the `WriteMessage`
function. It accepts any `io.Writer`, but typically this will be a `net.Conn`
to a remote node running a lings peer. Example syntax to request addresses
from a remote peer is:

```Go
	// Use the most recent protocol version supported by the package and the
	// main bitcoin network.
	pver := wire.ProtocolVersion
	lingsnet := wire.Mainnet

	// Create a new getaddr lings message.
	msg := wire.NewMsgGetAddr()

	// Writes a lings message msg to conn using the protocol version
	// pver, and the lings network lingsnet. The return is a possible
	// error.
	err := wire.WriteMessage(conn, msg, pver, lingsnet)
	if err != nil {
		// Log and handle the error
	}
```
