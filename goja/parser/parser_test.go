package parser

import (
	"errors"
	"regexp"
	"strings"
	"testing"

	"gitee.com/quant1x/pkg/goja/ast"
	"gitee.com/quant1x/pkg/goja/file"
	"gitee.com/quant1x/pkg/goja/token"
	"gitee.com/quant1x/pkg/goja/unistring"
)

func firstErr(err error) error {
	switch err := err.(type) {
	case ErrorList:
		return err[0]
	}
	return err
}

var matchBeforeAfterSeparator = regexp.MustCompile(`(?m)^[ \t]*---$`)

func testParse(src string) (parser *_parser, program *ast.Program, err error) {
	defer func() {
		if tmp := recover(); tmp != nil {
			switch tmp := tmp.(type) {
			case string:
				if strings.HasPrefix(tmp, "SyntaxError:") {
					parser = nil
					program = nil
					err = errors.New(tmp)
					return
				}
			}
			panic(tmp)
		}
	}()
	parser = newParser("", src)
	program, err = parser.parse()
	return
}

func TestParseFile(t *testing.T) {
	tt(t, func() {
		_, err := ParseFile(nil, "", `/abc/`, 0)
		is(err, nil)

		_, err = ParseFile(nil, "", `/(?!def)abc/`, IgnoreRegExpErrors)
		is(err, nil)

		_, err = ParseFile(nil, "", `/(?!def)abc/; return`, IgnoreRegExpErrors)
		is(err, "(anonymous): Line 1:15 Illegal return statement")
	})
}

func TestParseFunction(t *testing.T) {
	tt(t, func() {
		test := func(prm, bdy string, expect interface{}) *ast.FunctionLiteral {
			function, err := ParseFunction(prm, bdy)
			is(firstErr(err), expect)
			return function
		}

		test("a, b,c,d", "", nil)

		test("a, b;,c,d", "", "(anonymous): Line 1:15 Unexpected token ;")

		test("this", "", "(anonymous): Line 1:11 Unexpected token this")

		test("a, b, c, null", "", "(anonymous): Line 1:20 Unexpected token null")

		test("a, b,c,d", "return;", nil)

		test("a, b,c,d", "break;", "(anonymous): Line 2:1 Illegal break statement")

		test("a, b,c,d", "{}", nil)
	})
}

func TestParserErr(t *testing.T) {
	tt(t, func() {
		test := func(input string, expect interface{}) (*ast.Program, *_parser) {
			parser := newParser("", input)
			program, err := parser.parse()
			is(firstErr(err), expect)
			return program, parser
		}

		test("", nil)

		program, parser := test(`
        var abc;
        break; do {
        } while(true);
    `, "(anonymous): Line 3:9 Illegal break statement")
		{
			stmt := program.Body[1].(*ast.BadStatement)
			is(parser.position(stmt.From).Column, 9)
			is(parser.position(stmt.To).Column, 16)
			is(parser.slice(stmt.From, stmt.To), "break; ")
		}

		s := string([]byte{0x22, 0x25, 0x21, 0x63, 0x28, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x3d, 0x25, 0x63, 0x25, 0x9c, 0x29, 0x25, 0x21, 0x5c, 0x28, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x3d, 0x5c, 0xe2, 0x80, 0xa9, 0x29, 0x78, 0x39, 0x63, 0x22})
		test(s, `(anonymous): Line 1:16 Invalid UTF-8 character`)

		test("{", "(anonymous): Line 1:2 Unexpected end of input")

		test("}", "(anonymous): Line 1:1 Unexpected token }")

		test("3ea", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test("3in", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test("3in []", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test("3e", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test("3e+", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test("3e-", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test("3x", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test("3x0", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test("0x", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test("09", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test("018", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test("01.0", "(anonymous): Line 1:3 Unexpected number")

		test(".0.9", "(anonymous): Line 1:3 Unexpected number")

		test("0o3e1", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test("01a", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test("0x3in[]", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test("\"Hello\nWorld\"", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test("\u203f = 10", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test("x\\", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test("x\\\\", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test("x\\u005c", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test("x\\u002a", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test("x\\\\u002a", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test("/\n", "(anonymous): Line 1:1 Invalid regular expression: missing /")

		test("0 = 1", "(anonymous): Line 1:1 Invalid left-hand side in assignment")

		test("func() = 1", "(anonymous): Line 1:1 Invalid left-hand side in assignment")

		test("(1 + 1) = 2", "(anonymous): Line 1:2 Invalid left-hand side in assignment")

		test("1++", "(anonymous): Line 1:2 Invalid left-hand side in assignment")

		test("1--", "(anonymous): Line 1:2 Invalid left-hand side in assignment")

		test("--1", "(anonymous): Line 1:1 Invalid left-hand side in assignment")

		test("for((1 + 1) in abc) def();", "(anonymous): Line 1:1 Invalid left-hand side in for-in or for-of")

		test("[", "(anonymous): Line 1:2 Unexpected end of input")

		test("[,", "(anonymous): Line 1:3 Unexpected end of input")

		test("1 + {", "(anonymous): Line 1:6 Unexpected end of input")

		test("1 + { abc:abc", "(anonymous): Line 1:14 Unexpected end of input")

		test("1 + { abc:abc,", "(anonymous): Line 1:15 Unexpected end of input")

		test("var abc = /\n/", "(anonymous): Line 1:11 Invalid regular expression: missing /")

		test("var abc = \"\n", "(anonymous): Line 1:11 Unexpected token ILLEGAL")

		test("var if = 0", "(anonymous): Line 1:5 Unexpected token if")

		test("abc + 0 = 1", "(anonymous): Line 1:1 Invalid left-hand side in assignment")

		test("+abc = 1", "(anonymous): Line 1:1 Invalid left-hand side in assignment")

		test("1 + (", "(anonymous): Line 1:6 Unexpected end of input")

		test("\n\n\n{", "(anonymous): Line 4:2 Unexpected end of input")

		test("\n/* Some multiline\ncomment */\n)", "(anonymous): Line 4:1 Unexpected token )")

		test("+1 ** 2", "(anonymous): Line 1:4 Unary operator used immediately before exponentiation expression. Parenthesis must be used to disambiguate operator precedence")
		test("typeof 1 ** 2", "(anonymous): Line 1:10 Unary operator used immediately before exponentiation expression. Parenthesis must be used to disambiguate operator precedence")

		// TODO
		//{ set 1 }
		//{ get 2 }
		//({ set: s(if) { } })
		//({ set s(.) { } })
		//({ set: s() { } })
		//({ set: s(a, b) { } })
		//({ get: g(d) { } })
		//({ get i() { }, i: 42 })
		//({ i: 42, get i() { } })
		//({ set i(x) { }, i: 42 })
		//({ i: 42, set i(x) { } })
		//({ get i() { }, get i() { } })
		//({ set i(x) { }, set i(x) { } })

		test("function abc(if) {}", "(anonymous): Line 1:14 Unexpected token if")

		test("function abc(true) {}", "(anonymous): Line 1:14 Unexpected token true")

		test("function abc(false) {}", "(anonymous): Line 1:14 Unexpected token false")

		test("function abc(null) {}", "(anonymous): Line 1:14 Unexpected token null")

		test("function null() {}", "(anonymous): Line 1:10 Unexpected token null")

		test("function true() {}", "(anonymous): Line 1:10 Unexpected token true")

		test("function false() {}", "(anonymous): Line 1:10 Unexpected token false")

		test("function if() {}", "(anonymous): Line 1:10 Unexpected token if")

		test("a b;", "(anonymous): Line 1:3 Unexpected identifier")

		test("if.a", "(anonymous): Line 1:3 Unexpected token .")

		test("a if", "(anonymous): Line 1:3 Unexpected token if")

		test("a class", "(anonymous): Line 1:3 Unexpected token class")

		test("break\n", "(anonymous): Line 1:1 Illegal break statement")

		test("break 1;", "(anonymous): Line 1:7 Unexpected number")

		test("for (;;) { break 1; }", "(anonymous): Line 1:18 Unexpected number")

		test("continue\n", "(anonymous): Line 1:1 Illegal continue statement")

		test("continue 1;", "(anonymous): Line 1:10 Unexpected number")

		test("for (;;) { continue 1; }", "(anonymous): Line 1:21 Unexpected number")

		test("throw", "(anonymous): Line 1:1 Unexpected end of input")

		test("throw;", "(anonymous): Line 1:6 Unexpected token ;")

		test("throw \n", "(anonymous): Line 1:1 Unexpected end of input")

		test("for (var abc, def in {});", "(anonymous): Line 1:19 Unexpected token in")

		test("for ((abc in {});;);", nil)

		test("for ((abc in {}));", "(anonymous): Line 1:17 Unexpected token )")

		test("for (+abc in {});", "(anonymous): Line 1:1 Invalid left-hand side in for-in or for-of")

		test("if (false)", "(anonymous): Line 1:11 Unexpected end of input")

		test("if (false) abc(); else", "(anonymous): Line 1:23 Unexpected end of input")

		test("do", "(anonymous): Line 1:3 Unexpected end of input")

		test("while (false)", "(anonymous): Line 1:14 Unexpected end of input")

		test("for (;;)", "(anonymous): Line 1:9 Unexpected end of input")

		test("with (abc)", "(anonymous): Line 1:11 Unexpected end of input")

		test("try {}", "(anonymous): Line 1:1 Missing catch or finally after try")

		test("try {} catch () {}", "(anonymous): Line 1:15 Unexpected token )")

		test("\u203f = 1", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		// TODO
		// const x = 12, y;
		// const x, y = 12;
		// const x;
		// if(true) let a = 1;
		// if(true) const  a = 1;

		test(`new abc()."def"`, "(anonymous): Line 1:11 Unexpected string")

		test("/*", "(anonymous): Line 1:3 Unexpected end of input")

		test("/**", "(anonymous): Line 1:4 Unexpected end of input")

		test("/*\n\n\n", "(anonymous): Line 4:1 Unexpected end of input")

		test("/*\n\n\n*", "(anonymous): Line 4:2 Unexpected end of input")

		test("/*abc", "(anonymous): Line 1:6 Unexpected end of input")

		test("/*abc  *", "(anonymous): Line 1:9 Unexpected end of input")

		test("\n]", "(anonymous): Line 2:1 Unexpected token ]")

		test("\r\n]", "(anonymous): Line 2:1 Unexpected token ]")

		test("\n\r]", "(anonymous): Line 3:1 Unexpected token ]")

		test("//\r\n]", "(anonymous): Line 2:1 Unexpected token ]")

		test("//\n\r]", "(anonymous): Line 3:1 Unexpected token ]")

		test("/abc\\\n/", "(anonymous): Line 1:1 Invalid regular expression: missing /")

		test("//\r \n]", "(anonymous): Line 3:1 Unexpected token ]")

		test("/*\r\n*/]", "(anonymous): Line 2:3 Unexpected token ]")

		test("/*\r \n*/]", "(anonymous): Line 3:3 Unexpected token ]")

		test("\\\\", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test("\\u005c", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test("\\abc", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test("\\u0000", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test("\\u200c = []", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test("\\u200D = []", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test(`"\`, "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test(`"\u`, "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test("return", "(anonymous): Line 1:1 Illegal return statement")

		test("continue", "(anonymous): Line 1:1 Illegal continue statement")

		test("break", "(anonymous): Line 1:1 Illegal break statement")

		test("switch (abc) { default: continue; }", "(anonymous): Line 1:25 Illegal continue statement")

		test("do { abc } *", "(anonymous): Line 1:12 Unexpected token *")

		test("while (true) { break abc; }", "(anonymous): Line 1:16 Undefined label 'abc'")

		test("while (true) { continue abc; }", "(anonymous): Line 1:16 Undefined label 'abc'")

		test("abc: while (true) { (function(){ break abc; }); }", "(anonymous): Line 1:34 Undefined label 'abc'")

		test("abc: while (true) { (function(){ abc: break abc; }); }", nil)

		test("abc: while (true) { (function(){ continue abc; }); }", "(anonymous): Line 1:34 Undefined label 'abc'")

		test(`abc: if (0) break abc; else {}`, nil)

		test(`abc: if (0) { break abc; } else {}`, nil)

		test(`abc: if (0) { break abc } else {}`, nil)

		test("abc: while (true) { abc: while (true) {} }", "(anonymous): Line 1:21 Label 'abc' already exists")

		test(`if(0) { do { } while(0) } else { do { } while(0) }`, nil)

		test(`if(0) do { } while(0); else do { } while(0)`, nil)

		test("_: _: while (true) {]", "(anonymous): Line 1:4 Label '_' already exists")

		test("_:\n_:\nwhile (true) {]", "(anonymous): Line 2:1 Label '_' already exists")

		test("_:\n   _:\nwhile (true) {]", "(anonymous): Line 2:4 Label '_' already exists")

		test("function(){}", "(anonymous): Line 1:9 Unexpected token (")

		test("\n/*/", "(anonymous): Line 2:4 Unexpected end of input")

		test("/*/.source", "(anonymous): Line 1:11 Unexpected end of input")

		test("var class", "(anonymous): Line 1:5 Unexpected token class")

		test("var if", "(anonymous): Line 1:5 Unexpected token if")

		test("object Object", "(anonymous): Line 1:8 Unexpected identifier")

		test("[object Object]", "(anonymous): Line 1:9 Unexpected identifier")

		test("\\u0xyz", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test(`for (var abc, def in {}) {}`, "(anonymous): Line 1:19 Unexpected token in")

		test(`for (abc, def in {}) {}`, "(anonymous): Line 1:1 Invalid left-hand side in for-in or for-of")

		test(`for (var abc=def, ghi=("abc" in {}); true;) {}`, nil)

		{
			// Semicolon insertion

			test("this\nif (1);", nil)

			test("while (1) { break\nif (1); }", nil)

			test("throw\nif (1);", "(anonymous): Line 1:1 Illegal newline after throw")

			test("(function(){ return\nif (1); })", nil)

			test("while (1) { continue\nif (1); }", nil)

			test("debugger\nif (1);", nil)
		}

		{ // Reserved words

			test("class", "(anonymous): Line 1:6 Unexpected end of input")
			test("abc.class = 1", nil)
			test("var class;", "(anonymous): Line 1:5 Unexpected token class")

			test("const", "(anonymous): Line 1:6 Unexpected end of input")
			test("abc.const = 1", nil)
			test("var const;", "(anonymous): Line 1:5 Unexpected token const")

			test("enum", "(anonymous): Line 1:1 Unexpected reserved word")
			test("abc.enum = 1", nil)
			test("var enum;", "(anonymous): Line 1:5 Unexpected reserved word")

			test("export", "(anonymous): Line 1:1 Unexpected reserved word")
			test("abc.export = 1", nil)
			test("var export;", "(anonymous): Line 1:5 Unexpected reserved word")

			test("extends", "(anonymous): Line 1:1 Unexpected token extends")
			test("abc.extends = 1", nil)
			test("var extends;", "(anonymous): Line 1:5 Unexpected token extends")

			test("import", "(anonymous): Line 1:1 Unexpected reserved word")
			test("abc.import = 1", nil)
			test("var import;", "(anonymous): Line 1:5 Unexpected reserved word")

			test("super", "(anonymous): Line 1:1 'super' keyword unexpected here")
			test("abc.super = 1", nil)
			test("var super;", "(anonymous): Line 1:5 Unexpected token super")
			test(`
			obj = {
			  aaa: 1
			  bbb: "string"
			};`, "(anonymous): Line 4:6 Unexpected identifier")
			test("{}", nil)
			test("{a: 1}", nil)
			test("{a: 1,}", "(anonymous): Line 1:7 Unexpected token }")
			test("{a: 1, b: 2}", "(anonymous): Line 1:9 Unexpected token :")
			test("{a: 1, b: 2,}", "(anonymous): Line 1:9 Unexpected token :")
			test(`let f = () => new import('');`, "(anonymous): Line 1:19 Unexpected reserved word")

		}

		{ // Reserved words (strict)

			test(`implements`, nil)
			test(`abc.implements = 1`, nil)
			test(`var implements;`, nil)

			test(`interface`, nil)
			test(`abc.interface = 1`, nil)
			test(`var interface;`, nil)

			test(`let`, nil)
			test(`abc.let = 1`, nil)
			test(`var let;`, nil)

			test(`package`, nil)
			test(`abc.package = 1`, nil)
			test(`var package;`, nil)

			test(`private`, nil)
			test(`abc.private = 1`, nil)
			test(`var private;`, nil)

			test(`protected`, nil)
			test(`abc.protected = 1`, nil)
			test(`var protected;`, nil)

			test(`public`, nil)
			test(`abc.public = 1`, nil)
			test(`var public;`, nil)

			test(`static`, nil)
			test(`abc.static = 1`, nil)
			test(`var static;`, nil)

			test(`yield`, nil)
			test(`abc.yield = 1`, nil)
			test(`var yield;`, nil)
		}
		test(`0, { get a(param = null) {} };`, "(anonymous): Line 1:11 Getter must not have any formal parameters.")
		test(`let{f(`, "(anonymous): Line 1:7 Unexpected end of input")
		test("`", "(anonymous): Line 1:2 Unexpected end of input")
		test(" `", "(anonymous): Line 1:3 Unexpected end of input")
		test("` ", "(anonymous): Line 1:3 Unexpected end of input")
		test(`var{..(`, "(anonymous): Line 1:7 Unexpected token ILLEGAL")
		test(`var{get..(`, "(anonymous): Line 1:10 Unexpected token ILLEGAL")
		test(`var{set..(`, "(anonymous): Line 1:10 Unexpected token ILLEGAL")
		test(`(0 ?? 0 || true)`, "(anonymous): Line 1:9 Logical expressions and coalesce expressions cannot be mixed. Wrap either by parentheses")
		test(`(a || b ?? c)`, "(anonymous): Line 1:9 Logical expressions and coalesce expressions cannot be mixed. Wrap either by parentheses")
		test(`2 ?? 2 && 3 + 3`, "(anonymous): Line 1:3 Logical expressions and coalesce expressions cannot be mixed. Wrap either by parentheses")
		test(`
		class C {
            st\u0061tic m() {}
		}
		`, "(anonymous): Line 3:25 Unexpected identifier")
	})
}

func TestParser(t *testing.T) {
	tt(t, func() {
		test := func(source string, chk interface{}) *ast.Program {
			_, program, err := testParse(source)
			is(firstErr(err), chk)
			return program
		}
		test(`new (() => {});`, nil)

		test(`
            abc
            --
            []
        `, "(anonymous): Line 3:13 Invalid left-hand side in assignment")

		test(`
            abc--
            []
        `, nil)

		test("1\n[]\n", "(anonymous): Line 2:2 Unexpected token ]")

		test(`
            function abc() {
            }
            abc()
        `, nil)

		test("", nil)

		test("//", nil)

		test("/* */", nil)

		test("/** **/", nil)

		test("/*****/", nil)

		test("/*", "(anonymous): Line 1:3 Unexpected end of input")

		test("#", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test("/**/#", "(anonymous): Line 1:5 Unexpected token ILLEGAL")

		test("new +", "(anonymous): Line 1:5 Unexpected token +")

		program := test(";", nil)
		is(len(program.Body), 1)
		is(program.Body[0].(*ast.EmptyStatement).Semicolon, file.Idx(1))

		program = test(";;", nil)
		is(len(program.Body), 2)
		is(program.Body[0].(*ast.EmptyStatement).Semicolon, file.Idx(1))
		is(program.Body[1].(*ast.EmptyStatement).Semicolon, file.Idx(2))

		program = test("1.2", nil)
		is(len(program.Body), 1)
		is(program.Body[0].(*ast.ExpressionStatement).Expression.(*ast.NumberLiteral).Literal, "1.2")

		program = test("/* */1.2", nil)
		is(len(program.Body), 1)
		is(program.Body[0].(*ast.ExpressionStatement).Expression.(*ast.NumberLiteral).Literal, "1.2")

		program = test("\n", nil)
		is(len(program.Body), 0)

		test(`
            if (0) {
                abc = 0
            }
            else abc = 0
        `, nil)

		test("if (0) abc = 0 else abc = 0", "(anonymous): Line 1:16 Unexpected token else")

		test(`
            if (0) {
                abc = 0
            } else abc = 0
        `, nil)

		test(`
            if (0) {
                abc = 1
            } else {
            }
        `, nil)

		test(`
            do {
            } while (true)
        `, nil)

		test(`
            try {
            } finally {
            }
        `, nil)

		test(`
            try {
            } catch (abc) {
            } finally {
            }
        `, nil)

		test(`
            try {
            }
            catch (abc) {
            }
            finally {
            }
        `, nil)

		test(`try {} catch (abc) {} finally {}`, nil)

		test("try {} catch {}", nil)

		test(`
            do {
                do {
                } while (0)
            } while (0)
        `, nil)

		test(`
            (function(){
                try {
                    if (
                        1
                    ) {
                        return 1
                    }
                    return 0
                } finally {
                }
            })()
        `, nil)

		test("abc = ''\ndef", nil)

		test("abc = 1\ndef", nil)

		test("abc = Math\ndef", nil)

		test(`"\'"`, nil)

		test(`
            abc = function(){
            }
            abc = 0
        `, nil)

		test("abc.null = 0", nil)

		test("0x41", nil)

		test(`"\d"`, nil)

		test(`(function(){return this})`, nil)

		test(`
            Object.defineProperty(Array.prototype, "0", {
                value: 100,
                writable: false,
                configurable: true
            });
            abc = [101];
            abc.hasOwnProperty("0") && abc[0] === 101;
        `, nil)

		test(`new abc()`, nil)
		test(`new {}`, nil)

		test(`
            limit = 4
            result = 0
            while (limit) {
                limit = limit - 1
                if (limit) {
                }
                else {
                    break
                }
                result = result + 1
            }
        `, nil)

		test(`
            while (0) {
                if (0) {
                    continue
                }
            }
        `, nil)

		test("var \u0061\u0062\u0063 = 0", nil)

		// 7_3_1
		test("var test7_3_1\nabc = 66;", nil)
		test("var test7_3_1\u2028abc = 66;", nil)

		// 7_3_3
		test("//\u2028 =;", "(anonymous): Line 2:2 Unexpected token =")

		// 7_3_10
		test("var abc = \u2029;", "(anonymous): Line 2:1 Unexpected token ;")
		test("var abc = \\u2029;", "(anonymous): Line 1:11 Unexpected token ILLEGAL")
		test("var \\u0061\\u0062\\u0063 = 0;", nil)

		test("'", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		test("'\nstr\ning\n'", "(anonymous): Line 1:1 Unexpected token ILLEGAL")

		// S7.6_A4.3_T1
		test(`var $\u0030 = 0;`, nil)

		// S7.6.1.1_A1.1
		test(`switch = 1`, "(anonymous): Line 1:8 Unexpected token =")

		// S7.8.3_A2.1_T1
		test(`.0 === 0.0`, nil)

		// 7.8.5-1
		test("var regExp = /\\\rn/;", "(anonymous): Line 1:14 Invalid regular expression: missing /")

		// S7.8.5_A1.1_T2
		test("var regExp = /=/;", nil)

		// S7.8.5_A1.2_T1
		test("/*/", "(anonymous): Line 1:4 Unexpected end of input")

		// Sbp_7.9_A9_T3
		test(`
            do {
            ;
            } while (false) true
        `, nil)

		// S7.9_A10_T10
		test(`
            {a:1
            } 3
        `, nil)

		test(`
            abc
            ++def
        `, nil)

		// S7.9_A5.2_T1
		test(`
            for(false;false
            ) {
            break;
            }
        `, "(anonymous): Line 3:13 Unexpected token )")

		// S7.9_A9_T8
		test(`
            do {};
            while (false)
        `, "(anonymous): Line 2:18 Unexpected token ;")

		// S8.4_A5
		test(`
            "x\0y"
        `, nil)

		// S9.3.1_A6_T1
		test(`
            10e10000
        `, nil)

		// 10.4.2-1-5
		test(`
            "abc\
            def"
        `, nil)

		test("'\\\n'", nil)

		test("'\\\r\n'", nil)

		//// 11.13.1-1-1
		test("42 = 42;", "(anonymous): Line 1:1 Invalid left-hand side in assignment")
		test("s &^= 42;", "(anonymous): Line 1:4 Unexpected token ^=")

		// S11.13.2_A4.2_T1.3
		test(`
            abc /= "1"
        `, nil)

		// 12.1-1
		test(`
            try{};catch(){}
        `, "(anonymous): Line 2:13 Missing catch or finally after try")

		// 12.1-3
		test(`
            try{};finally{}
        `, "(anonymous): Line 2:13 Missing catch or finally after try")

		// S12.6.3_A11.1_T3
		test(`
            while (true) {
                break abc;
            }
        `, "(anonymous): Line 3:17 Undefined label 'abc'")

		// S15.3_A2_T1
		test(`var x / = 1;`, "(anonymous): Line 1:7 Unexpected token /")

		test(`
            function abc() {
                if (0)
                    return;
                else {
                }
            }
        `, nil)

		test("//\u2028 var =;", "(anonymous): Line 2:6 Unexpected token =")

		test(`
            throw
            {}
        `, "(anonymous): Line 2:13 Illegal newline after throw")

		// S7.6.1.1_A1.11
		test(`
            function = 1
        `, "(anonymous): Line 2:22 Unexpected token =")

		// S7.8.3_A1.2_T1
		test(`0e1`, nil)

		test("abc = 1; abc\n++", "(anonymous): Line 2:3 Unexpected end of input")

		// ---

		test("({ get abc() {} })", nil)

		test(`for (abc.def in {}) {}`, nil)

		test(`while (true) { break }`, nil)

		test(`while (true) { continue }`, nil)

		test(`abc=/^(?:(\w+:)\/{2}(\w+(?:\.\w+)*\/?)|(.{0,2}\/{1}))?([/.]*?(?:[^?]+)?\/)?((?:[^/?]+)\.(\w+))(?:\?(\S+)?)?$/,def=/^(?:(\w+:)\/{2})|(.{0,2}\/{1})?([/.]*?(?:[^?]+)?\/?)?$/`, nil)

		test(`(function() { try {} catch (err) {} finally {} return })`, nil)

		test(`0xde0b6b3a7640080.toFixed(0)`, nil)

		test(`/[^-._0-9A-Za-z\xb7\xc0-\xd6\xd8-\xf6\xf8-\u037d\u37f-\u1fff\u200c-\u200d\u203f\u2040\u2070-\u218f]/`, nil)

		test(`/[\u0000-\u0008\u000B-\u000C\u000E-\u001F\uD800-\uDFFF\uFFFE-\uFFFF]/`, nil)

		test("var abc = 1;\ufeff", nil)

		test("\ufeff/* var abc = 1; */", nil)

		test(`if (-0x8000000000000000<=abc&&abc<=0x8000000000000000) {}`, nil)

		test(`(function(){debugger;return this;})`, nil)

		test(`

        `, nil)

		test(`
            var abc = ""
            debugger
        `, nil)

		test(`
            var abc = /\[\]$/
            debugger
        `, nil)

		test(`
            var abc = 1 /
                2
            debugger
        `, nil)

		test("'ё\\\u2029'", nil)

		test(`[a, b] = [1, 2]`, nil)
		test(`({"a b": {}} = {})`, nil)

		test(`ref = (a, b = 39,) => {
		};`, nil)
		test(`(a,) => {}`, nil)

		test(`2 ?? (2 && 3) + 3`, nil)
		test(`(2 ?? 2) && 3 + 3`, nil)
		program = test(`a ?? b ?? c`, nil)
		is(len(program.Body), 1)
		is(program.Body[0].(*ast.ExpressionStatement).Expression.(*ast.BinaryExpression).Right.(*ast.Identifier).Name, "c")

		program = test(`
		class C {
			a
			b
			#c
			m() {
				return this.#c;
			}
		}
		`, nil)
		is(len(program.Body), 1)

		{
			program := test(`(-2)**53`, nil)
			st := program.Body[0].(*ast.ExpressionStatement).Expression.(*ast.BinaryExpression)
			is(st.Operator, token.EXPONENT)
			left := st.Left.(*ast.UnaryExpression)
			is(left.Operator, token.MINUS)
			op1 := left.Operand.(*ast.NumberLiteral)
			is(op1.Literal, "2")

			right := st.Right.(*ast.NumberLiteral)
			is(right.Literal, "53")
		}

	})
}

func TestParseDestruct(t *testing.T) {
	parser := newParser("", `({a: (a.b), ...spread,} = {})`)
	prg, err := parser.parse()
	if err != nil {
		t.Fatal(err)
	}
	_ = prg
}

func Test_parseStringLiteral(t *testing.T) {
	tt(t, func() {
		test := func(have string, want unistring.String) {
			parser := newParser("", have)
			parser.read()
			parser.read()
			_, res, err := parser.scanString(0, true)
			is(err, "")
			is(res, want)
		}

		test(`""`, "")
		test(`/=/`, "=")

		test("'1(\\\\d+)'", "1(\\d+)")

		test("'\\u2029'", "\u2029")

		test("'abc\\uFFFFabc'", "abc\uFFFFabc")

		test("'[First line \\\nSecond line \\\n Third line\\\n.     ]'",
			"[First line Second line  Third line.     ]")

		test("'\\u007a\\x79\\u000a\\x78'", "zy\nx")

		// S7.8.4_A4.2_T3
		test("'\\a'", "a")
		test("'\u0410'", "\u0410")

		// S7.8.4_A5.1_T1
		test("'\\0'", "\u0000")

		// S8.4_A5
		test("'\u0000'", "\u0000")

		// 15.5.4.20
		test("\"'abc'\\\n'def'\"", "'abc''def'")

		// 15.5.4.20-4-1
		test("\"'abc'\\\r\n'def'\"", "'abc''def'")

		// Octal
		test("'\\0'", "\000")
		test("'\\00'", "\000")
		test("'\\000'", "\000")
		test("'\\09'", "\0009")
		test("'\\009'", "\0009")
		test("'\\0009'", "\0009")
		test("'\\1'", "\001")
		test("'\\01'", "\001")
		test("'\\001'", "\001")
		test("'\\0011'", "\0011")
		test("'\\1abc'", "\001abc")

		test("'\\\u4e16'", "\u4e16")

		// err
		test = func(have string, want unistring.String) {
			parser := newParser("", have)
			parser.read()
			parser.read()
			_, res, err := parser.scanString(0, true)
			is(err, want)
			is(res, "")
		}

		test(`"\u"`, `invalid escape: \u: len("") != 4`)
		test(`"\u0"`, `invalid escape: \u: len("0") != 4`)
		test(`"\u00"`, `invalid escape: \u: len("00") != 4`)
		test(`"\u000"`, `invalid escape: \u: len("000") != 4`)

		test(`"\x"`, `invalid escape: \x: len("") != 2`)
		test(`"\x0"`, `invalid escape: \x: len("0") != 2`)
	})
}

func Test_parseNumberLiteral(t *testing.T) {
	tt(t, func() {
		test := func(input string, expect interface{}) {
			result, err := parseNumberLiteral(input)
			is(err, nil)
			is(result, expect)
		}

		test("0", 0)

		test("0x8000000000000000", float64(9.223372036854776e+18))
	})
}

func TestPosition(t *testing.T) {
	tt(t, func() {
		parser := newParser("", "// Lorem ipsum")

		// Out of range, idx0 (error condition)
		is(parser.slice(0, 1), "")
		is(parser.slice(0, 10), "")

		// Out of range, idx1 (error condition)
		is(parser.slice(1, 128), "")

		is(parser.str[0:0], "")
		is(parser.slice(1, 1), "")

		is(parser.str[0:1], "/")
		is(parser.slice(1, 2), "/")

		is(parser.str[0:14], "// Lorem ipsum")
		is(parser.slice(1, 15), "// Lorem ipsum")

		parser = newParser("", "(function(){ return 0; })")
		program, err := parser.parse()
		is(err, nil)

		var node ast.Node
		node = program.Body[0].(*ast.ExpressionStatement).Expression.(*ast.FunctionLiteral)
		is(node.Idx0(), file.Idx(2))
		is(node.Idx1(), file.Idx(25))
		is(parser.slice(node.Idx0(), node.Idx1()), "function(){ return 0; }")
		is(parser.slice(node.Idx0(), node.Idx1()+1), "function(){ return 0; })")
		is(parser.slice(node.Idx0(), node.Idx1()+2), "")
		is(node.(*ast.FunctionLiteral).Source, "function(){ return 0; }")

		node = program
		is(node.Idx0(), file.Idx(2))
		is(node.Idx1(), file.Idx(25))
		is(parser.slice(node.Idx0(), node.Idx1()), "function(){ return 0; }")

		parser = newParser("", "(function(){ return abc; })")
		program, err = parser.parse()
		is(err, nil)
		node = program.Body[0].(*ast.ExpressionStatement).Expression.(*ast.FunctionLiteral)
		is(node.(*ast.FunctionLiteral).Source, "function(){ return abc; }")
		node = program.Body[0].(*ast.ExpressionStatement).Expression.(*ast.FunctionLiteral).Body.List[0].(*ast.ReturnStatement)
		is(parser.slice(node.Idx0(), node.Idx1()), "return abc")

		parser = newParser("", "(function(){ return; })")
		program, err = parser.parse()
		is(err, nil)
		node = program.Body[0].(*ast.ExpressionStatement).Expression.(*ast.FunctionLiteral).Body.List[0].(*ast.ReturnStatement)
		is(parser.slice(node.Idx0(), node.Idx1()), "return")

		parser = newParser("", "true ? 1 : 2")
		program, err = parser.parse()
		is(err, nil)
		node = program.Body[0].(*ast.ExpressionStatement).Expression.(*ast.ConditionalExpression)
		is(parser.slice(node.Idx0(), node.Idx1()), "true ? 1 : 2")

		parser = newParser("", "a++")
		program, err = parser.parse()
		is(err, nil)
		node = program.Body[0].(*ast.ExpressionStatement).Expression.(*ast.UnaryExpression)
		is(parser.slice(node.Idx0(), node.Idx1()), "a++")

		parser = newParser("", "++a")
		program, err = parser.parse()
		is(err, nil)
		node = program.Body[0].(*ast.ExpressionStatement).Expression.(*ast.UnaryExpression)
		is(parser.slice(node.Idx0(), node.Idx1()), "++a")

		parser = newParser("", "xyz: for (i = 0; i < 10; i++) { if (i == 5) continue xyz; }")
		program, err = parser.parse()
		is(err, nil)
		node = program.Body[0].(*ast.LabelledStatement)
		is(parser.slice(node.Idx0(), node.Idx1()), "xyz: for (i = 0; i < 10; i++) { if (i == 5) continue xyz; }")
		node = program.Body[0].(*ast.LabelledStatement).Statement.(*ast.ForStatement)
		is(parser.slice(node.Idx0(), node.Idx1()), "for (i = 0; i < 10; i++) { if (i == 5) continue xyz; }")
		block := program.Body[0].(*ast.LabelledStatement).Statement.(*ast.ForStatement).Body.(*ast.BlockStatement)
		node = block.List[0].(*ast.IfStatement).Consequent.(*ast.BranchStatement)
		is(parser.slice(node.Idx0(), node.Idx1()), "continue xyz")

		parser = newParser("", "for (p in o) { break; }")
		program, err = parser.parse()
		is(err, nil)
		node = program.Body[0].(*ast.ForInStatement)
		is(parser.slice(node.Idx0(), node.Idx1()), "for (p in o) { break; }")
		node = program.Body[0].(*ast.ForInStatement).Body.(*ast.BlockStatement).List[0].(*ast.BranchStatement)
		is(parser.slice(node.Idx0(), node.Idx1()), "break")

		parser = newParser("", "while (i < 10) { i++; }")
		program, err = parser.parse()
		is(err, nil)
		node = program.Body[0].(*ast.WhileStatement)
		is(parser.slice(node.Idx0(), node.Idx1()), "while (i < 10) { i++; }")

		parser = newParser("", "do { i++; } while (i < 10 )")
		program, err = parser.parse()
		is(err, nil)
		node = program.Body[0].(*ast.DoWhileStatement)
		is(parser.slice(node.Idx0(), node.Idx1()), "do { i++; } while (i < 10 )")

		parser = newParser("", "with (1) {}")
		program, err = parser.parse()
		is(err, nil)
		node = program.Body[0].(*ast.WithStatement)
		is(parser.slice(node.Idx0(), node.Idx1()), "with (1) {}")

		parser = newParser("", `switch (a) {
	case 1: x--;
	case 2:
	default: x++;
}`)
		program, err = parser.parse()
		is(err, nil)
		node = program.Body[0].(*ast.SwitchStatement)
		is(parser.slice(node.Idx0(), node.Idx1()), `switch (a) {
	case 1: x--;
	case 2:
	default: x++;
}`)
	})
}

func TestExtractSourceMapLine(t *testing.T) {
	tt(t, func() {
		is(extractSourceMapLine(""), "")
		is(extractSourceMapLine("\n"), "")
		is(extractSourceMapLine(" "), "")
		is(extractSourceMapLine("1\n2\n3\n4\n"), "")

		src := `"use strict";
var x = {};
//# sourceMappingURL=delme.js.map`
		modSrc := `(function(exports, require, module) {` + src + `
})`
		is(extractSourceMapLine(modSrc), "//# sourceMappingURL=delme.js.map")
		is(extractSourceMapLine(modSrc+"\n\n\n\n"), "//# sourceMappingURL=delme.js.map")
	})
}

func TestSourceMapOptions(t *testing.T) {
	tt(t, func() {
		count := 0
		requestedPath := ""
		loader := func(p string) ([]byte, error) {
			count++
			requestedPath = p
			return nil, nil
		}
		src := `"use strict";
var x = {};
//# sourceMappingURL=delme.js.map`
		_, err := ParseFile(nil, "delme.js", src, 0, WithSourceMapLoader(loader))
		is(err, nil)
		is(count, 1)
		is(requestedPath, "delme.js.map")

		count = 0
		_, err = ParseFile(nil, "", src, 0, WithSourceMapLoader(loader))
		is(err, nil)
		is(count, 1)
		is(requestedPath, "delme.js.map")

		count = 0
		_, err = ParseFile(nil, "delme.js", src, 0, WithDisableSourceMaps)
		is(err, nil)
		is(count, 0)

		_, err = ParseFile(nil, "/home/user/src/delme.js", src, 0, WithSourceMapLoader(loader))
		is(err, nil)
		is(count, 1)
		is(requestedPath, "/home/user/src/delme.js.map")

		count = 0
		_, err = ParseFile(nil, "https://site.com/delme.js", src, 0, WithSourceMapLoader(loader))
		is(err, nil)
		is(count, 1)
		is(requestedPath, "https://site.com/delme.js.map")
	})
}

func TestParseTemplateCharacters(t *testing.T) {
	parser := newParser("", "`test\\\r\\\n${a}`")
	parser.next()
	if parser.token != token.BACKTICK {
		t.Fatalf("Token: %s", parser.token)
	}
	checkParseTemplateChars := func(expectedLiteral string, expectedParsed unistring.String, expectedFinished, expectParseErr, expectErr bool) {
		literal, parsed, finished, parseErr, err := parser.parseTemplateCharacters()
		if err != "" != expectErr {
			t.Fatal(err)
		}
		if literal != expectedLiteral {
			t.Fatalf("Literal: %q", literal)
		}
		if parsed != expectedParsed {
			t.Fatalf("Parsed: %q", parsed)
		}
		if finished != expectedFinished {
			t.Fatal(finished)
		}
		if parseErr != "" != expectParseErr {
			t.Fatalf("parseErr: %v", parseErr)
		}
	}
	checkParseTemplateChars("test\\\n\\\n", "test", false, false, false)
	parser.next()
	parser.expect(token.IDENTIFIER)
	if len(parser.errors) > 0 {
		t.Fatal(parser.errors)
	}
	if parser.token != token.RIGHT_BRACE {
		t.Fatal("Expected }")
	}
	if len(parser.errors) > 0 {
		t.Fatal(parser.errors)
	}
	checkParseTemplateChars("", "", true, false, false)
	if parser.chr != -1 {
		t.Fatal("Expected EOF")
	}
}

func TestParseTemplateLiteral(t *testing.T) {
	parser := newParser("", "f()\n`test${a}`")
	prg, err := parser.parse()
	if err != nil {
		t.Fatal(err)
	}
	if st, ok := prg.Body[0].(*ast.ExpressionStatement); ok {
		if expr, ok := st.Expression.(*ast.TemplateLiteral); ok {
			if expr.Tag == nil {
				t.Fatal("tag is nil")
			}
			if idx0 := expr.Tag.Idx0(); idx0 != 1 {
				t.Fatalf("Tag.Idx0(): %d", idx0)
			}
			if expr.OpenQuote != 5 {
				t.Fatalf("OpenQuote: %d", expr.OpenQuote)
			}
			if expr.CloseQuote != 14 {
				t.Fatalf("CloseQuote: %d", expr.CloseQuote)
			}
			if l := len(expr.Elements); l != 2 {
				t.Fatalf("len elements: %d", l)
			}
			if l := len(expr.Expressions); l != 1 {
				t.Fatalf("len expressions: %d", l)
			}
		} else {
			t.Fatal(st)
		}
	} else {
		t.Fatal(prg.Body[0])
	}
}

func TestParseTemplateLiteralWithTail(t *testing.T) {
	parser := newParser("", "f()\n`test${a}tail` ")
	prg, err := parser.parse()
	if err != nil {
		t.Fatal(err)
	}
	if st, ok := prg.Body[0].(*ast.ExpressionStatement); ok {
		if expr, ok := st.Expression.(*ast.TemplateLiteral); ok {
			if expr.CloseQuote != 18 {
				t.Fatalf("CloseQuote: %d", expr.CloseQuote)
			}
		} else {
			t.Fatal(st)
		}
	} else {
		t.Fatal(prg.Body[0])
	}
}
