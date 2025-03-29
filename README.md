
# WebCrawler

一個具備文章分類能力的分散式網路爬蟲系統，支援 JWT 身份驗證與 RabbitMQ 任務分派，使用 Python 預訓練模型進行內容分類，將結果儲存於 PostgreSQL，整個系統可部署於 Docker 環境中。此專案目標為實作一個可擴充、可維護的後端系統，並展示我作為後端工程師的技術整合與架構能力。

## 專案功能

- 使用者給定網址後自動爬取該頁面內的文章內容
- 透過 RabbitMQ 實現任務分散式處理
- 使用 Python 的預訓練模型對文章內容進行分類
- 分類後的文章資料儲存至 PostgreSQL 資料庫
- 透過 JWT 驗證使用者身份，確保 API 存取安全
- 使用 Docker 與 Docker Compose 進行容器化部署

## 技術棧

| 類別         | 技術                             |
|--------------|----------------------------------|
| 語言         | Go, Python                       |
| 爬蟲         | Go net/http                      |
| 任務佇列     | RabbitMQ                         |
| 模型分類     | Python + 預訓練分類模型           |
| 身分認證     | JWT (JSON Web Token)             |
| 資料庫       | PostgreSQL                       |
| 容器化       | Docker, Docker Compose           |

## 系統架構圖（簡述）

```
使用者 → REST API（含 JWT） → RabbitMQ（任務佇列）→ Worker（爬蟲 + 分類）→ PostgreSQL
                                                       ↳ Python model 進行分類
```

## 快速開始

### 前置需求

- Docker / Docker Compose
- Python 3.x
- Go 1.20+

### 環境設定

1. 複製專案：
```bash
git clone https://github.com/yourname/WebCrawler.git
cd WebCrawler
```

2. 設定 `.env` 檔案

3. 下載預訓練模型至本地
```
python3 -c "
from transformers import AutoModelForSequenceClassification, AutoTokenizer;

model_name = 'joeddav/xlm-roberta-large-xnli'
model = AutoModelForSequenceClassification.from_pretrained(model_name)
tokenizer = AutoTokenizer.from_pretrained(model_name)

model.save_pretrained('app/model')
tokenizer.save_pretrained('app/model')
"
```

4. 啟動服務：
```bash
docker-compose up --build
```

5. API 使用方式：
- `http://localhost:<PORT>` （PORT 由 `.env` 中的設定決定，例如 `PORT=8080`）

### RabbitMQ 管理介面
- 服務啟動後可透過 [http://localhost:15672](http://localhost:15672) 進入 RabbitMQ 管理介面
- 預設帳密：`guest/guest`

## 專案結構

```
WebCrawler/
├── app/                   # 分類模型的 tokenizer 設定
│   └── model/
├── classifier/            # Python 類別分類 API（含 Dockerfile）
├── config/                # Go 設定檔
├── crawl/                 # 爬蟲邏輯
├── handler/               # API handler（user/token/article 等）
├── internal/              # 認證、資料庫、pubsub 模組
│   ├── auth/
│   ├── database/
│   ├── models/
│   └── pubsub/
├── service/               # 核心邏輯（分類、檢查 robots、萃取文章）
├── sql/                   # SQL schema 與 queries
│   ├── schema/
│   └── queries/
├── static/                # HTML 靜態頁面
├── tests/                 # 單元測試
├── Dockerfile             # 主服務 Docker 設定
├── docker-compose.yml     # 多容器編排
├── entrypoint.sh          # Docker 進入點 script
└── README.md              # 專案說明文件

```

## 待辦與改進方向

- [ ] 支援更多網站的文章擷取格式（ex: Medium, Reddit）
- [ ] 改進分類模型效能與準確度
- [ ] 加入前端管理介面可視化爬取結果
- [ ] 支援多語言文章的分類與翻譯
