package goja

import (
	"errors"
	"gitee.com/quant1x/pkg/goja/unistring"
	"io"
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"
)

const hexUpper = "0123456789ABCDEF"

var (
	parseFloatRegexp = regexp.MustCompile(`^([+-]?(?:Infinity|[0-9]*\.?[0-9]*(?:[eE][+-]?[0-9]+)?))`)
)

func (r *Runtime) builtin_isNaN(call FunctionCall) Value {
	if math.IsNaN(call.Argument(0).ToFloat()) {
		return valueTrue
	} else {
		return valueFalse
	}
}

func (r *Runtime) builtin_parseInt(call FunctionCall) Value {
	str := call.Argument(0).toString().toTrimmedUTF8()
	radix := int(toInt32(call.Argument(1)))
	v, _ := parseInt(str, radix)
	return v
}

func (r *Runtime) builtin_parseFloat(call FunctionCall) Value {
	m := parseFloatRegexp.FindStringSubmatch(call.Argument(0).toString().toTrimmedUTF8())
	if len(m) == 2 {
		if s := m[1]; s != "" && s != "+" && s != "-" {
			switch s {
			case "+", "-":
			case "Infinity", "+Infinity":
				return _positiveInf
			case "-Infinity":
				return _negativeInf
			default:
				f, err := strconv.ParseFloat(s, 64)
				if err == nil || isRangeErr(err) {
					return floatToValue(f)
				}
			}
		}
	}
	return _NaN
}

func (r *Runtime) builtin_isFinite(call FunctionCall) Value {
	f := call.Argument(0).ToFloat()
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return valueFalse
	}
	return valueTrue
}

func (r *Runtime) _encode(uriString String, unescaped *[256]bool) String {
	reader := uriString.Reader()
	utf8Buf := make([]byte, utf8.UTFMax)
	needed := false
	l := 0
	for {
		rn, _, err := reader.ReadRune()
		if err != nil {
			if err != io.EOF {
				panic(r.newError(r.getURIError(), "Malformed URI"))
			}
			break
		}

		if rn >= utf8.RuneSelf {
			needed = true
			l += utf8.EncodeRune(utf8Buf, rn) * 3
		} else if !unescaped[rn] {
			needed = true
			l += 3
		} else {
			l++
		}
	}

	if !needed {
		return uriString
	}

	buf := make([]byte, l)
	i := 0
	reader = uriString.Reader()
	for {
		rn, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}

		if rn >= utf8.RuneSelf {
			n := utf8.EncodeRune(utf8Buf, rn)
			for _, b := range utf8Buf[:n] {
				buf[i] = '%'
				buf[i+1] = hexUpper[b>>4]
				buf[i+2] = hexUpper[b&15]
				i += 3
			}
		} else if !unescaped[rn] {
			buf[i] = '%'
			buf[i+1] = hexUpper[rn>>4]
			buf[i+2] = hexUpper[rn&15]
			i += 3
		} else {
			buf[i] = byte(rn)
			i++
		}
	}
	return asciiString(buf)
}

func (r *Runtime) _decode(sv String, reservedSet *[256]bool) String {
	s := sv.String()
	hexCount := 0
	for i := 0; i < len(s); {
		switch s[i] {
		case '%':
			if i+2 >= len(s) || !ishex(s[i+1]) || !ishex(s[i+2]) {
				panic(r.newError(r.getURIError(), "Malformed URI"))
			}
			c := unhex(s[i+1])<<4 | unhex(s[i+2])
			if !reservedSet[c] {
				hexCount++
			}
			i += 3
		default:
			i++
		}
	}

	if hexCount == 0 {
		return sv
	}

	t := make([]byte, len(s)-hexCount*2)
	j := 0
	isUnicode := false
	for i := 0; i < len(s); {
		ch := s[i]
		switch ch {
		case '%':
			c := unhex(s[i+1])<<4 | unhex(s[i+2])
			if reservedSet[c] {
				t[j] = s[i]
				t[j+1] = s[i+1]
				t[j+2] = s[i+2]
				j += 3
			} else {
				t[j] = c
				if c >= utf8.RuneSelf {
					isUnicode = true
				}
				j++
			}
			i += 3
		default:
			if ch >= utf8.RuneSelf {
				isUnicode = true
			}
			t[j] = ch
			j++
			i++
		}
	}

	if !isUnicode {
		return asciiString(t)
	}

	us := make([]rune, 0, len(s))
	for len(t) > 0 {
		rn, size := utf8.DecodeRune(t)
		if rn == utf8.RuneError {
			if size != 3 || t[0] != 0xef || t[1] != 0xbf || t[2] != 0xbd {
				panic(r.newError(r.getURIError(), "Malformed URI"))
			}
		}
		us = append(us, rn)
		t = t[size:]
	}
	return unicodeStringFromRunes(us)
}

func ishex(c byte) bool {
	switch {
	case '0' <= c && c <= '9':
		return true
	case 'a' <= c && c <= 'f':
		return true
	case 'A' <= c && c <= 'F':
		return true
	}
	return false
}

func unhex(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	}
	return 0
}

func (r *Runtime) builtin_decodeURI(call FunctionCall) Value {
	uriString := call.Argument(0).toString()
	return r._decode(uriString, &uriReservedHash)
}

func (r *Runtime) builtin_decodeURIComponent(call FunctionCall) Value {
	uriString := call.Argument(0).toString()
	return r._decode(uriString, &emptyEscapeSet)
}

func (r *Runtime) builtin_encodeURI(call FunctionCall) Value {
	uriString := call.Argument(0).toString()
	return r._encode(uriString, &uriReservedUnescapedHash)
}

func (r *Runtime) builtin_encodeURIComponent(call FunctionCall) Value {
	uriString := call.Argument(0).toString()
	return r._encode(uriString, &uriUnescaped)
}

func (r *Runtime) builtin_escape(call FunctionCall) Value {
	s := call.Argument(0).toString()
	var sb strings.Builder
	l := s.Length()
	for i := 0; i < l; i++ {
		r := s.CharAt(i)
		if r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z' || r >= '0' && r <= '9' ||
			r == '@' || r == '*' || r == '_' || r == '+' || r == '-' || r == '.' || r == '/' {
			sb.WriteByte(byte(r))
		} else if r <= 0xff {
			sb.WriteByte('%')
			sb.WriteByte(hexUpper[r>>4])
			sb.WriteByte(hexUpper[r&0xf])
		} else {
			sb.WriteString("%u")
			sb.WriteByte(hexUpper[r>>12])
			sb.WriteByte(hexUpper[(r>>8)&0xf])
			sb.WriteByte(hexUpper[(r>>4)&0xf])
			sb.WriteByte(hexUpper[r&0xf])
		}
	}
	return asciiString(sb.String())
}

func (r *Runtime) builtin_unescape(call FunctionCall) Value {
	s := call.Argument(0).toString()
	l := s.Length()
	var asciiBuf []byte
	var unicodeBuf []uint16
	_, u := devirtualizeString(s)
	unicode := u != nil
	if unicode {
		unicodeBuf = make([]uint16, 1, l+1)
		unicodeBuf[0] = unistring.BOM
	} else {
		asciiBuf = make([]byte, 0, l)
	}
	for i := 0; i < l; {
		r := s.CharAt(i)
		if r == '%' {
			if i <= l-6 && s.CharAt(i+1) == 'u' {
				c0 := s.CharAt(i + 2)
				c1 := s.CharAt(i + 3)
				c2 := s.CharAt(i + 4)
				c3 := s.CharAt(i + 5)
				if c0 <= 0xff && ishex(byte(c0)) &&
					c1 <= 0xff && ishex(byte(c1)) &&
					c2 <= 0xff && ishex(byte(c2)) &&
					c3 <= 0xff && ishex(byte(c3)) {
					r = uint16(unhex(byte(c0)))<<12 |
						uint16(unhex(byte(c1)))<<8 |
						uint16(unhex(byte(c2)))<<4 |
						uint16(unhex(byte(c3)))
					i += 5
					goto out
				}
			}
			if i <= l-3 {
				c0 := s.CharAt(i + 1)
				c1 := s.CharAt(i + 2)
				if c0 <= 0xff && ishex(byte(c0)) &&
					c1 <= 0xff && ishex(byte(c1)) {
					r = uint16(unhex(byte(c0))<<4 | unhex(byte(c1)))
					i += 2
				}
			}
		}
	out:
		if r >= utf8.RuneSelf && !unicode {
			unicodeBuf = make([]uint16, 1, l+1)
			unicodeBuf[0] = unistring.BOM
			for _, b := range asciiBuf {
				unicodeBuf = append(unicodeBuf, uint16(b))
			}
			asciiBuf = nil
			unicode = true
		}
		if unicode {
			unicodeBuf = append(unicodeBuf, r)
		} else {
			asciiBuf = append(asciiBuf, byte(r))
		}
		i++
	}
	if unicode {
		return unicodeString(unicodeBuf)
	}

	return asciiString(asciiBuf)
}

func createGlobalObjectTemplate() *objectTemplate {
	t := newObjectTemplate()
	t.protoFactory = func(r *Runtime) *Object {
		return r.global.ObjectPrototype
	}

	t.putStr("Object", func(r *Runtime) Value { return valueProp(r.getObject(), true, false, true) })
	t.putStr("Function", func(r *Runtime) Value { return valueProp(r.getFunction(), true, false, true) })
	t.putStr("Array", func(r *Runtime) Value { return valueProp(r.getArray(), true, false, true) })
	t.putStr("String", func(r *Runtime) Value { return valueProp(r.getString(), true, false, true) })
	t.putStr("Number", func(r *Runtime) Value { return valueProp(r.getNumber(), true, false, true) })
	t.putStr("RegExp", func(r *Runtime) Value { return valueProp(r.getRegExp(), true, false, true) })
	t.putStr("Date", func(r *Runtime) Value { return valueProp(r.getDate(), true, false, true) })
	t.putStr("Boolean", func(r *Runtime) Value { return valueProp(r.getBoolean(), true, false, true) })
	t.putStr("Proxy", func(r *Runtime) Value { return valueProp(r.getProxy(), true, false, true) })
	t.putStr("Reflect", func(r *Runtime) Value { return valueProp(r.getReflect(), true, false, true) })
	t.putStr("Error", func(r *Runtime) Value { return valueProp(r.getError(), true, false, true) })
	t.putStr("AggregateError", func(r *Runtime) Value { return valueProp(r.getAggregateError(), true, false, true) })
	t.putStr("TypeError", func(r *Runtime) Value { return valueProp(r.getTypeError(), true, false, true) })
	t.putStr("ReferenceError", func(r *Runtime) Value { return valueProp(r.getReferenceError(), true, false, true) })
	t.putStr("SyntaxError", func(r *Runtime) Value { return valueProp(r.getSyntaxError(), true, false, true) })
	t.putStr("RangeError", func(r *Runtime) Value { return valueProp(r.getRangeError(), true, false, true) })
	t.putStr("EvalError", func(r *Runtime) Value { return valueProp(r.getEvalError(), true, false, true) })
	t.putStr("URIError", func(r *Runtime) Value { return valueProp(r.getURIError(), true, false, true) })
	t.putStr("GoError", func(r *Runtime) Value { return valueProp(r.getGoError(), true, false, true) })

	t.putStr("eval", func(r *Runtime) Value { return valueProp(r.getEval(), true, false, true) })

	t.putStr("Math", func(r *Runtime) Value { return valueProp(r.getMath(), true, false, true) })
	t.putStr("JSON", func(r *Runtime) Value { return valueProp(r.getJSON(), true, false, true) })
	addTypedArrays(t)
	t.putStr("Symbol", func(r *Runtime) Value { return valueProp(r.getSymbol(), true, false, true) })
	t.putStr("WeakSet", func(r *Runtime) Value { return valueProp(r.getWeakSet(), true, false, true) })
	t.putStr("WeakMap", func(r *Runtime) Value { return valueProp(r.getWeakMap(), true, false, true) })
	t.putStr("Map", func(r *Runtime) Value { return valueProp(r.getMap(), true, false, true) })
	t.putStr("Set", func(r *Runtime) Value { return valueProp(r.getSet(), true, false, true) })
	t.putStr("Promise", func(r *Runtime) Value { return valueProp(r.getPromise(), true, false, true) })

	t.putStr("globalThis", func(r *Runtime) Value { return valueProp(r.globalObject, true, false, true) })
	t.putStr("NaN", func(r *Runtime) Value { return valueProp(_NaN, false, false, false) })
	t.putStr("undefined", func(r *Runtime) Value { return valueProp(_undefined, false, false, false) })
	t.putStr("Infinity", func(r *Runtime) Value { return valueProp(_positiveInf, false, false, false) })

	t.putStr("isNaN", func(r *Runtime) Value { return r.methodProp(r.builtin_isNaN, "isNaN", 1) })
	t.putStr("parseInt", func(r *Runtime) Value { return valueProp(r.getParseInt(), true, false, true) })
	t.putStr("parseFloat", func(r *Runtime) Value { return valueProp(r.getParseFloat(), true, false, true) })
	t.putStr("isFinite", func(r *Runtime) Value { return r.methodProp(r.builtin_isFinite, "isFinite", 1) })
	t.putStr("decodeURI", func(r *Runtime) Value { return r.methodProp(r.builtin_decodeURI, "decodeURI", 1) })
	t.putStr("decodeURIComponent", func(r *Runtime) Value { return r.methodProp(r.builtin_decodeURIComponent, "decodeURIComponent", 1) })
	t.putStr("encodeURI", func(r *Runtime) Value { return r.methodProp(r.builtin_encodeURI, "encodeURI", 1) })
	t.putStr("encodeURIComponent", func(r *Runtime) Value { return r.methodProp(r.builtin_encodeURIComponent, "encodeURIComponent", 1) })
	t.putStr("escape", func(r *Runtime) Value { return r.methodProp(r.builtin_escape, "escape", 1) })
	t.putStr("unescape", func(r *Runtime) Value { return r.methodProp(r.builtin_unescape, "unescape", 1) })

	// TODO: Annex B

	t.putSym(SymToStringTag, func(r *Runtime) Value { return valueProp(asciiString(classGlobal), false, false, true) })

	return t
}

var globalObjectTemplate *objectTemplate
var globalObjectTemplateOnce sync.Once

func getGlobalObjectTemplate() *objectTemplate {
	globalObjectTemplateOnce.Do(func() {
		globalObjectTemplate = createGlobalObjectTemplate()
	})
	return globalObjectTemplate
}

func (r *Runtime) getEval() *Object {
	ret := r.global.Eval
	if ret == nil {
		ret = r.newNativeFunc(r.builtin_eval, "eval", 1)
		r.global.Eval = ret
	}
	return ret
}

func digitVal(d byte) int {
	var v byte
	switch {
	case '0' <= d && d <= '9':
		v = d - '0'
	case 'a' <= d && d <= 'z':
		v = d - 'a' + 10
	case 'A' <= d && d <= 'Z':
		v = d - 'A' + 10
	default:
		return 36
	}
	return int(v)
}

// ECMAScript compatible version of strconv.ParseInt
func parseInt(s string, base int) (Value, error) {
	var n int64
	var err error
	var cutoff, maxVal int64
	var sign bool
	i := 0

	if len(s) < 1 {
		err = strconv.ErrSyntax
		goto Error
	}

	switch s[0] {
	case '-':
		sign = true
		s = s[1:]
	case '+':
		s = s[1:]
	}

	if len(s) < 1 {
		err = strconv.ErrSyntax
		goto Error
	}

	// Look for hex prefix.
	if s[0] == '0' && len(s) > 1 && (s[1] == 'x' || s[1] == 'X') {
		if base == 0 || base == 16 {
			base = 16
			s = s[2:]
		}
	}

	switch {
	case len(s) < 1:
		err = strconv.ErrSyntax
		goto Error

	case 2 <= base && base <= 36:
	// valid base; nothing to do

	case base == 0:
		// Look for hex prefix.
		switch {
		case s[0] == '0' && len(s) > 1 && (s[1] == 'x' || s[1] == 'X'):
			if len(s) < 3 {
				err = strconv.ErrSyntax
				goto Error
			}
			base = 16
			s = s[2:]
		default:
			base = 10
		}

	default:
		err = errors.New("invalid base " + strconv.Itoa(base))
		goto Error
	}

	// Cutoff is the smallest number such that cutoff*base > maxInt64.
	// Use compile-time constants for common cases.
	switch base {
	case 10:
		cutoff = math.MaxInt64/10 + 1
	case 16:
		cutoff = math.MaxInt64/16 + 1
	default:
		cutoff = math.MaxInt64/int64(base) + 1
	}

	maxVal = math.MaxInt64
	for ; i < len(s); i++ {
		if n >= cutoff {
			// n*base overflows
			return parseLargeInt(float64(n), s[i:], base, sign)
		}
		v := digitVal(s[i])
		if v >= base {
			break
		}
		n *= int64(base)

		n1 := n + int64(v)
		if n1 < n || n1 > maxVal {
			// n+v overflows
			return parseLargeInt(float64(n)+float64(v), s[i+1:], base, sign)
		}
		n = n1
	}

	if i == 0 {
		err = strconv.ErrSyntax
		goto Error
	}

	if sign {
		n = -n
	}
	return intToValue(n), nil

Error:
	return _NaN, err
}

func parseLargeInt(n float64, s string, base int, sign bool) (Value, error) {
	i := 0
	b := float64(base)
	for ; i < len(s); i++ {
		v := digitVal(s[i])
		if v >= base {
			break
		}
		n = n*b + float64(v)
	}
	if sign {
		n = -n
	}
	// We know it can't be represented as int, so use valueFloat instead of floatToValue
	return valueFloat(n), nil
}

var (
	uriUnescaped             [256]bool
	uriReserved              [256]bool
	uriReservedHash          [256]bool
	uriReservedUnescapedHash [256]bool
	emptyEscapeSet           [256]bool
)

func init() {
	for _, c := range "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_.!~*'()" {
		uriUnescaped[c] = true
	}

	for _, c := range ";/?:@&=+$," {
		uriReserved[c] = true
	}

	for i := 0; i < 256; i++ {
		if uriUnescaped[i] || uriReserved[i] {
			uriReservedUnescapedHash[i] = true
		}
		uriReservedHash[i] = uriReserved[i]
	}
	uriReservedUnescapedHash['#'] = true
	uriReservedHash['#'] = true
}
