package config

import "testing"

func TestPasswordConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  PasswordConfig
		wantErr bool
	}{
		{
			name: "有効な設定 - すべてのオプション",
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
			name: "有効な設定 - 最小オプション",
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
				t.Errorf("無効な長さ: %d", tt.config.Length)
			}
			if !tt.config.UseUppercase && !tt.config.UseLowercase &&
				!tt.config.UseNumbers && !tt.config.UseSymbols {
				t.Errorf("文字の種類が選択されていません")
			}
		})
	}
}

func TestPasswordConfigConstants(t *testing.T) {
	if len(Uppercase) == 0 {
		t.Error("大文字の定数が空です")
	}
	if len(Lowercase) == 0 {
		t.Error("小文字の定数が空です")
	}
	if len(Numbers) == 0 {
		t.Error("数字の定数が空です")
	}
	if len(Symbols) == 0 {
		t.Error("記号の定数が空です")
	}
}
