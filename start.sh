# #!/bin/bash

# # 啟動 Python API
# echo "🚀 啟動 Python API..."
# cd classifier
# # pip install -r requirements.txt
# nohup python3 classifier_api.py > classifier.log 2>&1 &
# cd ..

# # 啟動 Go Web Crawler
# echo "🚀 啟動 Go Web Crawler..."
# go run . "https://www.bbc.com" 150 100 3


#!/bin/bash
# 啟動 Python API
echo "🚀 啟動 Python API..."
cd classifier
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
        echo "✅ API 已準備就緒！分類結果: $response"
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
go run . "https://www.bbc.com" 150 100 3


