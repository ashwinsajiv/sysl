package parse

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFindTyperef(t *testing.T){
 input :=
`Foo:
    !type Bar:
        s string
        !type Y:
            hello int
            !type z:
                hello bool
                goodbye int
    !union U:
        Bar
        Y
`
fs := afero.NewMemMapFs()
	_ = afero.WriteFile(fs, "foo.sysl", []byte (input), os.ModePerm)

	mod, _ := NewParser().Parse("foo.sysl", fs)

	res := findTypeRef([]string{}, []string{"Foo1", "Bar"}, mod)
	res1 := findTypeRef([]string{}, []string{"", "Bar"}, mod)
	res2 := findTypeRef([]string{}, []string{"Foo", "Bar"}, mod)
	res3 := findTypeRef([]string{}, []string{"Foo", "U"}, mod)
	res4 := findTypeRef([]string{}, []string{"Foo", "Bar", "Y"}, mod)
	res5 := findTypeRef([]string{}, []string{"Foo", "Bar", "Y", "z"}, mod)
	assert.Nil(t, res)
	assert.Nil(t, res1)
	assert.NotNil(t, res2)
	assert.NotNil(t, res3)
	assert.NotNil(t, res4)
	assert.NotNil(t, res5)
}
