package mock

import (
	"errors"
	"fmt"
	"reflect"
)

type user struct {
	name string
}

// An interface to be mocked.
type api interface {
	fetch(*user) error
}

// A 2nd interface to be mocked.
type database interface {
	save(*user) error
}

// Mock implementation of the api interface.
type mockapi struct{ *Context }

func (m *mockapi) fetch(u *user) error {
	return m.Got(Call{Func: "fetch", In: Vs{u}, Set: Vs{u}}).ErrorAt(0)
}

// Mock implementation of the datebase interface.
type mockdb struct{ *Context }

func (m *mockdb) save(u *user) error {
	return m.Got(Call{Func: "save", In: Vs{u}}).ErrorAt(0)
}

// The function to be tested has two dependencies that can be mocked.
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
			FN("fetch", &user{}).OUT(clerr),
		},
		err: clerr,
	}, {
		calls: []Call{
			FN("fetch", &user{}).SET(user{"Joe"}),
			FN("save", &user{"Joe"}).OUT(dberr),
		},
		err: dberr,
	}, {
		calls: []Call{
			FN("fetch", &user{}).SET(user{"John"}),
			FN("save", &user{"John"}),
		},
		want: &user{"John"},
	}}

	for _, tt := range tests {
		// Construct a new *Context with the list of expected calls.
		m := Wants(tt.calls)

		// Have the mock dependencies use the *Context instance.
		cl := &mockapi{m}
		db := &mockdb{m}

		// Execute the function under test.
		got, err := dostuff(cl, db)

		// Compare the output.
		if err != tt.err {
			fmt.Println(err)
		}
		if !reflect.DeepEqual(got, tt.want) {
			fmt.Printf("user got=%v, want=%v\n", got, tt.want)
		}

		// Check that *Context received all the expected calls.
		if err := m.Err(); err != nil {
			fmt.Println(err)
		}
	}

	// Output:
}
