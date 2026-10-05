package main

import "testing"

func TestGraphVersionPattern(t *testing.T) {
	for v, want := range map[string]bool{"v26.0": true, "v9.1": true, "26.0": false, "v26": false, "v26.0/../x": false, "": false} {
		if got := graphVersionPattern.MatchString(v); got != want {
			t.Errorf("%q: got %v, want %v", v, got, want)
		}
	}
	if graphAPIVersion != "v26.0" {
		t.Errorf("default graphAPIVersion = %q, want v26.0", graphAPIVersion)
	}
}
