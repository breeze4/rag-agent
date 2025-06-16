package config

import (
	"log/slog"
	"os"
	"strconv"
)

type Config struct {
	LLMProvider    string
	ClaudeAPIKey   string
	GeminiAPIKey   string
	OpenAIAPIKey   string
	Port           int
	ChromaURL      string
	DBPath         string
	UploadDir      string
}

func Load() *Config {
	portStr := getEnv("PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		slog.Error("Invalid PORT value", "port", portStr, "error", err)
		port = 8080
	}

	config := &Config{
		LLMProvider:    getEnv("LLM_PROVIDER", "claude"),
		ClaudeAPIKey:   getEnv("CLAUDE_API_KEY", ""),
		GeminiAPIKey:   getEnv("GEMINI_API_KEY", ""),
		OpenAIAPIKey:   getEnv("OPENAI_API_KEY", ""),
		Port:           port,
		ChromaURL:      getEnv("CHROMA_URL", "http://localhost:8000"),
		DBPath:         getEnv("DB_PATH", "./data/rag.db"),
		UploadDir:      getEnv("UPLOAD_DIR", "./data/uploads"),
	}

	slog.Info("Configuration loaded", 
		"llm_provider", config.LLMProvider,
		"port", config.Port,
		"chroma_url", config.ChromaURL,
		"db_path", config.DBPath,
		"upload_dir", config.UploadDir,
	)

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}