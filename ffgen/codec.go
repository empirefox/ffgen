package ffgen

import (
	"reflect"
)

type PermitterGetter interface {
	GetPermitter(reflect.Type) (Permitter, bool)
}

type PermitterValidatorGetter interface {
	GetPermitterValidator(reflect.Type) (PermitterValidator, bool)
}

type Permitter interface {
	IsPermitted(field string) bool
	HasPermitted() bool
}

type PermitterValidator interface {
	Permitter
	ValidateField(field string, value interface{}) error
}

// Marshaler
type Marshaler interface {
	MarshalPermittedJSON(getter PermitterGetter) ([]byte, error)
}

type MarshalerWrapper struct {
	getter    PermitterGetter
	marshaler Marshaler
}

func NewMarshalerWrapper(marshaler Marshaler, getter PermitterGetter) *MarshalerWrapper {
	return &MarshalerWrapper{
		getter:    getter,
		marshaler: marshaler,
	}
}

func (w *MarshalerWrapper) MarshalJSON() ([]byte, error) {
	return w.marshaler.MarshalPermittedJSON(w.getter)
}

// Unmarshaler
type Unmarshaler interface {
	UnmarshalPermittedJSON(data []byte, getter PermitterValidatorGetter, unmarshaled *Unmarshaled) error
}

type UnmarshalerWrapper struct {
	getter      PermitterValidatorGetter
	unmarshaler Unmarshaler
	Unmarshaled *Unmarshaled
}

func NewUnmarshalerWrapper(unmarshaler Unmarshaler, getter PermitterValidatorGetter) *UnmarshalerWrapper {
	return &UnmarshalerWrapper{
		getter:      getter,
		unmarshaler: unmarshaler,
		Unmarshaled: new(Unmarshaled),
	}
}

func (w *UnmarshalerWrapper) UnmarshalJSON(data []byte) error {
	return w.unmarshaler.UnmarshalPermittedJSON(data, w.getter, w.Unmarshaled)
}

// Unmarshaled Only use Fields
// TODO support struct type
type Unmarshaled struct {
	Fields  map[string]interface{}
	Structs map[string]*Unmarshaled
	Slices  map[string][]*Unmarshaled
	StrMaps map[string]map[string]*Unmarshaled
	IntMaps map[string]map[int64]*Unmarshaled
}

func (u *Unmarshaled) Set(k string, v interface{}) {
	if u != nil {
		if u.Fields == nil {
			u.Fields = make(map[string]interface{})
		}
		u.Fields[k] = v
	}
}
