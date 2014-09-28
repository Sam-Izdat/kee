package kee

import (
    "math/big"
)

// APIID (no standard)
type apiid struct {
    slc []byte
    b58 string
    bigInt *big.Int
}

type apiidConfig struct {
    Cache bool
}

var apiidOptions = apiidConfig {
    Cache: true,            // Cache APIID strings, ignore new options
}

type apiidCtrl struct {
    Options *apiidConfig
}

// Takes string rep. of arbitrary precision integer and return APIID "object"
func (c apiidCtrl) FromString(s string) apiid {
    i := new(big.Int)
    i.SetString(s, 10)
    return apiid{slc: i.Bytes(), bigInt: i}
}

// Takes 64-bit integer and return APIID "object"
func (c apiidCtrl) FromInt(fpi uint64) apiid {
    i := new(big.Int)
    i.SetUint64(fpi)
    return apiid{slc: i.Bytes(), bigInt: i}
}

// Takes big integer and return APIID "object"
func (c apiidCtrl) FromBigInt(api *big.Int) apiid {
    i := new(big.Int)
    i.Abs(api)
    return apiid{slc: i.Bytes(), bigInt: i}
}

// Takes an arbitrary length byte slice and returns APIID "object"
func (c apiidCtrl) Set(slc []byte) apiid {
    i := new(big.Int)
    i.SetBytes(slc)
    return apiid{slc: i.Bytes(), bigInt: i}
}

// Decode APIID from base 58 string and returns APIID "object"
func (c apiidCtrl) Decode(s string) (apiid, error) {
    i, err := b58ToBigInt([]byte(s))
    return apiid{slc: i.Bytes(), bigInt: i}, err
}

// -- Produce --

// Alias for apiid.B58()
func (id apiid) String() string {
    return id.B58()
}

// Returns FPIID as slice
func (id apiid) Slc() []byte {
    return id.slc
}

// Returns FPIID as array
func (id apiid) BigInt() (res *big.Int) {
    if id.slc == nil || len(id.slc) == 0    { return new(big.Int) }
    return id.bigInt
}

// Returns FPIID as array
func (id apiid) B58() (res string) {
    if id.slc == nil || len(id.slc) == 0    { return "" }
    if apiidOptions.Cache && id.b58 != "" { return id.b58 }
    id.b58 = string(bigIntToB58(nil, id.bigInt))
    return id.b58
}