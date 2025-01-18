package handler

import (
	"embed"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/okamyuji/PasswordGenerator/internal/config"
)

// パスワード生成のコントラクトを定義するインターフェース
type PasswordGeneratorInterface interface {
	Generate(cfg config.PasswordConfig) (string, error)
}

// テンプレートレンダリングを抽象化するインターフェース
type TemplateRendererInterface interface {
	ExecuteTemplate(w http.ResponseWriter, name string, data interface{}) error
}

// html/templateとembed.FSを使用してTemplateRendererInterfaceを実装
type EmbedFSTemplateRenderer struct {
	tmpl *template.Template
}

func NewEmbedFSTemplateRenderer(content embed.FS) *EmbedFSTemplateRenderer {
	return &EmbedFSTemplateRenderer{
		tmpl: template.Must(template.ParseFS(content, "templates/*.html")),
	}
}

func (r *EmbedFSTemplateRenderer) ExecuteTemplate(w http.ResponseWriter, name string, data interface{}) error {
	return r.tmpl.ExecuteTemplate(w, name, data)
}

// インターフェースに依存する、具象実装ではないPasswordHandler
type PasswordHandler struct {
	renderer  TemplateRendererInterface
	generator PasswordGeneratorInterface
}

// 依存性注入を使用して新しいPasswordHandlerを作成
func NewPasswordHandler(
	renderer TemplateRendererInterface,
	generator PasswordGeneratorInterface,
) *PasswordHandler {
	return &PasswordHandler{
		renderer:  renderer,
		generator: generator,
	}
}

func (h *PasswordHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.handleGet(w, r)
		return
	}
	if r.Method == http.MethodPost {
		h.handlePost(w, r)
		return
	}
	http.Error(w, "メソッドは許可されていません", http.StatusMethodNotAllowed)
}

func (h *PasswordHandler) handleGet(w http.ResponseWriter, _ *http.Request) {
	// CSRFトークンを取得
	csrfToken := os.Getenv("CSRF_TOKEN")

	// テンプレートに渡すデータを準備
	data := map[string]string{
		"CSRFToken": csrfToken,
	}

	if err := h.renderer.ExecuteTemplate(w, "index.html", data); err != nil {
		slog.Error("テンプレート実行エラー", "error", err)
		http.Error(w, "内部サーバーエラー", http.StatusInternalServerError)
	}
}

func (h *PasswordHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "無効なフォームデータ", http.StatusBadRequest)
		return
	}

	length, _ := strconv.Atoi(r.FormValue("length"))
	// 長さのバリデーション
	if length <= 0 {
		http.Error(w, "無効な長さ", http.StatusBadRequest)
		return
	}

	pwdConfig := config.PasswordConfig{
		Length:        length,
		UseUppercase:  r.Form.Get("uppercase") == "true",
		UseLowercase:  r.Form.Get("lowercase") == "true",
		UseNumbers:    r.Form.Get("numbers") == "true",
		UseSymbols:    r.Form.Get("symbols") == "true",
		CustomSymbols: strings.TrimSpace(r.Form.Get("customSymbols")),
	}

	password, err := h.generator.Generate(pwdConfig)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte(password)); err != nil {
		slog.Error("パスワードの書き込みに失敗", "error", err)
		http.Error(w, "内部サーバーエラー", http.StatusInternalServerError)
	}
}
