package petshopmodel

import (
	"testing"

	"github.com/anz-bank/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBreed(t *testing.T) {
	t.Parallel()

	lifespan := decimal.MustParse64("12.25")
	weight := decimal.MustParse64("33")

	m := NewPetShopModel()

	m, lab, err := m.GetBreed().Insert().
		WithBreedName("Labrador").
		WithNumLegs(4).
		WithLegRank(0).
		WithSpecies("Dog").
		WithAvgLifespan(lifespan).
		WithAvgWeight(weight).
		Apply()
	require.NoError(t, err)
	assert.EqualValues(t, "Labrador", *lab.BreedName())
	assert.EqualValues(t, 4, *lab.NumLegs())
	assert.EqualValues(t, 0, *lab.LegRank())
	assert.EqualValues(t, "Dog", *lab.Species())
	assert.EqualValues(t, lifespan, *lab.AvgLifespan())
	assert.EqualValues(t, weight, *lab.AvgWeight())

	assert.True(t, lab.breedPK.Equal(lab.breedPK))
}

func TestBreedIterator(t *testing.T) {
	t.Parallel()

	m := NewPetShopModel()
	m, _, err := m.GetBreed().Insert().WithBreedName("Labrador").WithSpecies("Dog").Apply()
	require.NoError(t, err)
	m, _, err = m.GetBreed().Insert().WithBreedName("Birman").WithSpecies("Cat").Apply()
	require.NoError(t, err)
	m, _, err = m.GetBreed().Insert().WithBreedName("Goldfish").WithSpecies("Fish").Apply()
	require.NoError(t, err)

	breedSpecies := map[string]string{}
	for i := m.GetBreed().Iterator(); i.MoveNext(); {
		breed := i.Current()
		breedSpecies[*breed.BreedName()] = *breed.Species()
	}
	assert.Equal(t, map[string]string{"Labrador": "Dog", "Birman": "Cat", "Goldfish": "Fish"}, breedSpecies)
}
