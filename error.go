package mock

import (
	"encoding/json"
	"fmt"
)

type BadNumCallError struct {
	got, want int
}

func (e *BadNumCallError) Error() string {
	return fmt.Sprintf("mock: Wrong number of calls; got %d, want %d.", e.got, e.want)
}

type BadFuncCallError struct {
	got, want string
}

func (e *BadFuncCallError) Error() string {
	return fmt.Sprintf("mock: Wrong func call; got %q, want %q.", e.got, e.want)
}

type BadCallInputError struct {
	fn        string
	got, want Vs
}

func (e *BadCallInputError) Error() string {
	got, err := json.Marshal(e.got)
	if err != nil {
		panic(fmt.Sprintf("mock: error marshaling got Vs %+v to json.\n %v", e.got, err))
	}
	want, err := json.Marshal(e.want)
	if err != nil {
		panic(fmt.Sprintf("mock: error marshaling want Vs %+v to json.\n %v", e.want, err))
	}
	return fmt.Sprintf("mock: %q Call In\n got: %s\nwant: %s", e.fn, got, want)
}

type BadCallSetLenError struct {
	fn        string
	got, want Vs
}

func (e *BadCallSetLenError) Error() string {
	return fmt.Sprintf("mock: %q Call Set\n got: %v\nwant: %v", e.fn, e.got, e.want)
}
