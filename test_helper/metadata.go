package test_helper

import (
	"bytes"
	"encoding/gob"
)

type Meta struct {
	DeviceId int
}

func NewMetaFromBytes(d []byte) (Meta, error) {
	var m Meta
	buf := bytes.NewBuffer(d)
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(&m)
	return m, err
}

func (m *Meta) ToSlice() ([]byte, error) {
	buf := &bytes.Buffer{}
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(m)
	return buf.Bytes(), err
}

func Size() int {
	m := Meta{1}
	s, _ := m.ToSlice()
	return len(s)
}
