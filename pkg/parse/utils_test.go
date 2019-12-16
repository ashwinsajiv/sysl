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
		foo <: string`
fs := afero.NewMemMapFs()
	_ = afero.WriteFile(fs, "foo.sysl", []byte (input), os.ModePerm)

	mod, _ := NewParser().Parse("foo.sysl", fs)

	res := findTypeRef([]string{}, []string{"Foo", "Bar"}, mod)
	assert.NotNil(t, res)
}
