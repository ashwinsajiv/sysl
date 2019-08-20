package relgom

import (
	. "github.com/anz-bank/sysl/sysl2/codegen/golang" //nolint:golint,stylecheck
)

func (g *entityGenerator) goPKTypeDecl() Decl {
	return Types(TypeSpec{
		Name: *I(g.pkName),
		Type: Struct(g.goFieldsForSyslAttrDefs(g.isPkeyAttr, false, nil)...),
	}).WithDoc(Commentf("// %s is the Key for %s.", g.pkName, g.tname))
}
