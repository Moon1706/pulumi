// Copyright 2017 Pulumi, Inc. All rights reserved.

package eval

import (
	"github.com/pulumi/coconut/pkg/compiler/ast"
	"github.com/pulumi/coconut/pkg/compiler/symbols"
	"github.com/pulumi/coconut/pkg/diag"
	"github.com/pulumi/coconut/pkg/eval/rt"
	"github.com/pulumi/coconut/pkg/tokens"
)

// Invoker implements an intrinsic function's functionality.
type Invoker func(intrin *Intrinsic, e *evaluator, this *rt.Object, args []*rt.Object) *rt.Unwind

// Intrinsics contains the set of runtime functions that are callable by name through the Coconut standard library
// package.  Their functionality is implemented in the runtime because CocoIL cannot express the concepts they require
// to get their job done.  This includes things like dynamic introspection, invocation, and more.
var Intrinsics map[tokens.Token]Invoker

func init() {
	Intrinsics = map[tokens.Token]Invoker{
		"coconut:runtime:isFunction":    isFunction,
		"coconut:runtime:dynamicInvoke": dynamicInvoke,
	}
}

// Intrinsic is a special intrinsic function whose behavior is implemented by the runtime.
type Intrinsic struct {
	Node    diag.Diagable // the contextual node representing the place where this intrinsic got created.
	Func    ast.Function  // the underlying function's node (before mapping to an intrinsic).
	Nm      tokens.Name
	Tok     tokens.Token
	Sig     *symbols.FunctionType
	Invoker Invoker
}

var _ symbols.Function = (*Intrinsic)(nil)

func (node *Intrinsic) Symbol()                          {}
func (node *Intrinsic) Name() tokens.Name                { return node.Nm }
func (node *Intrinsic) Token() tokens.Token              { return node.Tok }
func (node *Intrinsic) Special() bool                    { return false }
func (node *Intrinsic) SpecialModInit() bool             { return false }
func (node *Intrinsic) Tree() diag.Diagable              { return node.Func }
func (node *Intrinsic) Function() ast.Function           { return node.Func }
func (node *Intrinsic) Signature() *symbols.FunctionType { return node.Sig }
func (node *Intrinsic) String() string                   { return string(node.Name()) }

// Invoke calls the intrinsic routine, passing additional context from the intrinsic struct itself.
func (node *Intrinsic) Invoke(e *evaluator, this *rt.Object, args []*rt.Object) *rt.Unwind {
	return node.Invoker(node, e, this, args)

}

// NewIntrinsic returns a new intrinsic function symbol with the given information.
func NewIntrinsic(tree diag.Diagable, fnc ast.Function, tok tokens.Token, nm tokens.Name,
	sig *symbols.FunctionType, invoker Invoker) *Intrinsic {
	return &Intrinsic{
		Node:    tree,
		Func:    fnc,
		Nm:      nm,
		Tok:     tok,
		Sig:     sig,
		Invoker: invoker,
	}
}

// MaybeIntrinsic checks whether the given symbol is an intrinsic and, if so, swaps it out with the actual runtime
// implementation of that intrinsic.  If the symbol is not an intrinsic, the original symbol is simply returned.
func MaybeIntrinsic(tree diag.Diagable, sym symbols.Symbol) symbols.Symbol {
	switch s := sym.(type) {
	case *Intrinsic:
		// Already an intrinsic; do not swap it out.
	case symbols.Function:
		// If this is a function whose token we recognize, create a new intrinsic symbol.  Note that we do not currently
		// cache these symbols because of the need to associate the AST node with the resulting symbol.
		tok := s.Token()
		if invoker, isintrinsic := Intrinsics[tok]; isintrinsic {
			sym = NewIntrinsic(tree, s.Function(), tok, tok.Name(), s.Signature(), invoker)
		}
	}
	return sym
}
