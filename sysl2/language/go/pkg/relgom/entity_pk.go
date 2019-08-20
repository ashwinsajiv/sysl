package relgom

import (
	. "github.com/anz-bank/sysl/sysl2/codegen/golang" //nolint:golint,stylecheck
)

const pkRecv = "k"

func (g *entityGenerator) goAppendPKDecls(decls []Decl) []Decl {
	if g.haveKeys {
		decls = append(decls,
			g.goPKTypeDecl(),
			g.goPKSetableHash(),
			g.goPKSetableEqual(),
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

func (g *entityGenerator) pkMethod(f FuncDecl) *FuncDecl {
	f.Recv = Fields(Field{
		Names: Idents(pkRecv),
		Type:  Star(I(g.pkName)),
	}).Parens()
	return &f
}

func (g *entityGenerator) goPKSetableHash() Decl {
	args := []Expr{I("i")}
	for _, field := range g.pkFields() {
		args = append(args, Dot(I(pkRecv), field.Names[0].Name.Text))
	}
	return g.pkMethod(FuncDecl{
		Name: *I("Hash"),
		Type: FuncType{
			Params:  *Fields(Field{Names: Idents("i"), Type: I("uint32")}),
			Results: Fields(Field{Type: I("uint32")}),
		},
		Body: Block(
			Return(Call(g.relgomlib("Hash"), args...)),
		),
	})
}

func (g *entityGenerator) goPKSetableEqual() Decl {
	return g.pkMethod(FuncDecl{
		Name: *I("Equal"),
		Type: FuncType{
			Params:  *Fields(Field{Names: Idents("i"), Type: Composite(I("interface"))}),
			Results: Fields(Field{Type: I("bool")}),
		},
		Body: Block(
			If(
				Init("l", "ok")(Assert(I("i"), Star(I(g.pkName)))),
				I("ok"),
				Return(Binary(I(pkRecv), "==", I("l"))),
			),
			Return(I("false")),
		),
	})
}
