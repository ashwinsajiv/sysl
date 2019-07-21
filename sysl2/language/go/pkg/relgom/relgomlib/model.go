package relgomlib

import (
	"encoding/json"

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
	if relation, has := b.relations.Get(key); has && relation.(*seq.HashMap).Size() > 0 {
		b.m[name] = relation.(json.Marshaler)
	}
}

type RelationMapExtractor struct {
	m         map[string]json.Unmarshaler
	relations *seq.HashMap
}

func NewRelationMapExtractor(relations *seq.HashMap) RelationMapExtractor {
	return RelationMapExtractor{map[string]json.Unmarshaler{}, relations}
}

func (b RelationMapExtractor) Map() *map[string]json.Unmarshaler {
	return &b.m
}

func (b RelationMapExtractor) Relations() *seq.HashMap {
	return b.relations
}

func (b RelationMapExtractor) Set(name string, value json.Unmarshaler) {
	b.m[name] = value
}

func (b RelationMapExtractor) Extract(name string, key interface{}) {
	b.relations.Set(key, b.m[name])
}
