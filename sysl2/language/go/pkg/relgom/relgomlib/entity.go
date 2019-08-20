package relgomlib

import (
	"fmt"
	"hash/crc32"
	"math"
	"strings"

	"github.com/mediocregopher/seq"
	"github.com/pkg/errors"
)

type EntityTypeStaticMetadata struct {
	PKMask       []uint64
	RequiredMask []uint64
}

// UpdateMaskForFieldButPanicIfAlreadySet checks the mask and panics if the field was already set.
func UpdateMaskForFieldButPanicIfAlreadySet(entityMask *uint64, fieldMask uint64) {
	if *entityMask&fieldMask != 0 {
		panic(errors.New("field already set"))
	}
	*entityMask |= fieldMask
}

// PanicIfRequiredFieldsNotSet checks the mask and panics if any required fields were not set.
// It takes a fieldsCommaList instead of a slice to ensure efficiency of successful scenarios.
func PanicIfRequiredFieldsNotSet(entityMasks []uint64, requiredFieldsMasks []uint64, fieldsCommalist string) {
	var fields []string
	var missingFields []string
	for i, entityMask := range entityMasks {
		if entityMask&requiredFieldsMasks[i] != requiredFieldsMasks[i] {
			gap := ^entityMask & requiredFieldsMasks[i]
			if fields == nil {
				fields = strings.Split(fieldsCommalist, ",")
			}
			for i, field := range fields {
				if gap&(uint64(1)<<uint(i)) != 0 {
					missingFields = append(missingFields, field)
				}
			}
		}
	}
	if len(missingFields) != 0 {
		quantifier := "field"
		if len(missingFields) != 1 {
			quantifier += "s"
		}
		panic(errors.Errorf("required field%s %s not set", quantifier, strings.Join(missingFields, ", ")))
	}
}

func Hash(i uint32, values ...interface{}) uint32 {
	for _, v := range values {
		i = hash(v, i)
	}
	return i
}

func hash(v interface{}, i uint32) uint32 {
	switch vt := v.(type) {
	case seq.Setable:
		return vt.Hash(i)
	case uint:
		return hashUint64(uint64(vt), i)
	case uint8:
		return hashUint32(uint32(vt), i)
	case uint16:
		return hashUint32(uint32(vt), i)
	case uint32:
		return hashUint32(vt, i)
	case uint64:
		return hashUint64(vt, i)
	case int:
		return hashUint64(uint64(vt), i)
	case int8:
		return hashUint32(uint32(vt), i)
	case int16:
		return hashUint32(uint32(vt), i)
	case int32:
		return hashUint32(uint32(vt), i)
	case int64:
		return hashUint64(uint64(vt), i)
	case float32:
		return hashUint32(math.Float32bits(vt), i)
	case float64:
		return hashUint64(math.Float64bits(vt), i)
	case string:
		return hashBytes([]byte(vt), i)
	case []rune:
		return hashBytes([]byte(string(vt)), i)
	case []byte:
		return hashBytes(vt, i)
	default:
		panic(fmt.Sprintf("unhashable: %T", vt))
	}
}

func hashUint32(v uint32, i uint32) uint32 {
	return 1928474805 * (v + i) // Random multiplier
}

func hashUint64(v uint64, i uint32) uint32 {
	return uint32(7216620393653828965 * (v + uint64(i))) // Random multiplier
}

func hashBytes(v []byte, i uint32) uint32 {
	return 1243383301 * (crc32.ChecksumIEEE(v) + i) // Random multiplier
}
