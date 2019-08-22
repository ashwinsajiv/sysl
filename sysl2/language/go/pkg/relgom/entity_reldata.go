package relgom

import (
	. "github.com/anz-bank/sysl/sysl2/codegen/golang" //nolint:golint,stylecheck
)

const relationDataRecv = "d"

func (g *entityGenerator) goAppendRelationDataDecls(decls []Decl) []Decl {
	return append(decls,
		g.relationDataTypeDecl(),
		g.relationDataSizeMethodDecl(),
		g.relationDataMarshalJSONMethodDecl(),
		g.relationDataUnmarshalJSONMethodDecl(),
	)
}

// type ${.name}RelationData struct {
//     set   *seq.HashMap
// }
func (g *entityGenerator) relationDataTypeDecl() Decl {
	return Types(TypeSpec{
		Name: *I(g.relationDataName),
		Type: Struct(
			Field{Names: Idents("set"), Type: Star(g.seq("HashMap"))},
		),
	}).WithDoc(Commentf("// %s represents a set of %s.", g.relationDataName, g.tname))
}

func (g *entityGenerator) relationDataMethod(f FuncDecl) *FuncDecl {
	f.Recv = Fields(Field{
		Names: Idents(relationDataRecv),
		Type:  I(g.relationDataName),
	}).Parens()
	return &f
}

func (g *entityGenerator) relationDataMethodPointerRecv(f FuncDecl) *FuncDecl {
	f.Recv = Fields(Field{
		Names: Idents(relationDataRecv),
		Type:  Star(I(g.relationDataName)),
	}).Parens()
	return &f
}

func (g *entityGenerator) relationDataSizeMethodDecl() Decl {
	return g.relationDataMethod(FuncDecl{
		Doc:  Comments(Commentf("// Size returns the number of tuples in %s.", relationDataRecv)),
		Name: *I("Size"),
		Type: FuncType{
			Results: Fields(Field{Type: I("uint64")}),
		},
		Body: Block(
			Return(Call(Dot(I(relationDataRecv), "set", "Size"))),
		),
	})
}

// func (r *${.name}RelationData) MarshalJSON() ([]byte, error) {
//     a := make([]${.name}Data, 0, r.set.Size())
//     for kv, m, has := r.set.FirstRestKV(); has; kv, m, has = r.set.FirstRestKV() {
//         a = append(a, kv.Val.(${.name}Data))
//     }
//     return json.Marshal(a)
// }
func (g *entityGenerator) relationDataMarshalJSONMethodDecl() Decl {
	return g.relationDataMethod(*marshalJSONMethodDecl(
		Init("a")(Call(I("make"),
			&ArrayType{Elt: Star(I(g.dataName))},
			Int(0),
			Call(Dot(I(relationDataRecv), "set", "Size")),
		)),
		&ForStmt{
			Init: Init("kv", "m", "has")(Call(Dot(I(relationDataRecv), "set", "FirstRestKV"))),
			Cond: I("has"),
			Post: Assign(I("kv"), I("m"), I("has"))("=")(Call(Dot(I("m"), "FirstRestKV"))),
			Body: *Block(
				Append(I("a"), Assert(Dot(I("kv"), "Val"), Star(I(g.dataName)))),
			),
		},
		Return(Call(g.json("Marshal"), I("a"))),
	))
}

// func (r *${.name}RelationData) UnmarshalJSON(data []byte) error {
//     a := []${.name}Data{}
//     if err := json.Unmarshal(data, &a); err != nil {
//         return err
//     }
//     set := seq.NewHashMap()
//     for _, e := range a {
//         set, _ = set.Set(e.${.name}PK, e)
//     }
//     *d = ${.name}RelationData{set}
//     return nil
// }
func (g *entityGenerator) relationDataUnmarshalJSONMethodDecl() Decl {
	var i, key Expr
	if g.haveKeys {
		i, key = I("_"), Dot(I("e"), g.pkName)
	} else {
		i, key = I("i"), I("i")
	}
	return g.relationDataMethodPointerRecv(*unmarshalJSONMethodDecl(
		Init("a")(&ArrayType{Elt: Composite(Star(I(g.dataName)))}),
		If(
			Init("err")(Call(g.json("Unmarshal"), I("data"), AddrOf(I("a")))),
			Binary(I("err"), "!=", Nil()),
			Return(I("err")),
		),
		Init("set")(Call(g.seq("NewHashMap"))),
		Range(i, I("e"), ":=", I("a"),
			Assign(I("set"), I("_"))("=")(Call(Dot(I("set"), "Set"), key, I("e"))),
		),
		Assign(Star(I(relationDataRecv)))("=")(Composite(I(g.relationDataName), I("set"))),
		Return(Nil()),
	))
}
