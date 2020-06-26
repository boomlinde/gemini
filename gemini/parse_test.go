package gemini

import (
	"bytes"
	"testing"
)

type tLine struct {
	rawtext       string
	expectType    LineType
	expectDisplay string
	expectLink    string
	skip          bool
}

func TestItemize(t *testing.T) {
	expect := []tLine{
		tLine{
			rawtext:       "# Hello",
			expectType:    H1Line,
			expectDisplay: "Hello",
		},
		tLine{
			rawtext:       "##Hello",
			expectType:    H2Line,
			expectDisplay: "Hello",
		},
		tLine{
			rawtext:       "Hello",
			expectType:    TextLine,
			expectDisplay: "Hello",
		},
		tLine{
			rawtext:       "   Hello",
			expectType:    TextLine,
			expectDisplay: "Hello",
		},
		tLine{
			rawtext:       "=>gopher://google.com",
			expectType:    LinkLine,
			expectDisplay: "gopher://google.com",
			expectLink:    "gopher://google.com",
		},
		tLine{
			rawtext:       "=> gopher://google.com",
			expectType:    LinkLine,
			expectDisplay: "gopher://google.com",
			expectLink:    "gopher://google.com",
		},
		tLine{
			rawtext:       "=>\t\t\tgopher://google.com\t\t\tgopher test!",
			expectType:    LinkLine,
			expectDisplay: "gopher test!",
			expectLink:    "gopher://google.com",
		},
		tLine{
			rawtext:       "=> gopher://google.com        \t     gopher test!",
			expectType:    LinkLine,
			expectDisplay: "gopher test!",
			expectLink:    "gopher://google.com",
		},
		tLine{
			rawtext: "```",
			skip:    true,
		},
		tLine{
			rawtext:    "=> gopher://google.com gopher test!",
			expectType: PreLine,
		},
		tLine{
			rawtext:    "=> gopher://google.com        \t     gopher test!",
			expectType: PreLine,
		},
		tLine{
			rawtext: "```",
			skip:    true,
		},
		tLine{
			rawtext:       "   Hello there",
			expectType:    TextLine,
			expectDisplay: "Hello there",
		},
		tLine{
			rawtext:       "=>gopher://google.com",
			expectType:    LinkLine,
			expectDisplay: "gopher://google.com",
			expectLink:    "gopher://google.com",
		},
		tLine{
			rawtext:       "=>",
			expectType:    LinkLine,
			expectDisplay: "",
			expectLink:    "",
		},
	}

	buf := []byte{}
	for _, l := range expect {
		buf = append(buf, []byte(l.rawtext)...)
		buf = append(buf, []byte("\r\n")...)
	}
	r := bytes.NewReader(buf)
	parsed, err := Itemize(r)
	if err != nil {
		t.Fatalf("failed to itemize: %v", err)
	}

	skipped := []tLine{}
	for _, l := range expect {
		if !l.skip {
			skipped = append(skipped, l)
		}
	}

	if len(parsed) != len(skipped) {
		t.Fatalf("expected parsed #lines to equal input #lines")
	}

	for i, l := range skipped {
		if l.expectType != parsed[i].Type {
			t.Errorf("line +%d: unexpected type", i)
		}
		if l.rawtext != parsed[i].Raw {
			t.Errorf("line +%d: unexpected Raw", i)
		}
		if l.expectDisplay != parsed[i].Display {
			t.Errorf("line +%d: unexpected Display", i)
		}
		if l.expectLink != parsed[i].Link {
			t.Errorf("line +%d: unexpected Link", i)
		}
	}
}
