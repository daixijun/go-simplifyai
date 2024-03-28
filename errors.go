package simplifyai

import (
	"errors"
)

var (
	ErrUnauthorized     = errors.New("没有提供正确的 API 密钥") // 401
	ErrInvalidParameter = errors.New("参数错误")           // 400
	ErrInternalError    = errors.New("服务器端内部错误")       // 500
	ErrNotFound         = errors.New("指定的任务不存在")       // 404
	ErrUnpaid           = errors.New("积分余额不足")         // 402
)

var ErrorMap = map[int]error{
	400: ErrInvalidParameter,
	401: ErrUnauthorized,
	402: ErrUnpaid,
	404: ErrNotFound,
	500: ErrInternalError,
}
