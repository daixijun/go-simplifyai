package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/daixijun/go-simplifyai"
)

var (
	apiKey   = os.Getenv("SIMPLIFYAI_API_KEY")
	fromLang = "English"
	toLang   = "Simplified Chinese"
	glossary = ""
	filePath = "source.pdf"
)

func main() {
	client := simplifyai.NewClient(apiKey, nil)
	ctx := context.Background()
	req := &simplifyai.CreateTranslationTaskRequest{
		FromLang: fromLang,
		ToLang:   toLang,
		Glossary: glossary,
		File:     filePath,
	}
	task, err := client.CreateTranslationTask(ctx, req)
	if err != nil {
		log.Fatalf("CreateTranslationTask() error = %v", err)
	}
	log.Printf("task: %+v", task)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	var res *simplifyai.QueryTranslationTaskResponse
	for range ticker.C {
		res, err = client.QueryTranslationTask(ctx, task.TaskId)
		if err != nil {
			log.Fatalf("QueryTranslationTask() error = %v", err)
		}

		if res.Progress == 100 {
			break
		}
	}
	log.Printf("res: %+v", res)
}
