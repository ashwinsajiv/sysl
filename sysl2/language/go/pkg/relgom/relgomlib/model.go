package relgomlib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/mediocregopher/seq"
)

type modelMetadataKey int

const ModelMetadataKey modelMetadataKey = 0

// ModelMetadata holds extra data common to all relgom-generated models.
type ModelMetadata struct {
	LastID uint64
}

func NewID(relations *seq.HashMap) (*seq.HashMap, uint64) {
	val, _ := relations.Get(ModelMetadataKey)
	metadata, _ := val.(ModelMetadata)
	metadata.LastID++
	id := metadata.LastID
	relations, _ = relations.Set(ModelMetadataKey, metadata)
	return relations, id
}

type RelationMapBuilder struct {
	m         map[string]json.Marshaler
	relations *seq.HashMap
}

func NewRelationMapBuilder(relations *seq.HashMap) RelationMapBuilder {
	return RelationMapBuilder{map[string]json.Marshaler{}, relations}
}

func (b RelationMapBuilder) Map() map[string]json.Marshaler {
	return b.m
}

func (b RelationMapBuilder) Set(name string, key interface{}) {
	if relation, has := b.relations.Get(key); has && relation.(Relation).Size() > 0 {
		b.m[name] = relation.(json.Marshaler)
	}
}

type RelationMapExtractor struct {
	m         map[string]relationKeyData
	relations *seq.HashMap
}

type relationKeyData struct {
	key             interface{}
	relationDataPtr interface{}
}

func NewRelationMapExtractor(relations *seq.HashMap) RelationMapExtractor {
	return RelationMapExtractor{map[string]relationKeyData{}, relations}
}

func (b RelationMapExtractor) Set(name string, key, relationDataPtr interface{}) {
	b.m[name] = relationKeyData{key: key, relationDataPtr: relationDataPtr}
}

func (b RelationMapExtractor) UnmarshalRelationDataJSON(data []byte) (*seq.HashMap, error) {
	buf := bytes.NewBuffer(data)
	dec := json.NewDecoder(buf)
	if tok, err := dec.Token(); err != nil || tok != json.Delim('{') {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("must be an object")
	}
	for tok, err := dec.Token(); err == nil && tok != json.Delim('}'); tok, err = dec.Token() {
		name, _ := tok.(string)
		if kd, has := b.m[name]; has {
			if err := dec.Decode(kd.relationDataPtr); err != nil {
				return nil, err
			}
			// TODO: Can we avoid reflection?
			relationData := reflect.ValueOf(kd.relationDataPtr).Elem().Interface()
			b.relations, _ = b.relations.Set(kd.key, relationData)
		}
	}
	return b.relations, nil
}
