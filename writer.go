package config

import "bufio"
import "io"
import "strconv"
import "strings"

type Writer interface {
	Bool(name string, val bool) Writer
	CloseSection() Writer
	Comment(comment string) Writer
	Float(name string, val float64) Writer
	Flush()
	Int(name string, val int64) Writer
	Line() Writer
	Section(name string, val string) Writer
	String(name string, val string) Writer
}

type writer struct {
	writer *bufio.Writer
	level  int
}

func NewWriter(wr io.Writer) Writer {
	o := new(writer)
	o.writer = bufio.NewWriter(wr)
	o.level = 0
	return o
}

var _ Writer = new(writer)

func (this *writer) wComment() *writer {
	this.writer.WriteByte('#')
	return this
}

func (this *writer) wIndent() *writer {
	for i := 0; i < this.level; i++ {
		this.writer.WriteString("    ")
	}
	return this
}

func (this *writer) wLevelDown() *writer {
	this.level = this.level - 1
	return this
}

func (this *writer) wLevelUp() *writer {
	this.level = this.level + 1
	return this
}

func (this *writer) wLine() *writer {
	this.writer.WriteByte('\n')
	return this
}

func (this *writer) wName(name string) *writer {
	this.writer.WriteString(name)
	return this
}

func (this *writer) wSectionEnd() *writer {
	this.writer.WriteByte('}')
	return this
}

func (this *writer) wSectionStart() *writer {
	this.writer.WriteByte('{')
	return this
}

func (this *writer) wSpace() *writer {
	this.writer.WriteByte(' ')
	return this
}

func (this *writer) wText(value string) *writer {
	this.writer.WriteString(value)
	return this
}

func (this *writer) wValue(value string) *writer {
	if this.isValueSafe(value) {
		this.writer.WriteString(value)
	} else {
		this.wValueEscaped(value)
	}
	return this
}

func (this *writer) wValueEscaped(value string) *writer {
	this.writer.WriteByte('"')
	for _, r := range value {
		if r == '"' {
			this.writer.WriteByte('\\')
		}
		this.writer.WriteRune(r)
	}
	this.writer.WriteByte('"')
	return this
}

func (this *writer) isValueSafe(value string) bool {
	for _, r := range value {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || (r == '_') || (r == '+') || (r == '-') || (r == '.')) {
			return false
		}
	}
	return true
}

func (this *writer) Bool(name string, val bool) Writer {
	return this.
		wIndent().
		wName(name).
		wSpace().
		wValue(strconv.FormatBool(val)).
		wLine()
}

func (this *writer) CloseSection() Writer {
	return this.
		wLevelDown().
		wIndent().
		wSectionEnd().
		wLine()
}

func (this *writer) Comment(comment string) Writer {
	if strings.IndexByte(comment, '\n') >= 0 {
		for _, commentLine := range strings.Split(comment, "\n") {
			this.Comment(commentLine)
		}
		return this
	} else {
		return this.
			wIndent().
			wComment().
			wSpace().
			wText(comment).
			wLine()
	}
}

func (this *writer) Float(name string, val float64) Writer {
	return this.
		wIndent().
		wName(name).
		wSpace().
		wValue(strconv.FormatFloat(val, 'g', -1, 64)).
		wLine()
}

func (this *writer) Flush() {
	this.writer.Flush()
}

func (this *writer) Int(name string, val int64) Writer {
	return this.
		wIndent().
		wName(name).
		wSpace().
		wValue(strconv.FormatInt(val, 10)).
		wLine()
}

func (this *writer) Line() Writer {
	return this.wLine()
}

func (this *writer) Section(name string, val string) Writer {
	return this.
		wIndent().
		wLevelUp().
		wName(name).
		wSpace().
		wValue(val).
		wSpace().
		wSectionStart().
		wLine()
}

func (this *writer) String(name string, val string) Writer {
	return this.
		wIndent().
		wName(name).
		wSpace().
		wValue(val).
		wLine()
}
