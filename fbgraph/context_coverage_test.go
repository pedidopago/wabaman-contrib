package fbgraph

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestNoContextFreeRequests keeps the package migrated. Every request the
// package issues must carry its caller's context: one built without it ignores
// cancellation, so the caller waits out the client timeout -- 120s on
// DefaultHTTPClient -- on work that is already doomed. Threading the context
// through was a breaking change for consumers, so a straggler reintroduced
// later costs them a second one. Building the request with context.Background()
// counts as a straggler: it compiles, it reads as migrated, and it discards the
// cancellation just as thoroughly as no context at all.
func TestNoContextFreeRequests(t *testing.T) {
	var offenders []string

	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("reading package directory: %v", err)
	}

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		file, err := parser.ParseFile(token.NewFileSet(), name, nil, 0)
		if err != nil {
			t.Fatalf("parsing %s: %v", name, err)
		}

		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			switch types.ExprString(call.Fun) {
			case "NewRequest", "http.NewRequest":
				// lowlevel.go's own http.NewRequest is the deprecated wrapper's
				// body; it is exported for consumers, so it stays. Only the rest
				// of the package is held to the rule.
				if name == "lowlevel.go" {
					return true
				}

				offenders = append(offenders, filepath.Join("fbgraph", name)+": builds a request with no context at all")
			case "NewRequestWithContext", "http.NewRequestWithContext":
				// Naming the callee is not enough: passing context.Background()
				// satisfies it while discarding the caller's cancellation, which
				// is the exact bug this package just spent a breaking change
				// removing from eleven call sites.
				if len(call.Args) == 0 {
					return true
				}
				// A deliberate detach still derives from the caller's context;
				// it is a decision, not a discarded one. Anything that does not
				// mention ctx at all is the straggler this looks for.
				arg := types.ExprString(call.Args[0])
				if !strings.Contains(arg, "ctx") {
					offenders = append(offenders,
						filepath.Join("fbgraph", name)+": passes "+arg+" instead of the caller's ctx")
				}
			}

			return true
		})
	}

	if len(offenders) > 0 {
		t.Errorf("these files build requests without a context, so cancellation cannot reach them: %v", offenders)
	}
}

// TestErrorAccessorsResetOnEntry keeps the LastGraphError/LastErrorRawBody
// invariant true rather than nearly true. errorFromResponse always refreshes the
// raw body but only sets lastGraphError when the payload carries an error code,
// so a method that does not clear on entry can report one call's Graph error
// paired with a later call's raw body -- a mismatch the caller cannot detect.
//
// Half of the methods honouring this is worse than none of them: callers start
// trusting an accessor that is only sometimes about the call they just made.
// This is a guard rather than a behavioural test because the property is "every
// method does it", which no single call can demonstrate.
func TestErrorAccessorsResetOnEntry(t *testing.T) {
	var missing []string

	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("reading package directory: %v", err)
	}

	files := map[string]*ast.File{}
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}

		file, err := parser.ParseFile(token.NewFileSet(), name, nil, 0)
		if err != nil {
			t.Fatalf("parsing %s: %v", name, err)
		}
		files[name] = file
	}

	calls := callGraph(files)
	reporting := closeOver(calls, reportingSeeds)

	// A method satisfies the invariant by clearing the state itself or by
	// delegating to something that clears it before issuing the request --
	// SendCallPermissionRequest goes through SendMessage, the template updaters
	// through post. Requiring a redundant second reset in those would be
	// pattern-matching the rule rather than applying it.
	//
	// The delegation half is position-blind: a method that issued its own
	// request and only afterwards called a resetter would satisfy this while
	// never clearing on entry. No such method exists, and detecting one needs
	// real flow analysis, so the direct check below carries the positional rule
	// and this carries reachability.
	directResetters := map[string]bool{}
	for _, file := range files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if ok && fn.Body != nil && resetsErrorState(fn) {
				directResetters[fn.Name.Name] = true
			}
		}
	}
	resetting := closeOver(calls, directResetters)

	for name, file := range files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			// Exported methods only. Unexported helpers like httpError are
			// called from inside a method that has already reset, and are
			// reached after the request: clearing there would wipe the very
			// error they exist to report.
			if !ok || fn.Body == nil || !isClientMethod(fn) || !fn.Name.IsExported() {
				continue
			}
			if !reporting[fn.Name.Name] {
				continue
			}
			if !resetting[fn.Name.Name] {
				missing = append(missing, name+":"+fn.Name.Name)
			}
		}
	}

	if len(missing) > 0 {
		t.Errorf("these methods report a Graph error but never clear the previous call's: %v", missing)
	}
}

func isClientMethod(fn *ast.FuncDecl) bool {
	if fn.Recv == nil || len(fn.Recv.List) != 1 {
		return false
	}

	return types.ExprString(fn.Recv.List[0].Type) == "*Client"
}

// reportingSeeds are the two helpers that publish into lastGraphError /
// lastErrorRawBody. Everything that reaches them, at any depth, is reporting a
// Graph error and therefore has to clear the previous call's first.
var reportingSeeds = map[string]bool{"errorFromResponse": true, "httpError": true}

// reportingFuncs resolves that transitively across the package. Matching direct
// calls only is not enough and neither is adding one more name: the guard first
// missed five methods that reach errorFromResponse through httpError, and then
// still missed GetWABAInfo, which reaches it through an unexported helper. A
// guard that reports green on the state it exists to prevent is worse than none.
func callGraph(files map[string]*ast.File) map[string]map[string]bool {
	calls := map[string]map[string]bool{}

	for _, file := range files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}

			callees := map[string]bool{}
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				// Matched by bare identifier: c.post and a package-level post
				// are one node here, as are same-named methods on different
				// receivers. Accurate for this package; it would over-mark
				// rather than under-mark if a name ever collides, which is the
				// safe direction for a guard.
				name := types.ExprString(call.Fun)
				if i := strings.LastIndex(name, "."); i >= 0 {
					name = name[i+1:]
				}
				callees[name] = true

				return true
			})
			calls[fn.Name.Name] = callees
		}
	}

	return calls
}

// closeOver marks every function that reaches a seed through any chain of calls.
func closeOver(calls map[string]map[string]bool, seeds map[string]bool) map[string]bool {
	marked := map[string]bool{}
	for name := range seeds {
		marked[name] = true
	}

	for changed := true; changed; {
		changed = false
		for fn, callees := range calls {
			if marked[fn] {
				continue
			}
			for callee := range callees {
				if marked[callee] {
					marked[fn] = true
					changed = true

					break
				}
			}
		}
	}

	return marked
}

// resetsErrorState reports whether fn clears lastGraphError before doing any
// work. Only the opening statements count: a reset after the request has been
// issued would clobber the error it is meant to report.
func resetsErrorState(fn *ast.FuncDecl) bool {
	var clearedError, clearedBody bool

	for _, stmt := range fn.Body.List {
		assign, ok := stmt.(*ast.AssignStmt)
		if !ok || len(assign.Lhs) != 1 {
			break
		}
		switch types.ExprString(assign.Lhs[0]) {
		case "c.lastGraphError":
			clearedError = true
		case "c.lastErrorRawBody":
			clearedBody = true
		}
	}

	// Both, or the pair can still be mismatched: a fresh nil error read
	// alongside a previous call's raw body.
	return clearedError && clearedBody
}
