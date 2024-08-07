# bech32m

A Go package for encoding and decoding bech32m addresses.

## Overview

This package provides functionality to work with bech32m addresses, including
encoding, decoding, and validation. It supports creating addresses from byte
arrays and ED25519 public keys.

## Installation

To install the package, use:

```sh
go get github.com/astriaorg/astria-cli-go/modules/bech32m
```

## Usage

### Importing the package

```go
import "github.com/astriaorg/astria-cli-go/modules/bech32m"
```

### Creating an address

From bytes:

```go
prefix := "cosmos"
var data [20]byte
// ... fill data ...
address, err := bech32m.EncodeFromBytes(prefix, data)
if err != nil {
    // Handle error
}
```

From ED25519 public key:

```go
prefix := "cosmos"
pubkey := ed25519.PublicKey(...)
address, err := bech32m.EncodeFromPublicKey(prefix, pubkey)
if err != nil {
    // Handle error
}
```

### Validating an address

```go
err := bech32m.Validate("astria1...")
if err != nil {
    // Address is invalid
}
```

### Decoding an address

```go
prefix, bytes, err := bech32m.DecodeFromString("astria1...")
if err != nil {
    // Handle error
}
```

### Working with Address struct

```go
// Get string representation
addressStr := address.String()

// Get prefix
prefix := address.Prefix()

// Get underlying bytes
bytes := address.Bytes()
```
