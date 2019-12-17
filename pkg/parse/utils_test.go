package parse

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
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

	var typerefTestsNil = []struct{
		ref []string
	}{
		{[]string{"Foo1", "Bar"}},
		{[]string{"", "Bar"}},
		{[]string{"", "Bar", "Y", "z"}},
	}

	for _, tt := range typerefTestsNil {
		require.Nil(t,findTypeRef([]string{}, tt.ref, mod) )
	}

	var typerefTestsNotNil = []struct{
		ref []string
	}{
		{[]string{"Foo", "Bar"}},
		{[]string{"Foo", "U"}},
		{[]string{"Foo", "Bar", "Y"}},
		{[]string{"Foo", "Bar", "Y", "z"}},
	}

	for _, tt := range typerefTestsNotNil {
		require.NotNil(t,findTypeRef([]string{}, tt.ref, mod) )
	}
}
