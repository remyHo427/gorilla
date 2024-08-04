package cpp

import "testing"

func TestNonDirectiveTokens(t *testing.T) {
	src := "int main(void) { return 0; }"
	p := NewParser(New(src))

	check(t, p, src)
}

func check(t *testing.T, p *Parser, expected string) {
	if expanded, err := p.Expand(); len(err) != 0 {
		for _, e := range err {
			t.Errorf(e.Error())
		}
	} else if expanded != expected {
		t.Errorf(`"%s" does not match expected, which is "%s"`,
			expanded, expected)
	}
}
