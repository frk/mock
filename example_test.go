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
	fetch(int, *user) error
}

// A 2nd interface to be mocked.
type database interface {
	save(*user) error
}

// Mock implementation of the api interface.
type mockapi struct{ *Context }

func (m *mockapi) fetch(id int, u *user) error {
	return m.Got(Call{Func: "fetch", In: Vs{id, u}, Set: Vs{u}}).ErrorAt(0)
}

// Mock implementation of the datebase interface.
type mockdb struct{ *Context }

func (m *mockdb) save(u *user) error {
	return m.Got(Call{Func: "save", In: Vs{u}}).ErrorAt(0)
}

// The function to be tested has two dependencies that can be mocked.
func dostuff(id int, cl api, db database) (*user, error) {
	u := &user{}
	if err := cl.fetch(id, u); err != nil {
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
		id    int
		calls []Call
		want  *user
		err   error
	}{{
		id: 111,
		calls: []Call{
			FN("fetch", 111, &user{}).OUT(clerr),
		},
		err: clerr,
	}, {
		id: 222,
		calls: []Call{
			FN("fetch", 222, &user{}).SET(user{"Joe"}),
			FN("save", &user{"Joe"}).OUT(dberr),
		},
		err: dberr,
	}, {
		id: 333,
		calls: []Call{
			FN("fetch", 333, &user{}).SET(user{"John"}),
			FN("save", &user{"John"}),
		},
		want: &user{"John"},
	}}

	for i, tt := range tests {
		// Construct a new *Context with the list of expected calls.
		m := Wants(tt.calls)

		// Have the mock dependencies use the *Context instance.
		cl := &mockapi{m}
		db := &mockdb{m}

		// Execute the function under test.
		got, err := dostuff(tt.id, cl, db)

		// Compare the output.
		if err != tt.err {
			fmt.Println(err)
		}
		if !reflect.DeepEqual(got, tt.want) {
			fmt.Printf("user got=%v, want=%v\n", got, tt.want)
		}

		// Check that *Context received all the expected calls.
		if err := m.Err(); err != nil {
			fmt.Println(i, err)
		}
	}

	// Output:
}
