package bizerr

// ErrorCode 定义错误码类型为字符串，使用大驼峰命名法
type ErrorCode string

// 定义常见的错误码
const (
	Unknown       ErrorCode = "Unknown"
	NotFound      ErrorCode = "NotFound"
	InvalidInput  ErrorCode = "InvalidInput"
	Unauthorized  ErrorCode = "Unauthorized"
	Forbidden     ErrorCode = "Forbidden"
	InternalError ErrorCode = "InternalError"
	RateLimited   ErrorCode = "RateLimited"
	// 可以根据需要添加更多错误码
)

// 预定义一些常用错误
var (
	ErrNotFound      = NewBusinessError(NotFound, "")
	ErrInvalidInput  = NewBusinessError(InvalidInput, "")
	ErrUnauthorized  = NewBusinessError(Unauthorized, "")
	ErrForbidden     = NewBusinessError(Forbidden, "")
	ErrInternalError = NewBusinessError(InternalError, "")
	RateLimitedError = NewBusinessError(RateLimited, "")
)
