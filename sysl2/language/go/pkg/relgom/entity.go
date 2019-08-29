package relgom

import (
	"strings"

	sysl "github.com/anz-bank/sysl/src/proto"
	. "github.com/anz-bank/sysl/sysl2/codegen/golang" //nolint:golint,stylecheck
	"github.com/anz-bank/sysl/sysl2/language/go/pkg/codegen"
	"github.com/anz-bank/sysl/sysl2/sysl/syslutil"
)

func maskForField(i int) Expr {
	return Index(Dot(I(builderRecv), "mask"), Int(i/64))
}

type entityGenerator struct {
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
	g := &entityGenerator{
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
		info := g.typeInfoForSyslType(nt.Type)
		if !info.opt {
			g.requiredMask[i/64] |= mask
		}
		if g.attrPatterns[nt.Name].Contains("autoinc") {
			g.autoinc.Insert(nt.Name)
			g.requiredMask[i/64] &= ^mask
		}
		if fkey := info.fkey; fkey != nil {
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
	decls = g.goAppendPKDecls(decls)
	decls = g.goAppendTupleDataDecls(decls)
	decls = g.goAppendTupleDecls(decls)
	decls = g.goAppendBuilderDecls(decls)
	decls = g.goAppendRelationDataDecls(decls)
	decls = g.goAppendRelationDecls(decls)
	decls = g.appendIterDecls(decls)

	return g.genSourceForDecls(g.tname, decls...)
}

func (g *entityGenerator) isPkeyAttr(nt syslutil.NamedType) bool {
	_, has := g.pkey[nt.Name]
	return has
}

func (g *entityGenerator) goFieldsForSyslAttrDefs(
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

func (g *entityGenerator) goFieldForSyslAttrDef(attrName string, attr *sysl.Type, export bool) Field {
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

type dotter = func(id string, ids ...string) Expr

func method(recv string, typ Expr, f func(recv string, dot dotter) FuncDecl) *FuncDecl {
	fd := f(recv, func(id string, ids ...string) Expr {
		return Dot(I(recv), id, ids...)
	})
	fd.Recv = Fields(Field{
		Names: Idents(recv),
		Type:  typ,
	}).Parens()
	return &fd
}
