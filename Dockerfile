# 使用するGoのバージョンを指定
FROM golang:1.22 as builder

# ビルド成果物を保存するディレクトリを作成
RUN mkdir /artifact

WORKDIR /workspace

COPY . .

# CGOを無効にしてバイナリをビルド
RUN CGO_ENABLED=0 go build -o /artifact/app ./cmd/server

FROM gcr.io/distroless/base-debian10

COPY --from=builder /artifact/app /app

# 静的ファイル用のディレクトリとファイルをコピー
COPY --from=builder /workspace/templates /templates
COPY --from=builder /workspace/authenticated /authenticated

# コンテナがリッスンするポートを指定
EXPOSE 8080

# コンテナ起動時に実行されるコマンド
CMD [ "/app" ]
