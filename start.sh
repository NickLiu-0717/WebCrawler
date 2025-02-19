#!/bin/bash

# 啟動 Python API
echo "🚀 啟動 Python API..."
cd classifier
# pip install -r requirements.txt
nohup python3 classifier_api.py > classifier.log 2>&1 &
cd ..

# 啟動 Go Web Crawler
echo "🚀 啟動 Go Web Crawler..."
go run . "https://www.bbc.com" 150 100 3
