# 使用 Golang 官方映像檔
FROM golang:1.23

# 設定工作目錄
WORKDIR /app

# 安裝 goose（SQL 遷移工具）
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# **明確加入 PATH**
ENV PATH="/go/bin:${PATH}"

# 複製 Go 相關檔案
COPY go.mod go.sum ./
RUN go mod tidy

# 複製程式碼
COPY . .

RUN chmod -R 755 /app/sql

# 編譯 Go 爬蟲
RUN go build -o crawler .

# 複製 entrypoint.sh
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# 使用 entrypoint.sh 確保 PostgreSQL 啟動後才執行 goose
CMD ["/app/entrypoint.sh"]
