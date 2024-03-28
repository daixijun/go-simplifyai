package simplifyai

type CreateTranslationTaskRequest struct {
	FromLang string // 源语言
	ToLang   string // 目标语言
	Glossary string // 术语表
	File     string // 文件路径
}

// CreateTranslationTaskResponse 创建翻译任务响应
type CreateTranslationTaskResponse struct {
	TaskId string `json:"taskId"` // 任务 ID, 创建任务时返回
	Price  int    `json:"price"`  // 此次翻译使用的积分
}

// TaskStatus 任务状态
type TaskStatus string

const (
	// TaskStatusUnpaid 未支付
	TaskStatusUnpaid TaskStatus = "Unpaid"
	// TaskStatusWaiting 等待中
	TaskStatusWaiting TaskStatus = "Waiting"
	// TaskStatusProcessing 处理中
	TaskStatusProcessing TaskStatus = "Processing"
	// TaskStatusCompleted 已完成
	TaskStatusCompleted TaskStatus = "Completed"
	// TaskStatusCancelled 已取消
	TaskStatusCancelled TaskStatus = "Cancelled"
	// TaskStatusTerminated 已终止
	TaskStatusTerminated TaskStatus = "Terminated"
	// TaskStatusNotSupported 不支持
	TaskStatusNotSupported TaskStatus = "NotSupported"
)

// QueryTranslationTaskResponse 查询翻译任务返回
type QueryTranslationTaskResponse struct {
	Status            TaskStatus `json:"status"`            // 任务状态
	Progress          int        `json:"progress"`          // 任务进度，0 到 100
	TranslatedFileURL string     `json:"translatedFileUrl"` // 译文文档的 URL，翻译完成才有。
	BilingualFileURL  string     `json:"bilingualFileUrl"`  // 双语文档的 URL，翻译完成才有。
	Price             int        `json:"price"`             // 此次翻译使用的积分
	TotalToken        int        `json:"totalToken"`        // 需要翻译的 token 总数
}

// StartTranslationTaskResponse 启动翻译任务返回
type StartTranslationTaskResponse struct {
	Status     TaskStatus `json:"status"`     // 任务状态
	Progress   int        `json:"progress"`   // 任务进度, 0 到 100
	Price      int        `json:"price"`      // 此次翻译使用的积分
	TotalToken int        `json:"totalToken"` // 需要翻译的 Token 总数
}
