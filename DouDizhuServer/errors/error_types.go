package errors

type ErrorCategory int

const (
	CategoryUnknown ErrorCategory = 0
	// 技术类错误
	CategoryDatabase   ErrorCategory = 1000
	CategoryNetwork    ErrorCategory = 2000
	CategoryThirdParty ErrorCategory = 3000
	// 业务类错误
	CategoryGameplay ErrorCategory = 10000
)

type ErrorCode string

const (
	CodeUnknown ErrorCode = "Unknown"
	// 通用错误
	CodeInvalidRequest ErrorCode = "INVALID_REQUEST"
	// Database
	CodeDBReadError  ErrorCode = "DB_READ_ERROR"
	CodeDBWriteError ErrorCode = "DB_WRITE_ERROR"
	// Network
	CodeNetworkReadError  ErrorCode = "NETWORK_READ_ERROR"
	CodeNetworkWriteError ErrorCode = "NETWORK_WRITE_ERROR"

	// Register
	CodeAccountExists         ErrorCode = "ACCOUNT_EXISTS"
	CodeAccountTooLong        ErrorCode = "ACCOUNT_TOO_LONG"
	CodeAccountTooShort       ErrorCode = "ACCOUNT_TOO_SHORT"
	CodeAccountFormatInvalid  ErrorCode = "ACCOUNT_FORMAT_INVALID"
	CodePasswordTooLong       ErrorCode = "PASSWORD_TOO_LONG"
	CodePasswordTooShort      ErrorCode = "PASSWORD_TOO_SHORT"
	CodePasswordFormatInvalid ErrorCode = "PASSWORD_FORMAT_INVALID"

	// Login
	CodeAccountNotExists ErrorCode = "ACCOUNT_NOT_EXISTS"
	CodePasswordWrong    ErrorCode = "PASSWORD_WRONG"
)
