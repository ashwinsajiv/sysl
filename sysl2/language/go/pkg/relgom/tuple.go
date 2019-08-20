package relgom

import (
	"strings"

	sysl "github.com/anz-bank/sysl/src/proto"
	. "github.com/anz-bank/sysl/sysl2/codegen/golang" //nolint:golint,stylecheck
	"github.com/anz-bank/sysl/sysl2/language/go/pkg/codegen"
	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
)

const (
	tupleDataRecv    = "d"
	relationDataRecv = "d"
	tupleRecv        = "t"
	relationRecv     = "r"
	builderRecv      = "b"
)

func maskForField(i int) Expr {
	return Index(Dot(I(builderRecv), "mask"), Int(i/64))
}

type tupleGenerator struct {
	*sourceGenerator
	*modelScope
	*commonModules
	relation         *sysl.Type_Relation
	namedAttrs       syslutil.NamedTypes
	nAttrs           int
	haveKeys         bool
	t                *sysl.Type
	tname            string
	pkName           string
	dataName         string
	relationDataName string
	builderName      string
	autoinc          syslutil.StrSet
	pkey             syslutil.StrSet
	patterns         syslutil.StrSet
	attrPatterns     map[string]syslutil.StrSet
	fkCount          map[string]int
	pkMask           []uint64
	requiredMask     []uint64
}

func genFileForSyslTypeDecl(s *modelScope, tname string, t *sysl.Type) error {
	relation := t.Type.(*sysl.Type_Relation_).Relation
	nAttrs := len(relation.AttrDefs)
	ename := ExportedName(tname)
	g := &tupleGenerator{
		sourceGenerator:  s.newSourceGenerator(),
		modelScope:       s,
		relation:         relation,
		namedAttrs:       syslutil.NamedTypesInSourceOrder(relation.AttrDefs),
		nAttrs:           nAttrs,
		haveKeys:         relation.PrimaryKey != nil,
		t:                t,
		tname:            ename,
		pkName:           NonExportedName(tname + "PK"),
		dataName:         NonExportedName(tname + "Data"),
		relationDataName: NonExportedName(tname + "RelationData"),
		builderName:      ename + "Builder",
		autoinc:          syslutil.MakeStrSet(),
		pkey:             syslutil.MakeStrSet(),
		patterns:         syslutil.MakeStrSetFromAttr("patterns", t.Attrs),
		attrPatterns:     map[string]syslutil.StrSet{},
		fkCount:          map[string]int{},
		pkMask:           make([]uint64, (nAttrs+63)/64),
		requiredMask:     make([]uint64, (nAttrs+63)/64),
	}
	g.commonModules = newCommonModules(g.sourceGenerator)

	if g.haveKeys {
		for _, pk := range g.relation.PrimaryKey.AttrName {
			g.pkey.Insert(pk)
		}
	}

	for i, nt := range g.namedAttrs {
		mask := uint64(1) << uint(i%64)
		patterns := syslutil.MakeStrSetFromAttr("patterns", nt.Type.Attrs)
		g.attrPatterns[nt.Name] = patterns
		if g.pkey.Contains(nt.Name) {
			g.pkMask[i/64] |= mask
		}
		if !nt.Type.Opt {
			g.requiredMask[i/64] |= mask
		}
		if g.attrPatterns[nt.Name].Contains("autoinc") {
			g.autoinc.Insert(nt.Name)
			g.requiredMask[i/64] &= ^mask
		}
		if fkey := g.typeInfoForSyslType(nt.Type).fkey; fkey != nil {
			path := strings.Join(fkey.Ref.Path, ".")
			if n, has := g.fkCount[path]; has {
				g.fkCount[path] = n + 1
			} else {
				g.fkCount[path] = 1
			}
		}
	}
	if len(g.autoinc) > 1 {
		panic("> 1 autoinc field not supported")
	}

	decls := []Decl{}

	if g.haveKeys {
		decls = append(decls,
			g.goKeyTypeDecl(),
		)
	}

	decls = append(decls,
		g.goTupleDataTypeDecl(),
		g.marshalTupleDataJSONFunc(),
		g.unmarshalTupleDataJSONFunc(),
		g.goTupleTypeDecl(),
	)

	for _, nt := range g.namedAttrs {
		decls = append(decls, g.goGetterFuncForSyslAttr(nt.Name, nt.Type))
	}

	decls = append(decls,
		g.goBuilderTypeDecl(),
	)

	for i, nt := range g.namedAttrs {
		if !g.autoinc.Contains(nt.Name) {
			decls = append(decls, g.goBuilderSetterFuncForSyslAttr(i, nt.Name, nt.Type))
		}
	}

	decls = append(decls,
		g.goTupleTypeStaticMetadataDecl(),
		g.goBuilderApplyDecl(),
		g.goRelationDataDecl(),
		g.marshalRelationDataJSONFunc(),
		g.unmarshalRelationDataJSONFunc(),
		g.goRelationDecl(),
	)

	if g.haveKeys {
		decls = append(decls,
			g.goRelationInsertMethod(),
			g.goRelationUpdateMethod(),
			g.goRelationDeleteMethod(),
			g.goRelationLookupMethod(),
		)
	}

	return g.genSourceForDecls(g.tname, decls...)
}

func (g *tupleGenerator) goKeyTypeDecl() Decl {
	return Types(TypeSpec{
		Name: *I(g.pkName),
		Type: Struct(g.goFieldsForSyslAttrDefs(g.isPkeyAttr, false, nil)...),
	}).WithDoc(Commentf("// %s is the Key for %s.", g.pkName, g.tname))
}

func (g *tupleGenerator) goTupleDataTypeDecl() Decl {
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

func (g *tupleGenerator) tupleDataMethod(f FuncDecl) *FuncDecl {
	f.Recv = Fields(Field{
		Names: Idents(tupleDataRecv),
		Type:  Star(I(g.dataName)),
	}).Parens()
	return &f
}

func (g *tupleGenerator) marshalTupleDataJSONFunc() Decl {
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

func (g *tupleGenerator) unmarshalTupleDataJSONFunc() Decl {
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

func (g *tupleGenerator) goTupleTypeDecl() Decl {
	return Types(TypeSpec{
		Name: *I(g.tname),
		Type: Struct(
			Field{Type: Star(I(g.dataName))},
			Field{Names: Idents("model"), Type: I(g.modelName)},
		),
	}).WithDoc(Commentf("// %s is the public representation tuple in the model.", g.tname))
}

func (g *tupleGenerator) tupleMethod(f FuncDecl) *FuncDecl {
	f.Recv = Fields(Field{
		Names: Idents(tupleRecv),
		Type:  I(g.tname),
	}).Parens()
	return &f
}

func (g *tupleGenerator) goGetterFuncForSyslAttr(attrName string, attr *sysl.Type) Decl {
	exp := ExportedID(attrName)
	nexp := NonExportedID(attrName)
	var field Expr = Dot(I(tupleRecv), nexp.Name.Text)
	typeInfo := g.typeInfoForSyslType(attr)
	if typeInfo.fkey == nil {
		return g.tupleMethod(FuncDecl{
			Doc:  Comments(Commentf("// %s gets the %s attribute from the %s.", exp.Name.Text, attrName, g.tname)),
			Name: *exp,
			Type: FuncType{
				Results: Fields(Field{Type: g.typeInfoForSyslType(attr).final}),
			},
			Body: Block(
				Return(field),
			),
		})
	}
	fpath := typeInfo.fkey.Ref.Path
	ambiguous := g.fkCount[strings.Join(fpath, ".")] > 1
	relation := ExportedID(fpath[0])
	exp2 := ExportedID(fpath[0])
	if ambiguous {
		exp2.Name.Text += "Via" + exp.Name.Text
	}
	if typeInfo.opt {
		field = Star(field)
	}
	return g.tupleMethod(FuncDecl{
		Doc: Comments(Commentf(
			"// %s gets the %s corresponding to the %s attribute from t.",
			exp2.Name.Text, fpath[0], attrName,
		)),
		Name: *exp2,
		Type: FuncType{
			Results: Fields(Field{Type: relation}),
		},
		Body: Block(
			Return(Star(Call(
				Dot(Call(Dot(I(tupleRecv), "model", "Get"+relation.Name.Text)), "Lookup"),
				field,
			))),
		),
	})
}

func (g *tupleGenerator) goBuilderTypeDecl() Decl {
	s := Struct(
		Field{Type: I(g.dataName)},
		Field{Names: Idents("model"), Type: I(g.modelName)},
		Field{Names: Idents("mask"), Type: ArrayN((g.nAttrs+63)/64, I("uint64"))},
		Field{Names: Idents("apply"), Type: g.applyFuncType()},
	)
	return Types(
		TypeSpec{Name: *ExportedID(g.builderName), Type: s},
	).WithDoc(Commentf("// %s builds an instance of %s in the model.", g.builderName, g.tname))
}

func (g *tupleGenerator) builderMethod(f FuncDecl) *FuncDecl {
	f.Recv = Fields(Field{
		Names: Idents(builderRecv),
		Type:  Star(I(g.builderName)),
	}).Parens()
	return &f
}

func (g *tupleGenerator) goBuilderSetterFuncForSyslAttr(i int, attrName string, attr *sysl.Type) Decl {
	exp := "With" + ExportedName(attrName)
	nexp := NonExportedName(attrName)
	typeInfo := g.typeInfoForSyslType(attr)
	updateMask := CallStmt(g.relgomlib("UpdateMaskForFieldButPanicIfAlreadySet"),
		AddrOf(maskForField(i)),
		Binary(Call(I("uint64"), Int(1)), "<<", Int(i%64)),
	)
	if typeInfo.fkey == nil {
		var value Expr = I("value")
		if g.requiredMask[i/64]&(uint64(1)<<uint(i%64)) == 0 {
			value = AddrOf(value)
		}
		return g.builderMethod(FuncDecl{
			Doc:  Comments(Commentf("// %s sets the %s attribute of the %s.", exp, attrName, g.builderName)),
			Name: *I(exp),
			Type: FuncType{
				Params:  *Fields(Field{Names: Idents("value"), Type: typeInfo.param}),
				Results: Fields(Field{Type: Star(I(g.builderName))}),
			},
			Body: Block(
				updateMask,
				Assign(Dot(I(builderRecv), nexp))("=")(value),
				Return(I(builderRecv)),
			),
		})
	}
	fpath := typeInfo.fkey.Ref.Path
	ambiguous := g.fkCount[strings.Join(fpath, ".")] > 1
	// relation := ExportedID(fpath[0])
	exp2 := "With" + ExportedName(fpath[0])
	if ambiguous {
		exp2 += "For" + ExportedName(attrName)
	}
	var field Expr = Dot(I("t"), NonExportedName(fpath[1]))
	if typeInfo.opt {
		field = AddrOf(field)
	}
	return g.builderMethod(FuncDecl{
		Doc:  Comments(Commentf("// %s sets the %s attribute of the %s from t.", exp, attrName, g.builderName)),
		Name: *I(exp2),
		Type: FuncType{
			Params:  *Fields(Field{Names: Idents("t"), Type: I(fpath[0])}),
			Results: Fields(Field{Type: Star(I(g.builderName))}),
		},
		Body: Block(
			updateMask,
			Assign(Dot(I(builderRecv), nexp))("=")(field),
			Return(I(builderRecv)),
		),
	})
}

func (g *tupleGenerator) goTupleTypeStaticMetadataDecl() Decl {
	pkMask := []Expr{}
	for _, u := range g.pkMask {
		pkMask = append(pkMask, Lit(u))
	}
	reqMask := []Expr{}
	for _, u := range g.requiredMask {
		reqMask = append(reqMask, Lit(u))
	}
	return Var(ValueSpec{
		Names: Idents(NonExportedName(g.tname + "StaticMetadata")),
		Values: []Expr{
			AddrOf(Composite(g.relgomlib("TupleTypeStaticMetadata"),
				KV(I("PKMask"), Composite(&ArrayType{Elt: I("uint64")}, pkMask...)),
				KV(I("RequiredMask"), Composite(&ArrayType{Elt: I("uint64")}, reqMask...)),
			)),
		},
	})
}

func (g *tupleGenerator) goBuilderApplyDecl() Decl {
	pkeys := make([]string, 0, len(g.relation.AttrDefs))
	for i, nt := range g.namedAttrs {
		if g.requiredMask[i/64]&(uint64(1)<<uint(i)) != 0 {
			pkeys = append(pkeys, nt.Name)
		} else {
			pkeys = append(pkeys, "")
		}
	}

	return g.builderMethod(FuncDecl{
		Doc:  Comments(Commentf("// Apply applies the built %s.", g.tname)),
		Name: *I("Apply"),
		Type: FuncType{
			Results: Fields(Field{Type: I(g.modelName)}, Field{Type: I(g.tname)}, Field{Type: I("error")}),
		},
		Body: Block(
			CallStmt(g.relgomlib("PanicIfRequiredFieldsNotSet"),
				Slice(Dot(I(builderRecv), "mask")),
				Dot(NonExportedID(g.tname+"StaticMetadata"), "RequiredMask"),
				String(strings.Join(pkeys, ",")),
			),
			Init("set", "err")(Call(Dot(I(builderRecv), "apply"), AddrOf(Dot(I(builderRecv), g.dataName)))),
			If(nil,
				Binary(I("err"), "!=", Nil()),
				Return(Composite(I(g.modelName)), Composite(I(g.tname)), I("err")),
			),
			Init("model", "_")(
				Call(Dot(I(builderRecv), "model", "relations", "Set"),
					NonExportedID(g.tname+"Key"),
					Composite(I(g.relationDataName), I("set")),
				),
			),
			Return(
				Composite(I(g.modelName), I("model")),
				Composite(I(g.tname), AddrOf(Dot(I(builderRecv), g.dataName)), Dot(I(builderRecv), "model")),
				Nil(),
			),
		),
	})
}

// type ${.name}RelationData struct {
//     set   *seq.HashMap
// }
func (g *tupleGenerator) goRelationDataDecl() Decl {
	return Types(TypeSpec{
		Name: *I(g.relationDataName),
		Type: Struct(
			Field{Names: Idents("set"), Type: Star(g.seq("HashMap"))},
		),
	}).WithDoc(Commentf("// %s represents a set of %s.", g.relationDataName, g.tname))
}

func (g *tupleGenerator) relationDataMethod(f FuncDecl) *FuncDecl {
	f.Recv = Fields(Field{
		Names: Idents(relationDataRecv),
		Type:  Star(I(g.relationDataName)),
	}).Parens()
	return &f
}

// func (r *${.name}RelationData) MarshalJSON() ([]byte, error) {
//     a := make([]${.name}Data, 0, r.set.Size())
//     for kv, m, has := r.set.FirstRestKV(); has; kv, m, has = r.set.FirstRestKV() {
//         a = append(a, kv.Val.(${.name}Data))
//     }
//     return json.Marshal(a)
// }
func (g *tupleGenerator) marshalRelationDataJSONFunc() Decl {
	return g.relationDataMethod(*marshalJSONMethodDecl(
		Init("a")(Call(I("make"),
			&ArrayType{Elt: I(g.dataName)},
			Int(0),
			Call(Dot(I(relationDataRecv), "set", "Size")),
		)),
		&ForStmt{
			Init: Init("kv", "m", "has")(Call(Dot(I(relationDataRecv), "set", "FirstRestKV"))),
			Cond: I("has"),
			Post: Assign(I("kv"), I("m"), I("has"))("=")(Call(Dot(I("m"), "FirstRestKV"))),
			Body: *Block(
				Append(I("a"), Assert(Dot(I("kv"), "Val"), I(g.dataName))),
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
func (g *tupleGenerator) unmarshalRelationDataJSONFunc() Decl {
	var i, key Expr
	if g.haveKeys {
		i, key = I("_"), Dot(I("e"), g.pkName)
	} else {
		i, key = I("i"), I("i")
	}
	return g.relationDataMethod(*unmarshalJSONMethodDecl(
		Init("a")(&ArrayType{Elt: Composite(I(g.dataName))}),
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

// ${ename := ExportedID(.name)}
// ${relation := `${ename}Relation`}
// // ${relation} represents a set of ${ename}.
// type ${relation} struct {
//     ${relation}Data
//     model PetShopModel
// }
func (g *tupleGenerator) goRelationDecl() Decl {
	relation := g.tname + "Relation"
	return Types(TypeSpec{
		Name: *ExportedID(relation),
		Type: Struct(
			Field{Type: I(g.relationDataName)},
			Field{Names: Idents("model"), Type: I(g.modelName)},
		),
	}).WithDoc(Commentf("// %s represents a set of %s.", relation, g.tname))
}

func (g *tupleGenerator) relationMethod(f FuncDecl) *FuncDecl {
	f.Recv = Fields(Field{
		Names: Idents(relationRecv),
		Type:  I(g.tname + "Relation"),
	}).Parens()
	return &f
}

func (g *tupleGenerator) goRelationInsertMethod() Decl {
	entity := I("t")
	modelSet := Dot(Call(Dot(I(relationRecv), "model", "Get"+g.tname)), "set")

	innerStmts := []Stmt{}
	if len(g.autoinc) > 0 {
		for _, nt := range g.namedAttrs {
			if g.autoinc.Contains(nt.Name) {
				t := g.typeInfoForSyslType(nt.Type).param
				innerStmts = append(innerStmts,
					Assign(Dot(I("t"), NonExportedName(nt.Name)))("=")(Call(t, I("id"))),
				)
			}
		}
	}
	innerStmts = append(innerStmts,
		Init("set", "_")(Call(Dot(modelSet, "Set"), Dot(entity, g.pkName), entity)),
		Return(I("set"), Nil()),
	)

	outerStmts := []Stmt{}
	var model Expr
	if len(g.autoinc) > 0 {
		outerStmts = append(outerStmts, Init("model", "id")(Call(Dot(I(relationRecv), "model", "newID"))))
		model = I("model")
	} else {
		model = Dot(I(relationRecv), "model")
	}
	outerStmts = append(outerStmts,
		Return(AddrOf(Composite(I(g.tname+"Builder"),
			KV(I("model"), model),
			KV(I("apply"), FuncT(*g.applyFuncType(), innerStmts...)),
		))),
	)

	return g.relationMethod(FuncDecl{
		Doc:  Comments(Commentf("// Insert creates a builder to insert a new %s.", g.tname)),
		Name: *I("Insert"),
		Type: FuncType{Results: Fields(Field{Type: Star(I(g.tname + "Builder"))})},
		Body: Block(outerStmts...),
	})
}

func (g *tupleGenerator) goRelationUpdateMethod() Decl {
	entity := I("t")
	modelSet := Dot(Call(Dot(I(relationRecv), "model", "Get"+g.tname)), "set")

	return g.relationMethod(FuncDecl{
		Doc:  Comments(Commentf("// Update creates a builder to update t in the model.")),
		Name: *I("Update"),
		Type: FuncType{
			Params:  *Fields(Field{Names: Idents("t"), Type: Star(I(g.tname))}),
			Results: Fields(Field{Type: Star(I(g.tname + "Builder"))}),
		},
		Body: Block(
			Init("b")(
				AddrOf(Composite(I(g.tname+"Builder"),
					KV(I(g.dataName), Star(Dot(I("t"), g.dataName))),
					KV(I("model"), Dot(I(relationRecv), "model")),
					KV(I("apply"), FuncT(*g.applyFuncType(),
						Init("set", "_")(Call(Dot(modelSet, "Set"), Dot(entity, g.pkName), entity)),
						Return(I("set"), Nil()),
					)),
				)),
			),
			CallStmt(I("copy"),
				Slice(Dot(I("b"), "mask")),
				Dot(NonExportedID(g.tname+"StaticMetadata"), "PKMask"),
			),
			Return(I("b")),
		),
	})
}

func (g *tupleGenerator) goRelationDeleteMethod() Decl {
	entity := I("t")
	modelSet := Dot(Call(Dot(I(relationRecv), "model", "Get"+g.tname)), "set")

	return g.relationMethod(FuncDecl{
		Doc:  Comments(Commentf("// Delete deletes t from the model.")),
		Name: *I("Delete"),
		Type: FuncType{
			Params: *Fields(Field{Names: Idents("t"), Type: Star(I(g.tname))}),
			Results: Fields(
				Field{Type: I(g.modelName)},
				Field{Type: I("error")},
			),
		},
		Body: Block(
			Init("set", "_")(Call(Dot(modelSet, "Del"), Dot(entity, g.pkName))),
			Return(Composite(I(g.modelName), I("set")), Nil()),
		),
	})
}

func (g *tupleGenerator) goRelationLookupMethod() Decl {
	fields := []Field{}
	kvs := []Expr{}
	for _, nt := range g.namedAttrs {
		if g.pkey.Contains(nt.Name) {
			fname := NonExportedName(nt.Name)
			fields = append(fields, Field{
				Names: Idents(fname),
				Type:  g.typeInfoForSyslType(nt.Type).param,
			})
			kvs = append(kvs, KV(I(fname), I(fname)))
		}
	}

	return g.relationMethod(FuncDecl{
		Doc:  Comments(Commentf("// Lookup searches %s by primary key.", g.tname)),
		Name: *I("Lookup"),
		Type: FuncType{
			Params:  *Fields(fields...),
			Results: Fields(Field{Type: Star(I(g.tname))}),
		},
		Body: Block(
			If(
				Init("t", "has")(Call(Dot(I(relationRecv), "set", "Get"),
					Composite(I(g.pkName), kvs...),
				)),
				I("has"),
				Return(AddrOf(Composite(I(g.tname),
					KV(I(g.dataName), Assert(I("t"), Star(I(g.dataName)))),
					KV(I("model"), Dot(I(relationRecv), "model")),
				))),
			),
			Return(Nil()),
		),
	})
}

func (g *tupleGenerator) goExportedStruct() *StructType {
	return Struct(g.goFieldsForSyslAttrDefs(
		syslutil.NamedTypeAll,
		true,
		func(nt syslutil.NamedType) map[string]string {
			return map[string]string{"json": nt.Name}
		},
	)...)
}

func (g *tupleGenerator) isPkeyAttr(nt syslutil.NamedType) bool {
	_, has := g.pkey[nt.Name]
	return has
}

func (g *tupleGenerator) goFieldsForSyslAttrDefs(
	include syslutil.NamedTypePredicate,
	export bool,
	computeTagMap func(nt syslutil.NamedType) map[string]string,
) []Field {
	fields := []Field{}
	for _, nt := range g.namedAttrs.Where(include) {
		field := g.goFieldForSyslAttrDef(nt.Name, nt.Type, export)
		if computeTagMap != nil {
			field.Tag = codegen.Tag(computeTagMap(nt))
		}
		fields = append(fields, field)
	}
	return fields
}

func (g *tupleGenerator) goFieldForSyslAttrDef(attrName string, attr *sysl.Type, export bool) Field {
	var id *Ident
	if export {
		id = ExportedID(attrName)
	} else {
		id = NonExportedID(attrName)
	}
	return Field{
		Names: []Ident{*id},
		Type:  g.typeInfoForSyslType(attr).final,
	}
}

func (g *tupleGenerator) applyFuncType() *FuncType {
	return &FuncType{
		Params:  *Fields(Field{Names: Idents("t"), Type: Star(I(g.dataName))}),
		Results: Fields(Field{Type: Star(g.seq("HashMap"))}, Field{Type: I("error")}),
	}
}
