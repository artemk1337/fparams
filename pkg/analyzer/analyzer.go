package analyzer

import (
	"fmt"
	"go/ast"
	"go/token"
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

type singleLineArgs struct {
	fset *token.FileSet
}

type fArgs struct {
	startPos token.Pos
	endPos   token.Pos
	args     []*ast.Field
}

func run(pass *analysis.Pass) (any, error) {
	sla := &singleLineArgs{fset: pass.Fset}

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

func (sla *singleLineArgs) checkFuncArgs(pass *analysis.Pass, fn *ast.FuncDecl) {
	var (
		args, rArgs *fArgs
	)

	if fn.Type.Params == nil || len(fn.Type.Params.List) == 0 {
		return // No arguments to check
	}

	if sla.checkFuncDecOneLine(fn) {
		return // function declaration in one line
	}

	args = &fArgs{
		startPos: fn.Type.Params.Pos() + 1,
		endPos:   fn.Type.Params.End() - 1,
		args:     fn.Type.Params.List,
	}
	if fn.Type.Results != nil {
		rArgs = &fArgs{
			startPos: fn.Type.Results.Pos() + 1,
			endPos:   fn.Type.Results.End() - 1,
			args:     fn.Type.Results.List,
		}
	}

	// check and replace args
	argsValid, rArgsValid := sla.validateFuncArgs(args, rArgs)

	if argsValid {
		args = nil
	}

	if rArgsValid {
		rArgs = nil
	}

	sla.reportMultiLineArg(pass, fn.Name.String(), args, rArgs)

	return
}

func (sla *singleLineArgs) checkFuncDecOneLine(fn *ast.FuncDecl) bool {
	bodyStartPos := sla.fset.Position(fn.Body.Pos())
	fnStartPos := sla.fset.Position(fn.Type.Pos())

	// TODO add extra check if need: bodyStartPos.Column < 120
	if fnStartPos.Line == bodyStartPos.Line {
		return true // function declaration in one line
	}

	return false
}

func (sla *singleLineArgs) validateFuncArgs(args *fArgs, rArgs *fArgs) (argsValid, rArgsValid bool) {
	argsValid, rArgsValid = true, true

	if args != nil && len(args.args) > 1 && !sla.validateFuncEachArgs(args) {
		argsValid = false
	}

	if rArgs != nil && len(rArgs.args) > 1 && !sla.validateFuncEachArgs(rArgs) {
		rArgsValid = false
	}

	return argsValid, rArgsValid
}

func (sla *singleLineArgs) validateFuncEachArgs(args *fArgs) bool {
	// prevPos start from "("
	// endPos end on ")"
	prevPos := sla.fset.Position(args.startPos)

	// iterate on each arg and check positions
	for _, arg := range args.args {
		argPos := sla.fset.Position(arg.Pos())
		if prevPos.Line == argPos.Line {
			return false
		}

		prevPos = argPos
	}

	// extra check for last arg
	if sla.fset.Position(args.args[len(args.args)-1].Pos()).Line == sla.fset.Position(args.endPos).Line {
		return false
	}

	return true
}

func (sla *singleLineArgs) reportMultiLineArg(
	pass *analysis.Pass,
	fnName string,
	args *fArgs,
	rArgs *fArgs,
) {
	if args == nil && rArgs == nil {
		return
	}

	msg := fmt.Sprintf(`the arguments of the function "%s" should start on a new line`, fnName)
	sla.reportV2(pass, msg, args, rArgs)
}

func (sla *singleLineArgs) reportV2(
	pass *analysis.Pass,
	msg string,
	args *fArgs,
	rArgs *fArgs,
) {
	var pos, end token.Pos

	suggestedFixes := make([]analysis.SuggestedFix, 0, 2)

	if args != nil {
		suggestedFixes = append(suggestedFixes, analysis.SuggestedFix{
			Message: msg,
			TextEdits: []analysis.TextEdit{{
				Pos:     args.startPos,
				End:     args.endPos,
				NewText: []byte("\n" + buildArgs(pass, args)),
			}},
		})
		pos = args.startPos
		end = args.endPos
	}

	if rArgs != nil {
		suggestedFixes = append(suggestedFixes, analysis.SuggestedFix{
			Message: msg,
			TextEdits: []analysis.TextEdit{{
				Pos:     rArgs.startPos,
				End:     rArgs.endPos,
				NewText: []byte("\n" + buildArgs(pass, rArgs)),
			}},
		})
		if !pos.IsValid() {
			pos = rArgs.startPos
		}
		end = rArgs.endPos
	}

	pass.Report(analysis.Diagnostic{
		Pos:            pos,
		End:            end,
		Message:        msg,
		SuggestedFixes: suggestedFixes,
	})
}

func buildArgs(pass *analysis.Pass, args *fArgs) string {
	argsSlice := make([]string, 0, len(args.args))
	for _, field := range args.args {
		astFieldType := pass.TypesInfo.TypeOf(field.Type)
		for _, arg := range field.Names {
			argsSlice = append(argsSlice, "\t"+arg.Name+" "+astFieldType.String()+",\n")
		}
	}

	return strings.Join(argsSlice, "")
}
