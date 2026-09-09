package audit

import (
	"bytes"
	"encoding/json"
	"sort"
)

// Canonicalize converts JSON bytes into RFC 8785 JCS canonical JSON.
func Canonicalize(data []byte) ([]byte, error) {
	var val any
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(&val); err != nil {
		return nil, err
	}
	return canonicalizeValue(val)
}

func canonicalizeValue(val any) ([]byte, error) {
	switch v := val.(type) {
	case map[string]any:
		return canonicalizeMap(v)
	case []any:
		return canonicalizeSlice(v)
	default:
		return json.Marshal(v)
	}
}

func canonicalizeMap(m map[string]any) ([]byte, error) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, k := range keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		kb, _ := json.Marshal(k)
		buf.Write(kb)
		buf.WriteByte(':')
		vb, err := canonicalizeValue(m[k])
		if err != nil {
			return nil, err
		}
		buf.Write(vb)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func canonicalizeSlice(s []any) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i, item := range s {
		if i > 0 {
			buf.WriteByte(',')
		}
		ib, err := canonicalizeValue(item)
		if err != nil {
			return nil, err
		}
		buf.Write(ib)
	}
	buf.WriteByte(']')
	return buf.Bytes(), nil
}
