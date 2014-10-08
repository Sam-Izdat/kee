package kee

import (
    "crypto/hmac"
    "crypto/sha1"
    "encoding/base32"
    "strings"
    "time"
    "regexp"
    "errors"
    "net/url"
)

// KTOTP type represents a secret capable of producing time-based one time passwords. (RFC 6238)
// It is exported only for reference and should be instantiated through its handler's methods.
type KTOTP struct {
    slc []byte
    b32 string
}

// TOTPConfig is the struct for TOTPOptions. It should only be used if  
// another handler with a different set of options is being created.
type TOTPConfig struct {
    LookAhead, LookBehind, B32Blocks int
    HyphB32 bool
}

// TOTPOptions defines the configuration used by the `kee.TOTP` handler.
// Options can also be changed through `kee.TOTP.Options`.
var TOTPOptions = TOTPConfig {
    LookAhead: 1,           // Allow passwords from n future 30-second blocks
    LookBehind: 1,          // Allow passwords from n previous 30-second blocks
    B32Blocks: 8,           // Secret length (change will invalidate stored pws)
    HyphB32: true,          // Hyphenate base 32 encoded secrets
}

// TOTPCtrl is a struct for the TOTP handler. 
// Unless another handler with different options is needed simply use instance `kee.TOTP`.
type TOTPCtrl struct {
    Options         *TOTPConfig
}

// New generates a new secret and returns KTOTP instance
func (c TOTPCtrl) New() KTOTP {
    bytes := make([]byte, 32)
    randomBits(bytes)
    return KTOTP{slc: bytes}
}

// Set loads an existing secret and returns KTOTP instance
func (c TOTPCtrl) Set(bytes []byte) KTOTP {
    bytesSlc := make([]byte, 32)
    copy(bytesSlc[:], bytes[:])
    return KTOTP{slc: bytesSlc}
}

// Decode takes base 32 encoded string of secret and returns KTOTP instance 
func (c TOTPCtrl) Decode(s string) (KTOTP, error) { 
    reg, err := regexp.Compile("[^A-Za-z0-9]+")
    if err != nil { return KTOTP{}, err }
    s = reg.ReplaceAllString(s, "")
    s = strings.ToUpper(s)
    if expLen := totpGetBlocks() * 4; len(s) != expLen {
        // forgiving case, but rejecting anything less
        return KTOTP{}, errors.New("secret length incorrect")
    }
    return KTOTP{b32: s}, nil // Conversion to byte value intentionally left for later
}

// MatchPasswords compares expected and received secrets, return true if they match, false if not
func (c TOTPCtrl) MatchPasswords(exp []uint32, rec uint32) bool { 
    for i := 0; i < len(exp); i++ {
        if exp[i] == rec { return true }
    }
    return false
}

// String is alias for B32()
func (id *KTOTP) String() string {
    return id.B32()
}

// Slc returns secret as slice. This method is only meant to be used immediately after 
// generating or loading a 32-byte secret to store it somewhere permanently.
func (id *KTOTP) Slc() []byte {
    return id.slc
}

// B32 returns base 32 encoded string representation of secret
func (id *KTOTP) B32() string {
    var res string
    if id.b32 != "" { res = id.b32 } else { 
        res = base32.StdEncoding.EncodeToString(id.slc) 
    }
    blocks := totpGetBlocks()
    res = res[0:blocks * 4]
    if TOTPOptions.HyphB32 { res = hyphenate(res, 4) }
    id.b32 = res
    return id.b32
}

// URI returns Uniform Resource Identifier with secret for QR code generation
func (id *KTOTP) URI(acct, issuer string) string {
    acct = url.QueryEscape(acct)
    issuer = url.QueryEscape(issuer)
    return "otpauth://totp/"+acct+"?secret="+id.B32()+"&issuer="+issuer
}

// The MIT License (MIT)
// Copyright (c) 2014 Robbie Vanbrabant

// MakePassword returns a slice of 6-digit time based passwords
func (id *KTOTP) MakePassword() ([]uint32, error) {
    // Value must always come from B32 string and not slice directly
    var sec string
    if id.b32 != "" { sec = id.b32 } else { sec = id.B32() }    // critical
    sec = strings.Replace(sec, "-", "", -1)                     // remove dashes
    key, err := base32.StdEncoding.DecodeString(sec)
    if err != nil {
        return []uint32{}, errors.New("failed to make password - decoding problem")
    }
    epochSeconds := time.Now().Unix()
    pwd := []uint32{0}

    pwd[0] = totpGetPassword(key, totpToBytes(epochSeconds/30))
    for i := int64(1); i <= int64(TOTPOptions.LookBehind); i++ {
        pwd = append(pwd, totpGetPassword(key, totpToBytes(epochSeconds/30 - i) ) )
    }
    for i := int64(1); i <= int64(TOTPOptions.LookAhead); i++ {
        pwd = append(pwd, totpGetPassword(key, totpToBytes(epochSeconds/30 + i) ) )
    }
    
    return pwd, nil
}

// --- Helpers ---

func totpGetBlocks() int {
    blocks := TOTPOptions.B32Blocks
    switch {
    case(blocks > 13):
        blocks = 13
    case(blocks < 4):
        blocks = 4
    }
    return blocks
}

func totpToBytes(value int64) []byte {
    var result []byte
    mask := int64(0xFF)
    shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
    for _, shift := range shifts {
        result = append(result, byte((value>>shift)&mask))
    }
    return result
}

func totpToUint32(bytes []byte) uint32 {
    return (uint32(bytes[0]) << 24) + (uint32(bytes[1]) << 16) +
        (uint32(bytes[2]) << 8) + uint32(bytes[3])
}

func totpGetPassword(key []byte, value []byte) uint32 {
    // sign the value using HMAC-SHA1
    hmacSha1 := hmac.New(sha1.New, key)
    hmacSha1.Write(value)
    hash := hmacSha1.Sum(nil)

    // We're going to use a subset of the generated hash.
    // Using the last nibble (half-byte) to choose the index to start from.
    // This number is always appropriate as it's maximum decimal 15, the hash will
    // have the maximum index 19 (20 bytes of SHA1) and we need 4 bytes.
    offset := hash[len(hash)-1] & 0x0F

    // get a 32-bit (4-byte) chunk from the hash starting at offset
    hashParts := hash[offset : offset+4]

    // ignore the most significant bit as per RFC 4226
    hashParts[0] = hashParts[0] & 0x7F

    number := totpToUint32(hashParts)

    // size to 6 digits
    // one million is the first number with 7 digits so the remainder
    // of the division will always return < 7 digits
    pwd := number % 1000000

    return pwd
}
