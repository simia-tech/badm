package badm

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

var types = map[string]Type{
	"uint64": func(value []byte) (string, error) {
		return fmt.Sprintf("%d", binary.BigEndian.Uint64(value)), nil
	},
	"string": func(value []byte) (string, error) {
		return string(value), nil
	},
	"hex": func(value []byte) (string, error) {
		return hex.EncodeToString(value), nil
	},
	"base64": func(value []byte) (string, error) {
		return base64.StdEncoding.EncodeToString(value), nil
	},
}

// Type defines a mapping of a byte-value to a string-representation.
type Type func([]byte) (string, error)

// TypeNames returns all possible type names.
func TypeNames() []string {
	names := []string{}
	for name := range types {
		names = append(names, name)
	}
	return names
}

// TypeFor returns the type with the provided name.
func TypeFor(name string) (Type, error) {
	if t, ok := types[name]; ok {
		return t, nil
	}
	return nil, fmt.Errorf("could not find type %s", name)
}

// RegisterType adds a new type to the type-list.
func RegisterType(name string, t Type) error {
	if _, found := types[name]; found {
		return fmt.Errorf("type %s is already registered", name)
	}
	types[name] = t
	return nil
}
