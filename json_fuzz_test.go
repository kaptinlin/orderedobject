package orderedobject

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func FuzzObjectJSON(f *testing.F) {
	for _, seed := range [][]byte{
		[]byte(`{}`),
		[]byte(`{"first":1,"second":2}`),
		[]byte(`{"a":1,"\u0061":2}`),
		[]byte(`{"value":"wrong"}`),
		[]byte(`{"array":[1,{"nested":true}],"object":{"key":"value"}}`),
		[]byte(`{"unterminated":`),
		[]byte(`{"valid":true} trailing`),
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		obj := New[any]().Set("old", "value").Set("keep", true)
		before := obj.Entries()
		err := obj.UnmarshalJSON(data)
		if err != nil {
			if diff := cmp.Diff(before, obj.Entries()); diff != "" {
				t.Fatalf("UnmarshalJSON() mutated receiver on error (-before +after):\n%s", diff)
			}
		} else {
			assertSuccessfulJSONDecode(t, data, obj)
		}

		typed := New[int]().Set("old", 1).Set("keep", 2)
		typedBefore := typed.Entries()
		if err := typed.UnmarshalJSON(data); err != nil {
			if diff := cmp.Diff(typedBefore, typed.Entries()); diff != "" {
				t.Fatalf("typed UnmarshalJSON() mutated receiver on error (-before +after):\n%s", diff)
			}
		}
	})
}

func assertSuccessfulJSONDecode(t *testing.T, input []byte, obj *Object[any]) {
	t.Helper()

	wantKeys, err := topLevelJSONKeys(input)
	if err != nil {
		t.Fatalf("topLevelJSONKeys() error after successful decode = %v", err)
	}
	if diff := cmp.Diff(wantKeys, obj.Keys()); diff != "" {
		t.Fatalf("decoded key order mismatch (-want +got):\n%s", diff)
	}

	seen := make(map[string]struct{}, obj.Len())
	for _, key := range obj.Keys() {
		if _, found := seen[key]; found {
			t.Fatalf("decoded duplicate key %q", key)
		}
		seen[key] = struct{}{}
	}

	encoded, err := obj.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() after successful decode error = %v", err)
	}
	if !json.Valid(encoded) || len(encoded) == 0 || encoded[0] != '{' {
		t.Fatalf("MarshalJSON() = %q, want one valid JSON object", encoded)
	}

	roundTrip, err := FromJSON[any](encoded)
	if err != nil {
		t.Fatalf("FromJSON(MarshalJSON()) error = %v", err)
	}
	if diff := cmp.Diff(obj.Entries(), roundTrip.Entries()); diff != "" {
		t.Fatalf("JSON round trip mismatch (-want +got):\n%s", diff)
	}
}

func topLevelJSONKeys(data []byte) ([]string, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	if _, err := dec.Token(); err != nil {
		return nil, err
	}

	keys := make([]string, 0)
	for dec.More() {
		key, err := dec.Token()
		if err != nil {
			return nil, err
		}
		keys = append(keys, key.(string))

		var value json.RawMessage
		if err := dec.Decode(&value); err != nil {
			return nil, err
		}
	}
	return keys, nil
}
