#!/bin/sh

# 設定模型路徑
MODEL_DIR="/app/model"

mkdir -p "$MODEL_DIR"

# 如果模型不存在，就下載
if [ ! -f "$MODEL_DIR/config.json" ]; then
  echo "Model not found, downloading..."
  python3 -c "
from transformers import AutoModelForSequenceClassification, AutoTokenizer
model_name = 'joeddav/xlm-roberta-large-xnli'
model = AutoModelForSequenceClassification.from_pretrained(model_name)
tokenizer = AutoTokenizer.from_pretrained(model_name)
model.save_pretrained('$MODEL_DIR')
tokenizer.save_pretrained('$MODEL_DIR')
"
else
  echo "Model already exists, skipping download."
fi

# 啟動 API
exec python3 classifier_api.py