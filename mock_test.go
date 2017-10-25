package mock

import (
	"errors"
	//"reflect"
	"testing"
)

type service interface {
	serve1()
	serve2(in1 string, in2 bool)
	serve3(in ...string) (int, error)
	serve4(in string, out *int) error
}

var _ service = &servicemock{}

type servicemock struct {
	*Context
}

func (m *servicemock) serve1() {
	m.Got(Call{Func: "serve1"})
}

func (m *servicemock) serve2(in1 string, in2 bool) {
	m.Got(Call{Func: "serve2", In: Vs{in1, in2}})
}

func (m *servicemock) serve3(in ...string) (int, error) {
	vs := Vs{}
	for _, s := range in {
		vs = append(vs, s)
	}

	out := m.Got(Call{Func: "serve3", In: vs})
	return out.IntAt(0), out.ErrorAt(1)
}

func (m *servicemock) serve4(in string, out *int) error {
	return m.Got(Call{Func: "serve4", In: Vs{in, out}, Set: Vs{out}}).ErrorAt(0)
}

func TestMock(t *testing.T) {
	_err := errors.New("E")
	_num := 5

	tests := []struct {
		name  string
		calls []Call
		exec  func(service, *testing.T)
	}{{
		name:  "no call",
		calls: []Call{},
		exec:  func(s service, t *testing.T) {},
	}, {
		name: "plain call",
		calls: []Call{
			{Func: "serve1"},
		},
		exec: func(s service, t *testing.T) {
			s.serve1()
		},
	}, {
		name: "plain call multi",
		calls: []Call{
			{Func: "serve1"},
			{Func: "serve1"},
			{Func: "serve1"},
		},
		exec: func(s service, t *testing.T) {
			s.serve1()
			s.serve1()
			s.serve1()
		},
	}, {
		name: "call with input",
		calls: []Call{
			{Func: "serve2", In: Vs{"hello", true}},
		},
		exec: func(s service, t *testing.T) {
			s.serve2("hello", true)
		},
	}, {
		name: "call with input multi",
		calls: []Call{
			{Func: "serve2", In: Vs{"foo", false}},
			{Func: "serve1"},
			{Func: "serve2", In: Vs{"bar", true}},
		},
		exec: func(s service, t *testing.T) {
			s.serve2("foo", false)
			s.serve1()
			s.serve2("bar", true)
		},
	}, {
		name: "call with input and output",
		calls: []Call{
			{Func: "serve3", In: Vs{"foo", "bar"}, Out: Vs{123, nil}},
		},
		exec: func(s service, t *testing.T) {
			got1, got2 := s.serve3("foo", "bar")
			want1, want2 := 123, error(nil)
			if got1 != want1 {
				t.Errorf("serve3 out1 got %d want %d", got1, want1)
			}
			if got2 != want2 {
				t.Errorf("serve3 out2 got %v want %v", got2, want2)
			}
		},
	}, {
		name: "call with input and output multi",
		calls: []Call{
			{Func: "serve3", In: Vs{"foo", "bar"}, Out: Vs{123, nil}},
			{Func: "serve1"},
			{Func: "serve2", In: Vs{"bar", true}},
			{Func: "serve3", In: Vs{"hello", "world", "!"}, Out: Vs{987, _err}},
		},
		exec: func(s service, t *testing.T) {
			got1, got2 := s.serve3("foo", "bar")
			want1, want2 := 123, error(nil)
			if got1 != want1 {
				t.Errorf("serve3 out1 got %d want %d", got1, want1)
			}
			if got2 != want2 {
				t.Errorf("serve3 out2 got %v want %v", got2, want2)
			}

			s.serve1()
			s.serve2("bar", true)

			got1, got2 = s.serve3("hello", "world", "!")
			want1, want2 = 987, _err
			if got1 != want1 {
				t.Errorf("serve3 out1 got %d want %d", got1, want1)
			}
			if got2 != want2 {
				t.Errorf("serve3 out2 got %v want %v", got2, want2)
			}
		},
	}, {
		name: "call with pointer indirection",
		calls: []Call{
			{Func: "serve4", In: Vs{"foobar", &_num}, Out: Vs{_err}, Set: Vs{557}},
		},
		exec: func(s service, t *testing.T) {
			num := _num
			got := s.serve4("foobar", &num)
			if num != 557 {
				t.Errorf("serve4 set got %d want %d", num, 557)
			}
			if got != _err {
				t.Errorf("serve4 out got %v want %v", got, _err)
			}

		},
	}}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			m := Wants(tt.calls)
			tt.exec(&servicemock{m}, t)

			if err := m.Err(); err != nil {
				t.Errorf("#%d: err %v", i, err)
			}
		})
	}
}
