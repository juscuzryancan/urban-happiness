package main

import "testing"

func TestOne(t *testing.T) {
	got := "Ryan"
	want := "Ryan"
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}
