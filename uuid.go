package kee

import(
    "encoding/base64"
    "encoding/base32"
    "encoding/ascii85"
    "strings"
    "errors"
    "fmt"
)

// UUID (RFC 4122)
type uuid struct {
    slc []byte
    hex string
    a85 string
    b64 string
    b32 string
    urn string
    url64 string
    url32 string
}

type uuidConfig struct {
    Cache, AllowInvalid bool
    MinVer, MaxVer uint8 
    PadB64, PadB32, WrapA85, HyphURL32 bool
}

var uuidOptions = uuidConfig {
    Cache: true,            // Cache UUID strings, ignore new options
    AllowInvalid: false,    // Allows setting of non-standard UUIDs
    MinVer: 1,              // Lowest UUID version allowed as valid
    MaxVer: 5,              // Highest UUID version allowed as valid
    PadB64: true,           // Add padding to base 64 encoded UUIDs
    PadB32: true,           // Add padding to base 32 encoded UUIDs
    WrapA85: false,         // Wrap ASCII 85 encoded UUIDs with <~ ~>
    HyphURL32: true,        // Hyphenate base 32 encoded URL UUIDs
}

type uuidCtrl struct {
    Options *uuidConfig
    NS map[string]string    // Namespaces
}

func (c uuidCtrl) newInst(bytes []byte, err error) (uuid, error) {
    res := uuid{slc: bytes}
    if err != nil { // A parsing or other unrecoverable error occured
        return uuid{}, err
    }
    if !uuidOptions.AllowInvalid && !res.IsValid() { 
        if len(res.slc) > 0 && res.Arr() == [16]byte{} { 
            // Allow NIL UUID but return error if no override
            return res, errors.New("nil UUID set")
        } 
        return uuid{}, errors.New("invalid UUID")
    }
    return res, nil
}

// Generates a random UUID and returns UUID "object" - alias for NewV4()
func (c uuidCtrl) New() uuid {
    res, _ := c.NewV4() // swallows errors but none should occur
    return res
}

// Takes a 16 byte array and returns UUID "object"
func (c uuidCtrl) Set(arr [16]byte) uuid {
    bytes := make([]byte, 16)
    bytes = arr[:]
    res, _ := c.newInst(bytes, nil)
    return res
}

// Decodes UUID from string and returns UUID "object"
func (c uuidCtrl) Decode(s string) (uuid, error) {
    var bytes []byte
    var err error
    switch len(s) {
    case 20: 
        bytes, err = c.fromA85(s)
    case 22:
        bytes, err = c.fromB64(s)
    case 24:
        if(s[:2] == "<~" && s[22:] == "~>") {
            bytes, err = c.fromA85(s)
        } else {
            bytes, err = c.fromB64(s)
        }
    case 26, 26+6: 
        bytes, err = c.fromB32(s)
    case 36, 36+9:
        bytes, err = c.fromHex(s)
    default:
        return c.newInst([]byte{}, errors.New("unrecognized UUID encoding"))
    }
    return c.newInst(bytes, err)
}

func (_ uuidCtrl) Match(ida, idb uuid) bool {
    return ida.Arr() == idb.Arr()
}


func (id uuid) IsValid() (valid bool) {
    if len(id.slc) != 16 { return false }
    ver := id.Version()
    if uint8(ver) < uuidOptions.MinVer || uint8(ver) > uuidOptions.MaxVer { 
        return false 
    }
    return true
}

// -- Produce --

// Alias for uuid.Hex()
func (id uuid) String() string {
    return id.Hex()
}

// Returns UUID as slice
func (id uuid) Slc() []byte {
    return id.slc
}

// Returns UUID as array
func (id uuid) Arr() (res [16]byte) {
    copy(res[:], id.slc[:])
    return 
}

// Generates canonical hex string representation, as in RFC 4122
func (id *uuid) Hex() string {
    if id.slc == nil || len(id.slc) == 0    { return "" }
    if uuidOptions.Cache && id.hex != ""    { return id.hex }
    u := id.slc
    id.hex = fmt.Sprintf(
        "%08x-%04x-%04x-%04x-%012x",
        u[:4], u[4:6], u[6:8], u[8:10], u[10:])
    return id.hex
}

// Generates ASCII 85 encoded string representation of UUID
func (id *uuid) A85() string {
    if id.slc == nil || len(id.slc) == 0    { return "" }
    if uuidOptions.Cache && id.a85 != ""    { return id.a85 }
    bytes := make([]byte, 20)
    ascii85.Encode(bytes, id.slc)
    if uuidOptions.WrapA85 {
        parts := []string{"<~", string(bytes[:]), "~>"}
        id.a85 = strings.Join(parts, "")
    } else { id.a85 = string(bytes) }    
    return id.a85
}
// Generates base 64 encoded string representation of UUID
func (id *uuid) B64() string {
    if id.slc == nil || len(id.slc) == 0    { return "" }
    if uuidOptions.Cache && id.b64 != ""    { return id.b64 }
    res := base64.StdEncoding.EncodeToString(id.slc)
    if !uuidOptions.PadB64 { res = res[0:22] }
    id.b64 = res
    return id.b64
}

// Generates base 32 encoded string representation of UUID
func (id *uuid) B32() string {
    if id.slc == nil || len(id.slc) == 0    { return "" }
    if uuidOptions.Cache && id.b32 != ""    { return id.b32 }
    res := base32.StdEncoding.EncodeToString(id.slc)
    if !uuidOptions.PadB32 { res = res[0:26] }
    id.b32 = res
    return id.b32
}

// Generates hex URN of UUID, as in RFC 2141
func (id *uuid) URN() string {
    if id.slc == nil || len(id.slc) == 0    { return "" }
    if uuidOptions.Cache && id.urn != ""    { return id.urn }
    res := []string{"urn:uuid:", id.Hex()}
    id.urn = strings.Join(res, "")
    return id.urn
}

// Generates a URL-safe base 64 representation UUID
func (id *uuid) URL64() string {
    var res string
    if id.slc == nil || len(id.slc) == 0    { return "" }
    if uuidOptions.Cache && id.url64 != ""  { return id.url64 }
    if uuidOptions.Cache && id.b64 != "" { res = id.b64 } else { res = id.B64() }
    id.url64 = b64ToURL64(res)
    return id.url64
}

// Generates a URL-safe base 32 representation of UUID with dashes
func (id *uuid) URL32() string {
    var res string
    if id.slc == nil || len(id.slc) == 0    { return "" }
    if uuidOptions.Cache && id.url32 != ""  { return id.url32 }
    if uuidOptions.Cache && id.b32 != "" { res = id.b32 } else { res = id.B32() }
    res = strings.Replace(res, "=", "", -1)
    if uuidOptions.HyphURL32 { res = hyphenate(res, 4) }
    id.url32 = res
    return id.url32
}


// -- Decode --

func (_ uuidCtrl) fromA85(s string) ([]byte, error) {
    if len(s) == 24 { s = s[2:22] }
    if len(s) != 20 {
        return []byte{}, errors.New("string of UUID ASCII 85 is wrong length")
    }
    dst, src := make([]byte, 16), make([]byte, 16)
    src = []byte(s)
    _, _, err := ascii85.Decode(dst, src, true)
    if err != nil { return []byte{}, err }
    return dst, nil
}

func (_ uuidCtrl) fromB64(s string) ([]byte, error) {
    s = url64ToB64(s)
    if len(s) == 22 { s = strings.Join([]string{s, "=="}, "") }
    if len(s) != 24 {
        return []byte{}, errors.New("string of UUID base 64 is wrong length")
    }
    dst, err := base64.StdEncoding.DecodeString(s)
    if err != nil { return []byte{}, err }
    return dst, nil
}

func (_ uuidCtrl) fromB32(s string) ([]byte, error) {
    s = strings.Replace(s, " ", "", -1)
    s = strings.Replace(s, "-", "", -1) 
    s = strings.Replace(s, "=", "", -1) 
    s = strings.ToUpper(s)
    if len(s) != 26 {
        return []byte{}, errors.New("string of UUID base 32 is wrong length")
    }
    s = strings.Join([]string{s, "======"}, "")
    dst, err := base32.StdEncoding.DecodeString(s)
    if err != nil { return []byte{}, err }
    return dst, nil
}

// Copyright 2011 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

func (_ uuidCtrl) fromHex(s string) ([]byte, error) {
    if len(s) == 36+9 {
        if strings.ToLower(s[:9]) != "urn:uuid:" {
            return []byte{}, errors.New("string of UUID URN is malformed") 
        }
        s = s[9:]
    } else if len(s) != 36 {
        return []byte{}, errors.New("string of UUID is wrong length") 
    }
    if s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
        return []byte{}, errors.New("canonical UUID string in wrong format") 
    }
    dst := make([]byte, 16)
    for i, x := range []int{
        0, 2, 4, 6,
        9, 11,
        14, 16,
        19, 21,
        24, 26, 28, 30, 32, 34} {
        v, ok := fromHexOctet(s[x:x+2])
        if !ok { return []byte{}, errors.New("bad octet or errant cosmic ray") }
        dst[i] = v
    }
    return dst, nil
}

// Variant returns the variant encoded in uuid.  It returns Invalid if
// uuid is invalid.
func (id uuid) Variant() uuidVariant {
    bytes := id.slc
    if len(bytes) != 16 {
        return uuidInvalid
    }
    switch {
    case (bytes[8] & 0xc0) == 0x80:
        return uuidRFC4122
    case (bytes[8] & 0xe0) == 0xc0:
        return uuidMicrosoft
    case (bytes[8] & 0xe0) == 0xe0:
        return uuidFuture
    default:
        return uuidReserved
    }
    panic("unreachable")
}

// Version returns the verison of uuid.  It returns false if uuid is not
// valid.
func (id uuid) Version() (uuidVersion) {
    bytes := id.slc
    if len(bytes) != 16 {
        return uuidVersion(0)
    }
    ver := uuidVersion(bytes[6] >> 4)
    return ver
}

// A Version represents UUIDs version.
type uuidVersion byte

// A Variant represents UUIDs variant.
type uuidVariant byte

// Constants returned by Variant.
const (
    uuidInvalid   = uuidVariant(iota)   // Invalid UUID
    uuidRFC4122                         // The variant specified in RFC4122
    uuidReserved                        // Reserved, NCS backward compatibility.
    uuidMicrosoft                       // Reserved, Microsoft Corporation backward compatibility.
    uuidFuture                          // Reserved for future definition.
)

func (v uuidVersion) String() string {
    return fmt.Sprintf("VERSION_%d", v)
}

func (v uuidVariant) String() string {
    switch v {
    case uuidRFC4122:
        return "RFC4122"
    case uuidReserved:
        return "Reserved"
    case uuidMicrosoft:
        return "Microsoft"
    case uuidFuture:
        return "Future"
    case uuidInvalid:
        return "Invalid"
    }
    return fmt.Sprintf("BadVariant%d", int(v))
}