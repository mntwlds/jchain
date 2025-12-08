package jchain

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

type parser struct {
	input    string
	len      int
	maxDepth int
	depth    int
}

func parseJSON(jsonStr string, maxDepth int) (res any, err error) {
	p := &parser{
		input:    jsonStr,
		len:      len(jsonStr),
		maxDepth: maxDepth,
	}

	defer func() {
		if r := recover(); r != nil {
			if pErr, ok := r.(parserError); ok {
				line, col := p.calculateLineCol(pErr.pos)
				err = fmt.Errorf("%s at line %d, column %d", pErr.msg, line, col)
			} else {
				err = fmt.Errorf("%v", r)
			}
			res = nil // Ensure result is nil on error
		}
	}()

	i := p.skipWhitespace(0)
	value, i := p.parseValue(i)
	i = p.skipWhitespace(i)

	if !p.checkOOB(i) {
		p.error(i, "Invalid JSON")
	}

	return value, nil
}

func (p *parser) error(pos int, msg string) {
	panic(parserError{pos: pos, msg: msg})
}

type parserError struct {
	pos int
	msg string
}

func (p *parser) checkOOB(i int) bool {
	return i >= p.len
}

func (p *parser) calculateLineCol(pos int) (int, int) {
	line := 1
	col := 1
	for i := 0; i < pos && i < p.len; i++ {
		if p.input[i] == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}
	return line, col
}

func (p *parser) parseValue(i int) (any, int) {
	if p.checkOOB(i) {
		p.error(i, "Invalid JSON")
	}

	switch p.input[i] {
	case '{':
		return p.parseObject(i)
	case '[':
		return p.parseArray(i)
	case '"':
		return p.parseString(i)
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return p.parseNumber(i)
	case 't':
		if i+4 > p.len || p.input[i:i+4] != "true" {
			p.error(i, "Invalid JSON")
		}
		return true, i + 4
	case 'f':
		if i+5 > p.len || p.input[i:i+5] != "false" {
			p.error(i, "Invalid JSON")
		}
		return false, i + 5
	case 'n':
		if i+4 > p.len || p.input[i:i+4] != "null" {
			p.error(i, "Invalid JSON")
		}
		return nil, i + 4
	default:
		p.error(i, "Invalid JSON")
	}
	return nil, i // Should be unreachable
}

func (p *parser) skipWhitespace(i int) int {
	for i < p.len && (p.input[i] == ' ' || p.input[i] == '\t' || p.input[i] == '\n' || p.input[i] == '\r') {
		i++
	}
	return i
}

func (p *parser) parseObject(i int) (map[string]any, int) {
	if p.maxDepth > 0 && p.depth >= p.maxDepth {
		p.error(i, "Maximum depth exceeded")
	}
	p.depth++
	defer func() { p.depth-- }()

	if i >= p.len || p.input[i] != '{' {
		p.error(i, "Invalid JSON")
	}
	i++
	jsonMap := make(map[string]any)
	i = p.skipWhitespace(i)
	if i < p.len && p.input[i] != '}' {
		for {
			var key string
			key, i = p.parseString(i)
			i = p.skipWhitespace(i)
			if i >= p.len || p.input[i] != ':' {
				p.error(i, "Invalid JSON")
			}
			i++
			i = p.skipWhitespace(i)
			var value any
			value, i = p.parseValue(i)
			i = p.skipWhitespace(i)
			jsonMap[key] = value
			if i >= p.len {
				p.error(i, "Invalid JSON")
			}
			if p.input[i] == ',' {
				i++
				i = p.skipWhitespace(i)
			} else if p.input[i] == '}' {
				i++
				break
			} else {
				p.error(i, "Invalid JSON")
			}
		}
	} else if i < p.len && p.input[i] == '}' {
		i++
	} else {
		p.error(i, "Invalid JSON")
	}
	return jsonMap, i
}

func (p *parser) parseArray(i int) (any, int) {
	if p.maxDepth > 0 && p.depth >= p.maxDepth {
		p.error(i, "Maximum depth exceeded")
	}
	p.depth++
	defer func() { p.depth-- }()

	if i >= p.len || p.input[i] != '[' {
		p.error(i, "Invalid JSON")
	}
	i++
	jsonArray := make([]any, 0)
	i = p.skipWhitespace(i)
	if i < p.len && p.input[i] != ']' {
		for {
			var value any
			value, i = p.parseValue(i)
			i = p.skipWhitespace(i)
			jsonArray = append(jsonArray, value)
			if i >= p.len {
				p.error(i, "Invalid JSON")
			}
			if p.input[i] == ',' {
				i++
				i = p.skipWhitespace(i)
			} else if p.input[i] == ']' {
				i++
				break
			} else {
				p.error(i, "Invalid JSON")
			}
		}
	} else if i < p.len && p.input[i] == ']' {
		i++
	} else {
		p.error(i, "Invalid JSON")
	}
	return jsonArray, i
}

func (p *parser) parseString(i int) (string, int) {
	if i >= p.len || p.input[i] != '"' {
		p.error(i, "Invalid JSON")
	}
	i++
	var sb strings.Builder
	if i >= p.len || p.input[i] != '"' {
		for i < p.len {
			if p.input[i] == '\\' {
				i++
				switch p.input[i] {
				case '"':
					sb.WriteByte('"')
					i++
				case '\\':
					sb.WriteByte('\\')
					i++
				case '/':
					sb.WriteByte('/')
					i++
				case 'b':
					sb.WriteByte('\b')
					i++
				case 'f':
					sb.WriteByte('\f')
					i++
				case 'n':
					sb.WriteByte('\n')
					i++
				case 'r':
					sb.WriteByte('\r')
					i++
				case 't':
					sb.WriteByte('\t')
					i++
				case 'u':
					i++
					if i+4 > p.len {
						p.error(i, "Invalid JSON")
					}
					val, err := strconv.ParseInt(p.input[i:i+4], 16, 32)
					if err != nil {
						p.error(i, "Invalid JSON")
					}
					if 0xD800 <= val && val <= 0xDBFF {
						i += 4
						if i+6 > p.len || p.input[i:i+2] != "\\u" {
							p.error(i, "Invalid JSON")
						}
						i += 2
						val2, err := strconv.ParseInt(p.input[i:i+4], 16, 32)
						if err != nil {
							p.error(i, "Invalid JSON")
						}
						if 0xDC00 <= val2 && val2 <= 0xDFFF {
							sb.WriteRune(utf16.DecodeRune(rune(val), rune(val2)))
							i += 4
						} else {
							p.error(i, "Invalid JSON")
						}
					} else if 0xDC00 <= val && val <= 0xDFFF {
						p.error(i, "Invalid JSON")
					} else {
						sb.WriteRune(rune(val))
						i += 4
					}
				default:
					p.error(i, "Invalid JSON")
				}
			} else if p.input[i] == '"' {
				break
			} else {
				val, size := utf8.DecodeRuneInString(p.input[i:])
				if val == utf8.RuneError && size == 1 {
					p.error(i, "Invalid JSON")
				}
				if val < 0x20 {
					p.error(i, "Invalid JSON")
				}
				sb.WriteRune(val)
				i += size
			}
		}
	}
	return sb.String(), i + 1
}

func (p *parser) parseNumber(i int) (any, int) {
	if p.checkOOB(i) {
		p.error(i, "Invalid JSON")
	}
	start := i
	isInt := true
	if p.input[i] == '-' {
		i++
	}
	if p.checkOOB(i) {
		p.error(i, "Invalid JSON")
	}
	if p.input[i] == '0' {
		i++
	} else {
		if p.input[i] >= '1' && p.input[i] <= '9' {
			i++
			if !p.checkOOB(i) {
				for i < p.len && p.input[i] >= '0' && p.input[i] <= '9' {
					i++
				}
			}
		} else {
			p.error(i, "Invalid JSON")
		}
	}
	if !p.checkOOB(i) && p.input[i] == '.' {
		isInt = false
		i++
		if p.checkOOB(i) {
			p.error(i, "Invalid JSON")
		}
		if p.input[i] < '0' || p.input[i] > '9' {
			p.error(i, "Invalid JSON")
		}
		for i < p.len && p.input[i] >= '0' && p.input[i] <= '9' {
			i++
		}
	}
	if !p.checkOOB(i) && (p.input[i] == 'e' || p.input[i] == 'E') {
		isInt = false
		i++
		if p.checkOOB(i) {
			p.error(i, "Invalid JSON")
		}
		if p.input[i] == '+' || p.input[i] == '-' {
			i++
		}
		if p.checkOOB(i) {
			p.error(i, "Invalid JSON")
		}
		if p.input[i] < '0' || p.input[i] > '9' {
			p.error(i, "Invalid JSON")
		}
		for i < p.len && p.input[i] >= '0' && p.input[i] <= '9' {
			i++
		}
	}

	temp := p.input[start:i]

	if isInt {
		number, err := strconv.ParseInt(temp, 10, 64)
		if err != nil {
			if numErr, ok := err.(*strconv.NumError); ok && numErr.Err == strconv.ErrRange {
				fNumber, fErr := strconv.ParseFloat(temp, 64)
				if fErr != nil {
					p.error(i, "Invalid number format")
				}
				return fNumber, i
			}
			p.error(i, "Invalid JSON")
		}
		return number, i
	} else {
		number, err := strconv.ParseFloat(temp, 64)
		if err != nil {
			p.error(i, "Invalid JSON")
		}
		return number, i
	}
}
