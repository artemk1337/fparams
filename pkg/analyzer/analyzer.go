package analyzer

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const (
	errMsgTotal   = `the parameters and return values of the function "%s" should be on separate lines`
	errMsgParams  = `the parameters of the function "%s" should be on separate lines`
	errMsgReturns = `the return values of the function "%s" should be on separate lines`
)

//nolint:gochecknoglobals
var flagSet flag.FlagSet

//nolint:gochecknoglobals
var (
	disableCheckFuncParams  bool
	disableCheckFuncReturns bool
	//maxLineLength           int
)

//nolint:gochecknoinits
func init() {
	flagSet.BoolVar(&disableCheckFuncParams, "disableCheckFuncParams", false, "disable check function params")
	flagSet.BoolVar(&disableCheckFuncReturns, "disableCheckFuncReturns", false, "disable check function returns")
	//flagSet.IntVar(&maxLineLength, "maxLineLength", 0, "maximum line length with function declaration")
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

// Params - extra model to store params and start/end position.
type Params struct {
	StartPos token.Pos
	EndPos   token.Pos
	Fields   []*ast.Field
}

func run(pass *analysis.Pass) (interface{}, error) {
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
	var params, returns *Params

	if (fn.Type.Params == nil || len(fn.Type.Params.List) == 0) &&
		(fn.Type.Results == nil || len(fn.Type.Results.List) == 0) {
		return // No arguments to check
	}

	if s.checkFuncInOneLine(fn) {
		return // function declaration in one line
	}

	// check exists params, flag and create input params struct
	if fn.Type.Params != nil && !disableCheckFuncParams {
		params = &Params{
			StartPos: fn.Type.Params.Pos() + 1,
			EndPos:   fn.Type.Params.End() - 1,
			Fields:   fn.Type.Params.List,
		}
	}
	// check exists results, flag and create return params struct
	if fn.Type.Results != nil && !disableCheckFuncReturns {
		returns = &Params{
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
}

func (s *fparams) checkFuncInOneLine(fn *ast.FuncDecl) bool {
	bodyStartPos := s.fset.Position(fn.Body.Pos())
	fnStartPos := s.fset.Position(fn.Type.Pos())

	// TODO add extra check: bodyStartPos.Column < 120
	return fnStartPos.Line == bodyStartPos.Line
}

func (s *fparams) validateFuncParams(params *Params, returns *Params) (paramsValid, returnsValid bool) {
	paramsValid, returnsValid = true, true

	if params != nil && len(params.Fields) > 0 && !s.validateFuncEachParam(params) {
		paramsValid = false
	}

	if returns != nil && len(returns.Fields) > 0 && !s.validateFuncEachParam(returns) {
		returnsValid = false
	}

	return paramsValid, returnsValid
}

func (s *fparams) validateFuncEachParam(params *Params) bool {
	// prevPos starts from "("
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
	// EndPos ends on ")"
	return s.fset.Position(params.Fields[len(params.Fields)-1].Pos()).Line != s.fset.Position(params.EndPos).Line
}

func (s *fparams) reportMultiLineParams(
	pass *analysis.Pass,
	fnName string,
	params *Params,
	returns *Params,
) {
	var errMsg string

	switch {
	case params != nil && returns != nil: // params and return values exists
		errMsg = fmt.Sprintf(errMsgTotal, fnName)
	case params != nil: // only params
		errMsg = fmt.Sprintf(errMsgParams, fnName)
	case returns != nil: // only return values
		errMsg = fmt.Sprintf(errMsgReturns, fnName)
	default:
		return
	}

	s.reportAndSuggest(pass, errMsg, params, returns)
}

func (s *fparams) reportAndSuggest(
	pass *analysis.Pass,
	msg string,
	params *Params,
	returns *Params,
) {
	var pos, end token.Pos

	// max size - 2; params and return values suggestion
	suggestedFixes := make([]analysis.SuggestedFix, 0, 2) //nolint:mnd

	// create params suggestion
	if params != nil {
		suggestedFixes = append(suggestedFixes, analysis.SuggestedFix{
			Message: msg,
			TextEdits: []analysis.TextEdit{{
				Pos:     params.StartPos,
				End:     params.EndPos,
				NewText: []byte("\n" + buildParamsString(params)),
			}},
		})
		// set pos and end
		pos = params.StartPos
		end = params.EndPos
	}

	// create return values suggestion
	if returns != nil {
		suggestedFixes = append(suggestedFixes, analysis.SuggestedFix{
			Message: msg,
			TextEdits: []analysis.TextEdit{{
				Pos:     returns.StartPos,
				End:     returns.EndPos,
				NewText: []byte("\n" + buildParamsString(returns)),
			}},
		})
		// set pos if not set before
		if !pos.IsValid() {
			pos = returns.StartPos
		}
		// rewrite end
		end = returns.EndPos
	}

	pass.Report(analysis.Diagnostic{
		Pos:            pos,
		End:            end,
		Message:        msg,
		SuggestedFixes: suggestedFixes,
	})
}

func buildParamsString(params *Params) string {
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
