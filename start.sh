#!/bin/bash

# if ! command -v goose &> /dev/null; then
#     echo "🚨 goose 未安裝，正在安裝..."
#     go install github.com/pressly/goose/v3/cmd/goose@latest
# fi

# 檢查是否已安裝 PostgreSQL
# if ! command -v psql &> /dev/null; then
#     echo "🚨 PostgreSQL 未安裝，正在安裝..."
    
#     # 根據作業系統安裝 PostgreSQL
#     if [[ "$OSTYPE" == "linux-gnu"* ]]; then
#         sudo apt update
#         sudo apt install -y postgresql postgresql-contrib
#     elif [[ "$OSTYPE" == "darwin"* ]]; then
#         brew install postgresql
#     else
#         echo "❌ 無法自動安裝 PostgreSQL，請手動安裝！"
#         exit 1
#     fi
# fi

# 確保 PostgreSQL 服務正在運行
# if ! pg_isready -q; then
#     echo "🔄 啟動 PostgreSQL 服務..."
#     sudo service postgresql start
# fi

# 等待 PostgreSQL 完全啟動
# sleep 2

# 設定資料庫名稱與使用者
DB_NAME="crawler"
DB_USER="postgres"
DB_PASS="postgres"

# 檢查資料庫是否存在
# DB_EXIST=$(sudo -u postgres psql -tAc "SELECT 1 FROM pg_database WHERE datname='$DB_NAME'")

# if [[ "$DB_EXIST" != "1" ]]; then
#     echo "📦 建立資料庫 $DB_NAME..."
#     sudo -u postgres psql <<EOF
# CREATE DATABASE $DB_NAME;
# CREATE USER $DB_USER postgres PASSWORD '$DB_PASS';
# EOF
# else
#     echo "✅ 資料庫 $DB_NAME 已存在，跳過建立步驟。"
# fi




# cd sql/schema
# goose postgres "postgres://$DB_USER:$DB_PASS@localhost:5432/$DB_NAME" up
# cd ../..

# 啟動 Python API
echo "🚀 啟動 Python API..."
cd classifier
# python3 -m pip install -r requirements.txt
nohup python3 classifier_api.py > classifier.log 2>&1 &
cd ..

# 等待 Python API 完全啟動
echo "⌛ 等待 Python API 啟動並能夠正確分類..."
API_URL="http://127.0.0.1:8000/classify/"
MAX_RETRIES=30
SLEEP_TIME=2

TEST_JSON='{"title": "Stock market reaches record high"}'

for ((i=1; i<=MAX_RETRIES; i++)); do
    response=$(curl --silent --header "Content-Type: application/json" \
                    --request POST --data "$TEST_JSON" "$API_URL" | jq -r '.catagory')
    
    if [[ "$response" != "null" && "$response" != "" ]]; then
        echo "✅ API 已準備就緒！"
        break
    fi

    echo "⏳ API 尚未完全準備好，重試中 ($i/$MAX_RETRIES)..."
    sleep $SLEEP_TIME
done

# 確保 API 真的可用，否則退出腳本
if [[ "$response" == "null" || "$response" == "" ]]; then
    echo "❌ API 啟動失敗，請檢查 classifier.log 記錄錯誤訊息"
    exit 1
fi

# 啟動 Go Web Crawler
echo "🚀 啟動 Go Web Crawler..."
go run .


