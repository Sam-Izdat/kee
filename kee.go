package kee

import (
    "regexp"
    "bytes"
    "text/template"
)

var (
    UUID uuidCtrl
    FPIID fpiidCtrl
    APIID apiidCtrl
    TOTP totpCtrl
    JUMBLE jumCtrl
)

func init() {
    // Handler for creating Universally Unique Identifiers
    UUID = uuidCtrl{
        &uuidOptions,
        map[string]string{ // Namespaces for Version 3 and 5
            "DNS":     "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
            "URL":     "6ba7b811-9dad-11d1-80b4-00c04fd430c8",
            "OID":     "6ba7b812-9dad-11d1-80b4-00c04fd430c8",
            "X500":    "6ba7b814-9dad-11d1-80b4-00c04fd430c8",
            "NIL":     "00000000-0000-0000-0000-000000000000",
        },
    }

    // Handler for creating Fixed Precision Integer Identifiers
    FPIID = fpiidCtrl{&fpiidOptions}

    // Handler for creating Arbitrary Precision Integer Identifiers
    APIID = apiidCtrl{&apiidOptions}

    // Handler for One-time Time Based Passwords
    TOTP  = totpCtrl{&totpOptions}

    // Handler for word-jumble identifiers
    JUMBLE = jumCtrl{ 
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

// Generic ID type for custom identifiers
type GenericID struct {
    idStr string
    idMap map[string]string
}


// Returns canonical string representation
func (id GenericID) String() string {
    return id.idStr
}

// Returns a map of specified ID values
func (id GenericID) Map() map[string]string {
    return id.idMap
}

// Parses s using supplied regexp and returns custom ID instance
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

// Composes m using supplied template and returns custom ID instance
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

// Returns a custom ID handler with provided pattern and template
func NewHandler(repat string, tmpl string) handler {
    return handler{repat, tmpl}
}