package main

import (
	"crypto/rand"
	"embed"
	"encoding/base64"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/okamyuji/PasswordGenerator/internal/generator"
	"github.com/okamyuji/PasswordGenerator/internal/handler"
	"github.com/okamyuji/PasswordGenerator/internal/middleware"
)

//go:embed templates/* static/* static/img/* static/css/* static/js/*
var content embed.FS

func main() {
	// 環境変数の設定
	if os.Getenv("APP_ENV") == "" {
		os.Setenv("APP_ENV", "development")
	}

	// CSRFトークンを生成して設定
	csrfToken := generateCSRFToken()
	os.Setenv("CSRF_TOKEN", csrfToken)

	// 構造化ロギングをセットアップ
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// ファイルシステムと静的コンテンツを埋め込み
	content, err := embedContent()
	if err != nil {
		logger.Error("コンテンツの埋め込みに失敗", "error", err)
		os.Exit(1)
	}

	// セキュリティミドルウェア
	securityMiddleware := middleware.NewSecurityMiddleware()

	// テンプレートレンダラー
	templateRenderer := handler.NewEmbedFSTemplateRenderer(content)

	// パスワードジェネレーター
	passwordGenerator := generator.New()

	// 依存性注入を使用したパスワードハンドラー
	passwordHandler := handler.NewPasswordHandler(templateRenderer, passwordGenerator)

	// ヘルスチェックエンドポイント
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			logger.Error("ヘルスチェックレスポンスの書き込みに失敗", "error", err)
		}
	})

	// ミドルウェアを使用したメインのパスワード生成ハンドラー
	http.HandleFunc("/", securityMiddleware.Middleware(passwordHandler.Handle))

	// セキュリティヘッダー付きの静的ファイル配信
	fs := http.FileServer(http.FS(content))
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		securityMiddleware.Middleware(func(w http.ResponseWriter, r *http.Request) {
			fs.ServeHTTP(w, r)
		})(w, r)
	})

	// サーバー構成
	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      nil,
	}

	// サーバー起動
	logger.Info("サーバー起動中", "address", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		logger.Error("サーバー起動失敗", "error", err)
		os.Exit(1)
	}
}

// generateCSRFToken は暗号学的に安全なランダムトークンを生成
func generateCSRFToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		slog.Error("CSRFトークンの生成に失敗", "error", err)
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

// embedContent は静的およびテンプレートファイルの埋め込みを処理
func embedContent() (embed.FS, error) {
	// 元の実装と同様、エラー処理付き
	return content, nil
}
