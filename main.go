package main

import (
	"crypto/rand"
	"embed"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

//go:embed templates/*
var content embed.FS

type PasswordConfig struct {
	Length        int    `json:"length"`
	UseUppercase  bool   `json:"useUppercase"`
	UseLowercase  bool   `json:"useLowercase"`
	UseNumbers    bool   `json:"useNumbers"`
	UseSymbols    bool   `json:"useSymbols"`
	CustomSymbols string `json:"customSymbols"`
}

const (
	uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowercase = "abcdefghijklmnopqrstuvwxyz"
	numbers   = "0123456789"
	symbols   = "!@#$%^&*()_+-=[]{}|;:,.<>?"
)

func generatePassword(config PasswordConfig) string {
	// 選択された文字セットの準備
	var charsets []string
	if config.UseUppercase {
		charsets = append(charsets, uppercase)
	}
	if config.UseLowercase {
		charsets = append(charsets, lowercase)
	}
	if config.UseNumbers {
		charsets = append(charsets, numbers)
	}
	if config.UseSymbols {
		if config.CustomSymbols != "" {
			charsets = append(charsets, config.CustomSymbols)
		} else {
			charsets = append(charsets, symbols)
		}
	}

	if len(charsets) == 0 || config.Length <= 0 {
		return ""
	}

	// 各文字セットから最低1文字を選択
	result := make([]byte, config.Length)
	used := make(map[int]bool)

	// 各文字セットから1文字ずつ必ず選択
	for i, charset := range charsets {
		if i >= config.Length {
			break
		}
		randomByte := make([]byte, 1)
		if _, err := rand.Read(randomByte); err != nil {
			return ""
		}
		pos := int(randomByte[0]) % config.Length
		// 既に使用済みの位置の場合、空いている位置を探す
		for used[pos] {
			pos = (pos + 1) % config.Length
		}
		result[pos] = charset[int(randomByte[0])%len(charset)]
		used[pos] = true
	}

	// 残りの文字を全文字セットからランダムに選択
	allChars := strings.Join(charsets, "")
	for i := 0; i < config.Length; i++ {
		if !used[i] {
			randomByte := make([]byte, 1)
			if _, err := rand.Read(randomByte); err != nil {
				return ""
			}
			result[i] = allChars[int(randomByte[0])%len(allChars)]
		}
	}

	// 生成されたパスワードをランダムにシャッフル
	for i := len(result) - 1; i > 0; i-- {
		randomByte := make([]byte, 1)
		if _, err := rand.Read(randomByte); err != nil {
			return ""
		}
		j := int(randomByte[0]) % (i + 1)
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}

func main() {
	tmpl := template.Must(template.ParseFS(content, "templates/*.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			tmpl.ExecuteTemplate(w, "index.html", nil)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		length, _ := strconv.Atoi(r.FormValue("length"))
		if length != 8 && length != 12 {
			length = 12 // デフォルト値
		}

		// デバッグ出力
		log.Printf("Raw form data: %+v", r.Form)

		config := PasswordConfig{
			Length:        length,
			UseUppercase:  r.Form.Get("uppercase") == "true",
			UseLowercase:  r.Form.Get("lowercase") == "true",
			UseNumbers:    r.Form.Get("numbers") == "true",
			UseSymbols:    r.Form.Get("symbols") == "true",
			CustomSymbols: strings.TrimSpace(r.Form.Get("customSymbols")),
		}

		password := generatePassword(config)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(password))
	})

	slog.Info("Server starting at http://localhost:8080")
	slog.Error("Server starting at http://localhost:8080", "error", http.ListenAndServe(":8080", nil))
}
