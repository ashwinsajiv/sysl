package parse

import (
	"math"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	sysl "github.com/anz-bank/sysl/pkg/sysl"
)

func makeExpr(ctx *sysl.SourceContext) *sysl.Expr {
	return &sysl.Expr{
		SourceContext: ctx,
	}
}

type sourceCtxHelper struct {
	filename string
}

func (s sourceCtxHelper) Get(ctx *antlr.BaseParserRuleContext) *sysl.SourceContext {
	return s.get(ctx.GetStart(), ctx.GetStop(), false)
}

func (s sourceCtxHelper) GetWithText(ctx *antlr.BaseParserRuleContext) *sysl.SourceContext {
	return s.get(ctx.GetStart(), ctx.GetStop(), true)
}

func (s sourceCtxHelper) get(start, end antlr.Token, withText bool) *sysl.SourceContext {
	var text string
	if withText {
		is := start.GetInputStream()
		endChar := end.GetStop()
		if endChar == 0 {
			endChar = is.Size() - 1
		}
		text = is.GetText(start.GetStart(), endChar)
	}

	return &sysl.SourceContext{
		File: s.filename,
		Start: &sysl.SourceContext_Location{
			Line: int32(start.GetLine()),
			Col:  int32(start.GetColumn()),
		},
		End: &sysl.SourceContext_Location{
			Line: int32(end.GetLine()),
			Col:  int32(end.GetColumn()),
		},
		Text: text,
	}
}

func primitiveFromNativeDataType(native antlr.TerminalNode) (*sysl.Type_Primitive_, *sysl.Type_Constraint) {
	if native == nil {
		return nil, nil
	}
	text := strings.ToUpper(native.GetText())
	primitiveType := sysl.Type_Primitive(sysl.Type_Primitive_value[text])
	if primitiveType != sysl.Type_NO_Primitive {
		return &sysl.Type_Primitive_{Primitive: primitiveType}, nil
	}

	var constraint *sysl.Type_Constraint
	switch text {
	case "INT32":
		primitiveType = sysl.Type_Primitive(sysl.Type_Primitive_value["INT"])
		constraint = &sysl.Type_Constraint{
			Range: &sysl.Type_Constraint_Range{
				Min: &sysl.Value{
					Value: &sysl.Value_I{
						I: math.MinInt32,
					},
				},
				Max: &sysl.Value{
					Value: &sysl.Value_I{
						I: math.MaxInt32,
					},
				},
			},
		}

	case "INT64":
		primitiveType = sysl.Type_Primitive(sysl.Type_Primitive_value["INT"])
		constraint = &sysl.Type_Constraint{
			Range: &sysl.Type_Constraint_Range{
				Min: &sysl.Value{
					Value: &sysl.Value_I{
						I: math.MinInt64,
					},
				},
				Max: &sysl.Value{
					Value: &sysl.Value_I{
						I: math.MaxInt64,
					},
				},
			},
		}
	}
	return &sysl.Type_Primitive_{Primitive: primitiveType}, constraint
}

func findTypeRef(scope []string, path []string, module *sysl.Module) *sysl.Type {
	// implement :p
	return nil
}
