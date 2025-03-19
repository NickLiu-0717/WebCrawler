# 使用 Python 官方映像檔
FROM python:3.10

# 設定工作目錄
WORKDIR /app

# 設定 Hugging Face 模型快取目錄
ENV HF_HOME=/app/model

# 複製 Python 依賴
COPY requirements.txt .

# 安裝 Python 依賴
RUN pip install --no-cache-dir -r requirements.txt

# 預先下載 Hugging Face 模型
RUN python3 -c "from transformers import AutoModel; AutoModel.from_pretrained('joeddav/xlm-roberta-large-xnli')"

# 複製程式碼
COPY . .

# 執行 Python API
CMD ["python3", "classifier_api.py"]
