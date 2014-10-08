##UUIDs/GUIDs
UUIDs, as described by RFC 4122, can be represented as:
- Byte slice/array
- Canonical hexadecimal string
- Hex URN (RFC 2141) string
- ASCII 85 string
- Base 64 / URL-safe base 64 string
- Base 32 / URL-safe base 32 string

Conversion functions are methods of the `kee.KUUID` type, returned when a UUID is generated or assigned using any of these functions:

- `func UUID.New() KUUID`
- `func UUID.NewV1() (KUUID, error)`
- `func UUID.NewV3(id KUUID, data []byte) (KUUID, error)`
- `func UUID.NewV4() (KUUID, error)`
- `func UUID.NewV5(id KUUID, data []byte) (KUUID, error)`
- `func UUID.Set(arr [16]byte) (KUUID, error)`
- `func UUID.Decode(s string) (KUUID, error)`

###Generating
There are several ways to create UUIDs and different versions serve somewhat different purposes. If in doubt, stick with Version 4. Take a look at [IETF's RFC](http://www.ietf.org/rfc/rfc4122.txt) for complete details.
```go
    // Version 4 (Random) UUID  -- most common version
    id := kee.UUID.New()            // shorthand with no error value
        // ... OR:
    idv4, err := kee.UUID.NewV4()   // explicit with error value
    
    // Version 1 (Hardware ID + Clock) UUID
    idv1, err := kee.UUID.NewV1()
    
    // Version 3 (MD5) UUID
    data := []byte("The quick brown fox jumped over the lazy dog.")
    domain1 := kee.UUID.New()       // Version 4 for domain
    idv3 := kee.UUID.NewV3(domain1, data) 
        // ... OR:
    domain2 := kee.UUID.NS["DNS"]
    idv3 = kee.UUID.NewV3(domain2, data)
    
    // Version 5 (SHA1) UUID (much like V3)
    idv5 := kee.UUID.NewV5(domain1, data)
```
###Setting bytes
```go
    bytes := [16]byte{96, 109, 186, 97, 248, 73, 72, 5, 155, 122, 167, 157, 88, 212, 217, 94}
    id, err := kee.UUID.Set(bytes)
```
###Encoding
Encoding is straightforward. Different encodings have different advantages and drawbacks.
```go
    // Print hex
    fmt.Println(id)
    fmt.Println(id.Hex())
    
    // Print ASCII 85
    fmt.Println(id.A85())
    
    // Print base 64
    fmt.Println(id.B64())     // Standard encoding
    fmt.Println(id.URL64())   // URL-safe encoding
    
    // Print base 32
    fmt.Println(id.B32())     // Standard encoding
    fmt.Println(id.URL32())   // Formatted, no pad
    
    // Print URN
    fmt.Println(id.URN()) 
    
    // Print slice/array
    fmt.Println(id.Slc())
    fmt.Println(id.Arr()) 
```
###Decoding
Any string generated can be decoded, whatever the encoding. Checking the error value instead of ignoring it with an underscore is advised.
```go
    id1, _ := kee.UUID.Decode("urn:uuid:4769491a-7237-4e06-a60a-cc3098563df1")
    id2, _ := kee.UUID.Decode("R2lJGnI3TgamCswwmFY98Q")
    id3, _ := kee.UUID.Decode("I5UUSGTSG5HANJQKZQYJQVR56E======")
    id4, _ := kee.UUID.Decode(`7qkO5E]6_tV@(O$QrZB?`)
    
    fmt.Println(id1, id2)
    // => 4769491a-7237-4e06-a60a-cc3098563df1 4769491a-7237-4e06-a60a-cc3098563df1
    fmt.Println(id3, id4)
    // => 4769491a-7237-4e06-a60a-cc3098563df1 4769491a-7237-4e06-a60a-cc3098563df1
```
###Options
Options can be set with  `kee.UUID.Options`, e.g.

```go
    kee.UUID.Options.Pad64 = false
```
Here's what's available, with default values listed:
```
    Cache: true            // Cache UUID strings, ignore new options
    AllowInvalid: false    // Allows setting of non-standard UUIDs
    MinVer: 1              // Lowest UUID version allowed as valid
    MaxVer: 5              // Highest UUID version allowed as valid
    PadB64: true           // Add padding to base 64 encoded UUIDs
    PadB32: true           // Add padding to base 32 encoded UUIDs
    WrapA85: false         // Wrap ASCII 85 encoded UUIDs with <~ ~>
    HyphURL32: true        // Hyphenate base 32 encoded URL UUIDs
```