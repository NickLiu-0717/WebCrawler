from fastapi import FastAPI
from pydantic import BaseModel
import torch
from transformers import pipeline
import logging
logging.getLogger("transformers.modeling_utils").setLevel(logging.ERROR)

app = FastAPI()
MODEL_PATH = "/app/model"
class ClassifyRequest(BaseModel):
    title: str

classifier = pipeline(
    "zero-shot-classification",
    model=MODEL_PATH,
    device=-1,
    framework="pt",
    torch_dtype=torch.float16,
    tokenizer=MODEL_PATH
)

# article = "研究揭清冠一號可抑制A型流感 降低發炎改善症狀 ｜ 公視新聞網 PNN"
labels = ["technology", "sports", "politics", "society", "entertainment", "health"]

# 讓 `pipeline()` 執行推理
# result = classifier(article, candidate_labels=labels)

# 再次確認 GPU 使用情況
# print(result['labels'][0])

@app.post("/classify/")
async def classify_article(request: ClassifyRequest):
    result = classifier(request.title, candidate_labels=labels)
    return {"catagory": result['labels'][0], "score": result['scores'][0]}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)