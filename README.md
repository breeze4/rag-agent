# RAG Therapist

A Retrieval-Augmented Generation (RAG) chatbot service that processes PDF documents and provides intelligent responses using Claude or Gemini LLMs.

## Features

- PDF document upload and processing
- Vector-based document search using Chroma
- Support for multiple LLM providers (Claude, Gemini)
- SQLite database for document metadata
- RESTful API for chat interactions
- Automatic text chunking and embedding generation

## Prerequisites

- Go 1.23.0 or later
- Git

## Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd rag-therapist
   ```

2. **Install Go dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   
   Create a `.env` file in the project root:
   ```bash
   cp .env.example .env
   ```
   
   Edit `.env` and configure the following variables:
   ```env
   # LLM Provider (claude or gemini)
   LLM_PROVIDER=claude
   
   # API Keys
   ANTHROPIC_API_KEY=your_anthropic_api_key_here
   GOOGLE_API_KEY=your_google_api_key_here
   OPENAI_API_KEY=your_openai_api_key_here  # For embeddings
   
   # Server Configuration
   PORT=8080
   DATA_DIR=./data
   
   # Chroma Vector Store
   CHROMA_URL=http://localhost:8000
   ```

4. **Set up Chroma vector database**
   
   Using Docker:
   ```bash
   docker run -d --name chroma -p 8000:8000 chromadb/chroma:latest
   ```
   
   Or install locally:
   ```bash
   pip install chromadb
   chroma run --host 0.0.0.0 --port 8000
   ```

## Build

### Development Build
```bash
go build -o bin/rag-therapist ./cmd/server
```

### Production Build
```bash
go build -ldflags="-w -s" -o bin/rag-therapist ./cmd/server
```

### Cross-platform Build
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o bin/rag-therapist-linux ./cmd/server

# macOS
GOOS=darwin GOARCH=amd64 go build -o bin/rag-therapist-macos ./cmd/server

# Windows
GOOS=windows GOARCH=amd64 go build -o bin/rag-therapist.exe ./cmd/server
```

## Run

### Using Go Run (Development)
```bash
go run ./cmd/server
```

### Using Built Binary
```bash
./bin/rag-therapist
```

### Using Docker (Optional)
```bash
# Build Docker image
docker build -t rag-therapist .

# Run container
docker run -p 8080:8080 --env-file .env rag-therapist
```

## API Usage

### Upload a PDF Document
```bash
curl -X POST http://localhost:8080/upload \
  -F "file=@document.pdf"
```

Response:
```json
{
  "id": 1,
  "file_name": "document.pdf",
  "status": "pending",
  "uploaded_at": "2024-01-15T10:30:00Z"
}
```

### Chat with Documents
```bash
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "What is the main topic of the uploaded document?"
  }'
```

Response:
```json
{
  "response": "Based on the uploaded document, the main topic is...",
  "sources": [
    {
      "document_id": 1,
      "chunk_id": "chunk_123",
      "relevance_score": 0.95
    }
  ]
}
```

### List Documents
```bash
curl http://localhost:8080/documents
```

## Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `LLM_PROVIDER` | LLM provider (claude/gemini) | claude | Yes |
| `ANTHROPIC_API_KEY` | Anthropic API key for Claude | - | If using Claude |
| `GOOGLE_API_KEY` | Google API key for Gemini | - | If using Gemini |
| `OPENAI_API_KEY` | OpenAI API key for embeddings | - | Yes |
| `PORT` | HTTP server port | 8080 | No |
| `DATA_DIR` | Data storage directory | ./data | No |
| `CHROMA_URL` | Chroma vector store URL | http://localhost:8000 | No |

### Data Directory Structure
```
data/
├── rag-therapist.db    # SQLite database
└── documents/          # Uploaded PDF files
    ├── 20240115_103000_document1.pdf
    └── 20240115_104500_document2.pdf
```

## Development

### Running Tests
```bash
go test ./...
```

### Code Formatting
```bash
go fmt ./...
```

### Linting
```bash
# Install golangci-lint first
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2

# Run linter
golangci-lint run
```

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go         # Configuration management
│   ├── llm/                  # LLM client implementations
│   ├── rag/                  # RAG pipeline logic
│   └── storage/              # Database and file storage
│       ├── database.go
│       ├── document_repository.go
│       ├── file_storage.go
│       └── storage_service.go
├── pkg/
│   └── models/
│       └── document.go       # Data models
├── .env.example              # Environment variables template
├── go.mod                    # Go module definition
├── go.sum                    # Go module checksums
└── README.md                 # This file
```

## Troubleshooting

### Common Issues

1. **Chroma connection failed**
   ```
   Error: failed to connect to Chroma
   ```
   - Ensure Chroma is running on the configured URL
   - Check firewall settings
   - Verify CHROMA_URL environment variable

2. **API key errors**
   ```
   Error: unauthorized API request
   ```
   - Verify API keys are correctly set in .env
   - Check API key permissions and quotas

3. **Database locked**
   ```
   Error: database is locked
   ```
   - Ensure no other instances are running
   - Check file permissions on data directory

4. **Large file uploads**
   - Default max upload size is 32MB
   - Adjust server configuration if needed

### Logs

Application logs are written to stdout with structured JSON format:
```bash
./bin/rag-therapist 2>&1 | jq .
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run linting and tests
6. Submit a pull request

## License

[Add your license information here]