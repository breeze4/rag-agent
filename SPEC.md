# RAG Chatbot Service - Simple Spec

## What it does
A local chatbot that can answer questions using your uploaded PDF documents. You can switch between Claude and Gemini for the AI responses.

## Core Features
1. **Chat API** - Send messages, get responses with document context
2. **PDF Upload** - Upload PDFs, they get processed and stored  
3. **LLM Switching** - Toggle between Claude/Gemini via environment variable 
4. **Document Search** - Find relevant text chunks from your docs

## API Endpoints (MVP)
```
POST /chat
Body: {"message": "your question"}
Response: {"response": "ai answer", "sources": ["doc chunks"]}

POST /upload  
Body: PDF file
Response: {"status": "processing"}
```

## Tech Stack Decisions

### Web Framework
**Options**: Gin, Echo, or stdlib  
**Recommendation**: Gin (simple, fast, good docs)

### Vector Database  
**Options**: Chroma, Qdrant, or in-memory  
**Recommendation**: Chroma (easiest setup, good Go client)

### PDF Processing
**Options**: unidoc, pdfcpu, or external service  
**Recommendation**: unidoc (pure Go, handles most PDFs)

### Embeddings
**Options**: OpenAI API, local sentence-transformers, or Cohere  
**Recommendation**: OpenAI text-embedding-3-small (cheap, reliable)

### Storage
**Options**: SQLite, PostgreSQL, or just files  
**Recommendation**: SQLite (zero-config, good enough)

### LLM Clients
- **Common**: The service should have a thin interface to allow for swappable implementation of LLM functionality. Switch between models with `LLM_PROVIDER` environment variable `gemini` or `claude`
- **Claude**: Official Anthropic Go SDK or HTTP client
- **Gemini**: Google's generativeai-go library

## Simple Architecture
```
HTTP Request → Gin Router → Handler
                    ↓
                RAG Logic (search docs + call LLM)
                    ↓
            Chroma (vectors) + SQLite (metadata)
```

## Environment Config
```
CLAUDE_API_KEY=sk-...
GEMINI_API_KEY=...
OPENAI_API_KEY=sk-... (for embeddings)
PORT=8080
```

## MVP Success Criteria
1. Upload a PDF and ask questions about it
2. Switch between Claude/Gemini models
3. Get responses with relevant document chunks
4. Run locally with one command