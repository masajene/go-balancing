# ビルドステージ
FROM golang:1.20-alpine AS build

# 作業ディレクトリの設定
WORKDIR /app

# 依存パッケージのインストール
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download

# アプリケーションのビルド
COPY . .
RUN CGO_ENABLED=0 go build -o gateway

# 本番用イメージ
FROM scratch
COPY --from=build /app/gateway /gateway

# tomlファイルのコピー
COPY backends.toml /backends.toml

ENTRYPOINT ["/gateway"]