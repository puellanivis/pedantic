package pedantic

import (
	"regexp"
	"regexp/syntax"
	"testing"
)

func TestUnicode(t *testing.T) {
	utf8NonASCIIString := `^[` + utf8NonASCII + `]$`

	r, err := syntax.Parse(utf8NonASCIIString, syntax.Perl)
	if err != nil {
		t.Fatal("unexpected error ", err)
	}

	t.Log(r.Simplify())

	type test struct {
		b     []byte
		match bool
	}

	tests := []test{
		// Test unicode codepoints.
		{[]byte{0xc2, 0x80}, false},       // 0x0080: technically it should, but no.
		{[]byte{0xc2, 0xa0}, true},        // 0x00A0
		{[]byte{0xdf, 0xbf}, true},        // 0x07FF
		{[]byte{0xe0, 0xa0, 0x80}, true},  // 0x0800
		{[]byte{0xe0, 0xbf, 0xbf}, true},  // 0x0fff
		{[]byte{0xe1, 0x80, 0x80}, true},  // 0x1000
		{[]byte{0xec, 0xbf, 0xbf}, true},  // 0xcfff
		{[]byte{0xed, 0x80, 0x80}, true},  // 0xd000
		{[]byte{0xed, 0x9f, 0xbf}, true},  // 0xd7ff
		{[]byte{0xed, 0xa0, 0x80}, false}, // 0xd800
		// skip surrogates
		{[]byte{0xed, 0xbf, 0xbf}, false}, // 0xdfff
		{[]byte{0xee, 0x80, 0x80}, true},  // 0xe000
		{[]byte{0xef, 0xb7, 0x8f}, true},  // 0xfdcf
		{[]byte{0xef, 0xb7, 0x90}, false}, // 0xfdd0

		{[]byte{0xef, 0xb7, 0xaf}, false}, // 0xfdef
		{[]byte{0xef, 0xb7, 0xb0}, true},  // 0xfdf0
		{[]byte{0xef, 0xbf, 0xbc}, true},  // 0xfffc
		{[]byte{0xef, 0xbf, 0xbd}, false}, // 0xfffd

		{[]byte{0xef, 0xbf, 0xbf}, false},       // 0xffff
		{[]byte{0xf0, 0x90, 0x80, 0x80}, true},  // 0x010000
		{[]byte{0xf0, 0xbf, 0xbf, 0xbd}, true},  // 0x03fffd
		{[]byte{0xf0, 0xbf, 0xbf, 0xbe}, false}, // 0x03fffe
		{[]byte{0xf0, 0xbf, 0xbf, 0xbf}, false}, // 0x03ffff
		{[]byte{0xf1, 0x80, 0x80, 0x80}, true},  // 0x040000
		{[]byte{0xf3, 0xbf, 0xbf, 0xbd}, true},  // 0x0ffffd
		{[]byte{0xf3, 0xbf, 0xbf, 0xbe}, false}, // 0x0ffffe
		{[]byte{0xf3, 0xbf, 0xbf, 0xbf}, false}, // 0x0fffff
		{[]byte{0xf4, 0x80, 0x80, 0x80}, true},  // 0x100000
		{[]byte{0xf4, 0x8f, 0xbf, 0xbd}, true},  // 0x10fffd
		{[]byte{0xf4, 0x8f, 0xbf, 0xbe}, false}, // 0x10fffe
		{[]byte{0xf4, 0x8f, 0xbf, 0xbf}, false}, // 0x10ffff
		{[]byte{0xf4, 0x90, 0x80, 0x80}, false}, // 0x110000
	}

	utf8NonASCIIRegex := regexp.MustCompile(utf8NonASCIIString)
	for _, tt := range tests {
		got := utf8NonASCIIRegex.Match(tt.b)
		if got != tt.match {
			switch tt.match {
			case true:
				t.Errorf("expected %x (%q) to match, but it did not", tt.b, tt.b)
			case false:
				s := string(tt.b)
				r := []rune(s)
				t.Errorf("expected %x (%q) (%x) to not match, but it did", tt.b, tt.b, r)
			}
		}
	}
}