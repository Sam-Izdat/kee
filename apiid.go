package kee

import (
    "math/big"
)

// KAPIID type represents an arbitrary precision integer identifier.
// It is exported only for reference and should be instantiated through its handler's methods.
type KAPIID struct {
    slc []byte
    b58 string
    bigInt *big.Int
}

type apiidConfig struct {
    Cache bool
}

// APIIDOptions defines the configuration used by the `kee.APIID` handler.
// Options can also be changed through `kee.APIID.Options`.
var APIIDOptions = apiidConfig {
    Cache: true,            // Cache APIID strings, ignore new options
}

// APIIDCtrl is a struct for the APIID handler. 
// Unless another handler with different options is needed simply use instance `kee.APIID`.
type APIIDCtrl struct {
    Options *apiidConfig
}

// FromString takes string representation of arbitrary precision integer and return KAPIID instance
func (c APIIDCtrl) FromString(s string) KAPIID {
    i := new(big.Int)
    i.SetString(s, 10)
    return KAPIID{slc: i.Bytes(), bigInt: i}
}

// FromInt takes 64-bit integer and return KAPIID instance
func (c APIIDCtrl) FromInt(fpi uint64) KAPIID {
    i := new(big.Int)
    i.SetUint64(fpi)
    return KAPIID{slc: i.Bytes(), bigInt: i}
}

// FromBigInt takes math/big Int and return KAPIID instance
func (c APIIDCtrl) FromBigInt(api *big.Int) KAPIID {
    i := new(big.Int)
    i.Abs(api)
    return KAPIID{slc: i.Bytes(), bigInt: i}
}

// Set takes an arbitrary-length byte slice and returns KAPIID instance
func (c APIIDCtrl) Set(slc []byte) KAPIID {
    i := new(big.Int)
    i.SetBytes(slc)
    return KAPIID{slc: i.Bytes(), bigInt: i}
}

// Decode takes base 58 encoded string of APIID and returns KAPIID instance
func (c APIIDCtrl) Decode(s string) (KAPIID, error) {
    i, err := b58ToBigInt([]byte(s))
    return KAPIID{slc: i.Bytes(), bigInt: i}, err
}

// -- Produce --

// String is alias for B58()
func (id KAPIID) String() string {
    return id.B58()
}

// Slc returns APIID as slice
func (id KAPIID) Slc() []byte {
    return id.slc
}

// BigInt returns APIID as a math/big Int
func (id KAPIID) BigInt() (res *big.Int) {
    if id.slc == nil || len(id.slc) == 0    { return new(big.Int) }
    return id.bigInt
}

// B58 returns base 58 encoded string representation of APIID
func (id KAPIID) B58() (res string) {
    if id.slc == nil || len(id.slc) == 0    { return "" }
    if APIIDOptions.Cache && id.b58 != "" { return id.b58 }
    id.b58 = string(bigIntToB58(nil, id.bigInt))
    return id.b58
}