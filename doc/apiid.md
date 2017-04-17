## Arbitrary Precision Integer Identifiers
APIIDs can be positive integers of any size and are represented as base 10 or base 58 strings and `big.Int` only. They can be:

- Byte slice
- `big.Int` integer (from "math/big" of the standard library)
- Base 58 string

### Generating from integer
```go
    id1 := kee.APIID.FromString("12345678901234567890123456789")
        // ... OR:
    id2 := kee.APIID.FromInt(123456789)
        // ... OR:
    myInt := 12345
    id3 := kee.APIID.FromString(strconv.Itoa(myInt))   // Must be string

    fmt.Println(id1, "&", id3) // => KEm5phz2fXwaGwm6 & 4ER
```
### Setting bytes
```go
    // Must be slice
    id := kee.APIID.Set([]byte{11, 22, 33, 44}) // 185999660
```
### Encoding
```go
    fmt.Println(id)                     // Base 58
    fmt.Println(id.B58())               //   "
    fmt.Println(id.BigInt())            // big.Int
    fmt.Println(id.BigInt().String())   // String
    fmt.Println(id.Slc())               // Slice
```
### Decoding
```go
    id := kee.APIID.Decode("hridG") // 185999660
```
### Options
```
    Cache: true            // Cache APIID strings, ignore new options
```
