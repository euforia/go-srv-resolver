package resolver

import "testing"

func TestResolver(t *testing.T) {
	// Test against local consul instance.
	rsl := NewResolver(8600, "127.0.0.1")
	_, err := rsl.Lookup("consul.service.consul")
	if err != nil {
		t.Fatal(err)
	}
}
