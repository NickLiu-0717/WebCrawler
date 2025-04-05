# 使用 Python 官方映像檔
FROM python:3.10

# 設定工作目錄
WORKDIR /app

# 設定 Hugging Face 模型快取目錄
ENV HF_HOME=/app/model

# 複製 Python 依賴
COPY requirements.txt .

# 安裝 Python 套件
RUN pip install --no-cache-dir -r requirements.txt

# 複製程式碼與 entrypoint
COPY . .

# 設定 entrypoint.sh 可執行
RUN chmod +x /app/entrypoint.sh

# 改為使用 entrypoint（這會取代原本的 CMD）
ENTRYPOINT ["/app/entrypoint.sh"]
