// Code generated by "kvschema"; DO NOT EDIT.

package docs

import (
	"github.com/google/note-maps/kv"
)

// Txn provides entities, components, and indexes backed by a key-value store.
type Txn struct{ kv.Partitioned }

func New(t kv.Txn) Txn { return Txn{kv.Partitioned{t, 0}} }

// SetDocument sets the Document associated with e to v.
//
// Corresponding indexes are updated.
func (s Txn) SetDocument(e kv.Entity, v *Document) error {
	key := make(kv.Prefix, 8+2+8)
	s.Partition.EncodeAt(key)
	DocumentPrefix.EncodeAt(key[8:])
	e.EncodeAt(key[10:])
	var old Document
	if err := s.Get(key, old.Decode); err != nil {
		return err
	}
	if err := s.Set(key, v.Encode()); err != nil {
		return err
	}
	lek := len(key)
	kv.Entity(0).EncodeAt(key[10:])
	key = append(key, kv.Component(0).Encode()...)
	var (
		lik = len(key)
		es  kv.EntitySlice
	)

	// Update Title index
	key = key[:lek].AppendComponent(TitlePrefix)
	for _, iv := range old.IndexTitle() {
		key = append(key[:lik], iv.Encode()...)
		if err := s.Get(key, es.Decode); err != nil {
			return err
		}
		if es.Remove(e) {
			if err := s.Set(key, es.Encode()); err != nil {
				return err
			}
		}
	}
	for _, iv := range v.IndexTitle() {
		key = append(key[:lik], iv.Encode()...)
		if err := s.Get(key, es.Decode); err != nil {
			return err
		}
		if es.Insert(e) {
			if err := s.Set(key, es.Encode()); err != nil {
				return err
			}
		}
	}
	return nil
}

// GetDocument returns the Document associated with e.
//
// If no Document has been explicitly set for e, and GetDocument will return
// the result of decoding a Document from an empty slice of bytes.
func (s Txn) GetDocument(e kv.Entity) (Document, error) {
	var v Document
	vs, err := s.GetDocumentSlice([]kv.Entity{e})
	if len(vs) >= 1 {
		v = vs[0]
	}
	return v, err
}

// GetDocumentSlice returns a Document for each entity in es.
//
// If no Document has been explicitly set for an entity, and the result will
// be a Document that has been decoded from an empty slice of bytes.
func (s Txn) GetDocumentSlice(es []kv.Entity) ([]Document, error) {
	result := make([]Document, len(es))
	key := make(kv.Prefix, 8+2+8)
	s.Partition.EncodeAt(key)
	DocumentPrefix.EncodeAt(key[8:])
	for i, e := range es {
		e.EncodeAt(key[10:])
		err := s.Get(key, (&result[i]).Decode)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// AllDocumentEntities returns the first n entities that have a Document, beginning
// with the first entity greater than or equal to *start.
//
// A nil start value will be interpreted as a pointer to zero.
//
// A value of n less than or equal to zero will be interpretted as the largest
// possible value.
func (s Txn) AllDocumentEntities(start *kv.Entity, n int) (es []kv.Entity, err error) {
	return s.AllComponentEntities(DocumentPrefix, start, n)
}

// EntitiesMatchingDocumentTitle returns entities with Document values that return a matching kv.String from their IndexTitle method.
//
// The returned EntitySlice is already sorted.
func (s Txn) EntitiesMatchingDocumentTitle(v kv.String) (kv.EntitySlice, error) {
	key := make(kv.Prefix, 8+2+8+2)
	s.Partition.EncodeAt(key)
	DocumentPrefix.EncodeAt(key[8:])
	kv.Entity(0).EncodeAt(key[10:])
	TitlePrefix.EncodeAt(key[18:])
	key = append(key, v.Encode()...)
	var es kv.EntitySlice
	return es, s.Get(key, es.Decode)
}

// EntitiesByDocumentTitle returns entities with
// Document values ordered by the kv.String values from their
// IndexTitle method.
//
// Reading begins at cursor, and ends when the length of the returned Entity
// slice is less than n. When reading is not complete, cursor is updated such
// that using it in a subequent call to ByTitle would return next n
// entities.
func (s Txn) EntitiesByDocumentTitle(cursor *kv.IndexCursor, n int) (es []kv.Entity, err error) {
	return s.EntitiesByComponentIndex(DocumentPrefix, TitlePrefix, cursor, n)
}
