package parse

import (
	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/spf13/afero"
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
		want *sysl.Type
	}{
		{[]string{"Foo1", "Bar"}, nil},
		{[]string{"", "Bar"}, nil},
		{[]string{"", "Bar", "Y", "z"}, nil},
	}

	for _, tt := range typerefTestsNil {
		got := findTypeRef([]string{}, tt.ref, mod)
		if got != tt.want {
			t.Errorf("Typreref: expected differs from actual")
		}
	}

	var typerefTestsNotNil = []struct{
		ref []string
		want bool
	}{
		{[]string{"Foo", "Bar"}, false},
		{[]string{"Foo", "U"}, false},
		{[]string{"Foo", "Bar", "Y"}, false},
		{[]string{"Foo", "Bar", "Y", "z"}, false},
	}

	for _, tt := range typerefTestsNotNil {
		res := findTypeRef([]string{}, tt.ref, mod)
		var got bool
		if res == nil{
			got = true
		} else{
			got = false
		}
		if got != tt.want {
			t.Errorf("Typreref: expected differs from actual")
		}
	}
}
