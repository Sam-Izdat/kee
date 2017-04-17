## Fixed Precision Integer Identifiers
FPIIDs accept only unsigned 64 bit integers but, by default, they are converted to 32 or 16 bits whenever possible *for the purposes of base 64/32 string encoding*. This can be disabled by setting the `ShortStr` option to `false`. Slices, arrays and integers will always return as 8 bytes. FPIIDs can be represented as:

- Byte slice/array
- Unsigned 64-bit integer
- Base 64 / URL-safe base 64 string
- Base 32 / URL-safe base 32 string

### Generating from integer
```go
    id1 := kee.FPIID.New(555555555555555)
        // ...OR:
    var myInt int = 12345
    id2 := kee.FPIID.New(uint64(myInt))
    
    // String method defaults to URL-safe base 64
    fmt.Println(id1, "&", id2) // => 47iKW0b5AQA & OTA
```
### Setting bytes
```go
    id := kee.FPIID.Set([8]byte{255, 255, 255, 255, 255, 255, 255, 255})
```
### Encoding
```go
    fmt.Println(id.Int())       // Unsigned 64 bit int
    fmt.Println(id.B64())       // Base 64
    fmt.Println(id)             // URL-safe base 64
    fmt.Println(id.URL64())     //       "
    fmt.Println(id.B32())       // Base 32
    fmt.Println(id.URL32())     // Formatted, no pad base 32
    fmt.Println(id.Slc())       // Slice
    fmt.Println(id.Arr())       // Array
```
### Decoding
```go
    id1 := kee.FPIID.Decode("4O4I-UW2G-7EAQ-A")  // dashes allowed
    id2 := kee.FPIID.Decode("sct0AA==")          // padding optional
```
### Options
```
    Cache: true            // Cache FPIID strings, ignore new options
    ShortStr: true         // Try conversion to uint32/16 for strings
    PadB64: true           // Add padding to base 64 encoded FPIIDs
    PadB32: true           // Add padding to base 32 encoded FPIIDs
    HyphURL32: true        // Hyphenate base 32 encoded URL FPIIDs
```