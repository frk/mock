// Package mock provides a type named Context which can be used to mock out one or more dependencies.
package mock

import (
	"log"
	"reflect"
)

type skip int

// X can be used to skip setting a passed in argument by pointer indirection.
const X skip = 0

// Call represents a call to one of the mocked functions/methods.
type Call struct {
	// The name of the function/method that was or is expected to be called.
	Func string
	// The input values that were or are expected to be passed to the Func called.
	In Vs
	// The output values the called Func should return.
	Out Vs
	// The values to be set to the call's passed in arguments by pointer indirection.
	Set Vs
}

// FN is a convenience constructor that returns a new instance of Call.
func FN(fn string, in ...interface{}) Call {
	return Call{Func: fn, In: Vs(in)}
}

// OUT returns a copy of the receiver with its Out set to the given out parameters.
func (c Call) OUT(out ...interface{}) Call {
	c.Out = Vs(out)
	return c
}

// SET returns a copy of the receiver with its Set set to the given set parameters.
func (c Call) SET(set ...interface{}) Call {
	c.Set = Vs(set)
	return c
}

// match is a helper method that checks whether the given Call has
// the same func name and the same arguments as the receiver Call.
func (c Call) match(k Call) error {
	if c.Func != k.Func {
		return &BadFuncCallError{got: c.Func, want: k.Func}
	}
	if !reflect.DeepEqual(c.In, k.In) {
		return &BadCallInputError{fn: c.Func, got: c.In, want: k.In}
	}
	return nil
}

// Context represents an aggregate of all mocked calls executed during a test.
// An instance of *Context is intended to be used by one or more mocked dependencies.
type Context struct {
	want []Call
	got  []Call
	errs []error
}

// New allocates and returns a new Context.
func New() *Context {
	return &Context{}
}

// Want allocates and returns a new Context with the given Call registered
// as expected by that Context.
func Want(c Call) *Context {
	return Wants([]Call{c})
}

// Wants allocates and returns a new Context with the given list of Calls
// registered as expected by that Context in that particular order.
func Wants(cs []Call) *Context {
	ctx := New()
	ctx.want = append(ctx.want, cs...)
	return ctx
}

// Want registers the given Call c as a call that the Context expects to happen.
func (ctx *Context) Want(c Call) {
	ctx.want = append(ctx.want, c)
}

// Wants registers the given slice of Calls as calls that the Context expects
// to happen. The Calls will be expected to happend in the same order in which
// they are provided in the slice.
func (ctx *Context) Wants(cs []Call) {
	ctx.want = append(ctx.want, cs...)
}

// Got registers the Call received and returns the return value of a matching
// expected Call, if there is no matching expected call it will return nil.
func (ctx *Context) Got(got Call) Vs {
	return ctx.match(got).Out
}

// match returns the matching expected Call to the given actual Call.
func (ctx *Context) match(got Call) Call {
	var want Call
	if ln := len(ctx.got); ln < len(ctx.want) {
		want = ctx.want[ln]
		if err := want.match(got); err != nil {
			ctx.errs = append(ctx.errs, err)
		}
	}
	ctx.got = append(ctx.got, got)

	for i, val := range want.Set {
		if val == X {
			continue
		}

		wv := reflect.ValueOf(val)
		gv := reflect.ValueOf(got.Set[i])

		if gv.Kind() == reflect.Ptr && (gv.Elem().Type() == wv.Type()) {
			gv.Elem().Set(wv)
		} else {
			log.Println("no match:", gv.Elem().Type(), wv.Type())
		}
	}
	return want
}

// Err checks the state of the Context and returns an error if any of it's
// expectations failed. It is intedend to be called at the end of each test case.
func (ctx *Context) Err() error {
	if got, want := len(ctx.got), len(ctx.want); got != want {
		return &BadNumCallError{got: got, want: want}
	}
	if len(ctx.errs) > 0 {
		return ctx.errs[0]
	}
	return nil
}

// Vs is a helper type with a number of methods that are useful for the Context.
type Vs []interface{}

// ValueAt returns the value at the given index as an interface{}.
// If there's no value at the given index the returned value will be nil.
func (vs Vs) ValueAt(index int) interface{} {
	if len(vs) > index {
		return vs[index]
	}
	return nil
}

// BoolAt returns the value at the given index as a bool. If there's no value
// at the given index, or its type is not bool, the returned value will be false.
func (vs Vs) BoolAt(index int) bool {
	if len(vs) > index {
		if v, ok := vs[index].(bool); ok {
			return v
		}
	}
	return false
}

// IntAt returns the value at the given index as an int. If there's no value at
// the given index, or its type is not int, the returned value will be 0.
func (vs Vs) IntAt(index int) int {
	if len(vs) > index {
		if v, ok := vs[index].(int); ok {
			return v
		}
	}
	return 0
}

// Float64At returns the value at the given index as a float64. If there's no value
// at the given index, or its type is not float64, the returned value will be 0.
func (vs Vs) Float64At(index int) float64 {
	if len(vs) > index {
		if v, ok := vs[index].(float64); ok {
			return v
		}
	}
	return 0
}

// StringAt returns the value at the given index as a string. If there's no value
// at the given index, or its type is not string, the returned value will be "".
func (vs Vs) StringAt(index int) string {
	if len(vs) > index {
		if v, ok := vs[index].(string); ok {
			return v
		}
	}
	return ""
}

// BytesAt returns the value at the given index as a []byte. If there's no value
// at the given index, or its type is not []byte, the returned value will be nil.
func (vs Vs) BytesAt(index int) []byte {
	if len(vs) > index {
		if v, ok := vs[index].([]byte); ok {
			return v
		}
	}
	return nil
}

// ErrorAt returns the value at the given index as an error. If there's no value
// at the given index, or its type is not error, the returned value will be nil.
func (vs Vs) ErrorAt(index int) error {
	if len(vs) > index {
		if v, ok := vs[index].(error); ok {
			return v
		}
	}
	return nil
}