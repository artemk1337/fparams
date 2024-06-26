package analyzer

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "fargs",
		Doc:  "checks if function arguments are all on one line or each on a new line",
		Run:  run,
	}
}

type fargs struct {
	fset *token.FileSet
}

type FuncArgs struct {
	StartPos token.Pos
	EndPos   token.Pos
	Args     []*ast.Field
}

func run(pass *analysis.Pass) (any, error) {
	sla := &fargs{fset: pass.Fset}

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

func (s *fargs) checkFuncArgs(pass *analysis.Pass, fn *ast.FuncDecl) {
	var (
		args, rArgs *FuncArgs
	)

	if fn.Type.Params == nil || len(fn.Type.Params.List) == 0 {
		return // No arguments to check
	}

	if s.checkFuncDecOneLine(fn) {
		return // function declaration in one line
	}

	args = &FuncArgs{
		StartPos: fn.Type.Params.Pos() + 1,
		EndPos:   fn.Type.Params.End() - 1,
		Args:     fn.Type.Params.List,
	}
	if fn.Type.Results != nil {
		rArgs = &FuncArgs{
			StartPos: fn.Type.Results.Pos() + 1,
			EndPos:   fn.Type.Results.End() - 1,
			Args:     fn.Type.Results.List,
		}
	}

	// check and replace args
	argsValid, rArgsValid := s.validateFuncArgs(args, rArgs)

	if argsValid {
		args = nil
	}

	if rArgsValid {
		rArgs = nil
	}

	s.reportMultiLineArg(pass, fn.Name.String(), args, rArgs)

	return
}

func (s *fargs) checkFuncDecOneLine(fn *ast.FuncDecl) bool {
	bodyStartPos := s.fset.Position(fn.Body.Pos())
	fnStartPos := s.fset.Position(fn.Type.Pos())

	// TODO add extra check: bodyStartPos.Column < 120
	if fnStartPos.Line == bodyStartPos.Line {
		return true // function declaration in one line
	}

	return false
}

func (s *fargs) validateFuncArgs(args *FuncArgs, rArgs *FuncArgs) (argsValid, rArgsValid bool) {
	argsValid, rArgsValid = true, true

	if args != nil && len(args.Args) > 0 && !s.validateFuncEachArgs(args) {
		argsValid = false
	}

	if rArgs != nil && len(rArgs.Args) > 0 && !s.validateFuncEachArgs(rArgs) {
		rArgsValid = false
	}

	return argsValid, rArgsValid
}

func (s *fargs) validateFuncEachArgs(args *FuncArgs) bool {
	// prevPos start from "("
	// EndPos end on ")"
	prevPos := s.fset.Position(args.StartPos)

	// iterate on each arg and check positions
	for _, arg := range args.Args {
		argPos := s.fset.Position(arg.Pos())
		if prevPos.Line == argPos.Line {
			return false
		}

		prevPos = argPos
	}

	// extra check for last arg
	if s.fset.Position(args.Args[len(args.Args)-1].Pos()).Line == s.fset.Position(args.EndPos).Line {
		return false
	}

	return true
}

func (s *fargs) reportMultiLineArg(
	pass *analysis.Pass,
	fnName string,
	args *FuncArgs,
	rArgs *FuncArgs,
) {
	if args == nil && rArgs == nil {
		return
	}

	msg := fmt.Sprintf(`the arguments of the function "%s" should start on a new line`, fnName)
	s.reportV2(pass, msg, args, rArgs)
}

func (s *fargs) reportV2(
	pass *analysis.Pass,
	msg string,
	args *FuncArgs,
	rArgs *FuncArgs,
) {
	var pos, end token.Pos

	suggestedFixes := make([]analysis.SuggestedFix, 0, 2)

	if args != nil {
		suggestedFixes = append(suggestedFixes, analysis.SuggestedFix{
			Message: msg,
			TextEdits: []analysis.TextEdit{{
				Pos:     args.StartPos,
				End:     args.EndPos,
				NewText: []byte("\n" + buildArgs(pass, args)),
			}},
		})
		pos = args.StartPos
		end = args.EndPos
	}

	if rArgs != nil {
		suggestedFixes = append(suggestedFixes, analysis.SuggestedFix{
			Message: msg,
			TextEdits: []analysis.TextEdit{{
				Pos:     rArgs.StartPos,
				End:     rArgs.EndPos,
				NewText: []byte("\n" + buildArgs(pass, rArgs)),
			}},
		})
		if !pos.IsValid() {
			pos = rArgs.StartPos
		}
		end = rArgs.EndPos
	}

	pass.Report(analysis.Diagnostic{
		Pos:            pos,
		End:            end,
		Message:        msg,
		SuggestedFixes: suggestedFixes,
	})
}

func buildArgs(pass *analysis.Pass, args *FuncArgs) string {
	var builder strings.Builder

	for _, field := range args.Args {
		fieldType := typeExprToString(field.Type)
		for _, arg := range field.Names {
			builder.WriteString("\t" + arg.Name + " " + fieldType)
			if len(field.Names) > 1 || len(args.Args) > 1 {
				builder.WriteString(",\n")
			} else {
				builder.WriteString("\n")
			}
		}
	}

	return builder.String()
}

func typeExprToString(argType ast.Expr) string {
	return types.ExprString(argType)
}
