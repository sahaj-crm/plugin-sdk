package extensionapi

import "testing"

func TestDecodeEntityInput(t *testing.T) {
	ev := &Event{
		Input: []byte(`{"Record":{"ID":4,"Name":"Course"},"Extras":{"com.example.optin":{"opt_in":"true","tags":["a","b"]}}}`),
	}
	input, err := DecodeEntityInput[map[string]any](ev)
	if err != nil {
		t.Fatal(err)
	}
	if input.Record["Name"] != "Course" {
		t.Fatalf("unexpected record: %#v", input.Record)
	}
	values := input.Extras.Values("com.example.optin")
	if !AsBool(values["opt_in"]) {
		t.Fatalf("expected opt_in true, got %#v", values["opt_in"])
	}
	if got := AsStrings(values["tags"]); len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Fatalf("unexpected tags: %#v", got)
	}
}

func TestDecodeOutput(t *testing.T) {
	ev := &Event{Output: []byte(`{"ID":4,"Name":"Course"}`)}
	out, err := DecodeOutput[map[string]any](ev)
	if err != nil {
		t.Fatal(err)
	}
	if out["ID"] != float64(4) {
		t.Fatalf("unexpected output: %#v", out)
	}
}

func TestAsBool(t *testing.T) {
	cases := []struct {
		in   any
		want bool
	}{
		{true, true},
		{"true", true},
		{"TRUE", true},
		{"false", false},
		{float64(1), true},
		{nil, false},
	}
	for _, tc := range cases {
		if got := AsBool(tc.in); got != tc.want {
			t.Fatalf("AsBool(%#v) = %v, want %v", tc.in, got, tc.want)
		}
	}
}
