# Password Generator

## 概要

このPasswordGeneratorは、安全でカスタマイズ可能なパスワードを生成するWebアプリケーションです。ユーザーは、パスワードの長さや含める文字の種類を柔軟に設定できます。

## 機能

- パスワード長のカスタマイズ（推奨範囲: 8〜128文字）
- 以下の文字種の選択が可能:
    - 大文字アルファベット
    - 小文字アルファベット
    - 数字
    - 記号
- カスタム記号の追加オプション
- 暗号学的に安全な乱数生成
- Webインターフェースでのパスワード生成

## 技術スタック

- 言語: Go (Golang)
- ウェブフレームワーク: 標準ライブラリ`net/http`
- テスト: Go標準のテスティングフレームワーク
- 依存性注入: カスタム実装

## 前提条件

- Go 1.21以上
- 以下の開発ツール（推奨）
    - goimports
    - staticcheck
    - golangci-lint
    - typos
    - codespell

## インストール

1. リポジトリをクローン

    ```bash
    git clone https://github.com/your-username/PasswordGenerator.git
    cd PasswordGenerator
    ```

2. 依存関係のインストール

    ```bash
    go mod tidy
    ```

## アプリケーションの実行

```bash
go run cmd/server/main.go
```

サーバーは `http://localhost:8080` で起動します。

## テストの実行

### 全テストの実行

```bash
go test ./... --shuffle=on
```

### 特定パッケージのテスト

```bash
go test ./internal/generator
go test ./internal/handler
```

### テストカバレッジの取得

#### カバレッジレポートの生成

```bash
# テスト全体のカバレッジ
go test ./... --shuffle=on -cover

# 詳細なカバレッジレポート
go test ./... --shuffle=on -coverprofile=coverage.out

# HTML形式のカバレッジレポート
go tool cover -html=coverage.out
```

## コード品質チェック

プロジェクトには `lint.sh` スクリプトが含まれており、以下のチェックを実行できます：

```bash
# lintスクリプトの実行
sh lint.sh
```

チェック内容

- `goimports`: コードのフォーマットと import の整理
- `go vet`: 静的解析
- `staticcheck`: 追加の静的解析
- `golangci-lint`: 包括的なコード品質チェック
- `typos`: スペルチェック
- `codespell`: コード内のスペルミスチェック

### lintスクリプトの詳細

`lint.sh` スクリプトは以下の特徴を持っています：

- 必要なコマンドの存在をチェック
- コマンドが見つからない場合にエラーを表示
- 各ツールの実行結果に応じて終了ステータスを設定
- `goimports`で検出されたエラーは自動修正
- その他のエラーは手動修正が必要

#### スクリプト使用上の注意

- false positiveがある場合:
    - `typos`については `_typos.toml` を編集
    - `codespell`については `.codespellrc` を編集

## セキュリティに関する注意

- パスワード生成には `crypto/rand` を使用し、暗号学的に安全な乱数を生成
- 生成されるパスワードはランダム性が高く、予測困難

## プロジェクト構造

```shell
.
├── cmd
│   └── server
│       ├── main.go          # アプリケーションのエントリーポイント
│       └── main_test.go     # サーバー関連のテスト
├── internal
│   ├── config
│   │   └── password.go      # パスワード設定の定義
│   ├── generator
│   │   └── password.go      # パスワード生成ロジック
│   └── handler
│       └── password.go      # HTTPハンドラー
└── lint.sh                  # コード品質チェックスクリプト
```

## 貢献

1. Issueを確認
2. フォーク
3. 機能追加/バグ修正のブランチを作成
4. テストを追加
5. コードをプッシュ
6. プルリクエストを送信

## ライセンス

MITライセンス

## 免責事項

このアプリケーションは教育目的および個人利用を想定しています。重要なシステムのパスワード生成には、専門的なパスワードマネージャーの利用を推奨します。

## お問い合わせ

質問や提案があれば、Issueを開いてください。
