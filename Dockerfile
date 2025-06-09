# Build stage
FROM golang:1.23-alpine AS builder

# セキュリティのためにnon-rootユーザーを作成
RUN adduser -D -g '' appuser

# 作業ディレクトリを設定
WORKDIR /app

# 依存関係をコピーしてダウンロード
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# バイナリをビルド（静的リンク）
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Runtime stage
FROM alpine:latest

# セキュリティアップデートとCA証明書
RUN apk --no-cache add ca-certificates tzdata

# non-rootユーザーを作成
RUN adduser -D -g '' appuser

# 作業ディレクトリを設定（appuserのホームディレクトリ）
WORKDIR /home/appuser/

# ビルド済みバイナリをコピー
COPY --from=builder /app/main .

# 静的ファイルをコピー
COPY --from=builder /app/cmd/server/static ./static
COPY --from=builder /app/cmd/server/templates ./templates

# ファイルの所有者をappuserに変更
RUN chown -R appuser:appuser /home/appuser/

# non-rootユーザーに切り替え
USER appuser

# ポートを公開
EXPOSE 8080

# 環境変数を設定
ENV PORT=8080
ENV GIN_MODE=release

# ヘルスチェック
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

# バイナリを実行
CMD ["./main"]