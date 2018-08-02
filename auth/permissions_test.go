package auth

import "testing"

func TestFullUserScope(t *testing.T) {
	want := int32(140492255)
	if got := FullUserScope(); got != want {
		t.Errorf("Unexpected sum of all user scopes. Got: %v, Want: %v", got, want)
	}
}
