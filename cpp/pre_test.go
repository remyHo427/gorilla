package cpp

import "testing"

func TestSplicing(t *testing.T) {
	if out := pre("this\\\nthis"); out != "thisthis" {
		t.Errorf(`expect "thisthis", got="%s"`, out)
	}
}

func TestTrigraph(t *testing.T) {
	input := `??= ??/ ??' ??( ??) ??! ??< ??> ??-`

	if out := pre(input); out != `# \ ^ [ ] | { } ~` {
		t.Errorf(`expected "# \ ^ [ ] | { } ~", got="%s"`, out)
	}
}
