package orderedobject

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func FuzzObjectOperations(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{0})
	f.Add([]byte{0, 1})
	f.Add([]byte{0, 0, 1})                            // empty key
	f.Add([]byte{0, 1, 1, 0, 1, 2})                   // repeated Set
	f.Add([]byte{1, 3, 0})                            // delete missing
	f.Add([]byte{0, 1, 1, 0, 2, 2, 1, 1, 0, 0, 1, 3}) // delete and reinsert
	f.Add([]byte{0, 1, 1, 3, 1, 0, 4, 1, 0, 2, 1, 0}) // clone, copies, lookup

	f.Fuzz(func(t *testing.T, data []byte) {
		keys := [...]string{"", "a", "b", "missing"}
		obj := New[int]()
		model := newOperationModel()

		for i := 0; i+2 < len(data); i += 3 {
			op := data[i] % 5
			key := keys[int(data[i+1])%len(keys)]
			value := int(int8(data[i+2]))

			switch op {
			case 0:
				obj.Set(key, value)
				model.set(key, value)
			case 1:
				obj.Delete(key)
				model.delete(key)
			case 2:
				assertLookupMatchesModel(t, obj, model, key)
			case 3:
				assertCloneIndependent(t, obj)
			case 4:
				assertCollectedViewsIndependent(t, obj)
			}

			assertObjectMatchesModel(t, obj, model)
		}

		assertObjectMatchesModel(t, obj, model)
	})
}

type operationModel struct {
	keys   []string
	values map[string]int
}

func newOperationModel() *operationModel {
	return &operationModel{values: make(map[string]int)}
}

func (m *operationModel) set(key string, value int) {
	if _, found := m.values[key]; !found {
		m.keys = append(m.keys, key)
	}
	m.values[key] = value
}

func (m *operationModel) delete(key string) {
	if _, found := m.values[key]; !found {
		return
	}
	delete(m.values, key)
	for i := range m.keys {
		if m.keys[i] == key {
			m.keys = append(m.keys[:i], m.keys[i+1:]...)
			return
		}
	}
}

func (m *operationModel) entries() []Entry[int] {
	entries := make([]Entry[int], len(m.keys))
	for i, key := range m.keys {
		entries[i] = Entry[int]{Key: key, Value: m.values[key]}
	}
	return entries
}

func assertObjectMatchesModel(t *testing.T, obj *Object[int], model *operationModel) {
	t.Helper()

	want := model.entries()
	if diff := cmp.Diff(want, obj.Entries()); diff != "" {
		t.Fatalf("Entries() mismatch (-want +got):\n%s", diff)
	}
	if got := obj.Len(); got != len(want) {
		t.Fatalf("Len() = %d, want %d", got, len(want))
	}

	gotAll := make([]Entry[int], 0, obj.Len())
	for key, value := range obj.All() {
		gotAll = append(gotAll, Entry[int]{Key: key, Value: value})
	}
	if diff := cmp.Diff(want, gotAll); diff != "" {
		t.Fatalf("All() mismatch (-want +got):\n%s", diff)
	}

	for _, key := range [...]string{"", "a", "b", "missing"} {
		assertLookupMatchesModel(t, obj, model, key)
	}
}

func assertLookupMatchesModel(t *testing.T, obj *Object[int], model *operationModel, key string) {
	t.Helper()

	want, wantFound := model.values[key]
	got, gotFound := obj.Get(key)
	if got != want || gotFound != wantFound {
		t.Fatalf("Get(%q) = (%d, %v), want (%d, %v)", key, got, gotFound, want, wantFound)
	}
	if got := obj.Has(key); got != wantFound {
		t.Fatalf("Has(%q) = %v, want %v", key, got, wantFound)
	}
}

func assertCloneIndependent(t *testing.T, obj *Object[int]) {
	t.Helper()

	before := obj.Entries()
	clone := obj.Clone()
	if diff := cmp.Diff(before, clone.Entries()); diff != "" {
		t.Fatalf("Clone() mismatch (-want +got):\n%s", diff)
	}

	clone.Set("clone-only", 1)
	if got, found := clone.Get("clone-only"); !found || got != 1 {
		t.Fatalf("clone.Get(\"clone-only\") = (%d, %v), want (1, true)", got, found)
	}
	if len(before) > 0 {
		clone.Delete(before[0].Key)
	}
	if diff := cmp.Diff(before, obj.Entries()); diff != "" {
		t.Fatalf("Clone mutation changed original (-before +after):\n%s", diff)
	}
}

func assertCollectedViewsIndependent(t *testing.T, obj *Object[int]) {
	t.Helper()

	before := obj.Entries()
	entries := obj.Entries()
	keys := obj.Keys()
	values := obj.Values()
	if len(entries) > 0 {
		entries[0] = Entry[int]{Key: "changed", Value: 99}
		keys[0] = "changed"
		values[0] = 99
	}
	if diff := cmp.Diff(before, obj.Entries()); diff != "" {
		t.Fatalf("Collected view mutation changed object (-before +after):\n%s", diff)
	}
}
