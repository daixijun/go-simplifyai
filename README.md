# Golang sdk for [SimplifyAI](https://translate.simplifyai.cn/)

SimplifyAI 接口文档地址: <https://translate.simplifyai.cn/developer>

## 使用示例

Install

`go get github.com/daixijun/go-simplifyai`

使用:

```go
package main

import (
 "context"
 "log"
 "os"

 "github.com/daixijun/simplifyai"
)

var (
 apiKey   = os.Getenv("SIMPLIFYAI_API_KEY")
 fromLang = "English"
 toLang   = "Simplified Chinese"
 glossary = ""
 filePath = "source.pdf"
)

func main() {
    // 创建 Client
    client := simplifyai.NewClient(apiKey, nil)
    ctx := context.Background()
    req := &simplifyai.CreateTranslationTaskRequest{
    FromLang: fromLang,
    ToLang:   toLang,
    Glossary: glossary,
    File:     filePath,
    }

    // 创建翻译任务
    task, err := client.CreateTranslationTask(ctx, req)
    if err != nil {
    log.Fatalf("CreateTranslationTask() error = %v", err)
    }

    // 查询翻译任务
    res, err = client.QueryTranslationTask(ctx, task.TaskId)
    if err != nil {
        log.Fatalf("QueryTranslationTask() error = %v", err)
    }

    // 删除翻译任务
    res, err = client.DeleteTranslationTask(ctx, task.TaskId)
    if err != nil {
        log.Fatalf("DeleteTranslationTask() error = %v", err)
    }

    // 启动翻译任务
    res, err = client.StartTranslationTask(ctx, task.TaskId)
    if err != nil {
        log.Fatalf("StartTranslationTask() error = %v", err)
    }

    // 查询可用语言
    langs, err = client.ListAvailableLanguages(ctx)
    if err != nil {
        log.Fatalf("ListAvailableLanguages() error = %v", err)
    }
}

```
