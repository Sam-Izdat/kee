// Package kee simplifies generating, parsing, composing, encoding and decoding resource identifiers
package kee

import (
    "regexp"
    "bytes"
    "text/template"
)

var (
    // UUID handler for creating Universally Unique Identifiers
    UUID UUIDCtrl   

    // FPIID handler for creating Fixed Precision Integer Identifiers
    FPIID FPIIDCtrl 

    // APIID handler for creating Arbitrary Precision Integer Identifiers
    APIID APIIDCtrl 

    // TOTP handler for One-time Time Based Passwords
    TOTP TOTPCtrl   

    // JUMBLE handler for word-jumble identifiers
    JUMBLE JUMCtrl  
)

func init() {
    UUID = UUIDCtrl{
        &UUIDOptions,
        map[string]string{ // Namespaces for Version 3 and 5
            "DNS":     "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
            "URL":     "6ba7b811-9dad-11d1-80b4-00c04fd430c8",
            "OID":     "6ba7b812-9dad-11d1-80b4-00c04fd430c8",
            "X500":    "6ba7b814-9dad-11d1-80b4-00c04fd430c8",
            "NIL":     "00000000-0000-0000-0000-000000000000",
        },
    }    
    FPIID = FPIIDCtrl{&FPIIDOptions}
    APIID = APIIDCtrl{&APIIDOptions}
    TOTP  = TOTPCtrl{&TOTPOptions}
    JUMBLE = JUMCtrl{ 
        phrase: []jumWord{ 
            &jumAdjectives{}, 
            &jumNouns{}, 
            &jumVerbs{}, 
            &jumAdverbs{},
        },
    }
}

type handler struct {
    repat string
    tmpl string
}

// GenericID type is for custom identifiers
type GenericID struct {
    idStr string
    idMap map[string]string
}

// String returns canonical string representation of the ID
func (id GenericID) String() string {
    return id.idStr
}

// Map returns a map of ID values specfied by handler's regex
func (id GenericID) Map() map[string]string {
    return id.idMap
}

// Parses s using supplied regexp and returns GenericID instance
func (p handler) Parse(s string) (GenericID, error) {
    res := make(map[string]string)
    re, err := regexp.Compile(p.repat)
    if err != nil { return GenericID{}, err }
    names := re.SubexpNames()
    result := re.FindStringSubmatch(s)
    for k, v := range result {
        if k == 0 { continue }
        res[string(names[k])] = string(v)
    }

    inst := GenericID{
        idStr: s,
        idMap: res,
    }

    return inst, nil
}

// Composes m using supplied template and returns GenericID instance
func (p handler) Compose(m map[string]string) (GenericID, error) {
    var res string
    var buf bytes.Buffer

    t := template.New("t")
    t, err := t.Parse(p.tmpl)
    if err != nil { return GenericID{}, err }
    err = t.Execute(&buf, m)
    if err != nil { return GenericID{}, err }
    res = buf.String()

    inst := GenericID{
        idStr: res,
        idMap: m,
    }

    return inst, nil
}

// NewHandler returns a custom ID handler with provided pattern and template
func NewHandler(repat string, tmpl string) handler {
    return handler{repat, tmpl}
}