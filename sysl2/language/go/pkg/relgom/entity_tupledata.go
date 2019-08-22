package relgom

import (
	. "github.com/anz-bank/sysl/sysl2/codegen/golang" //nolint:golint,stylecheck
	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
)

const tupleDataRecv = "d"

func (g *entityGenerator) goAppendTupleDataDecls(decls []Decl) []Decl {
	return append(decls,
		g.goTupleDataTypeDecl(),
		g.marshalTupleDataJSONFunc(),
		g.unmarshalTupleDataJSONFunc(),
	)
}

func (g *entityGenerator) goTupleDataTypeDecl() Decl {
	fields := []Field{}
	if g.haveKeys {
		fields = append(fields, Field{Type: I(g.pkName)})
	}
	fields = append(fields,
		g.goFieldsForSyslAttrDefs(syslutil.NamedTypeNot(g.isPkeyAttr), false, nil)...,
	)
	return Types(TypeSpec{
		Name: *I(g.dataName),
		Type: Struct(fields...),
	}).WithDoc(Commentf("// %s is the internal representation of a tuple in the model.", g.dataName))
}

func (g *entityGenerator) tupleDataMethod(f FuncDecl) *FuncDecl {
	f.Recv = Fields(Field{
		Names: Idents(tupleDataRecv),
		Type:  Star(I(g.dataName)),
	}).Parens()
	return &f
}

func (g *entityGenerator) marshalTupleDataJSONFunc() Decl {
	kvs := []Expr{}
	for _, nt := range g.namedAttrs {
		kvs = append(kvs, &KeyValueExpr{
			Key:   ExportedID(nt.Name),
			Value: Dot(I(tupleDataRecv), NonExportedName(nt.Name)),
		})
	}

	return g.tupleDataMethod(*marshalJSONMethodDecl(
		Return(Call(g.json("Marshal"), Composite(g.goExportedStruct(), kvs...))),
	))
}

func (g *entityGenerator) unmarshalTupleDataJSONFunc() Decl {
	tmp := "u"

	pkKVs := []Expr{}
	kvs := []Expr{}
	if g.haveKeys {
		kvs = append(kvs, nil)
	}
	for _, nt := range g.namedAttrs {
		if g.isPkeyAttr(nt) {
			pkKVs = append(pkKVs, &KeyValueExpr{
				Key:   NonExportedID(nt.Name),
				Value: Dot(I(tmp), ExportedName(nt.Name)),
			})
		} else {
			kvs = append(kvs, &KeyValueExpr{
				Key:   NonExportedID(nt.Name),
				Value: Dot(I(tmp), ExportedName(nt.Name)),
			})
		}
	}
	if g.haveKeys {
		kvs[0] = &KeyValueExpr{
			Key:   I(g.pkName),
			Value: Composite(I(g.pkName), pkKVs...),
		}
	}

	return g.tupleDataMethod(*unmarshalJSONMethodDecl(
		&DeclStmt{Decl: Var(ValueSpec{
			Names: Idents(tmp),
			Type:  g.goExportedStruct(),
		})},
		If(
			Init("err")(Call(g.json("Unmarshal"), I("data"), AddrOf(I(tmp)))),
			Binary(I("err"), "!=", Nil()),
			Return(I("err")),
		),
		Assign(Star(I(tupleDataRecv)))("=")(Composite(I(g.dataName), kvs...)),
		Return(Nil()),
	))
}

func (g *entityGenerator) goExportedStruct() *StructType {
	return Struct(g.goFieldsForSyslAttrDefs(
		syslutil.NamedTypeAll,
		true,
		func(nt syslutil.NamedType) map[string]string {
			return map[string]string{"json": nt.Name + ",omitempty"}
		},
	)...)
}
