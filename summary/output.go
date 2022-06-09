package summary

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/fatih/color"
)

var (
	// DefaultPadding specifies the number of spaces between columns in a table.
	DefaultPadding = 2

	// DefaultWriter specifies the output io.Writer for the Table.Print method.
	DefaultWriter io.Writer = os.Stdout

	// DefaultHeaderFormatter specifies the default Formatter for the table header.
	DefaultHeaderFormatter Formatter

	// DefaultAllowFormatter specifies the default Formatter for the allow cells.
	DefaultAllowFormatter = color.New(color.FgGreen).SprintfFunc()
	// DefaultDenyFormatter specifies the default Formatter for the deny cells.
	DefaultDenyFormatter = color.New(color.FgRed).SprintfFunc()
	// DefaultAuditFormatter specifies the default Formatter for the audit cells.
	DefaultAuditFormatter = color.New(color.FgYellow).SprintfFunc()

	// DefaultWidthFunc specifies the default WidthFunc for calculating column widths
	DefaultWidthFunc WidthFunc = utf8.RuneCountInString
)

type Formatter func(string, ...interface{}) string

type WidthFunc func(string) int

type Table interface {
	WithHeaderFormatter(f Formatter) Table
	WithAllowFormatter(f Formatter) Table
	WithDenyFormatter(f Formatter) Table
	WithAuditFormatter(f Formatter) Table
	WithPadding(p int) Table
	WithWriter(w io.Writer) Table
	WithWidthFunc(f WidthFunc) Table

	AddRow(vals ...interface{}) Table
	SetRows(rows [][]string) Table
	Print()
}

func Heading(columnHeaders ...interface{}) Table {
	t := table{header: make([]string, len(columnHeaders))}

	t.WithPadding(DefaultPadding)
	t.WithWriter(DefaultWriter)
	t.WithHeaderFormatter(DefaultHeaderFormatter)
	t.WithAllowFormatter(DefaultAllowFormatter)
	t.WithDenyFormatter(DefaultDenyFormatter)
	t.WithAuditFormatter(DefaultAuditFormatter)
	t.WithWidthFunc(DefaultWidthFunc)

	for i, col := range columnHeaders {
		t.header[i] = fmt.Sprint(col)
	}

	return &t
}

type table struct {
	AllowFormatter  Formatter
	DenyFormatter   Formatter
	AuditFormatter  Formatter
	HeaderFormatter Formatter
	Padding         int
	Writer          io.Writer
	Width           WidthFunc

	header []string
	rows   [][]string
	widths []int
}

func (t *table) WithHeaderFormatter(f Formatter) Table {
	t.HeaderFormatter = f
	return t
}

func (t *table) WithAllowFormatter(f Formatter) Table {
	t.AllowFormatter = f
	return t
}

func (t *table) WithDenyFormatter(f Formatter) Table {
	t.DenyFormatter = f
	return t
}

func (t *table) WithAuditFormatter(f Formatter) Table {
	t.AuditFormatter = f
	return t
}

func (t *table) WithPadding(p int) Table {
	if p < 0 {
		p = 0
	}

	t.Padding = p
	return t
}

func (t *table) WithWriter(w io.Writer) Table {
	if w == nil {
		w = os.Stdout
	}

	t.Writer = w
	return t
}

func (t *table) WithWidthFunc(f WidthFunc) Table {
	t.Width = f
	return t
}

func (t *table) AddRow(vals ...interface{}) Table {
	row := make([]string, len(t.header))
	for i, val := range vals {
		if i >= len(t.header) {
			break
		}
		row[i] = fmt.Sprint(val)
	}
	t.rows = append(t.rows, row)

	return t
}

func (t *table) SetRows(rows [][]string) Table {
	t.rows = [][]string{}
	headerLength := len(t.header)

	for _, row := range rows {
		if len(row) > headerLength {
			t.rows = append(t.rows, row[:headerLength])
		} else {
			t.rows = append(t.rows, row)
		}
	}

	return t
}

func (t *table) Print() {
	if len(t.rows) == 0 {
		fmt.Println("No Data")
		return
	}
	format := strings.Repeat("%s", len(t.header)) + "\n"
	t.calculateWidths()

	t.printHeader(format)

	for _, row := range t.rows {
		t.printRow(format, row)
	}
}

func (t *table) printHeader(format string) {
	vals := t.applyWidths(t.header, t.widths)
	if t.HeaderFormatter != nil {
		txt := t.HeaderFormatter(format, vals...)
		fmt.Fprint(t.Writer, txt)
	} else {
		fmt.Fprintf(t.Writer, format, vals...)
	}
}

func (t *table) printRow(format string, row []string) {
	vals := t.applyWidths(row, t.widths)

	fmt.Fprintf(t.Writer, format, vals...)
}

func (t *table) calculateWidths() {
	t.widths = make([]int, len(t.header))
	for _, row := range t.rows {
		for i, v := range row {
			if w := t.Width(v) + t.Padding; w > t.widths[i] {
				t.widths[i] = w
			}
		}
	}

	for i, v := range t.header {
		if w := t.Width(v) + t.Padding; w > t.widths[i] {
			t.widths[i] = w
		}
	}
}

func (t *table) applyWidths(row []string, widths []int) []interface{} {
	out := make([]interface{}, len(row))
	allow := color.New(color.FgGreen).SprintfFunc()
	deny := color.New(color.FgRed).SprintfFunc()
	audit := color.New(color.FgYellow).SprintfFunc()
	for i, s := range row {
		switch s {
		case "ALLOW":
			s = allow(s)
		case "DENY", "BLOCK":
			s = deny(s)
		case "AUDIT":
			s = audit(s)
		}
		out[i] = s + t.lenOffset(s, widths[i])
	}
	return out
}

func (t *table) lenOffset(s string, w int) string {
	l := w - t.Width(s)
	if l <= 0 {
		return ""
	}
	return strings.Repeat(" ", l)
}
