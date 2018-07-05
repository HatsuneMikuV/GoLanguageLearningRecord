package sexpr

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"text/scanner"
)

/*
	encode
*/
// Marshal encodes a Go value in S-expression form.
func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := encode(&buf, reflect.ValueOf(v)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func encode(buf *bytes.Buffer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("nil")
	/*
		练习 12.3： 实现encode函数缺少的分支。
		将布尔类型编码为t和nil，
		浮点数编码为Go语言的格式，
		复数1+2i编码为#C(1.0 2.0)格式。
		接口编码为类型名和值对，
		例如("[]int" (1 2 3))，
		但是这个形式可能会造成歧义：
		reflect.Type.String方法对于不同的类型可能返回相同的结果。
	*/
	case reflect.Bool:
		if v.Bool() == true {
			buf.WriteString("t")
		}else  {
			buf.WriteString("nil")
		}
	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(buf, "%f", v.Float())
	case reflect.Complex64, reflect.Complex128:
		fmt.Fprintf(buf, "#C(%.1f %.1f)", real(v.Complex()), imag(v.Complex()))
	case reflect.Interface:
		fmt.Fprintf(buf, "(\"%T\"(%d))", v.Interface(),v.Interface())

	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%d", v.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fmt.Fprintf(buf, "%d", v.Uint())

	case reflect.String:
		fmt.Fprintf(buf, "%q", v.String())

	case reflect.Ptr:
		return encode(buf, v.Elem())

	case reflect.Array, reflect.Slice: // (value ...)
		buf.WriteByte('(')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			if err := encode(buf, v.Index(i)); err != nil {
				return err
			}
		}
		buf.WriteByte(')')

	case reflect.Struct: // ((name value) ...)
		buf.WriteByte('(')
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			fmt.Fprintf(buf, "(%s ", v.Type().Field(i).Name)
			if err := encode(buf, v.Field(i)); err != nil {
				return err
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')

	case reflect.Map: // ((key value) ...)
		buf.WriteByte('(')
		for i, key := range v.MapKeys() {
			if i > 0 {
				buf.WriteByte(' ')
			}
			buf.WriteByte('(')
			if err := encode(buf, key); err != nil {
				return err
			}
			buf.WriteByte(' ')
			if err := encode(buf, v.MapIndex(key)); err != nil {
				return err
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')

	default: //chan, func
		return fmt.Errorf("unsupported type: %s", v.Type())
	}
	return nil
}

/*
	decode
*/
//!+Unmarshal
// Unmarshal parses S-expression data and populates the variable
// whose address is in the non-nil pointer out.
func Unmarshal(data []byte, out interface{}) (err error) {
	lex := &lexer{scan: scanner.Scanner{Mode: scanner.GoTokens}}
	lex.scan.Init(bytes.NewReader(data))
	lex.next() // get the first token
	defer func() {
		// NOTE: this is not an example of ideal error handling.
		if x := recover(); x != nil {
			err = fmt.Errorf("error at %s: %v", lex.scan.Position, x)
		}
	}()
	read(lex, reflect.ValueOf(out).Elem())
	return nil
}

//!-Unmarshal

//!+lexer
type lexer struct {
	scan  scanner.Scanner
	token rune // the current token
}

func (lex *lexer) next()        { lex.token = lex.scan.Scan() }
func (lex *lexer) text() string { return lex.scan.TokenText() }

func (lex *lexer) consume(want rune) {
	if lex.token != want { // NOTE: Not an example of good error handling.
		panic(fmt.Sprintf("got %q, want %q", lex.text(), want))
	}
	lex.next()
}

//!-lexer

// The read function is a decoder for a small subset of well-formed
// S-expressions.  For brevity of our example, it takes many dubious
// shortcuts.
//
// The parser assumes
// - that the S-expression input is well-formed; it does no error checking.
// - that the S-expression input corresponds to the type of the variable.
// - that all numbers in the input are non-negative decimal integers.
// - that all keys in ((key value) ...) struct syntax are unquoted symbols.
// - that the input does not contain dotted lists such as (1 2 . 3).
// - that the input does not contain Lisp reader macros such 'x and #'x.
//
// The reflection logic assumes
// - that v is always a variable of the appropriate type for the
//   S-expression value.  For example, v must not be a boolean,
//   interface, channel, or function, and if v is an array, the input
//   must have the correct number of elements.
// - that v in the top-level call to read has the zero value of its
//   type and doesn't need clearing.
// - that if v is a numeric variable, it is a signed integer.

//!+read
func read(lex *lexer, v reflect.Value) {
	switch lex.token {
	case scanner.Float:
		i, _ := strconv.ParseFloat(lex.text(),64) // NOTE: ignoring errors
		v.SetFloat(float64(i))
		lex.next()
		return
	case scanner.Ident:
		// The only valid identifiers are
		// "nil" and struct field names.
		if lex.text() == "nil" {
			v.Set(reflect.Zero(v.Type()))
			lex.next()
			return
		}
	case scanner.String:
		s, _ := strconv.Unquote(lex.text()) // NOTE: ignoring errors
		v.SetString(s)
		lex.next()
		return
	case scanner.Int:
		i, _ := strconv.Atoi(lex.text()) // NOTE: ignoring errors
		v.SetInt(int64(i))
		lex.next()
		return
	case '(':
		lex.next()
		readList(lex, v)
		lex.next() // consume ')'
		return
	}
	panic(fmt.Sprintf("unexpected token %q", lex.text()))
}

//!-read

//!+readlist
func readList(lex *lexer, v reflect.Value) {
	switch v.Kind() {
	case reflect.Bool:
		i, _ := strconv.ParseBool(lex.text()) // NOTE: ignoring errors
		v.SetBool(bool(i))
		lex.next()

	case reflect.Interface:
		//item := reflect.New(v.Type())
		fmt.Print(v.Type())
		//read(lex, item)
		i := v.Interface()
		v.Set(reflect.ValueOf(i))
		lex.next()

	case reflect.Array: // (item ...)
		for i := 0; !endList(lex); i++ {
			read(lex, v.Index(i))
		}

	case reflect.Slice: // (item ...)
		for !endList(lex) {
			item := reflect.New(v.Type().Elem()).Elem()
			read(lex, item)
			v.Set(reflect.Append(v, item))
		}

	case reflect.Struct: // ((name value) ...)
		for !endList(lex) {
			lex.consume('(')
			if lex.token != scanner.Ident {
				panic(fmt.Sprintf("got token %q, want field name", lex.text()))
			}
			name := lex.text()
			lex.next()
			read(lex, v.FieldByName(name))
			lex.consume(')')
		}

	case reflect.Map: // ((key value) ...)
		v.Set(reflect.MakeMap(v.Type()))
		for !endList(lex) {
			lex.consume('(')
			key := reflect.New(v.Type().Key()).Elem()
			read(lex, key)
			value := reflect.New(v.Type().Elem()).Elem()
			read(lex, value)
			v.SetMapIndex(key, value)
			lex.consume(')')
		}

	default:
		panic(fmt.Sprintf("cannot decode list into %v", v.Type()))
	}
}

func endList(lex *lexer) bool {
	switch lex.token {
	case scanner.EOF:
		panic("end of file")
	case ')':
		return true
	}
	return false
}

/*
	pretty
*/

func MarshalIndent(v interface{}) ([]byte, error) {
	p := printer{width: margin}
	if err := pretty(&p, reflect.ValueOf(v)); err != nil {
		return nil, err
	}
	return p.Bytes(), nil
}

const margin = 80

type token struct {
	kind rune // one of "s ()" (string, blank, start, end)
	str  string
	size int
}

type printer struct {
	tokens []*token // FIFO buffer
	stack  []*token // stack of open ' ' and '(' tokens
	rtotal int      // total number of spaces needed to print stream

	bytes.Buffer
	indents []int
	width   int // remaining space
}

func (p *printer) string(str string) {
	tok := &token{kind: 's', str: str, size: len(str)}
	if len(p.stack) == 0 {
		p.print(tok)
	} else {
		p.tokens = append(p.tokens, tok)
		p.rtotal += len(str)
	}
}
func (p *printer) pop() (top *token) {
	last := len(p.stack) - 1
	top, p.stack = p.stack[last], p.stack[:last]
	return
}
func (p *printer) begin() {
	if len(p.stack) == 0 {
		p.rtotal = 1
	}
	t := &token{kind: '(', size: -p.rtotal}
	p.tokens = append(p.tokens, t)
	p.stack = append(p.stack, t) // push
	p.string("(")
}
func (p *printer) end() {
	p.string(")")
	p.tokens = append(p.tokens, &token{kind: ')'})
	x := p.pop()
	x.size += p.rtotal
	if x.kind == ' ' {
		p.pop().size += p.rtotal
	}
	if len(p.stack) == 0 {
		for _, tok := range p.tokens {
			p.print(tok)
		}
		p.tokens = nil
	}
}
func (p *printer) space() {
	last := len(p.stack) - 1
	x := p.stack[last]
	if x.kind == ' ' {
		x.size += p.rtotal
		p.stack = p.stack[:last] // pop
	}
	t := &token{kind: ' ', size: -p.rtotal}
	p.tokens = append(p.tokens, t)
	p.stack = append(p.stack, t)
	p.rtotal++
}
func (p *printer) print(t *token) {
	switch t.kind {
	case 's':
		p.WriteString(t.str)
		p.width -= len(t.str)
	case '(':
		p.indents = append(p.indents, p.width)
	case ')':
		p.indents = p.indents[:len(p.indents)-1] // pop
	case ' ':
		if t.size > p.width {
			p.width = p.indents[len(p.indents)-1] - 1
			fmt.Fprintf(&p.Buffer, "\n%*s", margin-p.width, "")
		} else {
			p.WriteByte(' ')
			p.width--
		}
	}
}
func (p *printer) stringf(format string, args ...interface{}) {
	p.string(fmt.Sprintf(format, args...))
}

func pretty(p *printer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Invalid:
		p.string("nil")

		/*
			练习 12.4： 修改encode函数，以上面的格式化形式输出S表达式
		*/
	case reflect.Bool:
		if v.Bool() == true {
			p.string("t")
		}else  {
			p.string("nil")
		}
	case reflect.Float32, reflect.Float64:
		p.stringf("%f", v.Float())
	case reflect.Complex64, reflect.Complex128:
		p.string("#C")
		p.begin()
		p.stringf("%.1f", real(v.Complex()))
		p.space()
		p.stringf("%.1f", imag(v.Complex()))
		p.end()
	case reflect.Interface:
		p.stringf("\"%T\"", v.Interface())
		p.begin()
		p.stringf( "%d", v.Interface())
		p.end()

	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		p.stringf("%d", v.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p.stringf("%d", v.Uint())

	case reflect.String:
		p.stringf("%q", v.String())

	case reflect.Array, reflect.Slice: // (value ...)
		p.begin()
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				p.space()
			}
			if err := pretty(p, v.Index(i)); err != nil {
				return err
			}
		}
		p.end()

	case reflect.Struct: // ((name value ...)
		p.begin()
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				p.space()
			}
			p.begin()
			p.string(v.Type().Field(i).Name)
			p.space()
			if err := pretty(p, v.Field(i)); err != nil {
				return err
			}
			p.end()
		}
		p.end()

	case reflect.Map: // ((key value ...)
		p.begin()
		for i, key := range v.MapKeys() {
			if i > 0 {
				p.space()
			}
			p.begin()
			if err := pretty(p, key); err != nil {
				return err
			}
			p.space()
			if err := pretty(p, v.MapIndex(key)); err != nil {
				return err
			}
			p.end()
		}
		p.end()

	case reflect.Ptr:
		return pretty(p, v.Elem())

	default: // chan, func
		return fmt.Errorf("unsupported type: %s", v.Type())
	}
	return nil
}


/*
	练习 12.5： 修改encode函数，用JSON格式代替S表达式格式。
	然后使用标准库提供的json.Unmarshal解码器来验证函数是正确的。

	基本上是一致，结果是正确的
*/

/*
	encode
*/
// Marshal encodes a Go value in S-expression form.
func Marshal_Json(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := encode_json(&buf, reflect.ValueOf(v)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func encode_json(buf *bytes.Buffer, v reflect.Value) error {
	switch v.Kind() {

	case reflect.Invalid:
		buf.WriteString("nil")

	case reflect.Bool:
		if v.Bool() == true {
			buf.WriteString("true")
		}else  {
			buf.WriteString("false")
		}

	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(buf, "%f", v.Float())

	case reflect.Complex64, reflect.Complex128:
		fmt.Fprintf(buf, "#C(%.1f %.1f)", real(v.Complex()), imag(v.Complex()))

	case reflect.Interface:
		fmt.Fprintf(buf, "%d", v.Interface())

	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%d", v.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fmt.Fprintf(buf, "%d", v.Uint())

	case reflect.String:
		fmt.Fprintf(buf, "%q", v.String())

	case reflect.Ptr:
		return encode_json(buf, v.Elem())

	case reflect.Array, reflect.Slice: // (value ...)
		buf.WriteByte('{')
		for i := 0; i < v.Len(); i++ {
			if err := encode_json(buf, v.Index(i)); err != nil {
				return err
			}
		}
		buf.WriteByte('}')

	case reflect.Struct: // ((name value) ...)
		buf.WriteByte('{')
		for i := 0; i < v.NumField(); i++ {
			fmt.Fprintf(buf, "\"%s\":", v.Type().Field(i).Name)
			if err := encode_json(buf, v.Field(i)); err != nil {
				return err
			}
			buf.WriteByte(',')
		}
		buf.WriteByte('}')

	case reflect.Map: // ((key value) ...)
		buf.WriteByte('{')
		for _, key := range v.MapKeys() {
			if err := encode_json(buf, key); err != nil {
				return err
			}
			buf.WriteByte(':')
			if err := encode_json(buf, v.MapIndex(key)); err != nil {
				return err
			}
		}
		buf.WriteByte('}')

	default: //chan, func
		return fmt.Errorf("unsupported type: %s", v.Type())
	}
	return nil
}
