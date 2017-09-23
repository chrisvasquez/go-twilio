package twilio

import "testing"

func TestSuccess(t *testing.T) {
	var success = true
	if success != true {
		t.Errorf("Failed a successful test.  Something really wrong here.")
	}
}

// For further implementation

