package gemini

import (
	"bytes"
	"testing"
)

func TestToHtml(t *testing.T) {
	join := func(s []string, sep string) []byte {
		out := []byte{}
		for _, v := range s {
			out = append(out, []byte(v+sep)...)
		}
		return out
	}
	expected := join([]string{
		"<h1>This is a test</h1>",
		"This is some plain text. Here&#39;s a list instead:<br>",
		"<br>",
		"<ul>",
		"<li>Item 1</li>",
		"<li>Item 2</li>",
		"<li>Item 3</li>",
		"</ul>",
		"<br>",
		"Here&#39;s some preformatted text:<br>",
		"<br>",
		"<pre>func main() {",
		"    doc:= []string{}",
		"}</pre>",
		"<br>",
		"Some quote:<br>",
		"<br>",
		"<i>I have never been so insulted!</i><br>",
		"<i>This is the worst I&#39;ve ever been treated!</i><br>",
		"<br>",
		"Finally, some links for you:<br>",
		"<br>",
		"<a href=\"gemini://example.com?t=7&h=9\">gemini://example.com?t=7&amp;h=9</a><br>",
		"<a href=\"gopher://text.garden\">My gopherhole</a><br>",
	}, "\n")

	serialized := join([]string{
		"# This is a test",
		"This is some plain text. Here's a list instead:",
		"",
		"* Item 1",
		"* Item 2",
		"* Item 3",
		"",
		"Here's some preformatted text:",
		"",
		"```",
		"func main() {",
		"    doc:= []string{}",
		"}",
		"```",
		"",
		"Some quote:",
		"",
		"> I have never been so insulted!",
		"> This is the worst I've ever been treated!",
		"",
		"Finally, some links for you:",
		"",
		"=> gemini://example.com?t=7&h=9",
		"=> gopher://text.garden My gopherhole",
	}, "\r\n")

	r := bytes.NewReader(serialized)
	lines, err := Itemize(r)
	if err != nil {
		t.Fatalf("failed to itemize: %v", err)
	}

	b := []byte{}
	w := bytes.NewBuffer(b)

	if err := ToHtml(lines, w); err != nil {
		t.Fatalf("failed to convert to html: %v", err)
	}

	if !bytes.Equal(expected, w.Bytes()) {
		t.Error("output not as expected")
	}
}
