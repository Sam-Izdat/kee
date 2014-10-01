package kee

import(
    "encoding/binary"
    "encoding/base64"
    "encoding/base32"
    "strings"
    "errors"
)

// KFPIID type represents a fixed precision integer identifier.
// It is exported only for reference and should be instantiated through its handler's methods.
type KFPIID struct {
    slc []byte
    b64 string
    b32 string
    url64 string
    url32 string
}

const ( // Maximum values for signed ints
    maxVal16 uint64 = 65535
    maxVal32 uint64 = 4294967295
    maxVal64 uint64 = 18446744073709551615
)

// FPIIDConfig is the struct for FPIIDOptions. It should only be used if  
// another handler with a different set of options is being created.
type FPIIDConfig struct {
    Cache, ShortStr bool
    PadB64, PadB32, HyphURL32 bool
}

// FPIIDOptions defines the configuration used by the `kee.FPIID` handler.
// Options can also be changed through `kee.FPIID.Options`.
var FPIIDOptions = FPIIDConfig {
    Cache: true,            // Cache FPIID strings, ignore new options
    ShortStr: true,         // Try conversion to uint32/16 for strings
    PadB64: true,           // Add padding to base 64 encoded FPIIDs
    PadB32: true,           // Add padding to base 32 encoded FPIIDs
    HyphURL32: true,        // Hyphenate base 32 encoded URL FPIIDs
}

// FPIIDCtrl is a struct for the APIID handler. 
// Unless another handler with different options is needed simply use instance `kee.FPIID`.
type FPIIDCtrl struct {
    Options *FPIIDConfig
}

// FromInt takes a [8]byte array and returns a KFPIID instance
func (c FPIIDCtrl) FromInt(id uint64) KFPIID {
    bytes := make([]byte, 8)
    binary.LittleEndian.PutUint64(bytes, uint64(id))
    return KFPIID{slc: bytes}
}

// Set takes an [8]byte array and returns a KFPIID instance
func (c FPIIDCtrl) Set(arr [8]byte) KFPIID {
    bytes := make([]byte, 8)
    bytes = arr[:]
    return KFPIID{slc: bytes}
}

// Decode takes encoded string of FPIID and returns KFPIID instance
func (c FPIIDCtrl) Decode(s string) (KFPIID, error) {
    var bytes []byte
    var err error
    tmp := strings.Replace(s, "=", "", -1)
    switch len(tmp) {
    case 3:     // B64 uint16 // len 4 with pad
        bytes, err = c.fromB64(s, 16, 4)
    case 4:     // B32 uint16 // len 8 with pad
        bytes, err = c.fromB32(s, 16, 8)
    case 6:     // B64 uint32 // len 8 with pad
        bytes, err = c.fromB64(s, 32, 8)
    case 7:     // B32 uint32 // len 8 with pad
        bytes, err = c.fromB32(s, 32, 8)
    case 11:    // B64 uint64 // len 12 with pad
        bytes, err = c.fromB64(s, 64, 12)
    case 13:    // B32 uint64 // len 16 with pad
        bytes, err = c.fromB32(s, 64, 16)
    default:
        return KFPIID{slc: []byte{}}, errors.New("unrecognized FPIID encoding")
    }
    return KFPIID{slc: bytes}, err
}

// -- Produce --

// String is alias for URL64()
func (id KFPIID) String() string {
    return id.URL64()
}

// Slc returns FPIID as slice
func (id KFPIID) Slc() []byte {
    return id.slc
}

// Arr returns FPIID as array
func (id KFPIID) Arr() (res [8]byte) {
    copy(res[:], id.slc[:])
    return 
}

// Int returns FPIID as unsigned 64 bit integer
func (id KFPIID) Int() (res uint64) {
    if id.slc == nil || len(id.slc) == 0 { return 0 }
    switch (len(id.slc)) {
    case 2:
        res = uint64(binary.LittleEndian.Uint16(id.slc))
    case 4:
        res = uint64(binary.LittleEndian.Uint32(id.slc))
    case 8:
        res = binary.LittleEndian.Uint64(id.slc)
    default:
        res = 0
    }    
    return
}

// B64 returns base 64 encoded string representation of FPIID
func (id *KFPIID) B64() string {
    if id.slc == nil || len(id.slc) == 0 { return "" }
    if FPIIDOptions.Cache && id.b64 != "" { return id.b64 }
    var bytes []byte;
    copy(bytes, id.slc) 
    if FPIIDOptions.ShortStr { bytes = fpiidTrimBytes(id) }
    res := base64.StdEncoding.EncodeToString(bytes)
    if !FPIIDOptions.PadB64 { res = strings.Replace(res, "=", "", -1) }
    id.b64 = res
    return id.b64
}

// B32 returns base 32 encoded string representation of FPIID
func (id *KFPIID) B32() string {
    if id.slc == nil || len(id.slc) == 0    { return "" }
    if FPIIDOptions.Cache && id.b32 != ""   { return id.b32 }
    var bytes []byte;
    copy(bytes, id.slc)
    if FPIIDOptions.ShortStr { bytes = fpiidTrimBytes(id) }
    res := base32.StdEncoding.EncodeToString(bytes)
    if !FPIIDOptions.PadB32 { res = strings.Replace(res, "=", "", -1) }
    id.b32 = res
    return id.b32
}

// URL64 returns URL-safe base 64 string representation FPIID
func (id *KFPIID) URL64() string {
    var res string
    if id.slc == nil || len(id.slc) == 0    { return "" }
    if FPIIDOptions.Cache && id.url64 != "" { return id.url64 }
    if FPIIDOptions.Cache && id.b64 != "" { res = id.b64 } else { res = id.B64() }
    id.url64 = b64ToURL64(res)
    return id.url64
}

// URL32 returns formatted, URL-safe base 32 string representation of FPIID
func (id *KFPIID) URL32() string {
    var res string
    if id.slc == nil || len(id.slc) == 0    { return "" }
    if FPIIDOptions.Cache && id.url32 != "" { return id.url32 }
    if FPIIDOptions.Cache && id.b32 != "" { res = id.b32 } else { res = id.B32() }
    res = strings.Replace(res, "=", "", -1)
    if FPIIDOptions.HyphURL32 { res = hyphenate(res, 4) }
    id.url32 = res
    return id.url32
}

// -- Decode --

func (_ FPIIDCtrl) fromB64(s string, bLen, sLen int) ([]byte, error) {
    bytes := make([]byte, bLen/8)
    pieces := []string{url64ToB64(s)}
    for ; len(s) < sLen; sLen-- {
        pieces = append(pieces, "=")
    }
    s = strings.Join(pieces, "")
    bytes, err := base64.StdEncoding.DecodeString(s)
    return bytes, err
}

func (_ FPIIDCtrl) fromB32(s string, bLen, sLen int) ([]byte, error) {
    bytes := make([]byte, bLen/8)
    s = strings.Replace(s, " ", "", -1)
    s = strings.Replace(s, "-", "", -1) 
    s = strings.Replace(s, "=", "", -1) 
    s = strings.ToUpper(s)
    pieces := []string{s}
    for ; len(s) < sLen; sLen-- {
        pieces = append(pieces, "=")
    }
    s = strings.Join(pieces, "")
    bytes, err := base32.StdEncoding.DecodeString(s)
    return bytes, err
}

// -- Helpers --

func fpiidTrimBytes(id *KFPIID) []byte {
    val := id.Int()
    switch {
    case (val <= maxVal16):
        tmp := make([]byte, 2)
        binary.LittleEndian.PutUint16(tmp, uint16(val))
        return tmp
    case (val <= maxVal32):
        tmp := make([]byte, 4)
        binary.LittleEndian.PutUint32(tmp, uint32(val))
        return tmp
    }
    return id.slc
}