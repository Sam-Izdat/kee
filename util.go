package kee

import(
    crand "crypto/rand"
    mrand "math/rand"
    "bytes"
    "strings"
    "io"
    "math/big"
    "strconv"
    "time"
)

// Substitutions for converting between URL component and standard base 64 
var b64chrs, urlchrs = 
        [3]string{"+", "/", "="}, 
        [3]string{"-", "_", ""}

// Removes URL-unsafe base 64 characters and returns safe URL component
func b64ToURL64(s string) string {
    for key, val := range b64chrs {
        s = strings.Replace(s, val, urlchrs[key], -1)
    }
    return s
}

// Converts URL component back to standard base 64 encoding
func url64ToB64(s string) string {
    for key, val := range urlchrs[0:2] { // skip empty string!
        s = strings.Replace(s, val, b64chrs[key], -1)
    }
    return s
}

// Inserts a dash every n characters
func hyphenate(s string, n int) string {
    os := strings.Split(s, "")
    var ns []string
    var buf bytes.Buffer
    for i, r := range os {
        buf.WriteString(r)
        if (i > 0 && (i+1)%n == 0) || i+1 == len(s) {
            ns = append(ns, buf.String())
            buf.Reset()
        }
    }
    return strings.Join(ns, "-")
}

// fromHexChar converts a hex character into its value and a success flag.
func fromHexChar(c byte) (byte, bool) {
    switch {
    case '0' <= c && c <= '9':
        return c - '0', true
    case 'a' <= c && c <= 'f':
        return c - 'a' + 10, true
    case 'A' <= c && c <= 'F':
        return c - 'A' + 10, true
    }
    return 0, false
}

func fromHexOctet(s string) (byte, bool) {
    a, ok := fromHexChar(s[0])
    if !ok {
        return 0, false
    }
    b, ok := fromHexChar(s[1])
    if !ok {
        return 0, false
    }
    return (a << 4) | b, true
}

// Returns integer in range
func randIntr(min, max int) int {
    r := mrand.New(mrand.NewSource(time.Now().UnixNano()))
    return int( r.Float64() * float64((max - min) + min) )
}

// Returns integer between 0 and n
func randIntn(n int) int {
    if n == 0 { return 0 }
    r := mrand.New(mrand.NewSource(time.Now().UnixNano()))
    return int( r.Intn(n) )
}

// Copyright 2011 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// completely fills slice b with random data.
func randomBits(b []byte) {
    reader := crand.Reader
    if _, err := io.ReadFull(reader, b); err != nil {
        panic(err.Error()) // rand should never fail
    }
}

// Copyright (c) 2012 Tommi Virtanen

// Package base58 implements a human-friendly base58 encoding.
//
// As opposed to base64 and friends, base58 is typically used to
// convert integers. You can use big.Int.SetBytes to convert arbitrary
// bytes to an integer first, and big.Int.Bytes the other way around.

const b58Alphabet = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"

var b58DecodeMap [256]byte

func init() {
    for i := 0; i < len(b58DecodeMap); i++ {
        b58DecodeMap[i] = 0xFF
    }
    for i := 0; i < len(b58Alphabet); i++ {
        b58DecodeMap[b58Alphabet[i]] = byte(i)
    }
}

type b58CorruptInputError int64

func (e b58CorruptInputError) Error() string {
    return "illegal base58 data at input byte " + strconv.FormatInt(int64(e), 10)
}

// Decode a big integer from the bytes. Returns an error on corrupt input.
func b58ToBigInt(src []byte) (*big.Int, error) {
    n := new(big.Int)
    radix := big.NewInt(58)
    for i := 0; i < len(src); i++ {
        b := b58DecodeMap[src[i]]
        if b == 0xFF {
            return nil, b58CorruptInputError(i)
        }
        n.Mul(n, radix)
        n.Add(n, big.NewInt(int64(b)))
    }
    return n, nil
}

// Encode encodes src, appending to dst. Be sure to use the returned
// new value of dst.
func bigIntToB58(dst []byte, src *big.Int) []byte {
    start := len(dst)
    n := new(big.Int)
    n.Set(src)
    radix := big.NewInt(58)
    zero := big.NewInt(0)

    for n.Cmp(zero) > 0 {
        mod := new(big.Int)
        n.DivMod(n, radix, mod)
        dst = append(dst, b58Alphabet[mod.Int64()])
    }

    for i, j := start, len(dst)-1; i < j; i, j = i+1, j-1 {
        dst[i], dst[j] = dst[j], dst[i]
    }
    return dst
}