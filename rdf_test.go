package gon3

import (
	"testing"
)

func Test_unescapeUChar(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "on escapes",
			s:    "simple",
			want: "simple",
		},
		{
			name: "empty",
			s:    "",
			want: "",
		},
		{
			name: "single short escape",
			s:    `ab\u0063de`,
			want: "abcde",
		},
		{
			name: "single long escape",
			s:    `ab\U00000063de`,
			want: "abcde",
		},
		{
			name: "escaped short escape",
			s:    `a \\user`,
			want: `a \\user`,
		},
		{
			name: "escaped long escape",
			s:    `a \\User`,
			want: `a \\User`,
		},
		{
			name: "escaped escape before a short escape",
			s:    `a \\\u0063lass`,
			want: `a \\class`,
		},
		{
			name: "leading short escape",
			s:    `\u0061bc`,
			want: "abc",
		},
		{
			name: "leading escaped escape before a short escape",
			s:    `\\\u0061bc`,
			want: `\\abc`,
		},
		{
			name: "multiple escaped escapes",
			s:    `a \\\\user`,
			want: `a \\\\user`,
		},
		{
			name: "multiple short escapes",
			s:    `an \u0061b\u0063 \\user`,
			want: `an abc \\user`,
		},
		{
			name: "ending with escaped escape",
			s:    `end\\u`,
			want: `end\\u`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("unescapeUChar(%s) = panic: %v", tt.s, r)
				}
			}()
			if got := unescapeUChar(tt.s); got != tt.want {
				t.Errorf("unescapeUChar(%s) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}
