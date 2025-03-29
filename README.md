
# WebCrawler

一個具備文章分類能力的分散式網路爬蟲系統，支援 JWT 身份驗證與 RabbitMQ 任務分派，使用 Python 預訓練模型進行內容分類，將結果儲存於 PostgreSQL，整個系統可部署於 Docker 環境中。此專案目標為實作一個可擴充、可維護的後端系統，並展示我作為後端工程師的技術整合與架構能力。

## 🔍 專案功能

- ✅ 使用者輸入網址後，自動爬取該頁面內的文章內容
- ✅ 透過 RabbitMQ 實現任務分散式處理
- ✅ 使用 Python 的預訓練模型對文章內容進行分類
- ✅ 分類後的文章資料儲存至 PostgreSQL 資料庫
- ✅ 透過 JWT 驗證使用者身份，確保 API 存取安全
- ✅ 使用 Docker 與 Docker Compose 進行容器化部署

## 🧰 技術棧

| 類別         | 技術                             |
|--------------|----------------------------------|
| 語言         | Go, Python                       |
| 爬蟲         | Go net/http, goquery             |
| 任務佇列     | RabbitMQ                         |
| 模型分類     | Python + 預訓練分類模型           |
| 身分認證     | JWT (JSON Web Token)             |
| 資料庫       | PostgreSQL                       |
| 容器化       | Docker, Docker Compose           |

## 🧱 系統架構圖（簡述）

```
使用者 → REST API（含 JWT） → RabbitMQ（任務佇列）→ Worker（爬蟲 + 分類）→ PostgreSQL
                                                       ↳ Python model 進行分類
```

## 🚀 快速開始

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

3. 啟動服務：
```bash
docker-compose up --build
```

4. API 使用方式：
- `POST /crawl`：提交一個網址進行爬蟲與分類
- 需帶入 JWT Token 作為認證

### RabbitMQ 管理介面
- 服務啟動後可透過 [http://localhost:15672](http://localhost:15672) 進入 RabbitMQ 管理介面
- 預設帳密：`guest/guest`

## 📁 專案結構

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

## 📌 待辦與改進方向

- [ ] 支援更多網站的文章擷取格式（ex: Medium, Reddit）
- [ ] 改進分類模型效能與準確度
- [ ] 加入前端管理介面可視化爬取結果
- [ ] 支援多語言文章的分類與翻譯

## 🙋‍♂️ 作者

本專案由我個人獨立開發，目的是用來展示我在後端開發、系統架構設計與 DevOps 的實務能力，並作為求職 Junior Backend Engineer 的作品集。

如果你對我的專案有任何建議或想交流，歡迎聯絡我！
