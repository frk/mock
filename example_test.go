package mock

import (
	"errors"
	"fmt"
	"reflect"
)

type user struct {
	name string
}

type api interface {
	fetch(*user) error
}

type database interface {
	save(*user) error
}

type mockapi struct{ *Context }

func (m *mockapi) fetch(u *user) error {
	return m.Got(Call{Func: "fetch", In: Vs{u}, Set: Vs{u}}).ErrorAt(0)
}

type mockdb struct{ *Context }

func (m *mockdb) save(u *user) error {
	return m.Got(Call{Func: "save", In: Vs{u}}).ErrorAt(0)
}

func dostuff(cl api, db database) (*user, error) {
	u := &user{}
	if err := cl.fetch(u); err != nil {
		return nil, err
	}

	// ...

	if err := db.save(u); err != nil {
		return nil, err
	}

	return u, nil
}

func ExampleContext() {

	clerr := errors.New("api error")
	dberr := errors.New("db error")

	tests := []struct {
		calls []Call
		want  *user
		err   error
	}{{
		calls: []Call{
			{Func: "fetch", In: Vs{&user{}}, Out: Vs{clerr}},
		},
		err: clerr,
	}, {
		calls: []Call{
			{Func: "fetch", In: Vs{&user{}}, Out: Vs{nil}, Set: Vs{user{"Joe"}}},
			{Func: "save", In: Vs{&user{"Joe"}}, Out: Vs{dberr}},
		},
		err: dberr,
	}, {
		calls: []Call{
			{Func: "fetch", In: Vs{&user{}}, Out: Vs{nil}, Set: Vs{user{"John"}}},
			{Func: "save", In: Vs{&user{"John"}}, Out: Vs{nil}},
		},
		want: &user{"John"},
	}}

	for _, tt := range tests {
		m := Wants(tt.calls)
		cl := &mockapi{m}
		db := &mockdb{m}

		got, err := dostuff(cl, db)
		if err != tt.err {
			fmt.Println(err)
		}
		if !reflect.DeepEqual(got, tt.want) {
			fmt.Printf("user got=%v, want=%v\n", got, tt.want)
		}

		if err := m.Err(); err != nil {
			fmt.Println(err)
		}
	}

	// Output:
}
