package mock

import (
	"fmt"
	"reflect"
)

type skip int

const X skip = 0

// Call represents a call to one of the mocked functions/methods.
type Call struct {
	Func string // The name of the function/method that was or is expected to be called.
	In   Vs     // The input values that were or are expected to be passed to the Func called.
	Out  Vs     // The output values the called Func should return.
	Set  Vs     // The values to be set to the call's passed in arguments by pointer indirection.
}

// match is a helper method that checks whether the given Call has
// the same func name and the same arguments as the receiver Call.
func (c Call) match(k Call) bool {
	return (c.Func == k.Func) && reflect.DeepEqual(c.In, k.In)
}

// Context
type Context struct {
	want []Call
	got  []Call
	res  []bool
	nth  int
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
	var ok bool
	var want Call
	if ln := len(ctx.got); ln < len(ctx.want) {
		if w := ctx.want[ln]; w.match(got) {
			want = w
			ok = true
		}
	}
	ctx.got = append(ctx.got, got)
	ctx.res = append(ctx.res, ok)

	for i, val := range want.Set {
		if val == X {
			continue
		}

		wv := reflect.ValueOf(val)
		gv := reflect.ValueOf(got.Set[i])

		if gv.Kind() == reflect.Ptr && (gv.Elem().Type() == wv.Type()) {
			gv.Elem().Set(wv)
		}
	}
	return want
}

func (ctx *Context) Err() error {
	if got, want := len(ctx.got), len(ctx.want); got != want {
		return &BadNumCallError{got: got, want: want}
	}
	for i, want := range ctx.want {
		got := ctx.got[i]

		if got.Func != want.Func {
			return &BadFuncCallError{got: got.Func, want: want.Func}
		}
		if !reflect.DeepEqual(got.In, want.In) {
			return &BadCallInputError{fn: got.Func, got: got.In, want: want.In}
		}
	}
	return nil
}

type Vs []interface{}

func (vs Vs) ValueAt(index int) interface{} {
	if len(vs) > index {
		return vs[index]
	}
	panic(fmt.Sprintf("mock.Vs.ValueAt: no value at index %d. (%#v)", index, vs))
	return nil
}

func (vs Vs) BoolAt(index int) bool {
	if len(vs) > index {
		if v, ok := vs[index].(bool); ok {
			return v
		}
	}
	panic(fmt.Sprintf("mock.Vs.BoolAt: no bool at index %d. (%#v)", index, vs))
	return false
}

func (vs Vs) IntAt(index int) int {
	if len(vs) > index {
		if v, ok := vs[index].(int); ok {
			return v
		}
	}
	panic(fmt.Sprintf("mock.Vs.IntAt: no int at index %d. (%#v)", index, vs))
	return 0
}

func (vs Vs) Float64At(index int) float64 {
	if len(vs) > index {
		if v, ok := vs[index].(float64); ok {
			return v
		}
	}
	panic(fmt.Sprintf("mock.Vs.Float64At: no float64 at index %d. (%#v)", index, vs))
	return 0
}

func (vs Vs) StringAt(index int) string {
	if len(vs) > index {
		if v, ok := vs[index].(string); ok {
			return v
		}
	}
	panic(fmt.Sprintf("mock.Vs.StringAt: no string at index %d. (%#v)", index, vs))
	return ""
}

func (vs Vs) ErrorAt(index int) error {
	if len(vs) > index {
		if v, ok := vs[index].(error); ok {
			return v
		}
	}
	panic(fmt.Sprintf("mock.Vs.ErrorAt: no error at index %d. (%#v)", index, vs))
	return nil
}