package relgom

import (
	. "github.com/anz-bank/sysl/sysl2/codegen/golang" //nolint:golint,stylecheck
	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
)

const (
	modelRecv = "m"
)

type modelGenerator struct {
	*sourceGenerator
	*modelScope
	*commonModules
	namedTypes syslutil.NamedTypes
}

func newModelGenerator(s *modelScope) *modelGenerator {
	g := &modelGenerator{
		sourceGenerator: s.newSourceGenerator(),
		modelScope:      s,
		namedTypes:      syslutil.NamedTypesInSourceOrder(s.model.Types),
	}
	g.commonModules = newCommonModules(g.sourceGenerator)
	return g
}

func (g *modelGenerator) genFileForSyslModel() error {

	decls := []Decl{
		// type ${.Model+"ModelKeys"} int
		Types(TypeSpec{
			Name: *NonExportedID(g.modelName + "ModelKeys"),
			Type: I("int"),
		}),
	}

	modelKeys := []ValueSpec{}
	for _, nt := range g.namedTypes {
		modelKeys = append(modelKeys, ValueSpec{
			Names: []Ident{*NonExportedID(nt.Name + "Key")},
		})
	}
	modelKeys[0].Type = NonExportedID(g.modelName + "ModelKeys")
	modelKeys[0].Values = []Expr{Iota()}
	decls = append(decls, Const(modelKeys...))

	decls = append(decls, Types(*TypeSpec{
		Name: *I(g.modelName),
		Type: Struct(Field{Names: Idents("relations"), Type: Star(g.seq("HashMap"))}),
	}.WithDoc(Commentf("// %s is the model.", g.modelName))))

	// // New<Model> creates a new <Model>.
	// func New<Model>() *<Model> {
	// 	   return &<Model>{}
	// }
	newName := "New" + g.modelName
	newDeclFunc := &FuncDecl{
		Doc:  Comments(Commentf("// %s creates a new %s.", newName, g.modelName)),
		Name: *I(newName),
		Type: FuncType{
			Results: Fields(Field{Type: Star(I(g.modelName))}),
		},
		Body: Block(
			Return(Unary("&", Composite(I(g.modelName),
				Call(Dot(I("seq"), "NewHashMap"),
					Unary("&", Composite(Dot(I("seq"), "KV"),
						KV(I("Key"), g.relgomlib("ModelMetadataKey")),
						KV(I("Val"), Composite(g.relgomlib("ModelMetadata"))),
					)),
				),
			))),
		),
	}
	decls = append(decls, newDeclFunc)

	// // <Relation> returns the model's <Relation> relation.
	// func (m *<Model>) Get<Relation>() *<Relation>Relation {
	//     return m.relations.GetVal
	// }
	for _, nt := range g.namedTypes {
		ename := ExportedName(nt.Name)
		sname := ename + "Relation"
		rdname := NonExportedName(nt.Name) + "RelationData"

		decls = append(decls,
			g.modelMethod(FuncDecl{
				Doc:  Comments(Commentf("// %s returns the model's %[1]s relation.", ename)),
				Name: *I("Get" + ename),
				Type: FuncType{Results: Fields(Field{Type: Star(I(sname))})},
				Body: Block(
					If(
						Init("relation", "has")(Call(Dot(I(modelRecv), "relations", "Get"),
							NonExportedID(nt.Name+"Key"),
						)),
						I("has"),
						Return(AddrOf(Composite(I(sname),
							Assert(I("relation"), I(rdname)),
							Star(I(modelRecv)),
						))),
					),
					Return(AddrOf(Composite(I(sname),
						Composite(I(rdname)),
						Star(I(modelRecv)),
					))),
				),
			}))
	}

	decls = append(decls,
		g.marshalJSONModelMethod(),
		g.unmarshalJSONModelMethod(),
		g.modelMethod(FuncDecl{
			Doc:  Comments(Commentf("// newID returns a new id for the model")),
			Name: *I("newID"),
			Type: FuncType{Results: Fields(Field{Type: I(g.modelName)}, Field{Type: I("uint64")})},
			Body: Block(
				Init("relations", "id")(Call(g.relgomlib("NewID"), Dot(I(modelRecv), "relations"))),
				Return(Composite(I(g.modelName), I("relations")), I("id")),
			)}),
	)

	return g.genSourceForDecls(g.modelName, decls...)
}

func (g *modelGenerator) marshalJSONModelMethod() *FuncDecl {
	stmts := []Stmt{
		Init("b")(Call(g.relgomlib("NewRelationMapBuilder"), Dot(I(modelRecv), "relations"))),
	}
	for _, nt := range g.namedTypes {
		stmts = append(stmts,
			CallStmt(Dot(I("b"), "Set"), String(nt.Name), NonExportedID(nt.Name+"Key")),
		)
	}
	stmts = append(stmts,
		Return(Call(g.json("Marshal"), Call(Dot(I("b"), "Map")))),
	)
	return g.modelMethod(*marshalJSONMethodDecl(stmts...))
}

func (g *modelGenerator) unmarshalJSONModelMethod() *FuncDecl {
	stmts := []Stmt{
		Init("e")(Call(g.relgomlib("NewRelationMapExtractor"), Dot(I(modelRecv), "relations"))),
	}
	for _, nt := range g.namedTypes {
		stmts = append(stmts,
			CallStmt(Dot(I("e"), "Set"), String(nt.Name), AddrOf(Composite(NonExportedID(nt.Name+"Data")))),
		)
	}
	stmts = append(stmts,
		If(
			Init("err")(Call(g.json("Unmarshal"), I("data"), Call(Dot(I("e"), "Map")))),
			Binary(I("err"), "!=", Nil()),
			Return(I("err")),
		),
	)
	for _, nt := range g.namedTypes {
		stmts = append(stmts,
			CallStmt(Dot(I("e"), "Extract"), String(nt.Name), NonExportedID(nt.Name+"Key")),
		)
	}
	stmts = append(stmts,
		Assign(Dot(I(modelRecv), "relations"))("=")(Call(Dot(I("e"), "Relations"))),
		Return(Nil()),
	)
	return g.modelMethod(*unmarshalJSONMethodDecl(stmts...))
}

func (g *modelGenerator) modelMethod(f FuncDecl) *FuncDecl {
	f.Recv = Fields(Field{
		Names: Idents(modelRecv),
		Type:  Star(I(g.modelName)),
	}).Parens()
	return &f
}
