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
			name: "all character types",
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
			name: "custom symbols",
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
			name: "numbers only",
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
			name: "invalid length",
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
			name: "no character types selected",
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
				t.Errorf("Generator.Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantLen {
				t.Errorf("Generator.Generate() length = %v, want %v", len(got), tt.wantLen)
			}
			if !tt.wantErr && !tt.validate(got) {
				t.Errorf("Generator.Generate() = %v, failed validation", got)
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

	// Generate multiple passwords to ensure they're different
	passwords := make(map[string]bool)
	for i := 0; i < 100; i++ {
		pass, err := g.Generate(config)
		if err != nil {
			t.Errorf("Generator.Generate() error = %v", err)
			continue
		}
		if passwords[pass] {
			t.Errorf("Generated duplicate password: %v", pass)
		}
		passwords[pass] = true
	}
}
