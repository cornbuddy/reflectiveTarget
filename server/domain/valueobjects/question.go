package valueobjects

import (
	"fmt"
	"strings"
)

type Question struct {
	ID

	Text string
}

type Questions []Question

func (q *Question) String() string {
	return fmt.Sprintf("{ID: %d, Text: '%s'}", q.ID, q.Text)
}

func (qs Questions) String() string {
	var sb strings.Builder
	for i, q := range qs {
		fmt.Fprintf(&sb, "%s", q.String())
		if i < len(qs)-1 {
			fmt.Fprintf(&sb, ", ")
		}
	}

	return fmt.Sprintf("[%s]", sb.String())
}
