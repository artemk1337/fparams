package analyzer

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"strings"
)

//nolint:gochecknoglobals
var flagSet flag.FlagSet

//nolint:gochecknoglobals
var (
	disableCheckFuncParams  bool
	disableCheckFuncReturns bool
)

//nolint:gochecknoinits
func init() {
	flagSet.BoolVar(&disableCheckFuncParams, "disableCheckFuncParams", false, "disable check function params")
	flagSet.BoolVar(&disableCheckFuncReturns, "disableCheckFuncReturns", false, "disable check function returns")
}

func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "fparams",
		Doc:  "checks if function params and returns are all on one line or each on a new line",
		Run:  run,
	}
}

type fparams struct {
	fset *token.FileSet
}

type FuncParams struct {
	StartPos token.Pos
	EndPos   token.Pos
	Fields   []*ast.Field
}

func run(pass *analysis.Pass) (any, error) {
	sla := &fparams{fset: pass.Fset}

	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			if fn, ok := n.(*ast.FuncDecl); ok {
				sla.checkFuncArgs(pass, fn)
			}

			return true
		})
	}

	return nil, nil
}

func (s *fparams) checkFuncArgs(pass *analysis.Pass, fn *ast.FuncDecl) {
	var (
		params, returns *FuncParams
	)

	if (fn.Type.Params == nil || len(fn.Type.Params.List) == 0) &&
		(fn.Type.Results == nil || len(fn.Type.Results.List) == 0) {
		return // No arguments to check
	}

	if s.checkFuncInOneLine(fn) {
		return // function declaration in one line
	}

	// check exists params, flag and create input params struct
	if fn.Type.Params != nil && !disableCheckFuncParams {
		params = &FuncParams{
			StartPos: fn.Type.Params.Pos() + 1,
			EndPos:   fn.Type.Params.End() - 1,
			Fields:   fn.Type.Params.List,
		}
	}
	// check exists results, flag and create return params struct
	if fn.Type.Results != nil && !disableCheckFuncReturns {
		returns = &FuncParams{
			StartPos: fn.Type.Results.Pos() + 1,
			EndPos:   fn.Type.Results.End() - 1,
			Fields:   fn.Type.Results.List,
		}
	}

	// check and replace params
	paramsValid, returnsValid := s.validateFuncParams(params, returns)

	if paramsValid {
		params = nil
	}

	if returnsValid {
		returns = nil
	}

	s.reportMultiLineParams(pass, fn.Name.String(), params, returns)

	return
}

func (s *fparams) checkFuncInOneLine(fn *ast.FuncDecl) bool {
	bodyStartPos := s.fset.Position(fn.Body.Pos())
	fnStartPos := s.fset.Position(fn.Type.Pos())

	// TODO add extra check: bodyStartPos.Column < 120
	if fnStartPos.Line == bodyStartPos.Line {
		return true // function declaration in one line
	}

	return false
}

func (s *fparams) validateFuncParams(params *FuncParams, returns *FuncParams) (paramsValid, returnsValid bool) {
	paramsValid, returnsValid = true, true

	if params != nil && len(params.Fields) > 0 && !s.validateFuncEachParam(params) {
		paramsValid = false
	}

	if returns != nil && len(returns.Fields) > 0 && !s.validateFuncEachParam(returns) {
		returnsValid = false
	}

	return paramsValid, returnsValid
}

func (s *fparams) validateFuncEachParam(params *FuncParams) bool {
	// prevPos start from "("
	// EndPos end on ")"
	prevPos := s.fset.Position(params.StartPos)

	// iterate on each param and check positions
	for _, arg := range params.Fields {
		for _, name := range arg.Names {
			namePos := s.fset.Position(name.Pos())
			if prevPos.Line == namePos.Line {
				return false
			}

			prevPos = namePos
		}
	}

	// extra check for last arg
	if s.fset.Position(params.Fields[len(params.Fields)-1].Pos()).Line == s.fset.Position(params.EndPos).Line {
		return false
	}

	return true
}

func (s *fparams) reportMultiLineParams(
	pass *analysis.Pass,
	fnName string,
	params *FuncParams,
	returns *FuncParams,
) {
	if params == nil && returns == nil {
		return
	}

	msg := fmt.Sprintf(`the parameters and returns of the function "%s" should start on a new line`, fnName)
	s.reportV2(pass, msg, params, returns)
}

func (s *fparams) reportV2(
	pass *analysis.Pass,
	msg string,
	params *FuncParams,
	returns *FuncParams,
) {
	var pos, end token.Pos

	suggestedFixes := make([]analysis.SuggestedFix, 0, 2)

	if params != nil {
		suggestedFixes = append(suggestedFixes, analysis.SuggestedFix{
			Message: msg,
			TextEdits: []analysis.TextEdit{{
				Pos:     params.StartPos,
				End:     params.EndPos,
				NewText: []byte("\n" + buildArgs(params)),
			}},
		})
		pos = params.StartPos
		end = params.EndPos
	}

	if returns != nil {
		suggestedFixes = append(suggestedFixes, analysis.SuggestedFix{
			Message: msg,
			TextEdits: []analysis.TextEdit{{
				Pos:     returns.StartPos,
				End:     returns.EndPos,
				NewText: []byte("\n" + buildArgs(returns)),
			}},
		})
		if !pos.IsValid() {
			pos = returns.StartPos
		}
		end = returns.EndPos
	}

	pass.Report(analysis.Diagnostic{
		Pos:            pos,
		End:            end,
		Message:        msg,
		SuggestedFixes: suggestedFixes,
	})
}

func buildArgs(params *FuncParams) string {
	var builder strings.Builder

	for _, field := range params.Fields {
		fieldType := typeExprToString(field.Type)
		for _, arg := range field.Names {
			builder.WriteString("\t" + arg.Name + " " + fieldType)

			if len(field.Names) > 1 || len(params.Fields) > 1 {
				builder.WriteString(",\n")
			} else {
				builder.WriteString("\n")
			}
		}
	}

	return builder.String()
}

func typeExprToString(paramType ast.Expr) string {
	return types.ExprString(paramType)
}
