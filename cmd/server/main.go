package main

import (
	"embed"
	"log/slog"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/okamyuji/PasswordGenerator/internal/generator"
	"github.com/okamyuji/PasswordGenerator/internal/handler"
)

//go:embed templates/* static/*
var content embed.FS

func main() {
	// MIMEタイプの設定
	if err := mime.AddExtensionType(".css", "text/css"); err != nil {
		slog.Error("Failed to add MIME type for .css", "error", err)
	}
	if err := mime.AddExtensionType(".js", "application/javascript"); err != nil {
		slog.Error("Failed to add MIME type for .js", "error", err)
	}

	// 静的ファイル用のハンドラ
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[1:] // 先頭の'/'を除去
		data, err := content.ReadFile(path)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		// Content-Typeの設定
		ext := filepath.Ext(path)
		if mimeType := mime.TypeByExtension(ext); mimeType != "" {
			w.Header().Set("Content-Type", mimeType)
		}

		if _, err := w.Write(data); err != nil {
			slog.Error("Failed to write data", "error", err)
		}
	})

	// Dependency Injectionを使用したパスワード生成ハンドラの作成
	templateRenderer := handler.NewEmbedFSTemplateRenderer(content)
	passwordGenerator := generator.New()
	passwordHandler := handler.NewPasswordHandler(templateRenderer, passwordGenerator)

	http.HandleFunc("/", passwordHandler.Handle)

	slog.Info("Server starting at http://localhost:8080")
	slog.Error("Failed to start server", "error", http.ListenAndServe(":8080", nil))
}
