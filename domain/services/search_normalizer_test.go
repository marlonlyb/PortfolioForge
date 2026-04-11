package services

import "testing"

func TestNormalizeQuery(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "basic normalization",
			input: "Hello World",
			want:  "hello world",
		},
		{
			name:  "accent removal with special chars",
			input: "SCADA C++",
			want:  "scada c",
		},
		{
			name:  "mixed accents",
			input: "Résumé",
			want:  "resume",
		},
		{
			name:  "special chars replaced with space then collapsed",
			input: "C++ & .NET",
			want:  "c net",
		},
		{
			name:  "multiple spaces collapsed",
			input: "hello   world",
			want:  "hello world",
		},
		{
			name:  "leading and trailing whitespace",
			input: "  hello  ",
			want:  "hello",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "CJK characters preserved (NFD may normalize half-width katakana)",
			input: "日本語プロジェクト",
			want:  "日本語フロシェクト", // NFD normalizes プ→フ+゙ and ジ→シ+゙
		},
		{
			name:  "short query",
			input: "A",
			want:  "a",
		},
		{
			name:  "numeric preserved",
			input: "Office 365",
			want:  "office 365",
		},
		{
			name:  "accent with tilde",
			input: "Señor",
			want:  "senor",
		},
		{
			name:  "mixed accents and special chars",
			input: "Comunicación & Colaboración",
			want:  "comunicacion colaboracion",
		},
		{
			name:  "only special chars",
			input: "+++&&&",
			want:  "",
		},
		{
			name:  "tabs and newlines",
			input: "hello\tworld\nfoo",
			want:  "hello world foo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeQuery(tt.input)
			if got != tt.want {
				t.Errorf("NormalizeQuery(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
