package gemini

import (
	"fmt"
	"html"
	"io"
)

func ToHtml(lines []Line, w io.Writer) error {
	var err error
	errw := func(b string, p ...interface{}) {
		if err != nil {
			return
		}
		_, err = w.Write([]byte(fmt.Sprintf(b, p...)))
	}

	lastType := LineType(-1)
	for i, l := range lines {
		nextType := LineType(-1)
		if i < len(lines)-1 {
			nextType = lines[i+1].Type
		}

		switch l.Type {
		case LinkLine:
			errw("<a href=\"%s\">%s</a><br>\n", l.Link, html.EscapeString(l.Display))
		case PreLine:
			if lastType != PreLine {
				errw("<pre>")
			}
			errw(html.EscapeString(l.Raw))
			if nextType != PreLine {
				errw("</pre>")
			}
			errw("\n")
		case TextLine:
			errw("%s<br>\n", html.EscapeString(l.Display))
		case ListLine:
			if lastType != ListLine {
				errw("<ul>\n")
			}
			errw("<li>%s</li>\n", html.EscapeString(l.Display))
			if nextType != ListLine {
				errw("</ul>\n")
			}
		case H1Line:
			errw("<h1>%s</h1>\n", html.EscapeString(l.Display))
		case H2Line:
			errw("<h2>%s</h2>\n", html.EscapeString(l.Display))
		case H3Line:
			errw("<h3>%s</h3>\n", html.EscapeString(l.Display))
		case QuoteLine:
			errw("<i>%s</i><br>\n", html.EscapeString(l.Display))
		default:
			return fmt.Errorf("Unknown Type: %d", l.Type)
		}
		lastType = l.Type
	}

	return err
}
