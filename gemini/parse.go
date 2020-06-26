package gemini

import (
	"bufio"
	"io"
	"strings"
)

type LineType int

const (
	LinkLine LineType = iota
	PreLine
	TextLine
	ListLine
	H1Line
	H2Line
	H3Line
	QuoteLine
)

type Line struct {
	Type    LineType
	Raw     string
	Display string
	Link    string
}

// Parse a text/gemini document into a list of individual lines
func Itemize(r io.Reader) ([]Line, error) {
	pref := false
	lines := []Line{}

	br := bufio.NewReader(r)

	for {
		rawline, err := br.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		// Remove line terminators
		rawline = strings.TrimRight(rawline, "\r\n")

		line := Line{Raw: rawline}

		if strings.HasPrefix(rawline, "```") {
			// Toggle preformatting and don't output
			pref = !pref
			continue
		} else if pref {
			line.Type = PreLine
		} else if strings.HasPrefix(rawline, "=>") {
			line.Type = LinkLine

			l := strings.TrimSpace(rawline[2:])
			splitchar := " "
			ftab := strings.IndexByte(l, '\t')
			fspace := strings.IndexByte(l, ' ')
			if fspace == -1 {
				fspace = 10000
			}
			if fspace == -1 || (ftab != -1 && ftab < fspace) {
				splitchar = "\t"
			}

			fields := strings.SplitN(l, splitchar, 2)
			if len(fields) == 1 {
				// Duplicate for display
				fields = append(fields, fields[0])
			}
			line.Link = strings.TrimSpace(fields[0])
			line.Display = strings.TrimSpace(fields[1])
		} else if strings.HasPrefix(rawline, "* ") {
			line.Type = ListLine
			line.Display = strings.TrimSpace(rawline[2:])
		} else if strings.HasPrefix(rawline, ">") {
			line.Type = QuoteLine
			line.Display = strings.TrimSpace(rawline[1:])
		} else if strings.HasPrefix(rawline, "###") {
			line.Type = H3Line
			line.Display = strings.TrimSpace(rawline[3:])
		} else if strings.HasPrefix(rawline, "##") {
			line.Type = H2Line
			line.Display = strings.TrimSpace(rawline[2:])
		} else if strings.HasPrefix(rawline, "#") {
			line.Type = H1Line
			line.Display = strings.TrimSpace(rawline[1:])
		} else {
			line.Type = TextLine
			line.Display = strings.TrimSpace(rawline)
		}

		lines = append(lines, line)
	}

	return lines, nil
}
