package client

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClientApp_Stop(t *testing.T) {
	// Prepare test data
	cfg := Config{}
	app := NewClientApp(cfg)

	// Mock deferred operations
	var calledOps []bool
	for i := 0; i < 3; i++ {
		index := i // Capture the loop variable
		op := func() error {
			calledOps[index] = true
			return nil
		}
		app.deferredOps = append(app.deferredOps, op)
		calledOps = append(calledOps, false)
	}

	app.deferredOps = append(app.deferredOps, func() error {
		calledOps[len(calledOps)-1] = true
		return errors.New("deferred operation error")
	})
	calledOps = append(calledOps, false)

	err := app.Stop()

	// Check if all deferred operations were called
	for i, called := range calledOps {
		if !called && i != 3 {
			t.Errorf("Deferred operation at index %d was not called", i)
		}
	}

	require.Errorf(t, err, "deferred operation error")

}

// TestServerApp_Stop_Error - test the Stop method when a deferred operation returns an error
func TestServerApp_Stop_Error(t *testing.T) {
	// Prepare test data
	cfg := Config{}
	app := NewClientApp(cfg)

	// Mock a deferred operation that returns an error
	errMsg := "deferred operation error"
	op := func() error {
		return errors.New(errMsg)
	}
	app.deferredOps = append(app.deferredOps, op)

	// Call the Stop method
	err := app.Stop()

	// Check the error returned by Stop method
	if err == nil {
		t.Error("expected an error from Stop method, but got nil")
	} else if err.Error() != errMsg {
		t.Errorf("expected error message '%s', but got '%s'", errMsg, err.Error())
	}
}
