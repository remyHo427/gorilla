package parse

import (
	"gorilla/lex"
	"testing"
)

type Pair struct {
	input  string
	output string
}

func check(t *testing.T, tt []Pair) {
	for i, test := range tt {
		l := lex.New(test.input)
		p := New(l)
		tree, err := p.Parse()

		for _, e := range err {
			t.Errorf(e.Error())
		}

		if len(p.err) > 0 {
			continue
		} else if len(tree) == 0 {
			t.Errorf("no ast produced")
			continue
		}

		if test.output != tree[0].String() {
			t.Errorf("expected \"%s\", got \"%s\" at tt[%d]",
				test.output, tree[0].String(), i)
		}
	}
}
