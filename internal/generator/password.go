package generator

import (
	"crypto/rand"
	"fmt"
	"strings"

	"github.com/okamyuji/PasswordGenerator/internal/config"
)

type Generator struct{}

func New() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(cfg config.PasswordConfig) (string, error) {
	// 最初にバリデーションを実行
	const MaxPasswordLength = 1000 // パスワードの最大長を設定
	
	if cfg.Length <= 0 {
		return "", fmt.Errorf("無効な長さ: %d", cfg.Length)
	}
	
	if cfg.Length > MaxPasswordLength {
		return "", fmt.Errorf("パスワード長が最大値を超えています: %d (最大: %d)", cfg.Length, MaxPasswordLength)
	}

	var charsets []string
	if cfg.UseUppercase {
		charsets = append(charsets, config.Uppercase)
	}
	if cfg.UseLowercase {
		charsets = append(charsets, config.Lowercase)
	}
	if cfg.UseNumbers {
		charsets = append(charsets, config.Numbers)
	}
	if cfg.UseSymbols {
		if cfg.CustomSymbols != "" {
			charsets = append(charsets, cfg.CustomSymbols)
		} else {
			charsets = append(charsets, config.Symbols)
		}
	}

	if len(charsets) == 0 {
		return "", fmt.Errorf("no character types selected")
	}

	// 選択された文字セットの準備
	if cfg.UseUppercase {
		charsets = append(charsets, config.Uppercase)
	}
	if cfg.UseLowercase {
		charsets = append(charsets, config.Lowercase)
	}
	if cfg.UseNumbers {
		charsets = append(charsets, config.Numbers)
	}
	if cfg.UseSymbols {
		if cfg.CustomSymbols != "" {
			charsets = append(charsets, cfg.CustomSymbols)
		} else {
			charsets = append(charsets, config.Symbols)
		}
	}

	if len(charsets) == 0 || cfg.Length <= 0 {
		return "", nil
	}

	result := make([]byte, cfg.Length)
	used := make(map[int]bool)

	// 各文字セットから1文字ずつ必ず選択
	for i, charset := range charsets {
		if i >= cfg.Length {
			break
		}
		randomByte := make([]byte, 1)
		if _, err := rand.Read(randomByte); err != nil {
			return "", err
		}
		pos := int(randomByte[0]) % cfg.Length
		for used[pos] {
			pos = (pos + 1) % cfg.Length
		}
		result[pos] = charset[int(randomByte[0])%len(charset)]
		used[pos] = true
	}

	// 残りの文字を全文字セットからランダムに選択
	allChars := strings.Join(charsets, "")
	for i := 0; i < cfg.Length; i++ {
		if !used[i] {
			randomByte := make([]byte, 1)
			if _, err := rand.Read(randomByte); err != nil {
				return "", err
			}
			result[i] = allChars[int(randomByte[0])%len(allChars)]
		}
	}

	// 生成されたパスワードをシャッフル
	for i := len(result) - 1; i > 0; i-- {
		randomByte := make([]byte, 1)
		if _, err := rand.Read(randomByte); err != nil {
			return "", err
		}
		j := int(randomByte[0]) % (i + 1)
		result[i], result[j] = result[j], result[i]
	}

	return string(result), nil
}
