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

func TestStrip(t *testing.T) {
	if out := pre("//abc\n"); out != " " {
		t.Errorf("single line comment not stripped")
	} else if out = pre("/*abc\nabc*/"); out != " " {
		t.Errorf("multi line comment not stripped")
	} else if out = pre("// /*abc\n //// /* // */"); out != "   " {
		t.Errorf("mixed single line and multi line comment not stripped")
	}
}
