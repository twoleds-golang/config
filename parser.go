package config

import "bufio"
import "bytes"
import "errors"
import "fmt"
import "io"
import "os"

func ParseFromBytes(data []byte) (cfg Config, err error) {
	return ParseFromReader(bytes.NewReader(data))
}

func ParseFromByteReader(reader io.ByteReader) (cfg Config, err error) {
	p := newParser(reader)
	if err := p.parse(); err != nil {
		return nil, err
	}
	return p.config(), nil
}

func ParseFromFile(file string) (cfg Config, err error) {
	r, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return ParseFromReader(r)
}

func ParseFromReader(reader io.Reader) (cfg Config, err error) {
	return ParseFromByteReader(bufio.NewReader(reader))
}

func ParseFromString(str string) (cfg Config, err error) {
	return ParseFromBytes([]byte(str))
}

type parser struct {
	reader   io.ByteReader
	bufName  []byte
	bufValue []byte
	state    parserState
	builder  Builder
	curLine  uint32
	curCol   uint32
	lastByte byte
}

type parserState uint16

const (
	parserBegin parserState = iota
	parserComment
	parserName
	parserValue
	parserValueEnd
	parserValueEscaped
	parserValueStart
)

func newParser(reader io.ByteReader) *parser {
	p := new(parser)
	p.reader = reader
	p.bufName = make([]byte, 0, 128)
	p.bufValue = make([]byte, 0, 128)
	p.state = parserBegin
	p.builder = NewBuilder()
	p.curLine = 1
	p.curCol = 1
	p.lastByte = 0
	return p
}

func (this *parser) config() Config {
	return this.builder.Config()
}

func (this *parser) parse() error {

	for {

		b, err := this.reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				if this.state != parserBegin {
					return err
				}
				return nil
			} else {
				return err
			}
		}

		switch this.state {
		case parserBegin:
			if b == '#' {
				this.state = parserComment
			} else if b == '}' {
				this.builder.CloseSection()
				return nil
			} else if (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_' {
				this.bufName = append(this.bufName, b)
				this.state = parserName
			} else if b != ' ' && b != '\t' && b != '\r' && b != '\n' {
				return errors.New(fmt.Sprintf("Wrong character on line %d at column %d", this.curLine, this.curCol))
			}
		case parserComment:
			if b == '\n' {
				this.state = parserBegin
			}
		case parserName:
			if (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_' {
				this.bufName = append(this.bufName, b)
			} else if b == ' ' || b == '\t' {
				this.state = parserValueStart
			} else {
				// TODO: Throw an error
				panic("Wrong configuration")
			}
		case parserValueStart:
			if (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_' || b == '+' || b == '-' || b == '.' {
				this.bufValue = append(this.bufValue, b)
				this.state = parserValue
			} else if b == '"' {
				this.state = parserValueEscaped
			} else if b == '\r' || b == '\n' {
				this.builder.String(string(this.bufName), "")
				this.bufName = this.bufName[:0]
				this.state = parserBegin
			} else if b == '#' {
				this.builder.String(string(this.bufName), "")
				this.bufName = this.bufName[:0]
				this.state = parserComment
			} else if b == '{' {
				this.builder.Section(string(this.bufName), "")
				if err = this.parse(); err != nil {
					return err
				}
				this.state = parserBegin
			} else if b != ' ' && b != '\t' {
				// TODO: Throw an error
				panic("Wrong configuration")
			}
		case parserValue:
			if (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_' || b == '+' || b == '-' || b == '.' {
				this.bufValue = append(this.bufValue, b)
				this.state = parserValue
			} else if b == ' ' || b == '\t' {
				this.state = parserValueEnd
			} else if b == '\r' || b == '\n' {
				this.builder.String(string(this.bufName), string(this.bufValue))
				this.bufName = this.bufName[:0]
				this.bufValue = this.bufValue[:0]
				this.state = parserBegin
			} else {
				// TODO: Throw an error
				panic("Wrong configuration")
			}
		case parserValueEscaped:
			if b == '"' && this.lastByte != '\\' {
				this.state = parserValueEnd
			} else if (b == '\\' && this.lastByte == '\\') || (b != '\\') {
				this.bufValue = append(this.bufValue, b)
				fmt.Println(string(this.bufValue))
			}
			this.lastByte = b
		case parserValueEnd:
			if b == '\r' || b == '\n' {
				this.builder.String(string(this.bufName), string(this.bufValue))
				this.bufName = this.bufName[:0]
				this.bufValue = this.bufValue[:0]
				this.state = parserBegin
			} else if b == '#' {
				this.builder.String(string(this.bufName), string(this.bufValue))
				this.bufName = this.bufName[:0]
				this.bufValue = this.bufValue[:0]
				this.state = parserBegin
			} else if b == '{' {
				this.builder.Section(string(this.bufName), "")
				if err = this.parse(); err != nil {
					return err
				}
				this.state = parserBegin
			} else if b != ' ' && b != '\t' {
				// TODO: Throw an error
				panic("Wrong configuration")
			}
		}

		if b == '\n' {
			this.curLine = this.curLine + 1
			this.curCol = 1
		} else {
			this.curCol = this.curCol + 1
		}

	}

}
