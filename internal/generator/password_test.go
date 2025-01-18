package generator

import (
	"strings"
	"testing"

	"github.com/okamyuji/PasswordGenerator/internal/config"
)

func TestGenerator_Generate(t *testing.T) {
	tests := []struct {
		name     string
		config   config.PasswordConfig
		wantLen  int
		wantErr  bool
		validate func(string) bool
	}{
		{
			name: "全ての文字種",
			config: config.PasswordConfig{
				Length:       12,
				UseUppercase: true,
				UseLowercase: true,
				UseNumbers:   true,
				UseSymbols:   true,
			},
			wantLen: 12,
			wantErr: false,
			validate: func(s string) bool {
				hasUpper := strings.ContainsAny(s, config.Uppercase)
				hasLower := strings.ContainsAny(s, config.Lowercase)
				hasNumber := strings.ContainsAny(s, config.Numbers)
				hasSymbol := strings.ContainsAny(s, config.Symbols)
				return hasUpper && hasLower && hasNumber && hasSymbol
			},
		},
		{
			name: "カスタム記号",
			config: config.PasswordConfig{
				Length:        8,
				UseUppercase:  true,
				UseLowercase:  true,
				UseSymbols:    true,
				CustomSymbols: "!@#",
			},
			wantLen: 8,
			wantErr: false,
			validate: func(s string) bool {
				hasUpper := strings.ContainsAny(s, config.Uppercase)
				hasLower := strings.ContainsAny(s, config.Lowercase)
				hasSymbol := strings.ContainsAny(s, "!@#")
				return hasUpper && hasLower && hasSymbol
			},
		},
		{
			name: "数字のみ",
			config: config.PasswordConfig{
				Length:     8,
				UseNumbers: true,
			},
			wantLen: 8,
			wantErr: false,
			validate: func(s string) bool {
				return strings.Trim(s, config.Numbers) == ""
			},
		},
		{
			name: "無効な長さ",
			config: config.PasswordConfig{
				Length:     0,
				UseNumbers: true,
			},
			wantLen: 0,
			wantErr: true,
			validate: func(s string) bool {
				return true
			},
		},
		{
			name: "文字種が選択されていない",
			config: config.PasswordConfig{
				Length: 8,
			},
			wantLen: 0,
			wantErr: true,
			validate: func(s string) bool {
				return true
			},
		},
	}

	g := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := g.Generate(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generator.Generate() エラー = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantLen {
				t.Errorf("Generator.Generate() 長さ = %v, want %v", len(got), tt.wantLen)
			}
			if !tt.wantErr && !tt.validate(got) {
				t.Errorf("Generator.Generate() = %v, 検証失敗", got)
			}
		})
	}
}

func TestGenerator_GenerateMultiple(t *testing.T) {
	g := New()
	config := config.PasswordConfig{
		Length:       12,
		UseUppercase: true,
		UseLowercase: true,
		UseNumbers:   true,
		UseSymbols:   true,
	}

	// 重複がないことを確認するために複数のパスワードを生成
	passwords := make(map[string]bool)
	for i := 0; i < 100; i++ {
		pass, err := g.Generate(config)
		if err != nil {
			t.Errorf("Generator.Generate() エラー = %v", err)
			continue
		}
		if passwords[pass] {
			t.Errorf("重複したパスワードを生成: %v", pass)
		}
		passwords[pass] = true
	}
}
