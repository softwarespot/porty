package ports

import (
	"strconv"
	"testing"

	testhelpers "github.com/softwarespot/porty/test-helpers"
)

func Test_New(t *testing.T) {
	p := New()

	// Should return an error if the username and app name doesn't exist
	up, err := p.Get("test-user", "test-app-name-1")
	testhelpers.AssertError(t, err)
	testhelpers.AssertEqual(t, up, userPortOnError)

	port, err := p.Next()
	testhelpers.AssertNoError(t, err)
	testhelpers.AssertEqual(t, port, 8000)

	// Should regsiter the next port, for a username and app name
	up, err = p.Register("test-user", "test-app-name-1")
	testhelpers.AssertNoError(t, err)
	testhelpers.AssertEqual(t, up.Port, 8000)

	up, err = p.Get("test-user", "test-app-name-1")
	testhelpers.AssertNoError(t, err)
	testhelpers.AssertEqual(t, up.Port, 8000)

	// Should get the next available port
	port, err = p.Next()
	testhelpers.AssertNoError(t, err)
	testhelpers.AssertEqual(t, port, 8001)

	// Should not register the next port, if the username and app name already exists
	up, err = p.Register("test-user", "test-app-name-1")
	testhelpers.AssertError(t, err)
	testhelpers.AssertEqual(t, up.Port, 0)

	// Should unregister the next port, for a username and app name
	up, err = p.Unregister("test-user", "test-app-name-1")
	testhelpers.AssertNoError(t, err)
	testhelpers.AssertEqual(t, up.Port, 8000)

	up, err = p.Get("test-user", "test-app-name-1")
	testhelpers.AssertError(t, err)
	testhelpers.AssertEqual(t, up, userPortOnError)

	// Should get the next available port
	port, err = p.Next()
	testhelpers.AssertNoError(t, err)
	testhelpers.AssertEqual(t, port, 8000)

	// Should register many ports
	for i := DefaultMinPort; i <= DefaultMaxPort; i++ {
		up, err := p.Register("test-user", strconv.Itoa(int(i)))
		testhelpers.AssertNoError(t, err)
		testhelpers.AssertEqual(t, i, up.Port)
	}
	for i := DefaultMinPort; i <= DefaultMaxPort; i++ {
		up, err := p.Get("test-user", strconv.Itoa(int(i)))
		testhelpers.AssertNoError(t, err)
		testhelpers.AssertEqual(t, i, up.Port)
	}

	// Should not register, when all the ports have been used
	up, err = p.Register("test-user", "test-app-name-1")
	testhelpers.AssertError(t, err)
	testhelpers.AssertEqual(t, up.Port, 0)
}
