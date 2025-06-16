# RAG Chatbot Implementation Tasks

When going through tasks, check each one off as it completes.

## Setup & Dependencies
- [x] Initialize Go module and basic project structure
- [x] Add dependencies: gin, chroma-go, unidoc, sqlite, anthropic-sdk-go, generativeai-go
- [x] Create .env file with API keys and config
- [x] Set up basic logging with slog

## Database & Storage
- [x] Create SQLite schema for documents metadata
- [x] Implement document storage (save uploaded PDFs to disk)
- [x] Create basic database operations (insert, query documents)

## Vector Store
- [x] Set up Chroma client connection
- [x] Create collection for document chunks
- [x] Implement add/search functions for vectors
- [x] Test basic vector operations

## LLM Interface
- [ ] Define common LLM interface (Chat method)
- [ ] Create LLM factory that reads LLM_PROVIDER env var
- [ ] Implement Claude client using Anthropic SDK
- [ ] Implement Gemini client using Google SDK
- [ ] Test both LLM clients work independently

## PDF Processing
- [ ] Implement PDF text extraction using unidoc
- [ ] Create text chunking (fixed-size chunks with overlap)
- [ ] Generate embeddings for chunks using OpenAI API
- [ ] Store chunks in Chroma with metadata

## RAG Logic
- [ ] Implement similarity search function
- [ ] Create RAG pipeline: query → search → context + prompt → LLM
- [ ] Test RAG pipeline with sample documents

## HTTP API
- [ ] Set up Gin server with basic routes
- [ ] Implement POST /upload handler (save PDF, process async)
- [ ] Implement POST /chat handler (RAG pipeline)
- [ ] Add basic error handling and JSON responses

## Integration & Testing
- [ ] Test full workflow: upload PDF → ask question → get response
- [ ] Test switching between Claude and Gemini via env var
- [ ] Add basic validation (file size, PDF format)
- [ ] Create simple README with setup instructions

## Tech debt
- [ ] Replace Chroma with a windows compatible tool