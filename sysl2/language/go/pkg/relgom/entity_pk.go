package relgom

import (
	. "github.com/anz-bank/sysl/sysl2/codegen/golang" //nolint:golint,stylecheck
)

func (g *entityGenerator) goAppendPKDecls(decls []Decl) []Decl {
	if g.haveKeys {
		decls = append(decls,
			g.goPKTypeDecl(),
		)
	}
	return decls
}

func (g *entityGenerator) goPKTypeDecl() Decl {
	return Types(TypeSpec{
		Name: *I(g.pkName),
		Type: Struct(g.pkFields()...),
	}).WithDoc(Commentf("// %s is the Key for %s.", g.pkName, g.tname))
}

func (g *entityGenerator) pkFields() []Field {
	return g.goFieldsForSyslAttrDefs(g.isPkeyAttr, false, nil)
}
