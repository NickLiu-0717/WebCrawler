
# WebCrawler

![code coverage badge](https://github.com/NickLiu-0717/Webcrawler/actions/workflows/ci.yml/badge.svg)
![Go version](https://img.shields.io/badge/go-1.23-blue)

A distributed web crawler system with article classification capabilities. It supports JWT authentication and task dispatch via RabbitMQ, utilizes a pre-trained Python model for content classification, stores the results in PostgreSQL, and is fully deployable in a Docker environment. This project aims to implement a scalable and maintainable backend system and showcase my backend development skills in integration and architecture design.

## Features

- Automatically crawls article content from a given URL
- Distributed task processing via RabbitMQ
- Content classification using a pre-trained Python model
- Stores classified articles in a PostgreSQL database
- JWT-based user authentication to ensure secure API access
- Containerized deployment with Docker and Docker Compose

## Tech Stack

| Category         | Technology                             |
|--------------|----------------------------------|
| Languages         | Go, Python                       |
| Crawler         | Go net/http                      |
| Task Queue     | RabbitMQ                         |
| Classification     | Python + Pre-trained model          |
| Authentication     | JWT (JSON Web Token)             |
| Database       | PostgreSQL                       |
| Containerization       | Docker, Docker Compose           |

## System Architecture (Overview)

```
User → REST API (with JWT) → RabbitMQ (task queue) → Worker (crawling + classification) → PostgreSQL
                                                    ↳ Python model for content classification
```

## Getting Started

### Prerequisites

- Docker / Docker Compose
- Python 3.x
- Go 1.23+

### Environment Setup

1. Clone the repository：
```bash
git clone https://github.com/yourname/WebCrawler.git
cd WebCrawler
```

2. Configure the `.env` file

3. Download the pre-trained model locally
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

4. Start the service：
```bash
docker-compose up --build
```

5.API Usage：
- `http://localhost:<PORT>` （Port is defined in your `.env` file, e.g.,`PORT=8080`）

### RabbitMQ Management Interface
- Once the service is up, access RabbitMQ at [http://localhost:15672](http://localhost:15672)
- Default credentials: `guest/guest`

## Project Structure

```
WebCrawler/
├── app/                   # Tokenizer and model config
│   └── model/
├── classifier/            # Python classification API (with Dockerfile)
├── config/                # Go config files
├── crawl/                 # Crawling logic
├── handler/               # API handler（user/token/article）
├── internal/              # Auth, database, pubsub modules
│   ├── auth/
│   ├── database/
│   ├── models/
│   └── pubsub/
├── service/               # Core logic (classification, robots check, content extraction)
├── sql/                   # SQL schema and  queries
│   ├── schema/
│   └── queries/
├── static/                # HTML static pages
├── tests/                 # Unit tests
├── Dockerfile             # Docker config for main service
├── docker-compose.yml     # Multi-container orchestration
├── entrypoint.sh          # Docker entrypoint script
└── README.md              # Project documentation

```

## TODO & Future Improvements

- [ ] Support more website formats (e.g., Medium, Reddit)
- [ ] Improve classification model performance and accuracy
- [ ] Add a frontend dashboard to visualize crawled results
- [ ] Support multilingual article classification and translation

## Contributing
Developers interested in this project are welcome to fork and submit a pull request (PR) to the `main` branch.
For major changes, it's recommended to open an Issue for discussion first.
Please follow the existing code style and include necessary tests where applicable.
