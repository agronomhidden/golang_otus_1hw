package hw02_unpack_string //nolint:golint,stylecheck

import (
	"testing"
	"github.com/stretchr/testify/require"
)

type test struct {
	input    string
	expected string
	err      error
}

func TestUnpack(t *testing.T) {
	for _, tst := range [...]test{
		{
			input:    "a4bc2d5e",
			expected: "aaaabccddddde",
		},
		{
			input:    "abccd",
			expected: "abccd",
		},
		{
			input:    "3abc",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "45",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "aaa10b",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "",
			expected: "",
		},
		{
			input:    "฿9",
			expected: "฿฿฿฿฿฿฿฿฿",
		},
		{
			input:    " 5",
			expected: "     ",
		},
		{
			input:    "5",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "0",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "a0",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    `a\f`,
			expected: "af",
		},
		{
			input:    "Ё3а4г1",
			expected: "ЁЁЁааааг",
		},
	} {
		result, err := Unpack(tst.input)
		require.Equal(t, tst.err, err)
		require.Equal(t, tst.expected, result)
	}
}

func TestUnpackWithEscape(t *testing.T) {
	//t.Skip() // Remove if task with asterisk completed

	for _, tst := range [...]test{
		{
			input:    `qwe\4\5`,
			expected: `qwe45`,
		},
		{
			input:    `qwe\45`,
			expected: `qwe44444`,
		},
		{
			input:    `qwe\\5`,
			expected: `qwe\\\\\`,
		},
		{
			input:    `qwe\\\3`,
			expected: `qwe\3`,
		},
		{
			input:    `\55`,
			expected: `55555`,
		},
		{
			input:    `\w`,
			expected: `w`,
		},
		{
			input:    `q\`,
			expected: ``,
			err:      ErrUnexpectedEnd,
		},
	} {
		result, err := Unpack(tst.input)
		require.Equal(t, tst.err, err)
		require.Equal(t, tst.expected, result)
	}
}
