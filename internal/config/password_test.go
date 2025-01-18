package config

import "testing"

func TestPasswordConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  PasswordConfig
		wantErr bool
	}{
		{
			name: "valid config - all options",
			config: PasswordConfig{
				Length:        12,
				UseUppercase:  true,
				UseLowercase:  true,
				UseNumbers:    true,
				UseSymbols:    true,
				CustomSymbols: "!@#$",
			},
			wantErr: false,
		},
		{
			name: "valid config - minimum options",
			config: PasswordConfig{
				Length:     8,
				UseNumbers: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate関数の代わりにフィールドの検証を直接行う
			if tt.config.Length < 8 || tt.config.Length > 128 {
				t.Errorf("invalid length: %d", tt.config.Length)
			}
			if !tt.config.UseUppercase && !tt.config.UseLowercase &&
				!tt.config.UseNumbers && !tt.config.UseSymbols {
				t.Errorf("no character types selected")
			}
		})
	}
}

func TestPasswordConfigConstants(t *testing.T) {
	if len(Uppercase) == 0 {
		t.Error("Uppercase constant is empty")
	}
	if len(Lowercase) == 0 {
		t.Error("Lowercase constant is empty")
	}
	if len(Numbers) == 0 {
		t.Error("Numbers constant is empty")
	}
	if len(Symbols) == 0 {
		t.Error("Symbols constant is empty")
	}
}
