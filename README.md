
![KEE dot GO](doc/logo.png)

...is a golang single-package library for dealing with resource identifiers and the wacky things they do. 

It handles UUIDs/GUIDs, integer identifiers and custom IDs with regular grammar. It also does TOTPs (time-based one time passwords), because that's vaguely related, fits the API and makes the name kind of a pun. Oh, and it generates incoherent nonsense too, if that's your thing.

## Why?

Because fumbling with identifiers and encodings, despite the extensive standard library, is a tedious, somewhat error-prone process that eventually tends to fill a code base with spaghetti. Also, random IDs like "SoppyClownGallopsAimlessly" are a *lot* more fun than using AUTO_INCREMENT.


## Install
Grab the package with: 

    $ go get github.com/Sam-Izdat/kee

[![Build Status](http://drone.io/github.com/Sam-Izdat/kee/status.png)](https://drone.io/github.com/Sam-Izdat/kee/latest) 
[![License MIT](http://img.shields.io/badge/license-MIT-red.svg?style=flat-square)](http://opensource.org/licenses/MIT)
[![GoDoc](http://img.shields.io/badge/doc-REFERENCE-blue.svg?style=flat-square)](https://godoc.org/github.com/Sam-Izdat/kee)

## Basic usage

```go
package main

import(
    "fmt"
    "github.com/Sam-Izdat/kee"
)

func main() {
    // Get a random UUID
    id := kee.UUID.New()
    
    // Print it out
    fmt.Println(id)         // => 2d0bbb67-f3f1-4632-9e27-ca3cd7265e22
}
```

## Making it useful

But wait, there's more...

### The universal and the global

```go
// Get a random UUID
idA := kee.UUID.New()

// Encode it in base 64, URL-safe base 64, ASCII 85, base 32
fmt.Println(idA.B64())      // => LQu7Z/PxRjKeJ8o81yZeIg==
fmt.Println(idA.URL64())    // => LQu7Z_PxRjKeJ8o81yZeIg
fmt.Println(idA.A85())      // => /IT1'oC5:*SgVZCf-XfJ
fmt.Println(idA.B32())      // => FUF3WZ7T6FDDFHRHZI6NOJS6EI======

// Decode another ID from formatted base 32
idB, _ := kee.UUID.Decode("MBW3-UYPY-JFEA-LG32-U6OV-RVGZ-LY")

// Get its raw [16]byte array
fmt.Println(idB.Arr())   
    // => [96 109 186 97 248 73 72 5 155 122 167 157 88 212 217 94]

// Set ID A to the value of ID B
idA = kee.UUID.Set(idB.Arr())

// Check if it's valid
fmt.Println(idB.IsValid())   // => true

// Determine its version and variant
vrs, vrn := idB.Version(), idB.Variant()
fmt.Println("UUID is", vrs, vrn) // => UUID is VERSION_4 RFC4122

// Compare the two
fmt.Println(kee.UUID.Match(idA, idB)) // => true
```
Available methods for UUID output: `Slc`, `Arr`, `Hex`, `A85`, `B64`, `URL64`, `B32`, `URL32` and `URN`.

The `Decode` method of the UUID handler accepts any valid string output listed above.

See documentation or source for UUIDs other than Version 4.

### The less cosmopolitan identifiers

```go
// Make a fixed precision integer identifier
idfa := kee.FPIID.FromInt(555555555555555)

// Variable parameters must be typed uint64
var myInt int = 12345
idfb := kee.FPIID.FromInt(uint64(myInt))

fmt.Println(idfa, "&", idfb)    // => 47iKW0b5AQA & OTA
fmt.Println(idfa.URL32())       // => 4O4I-UW2G-7EAQ-A

// Decode from base 64
idfc, _ := kee.FPIID.Decode("OTA")
fmt.Println(idfc.Int()) // => 12345

// Make an arbitrary-precision integer identifier
idaa := kee.APIID.FromString("654654654654654654654654")
idab := kee.APIID.FromInt(512)
fmt.Println(idaa, "&", idab) // => 8MJbCS5foMSAeC & 9Q

// Decode from base 58
idac, _ := kee.APIID.Decode("9Q")
fmt.Println(idac.BigInt()) // => 512

```
Available methods for FPIID output are: `Slc`, `Arr`, `Int`, `B64`, `URL64`, `B32`, and `URL32`.

Available methods for APIID output are: `Slc`, `BigInt`, and `B58`.

The `Decode` method of the FPIID/APIID handlers accepts any valid string output listed above.

### SplendidToucanVanishDarkly

```go
// ScrawlyFittersFlounderWhither
wut, _ := kee.JUMBLE.New(2,2,2,2)

fmt.Println(wut)
// => FlabbyAgnatesUnhorsedLeftwards

fmt.Println(wut.SampleSpace())
// => 755632777339440

wut, _ = kee.JUMBLE.New(1,4,2,0)
fmt.Println(wut, wut.SampleSpace())
// => GorgedLethalityComplied 72298351944
```
FatiguedPhalangeTriggingBackward? InboundClaptrapPreludesNudely.

### Multi-factor authentication

```go
// Generate a secret
newSecret := kee.TOTP.New()

// Get its [32]byte slice for later use
data := newSecret.Slc()   // only available after New() -- check for nil

// Set a previous secret
secret := kee.TOTP.Set(data)

// Print it in formatted base 32
fmt.Println(secret)       // => KRNX-EGXV-JZVR-GN6P-AF3D-LNO7-UGI7-YMX6

// ...or a URI to make a QR code
uri := secret.URI("Acct Name", "Issuer")
fmt.Println(uri) 
    // => otpauth://totp/Acct+Name?secret=[blabla]&issuer=Issuer

// Generate one-time passowrd(s)
expected, _ := secret.MakePassword()

// Compare it with the password received
received := uint32(123456)
openSafe := kee.TOTP.MatchPasswords(expected, received)

```
This is generally intended for mobile devices and works with the [Google Authenticator](https://play.google.com/store/apps/details?id=com.google.android.apps.authenticator2&hl=en) application.

## Advanced usage

An ID instance comes from the handler of its respective type:
- `UUID` for *Universally/Globally Unique Identifiers*
- `FPIID` for *Fixed Precision Integer Identifiers*
- `APIID` for *Arbitrary Precision Integer Identifiers*
- `TOTP` for *Time-based One-time Passwords*
- `JUMBLE` for *Gibberish*

These handlers share few common methods where appropriate -- at least in purpose, by convention:

- `New()` creates a *new* (meaning original) locally, globally, universally or multiversally unique identifier from scratch. Invoking this method may generate random data, increment a counter or query a database to return a unique identifier; it may take parameters, such as data for a hash function

- `Set()` takes a byte slice or a byte array representing an existing ID and manually assigns the ID instance its definitive value

- `Decode()` derives the instance's value from some formatted or unformatted encoding (like hex, base 32, ASCII 35, etc) representing an existing ID

You can extend functionality by writing your own handlers and two built-in methods are reserved just for such occasions:

- `Parse()` maps out data from the canonical string representation of some ID

- `Compose()` takes a map of that data and reassembles the ID according to its defined structure

While there's not much information to be gleaned from essentially random bits and incrementing integers, many identifiers can have a bit more to say. Let's write a handler for ISBNs.

```go
package main

import (
    "fmt"
    "strconv"
    "github.com/Sam-Izdat/kee"
)

// Define the capture groups with a regular expression pattern
var isbn13pat string = `(?P<label>ISBN|ISBN-13)+[ ]` + 
    `(?P<prefix>`       + `[0-9]*)+[- ]` + 
    `(?P<group>`        + `[0-9]*)+[- ]` +
    `(?P<publisher>`    + `[0-9]*)+[- ]` +
    `(?P<title>`        + `[0-9]*)+[- ]` + 
    `(?P<checksum>`     + `[0-9]*)`

// Define a standard-library template for its canonical format
var isbn13tmpl string = `ISBN-13 {{.prefix}}-{{.group}}-` + 
    `{{.publisher}}-{{.title}}-{{.checksum}}`

// Grab an ID handler
var ISBN = kee.NewHandler(isbn13pat, isbn13tmpl)

// The ID instances should embed the generic ID struct
type isbn struct {
    kee.GenericID
}

// Write any methods you need for your ID instances
// e.g. compute checksum and compare it to check digit
func (p isbn) Check() bool {
    var m map[string]string = p.Map()
    var digits string = (
        m["prefix"] + m["group"] + 
        m["publisher"] + m["title"])
    var sum, weight, tmp int
    var checkdigit, _ = strconv.Atoi(m["checksum"])
    if checkdigit == 0 { checkdigit = 10}
    for k, v := range digits {
        if (k + 1) % 2 == 0 {
            weight = 3
        } else { weight = 1 }
        tmp, _ = strconv.Atoi(string(v))
        sum += tmp * weight
    }
    return 10 - (sum % 10) == checkdigit
}

func main() {
    // Get generic ID instance
    ida, _ := ISBN.Parse("ISBN 978-0-306-40615-7") 
    idb := isbn{ida} // Embed generic ID in isbn instance

    // Test it out
    fmt.Println(ISBN.Compose(idb.Map())) 
        // => ISBN-13 978-0-306-40615-7 <nil>

    // Verify check digit
    fmt.Println(idb.Check()) // => true
}
```
What's provided is really just scaffolding for anyone wishing to follow the conventions above for convenience, consistency and improved code readability. More functionality may be added on later.

# Potential gotchas
- Encoded strings are cached to avoid re-encoding the same string every time it's requested. If you need to change the options *after* creating an ID (e.g. to remove padding or change formatting), turn off the appopriate `Cache` option and consider managing your own variables if performance is a factor.
```go
    id := kee.UUID.New()
    fmt.Println( id.URL32() ) // => Y46J-ZTZ4-3NFO-JBM5-4MLR-BLJR-3M
    kee.UUID.Options.HyphURL32 = false
    fmt.Println( id.URL32() ) // => Y46J-ZTZ4-3NFO-JBM5-4MLR-BLJR-3M
    kee.UUID.Options.Cache = false
    fmt.Println( id.URL32() ) // => Y46JZTZ43NFOJBM54MLRBLJR3M
```
- Version 2 (DCE Security) UUIDs have been axed on account of being unnatural, easy to misuse and generally ridiculous.
```go
    if _, err := kee.UUID.NewV2(); err != nil {
        fmt.Println(err) // => no
    }
```
- A lot of UUID/GUID implementations ignore the RFC spec, just fill 16 bytes with random porridge and call it a day. This porridge will generally be rejected unless the right nibbles just happen to identify it as something it probably isn't. If you absolutely need to accept this porridge, set the `AcceptInvalid` option to `true`.
```go
    if _, err := kee.UUID.Decode("12341324-1234-0000-0000-123412341234"); err != nil {
        fmt.Println(err) // => Invalid UUID
    }
```
- Arbitrary-precision integers are really arbitrary in size just as their encoded versions are in length. If you plan on counting past eighteen quintillion or dividing by zero please mind the constraints and fasten appropriate protective head gear.
```go
    id := kee.APIID.FromString("18446744073709551615") 
    fmt.Println(id.BigInt().Uint64())  // => 18446744073709551615
    id = kee.APIID.FromString("18446744073709551616") // uint64 will overflow
    fmt.Println(id.BigInt().Uint64())  // => 0
    id = kee.APIID.FromString("28446744073709551616") 
    fmt.Println(id.BigInt().Uint64())  // => 10000000000000000000
    fmt.Println(id.BigInt())           // => 28446744073709551616
```
- Check that the server time is exactly correct when using time-based passwords.
- Some JUMBLE words are quite long. If storing phrases in a database, allow for at least 100 characters. Unless words are omitted, the sample space is usually very large but still well below a UUID; checking for collisions is a good idea.
- Please, dear god, don't try to parse email addresses with regular expressions. You have been warned.

# What still needs doin'
- Complete unit tests
- Some base set of built-in regular expressions for parsing common identifiers

# Attributions
Code or content anonymously pinched from:

- Tommi Virtanen: [base58](https://github.com/tv42/base58)
- Robbie Vanbrabant: [two-factor-auth](https://github.com/robbiev/two-factor-auth)
- The [go-uuid package](https://code.google.com/p/go-uuid/)
- Ashley Bovan: [word lists](http://www.ashley-bovan.co.uk/words/partsofspeech.html)
- [WPZOOM](http://www.wpzoom.com/): key icon used above (Attribution-Share Alike 3.0)

# License

MIT with some BSD-licensed snippets
