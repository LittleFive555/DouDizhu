package errordef

type GameError struct {
	Category     ErrorCategory
	Code         ErrorCode
	ClientMsg    string
	InternalMsg  string
	WrappedError error
}

func (e *GameError) Error() string {
	return e.InternalMsg
}

func AsGameError(err error) *GameError {
	if err == nil {
		return nil
	}
	if e, ok := err.(*GameError); ok {
		return e
	}
	return NewUnknownError(err)
}

func NewUnknownError(wrappedError error) *GameError {
	return newGameErrorWithWrappedError(CategoryUnknown, CodeUnknown, wrappedError)
}

func NewDatabaseError(code ErrorCode, wrappedError error) *GameError {
	return newGameErrorWithWrappedError(CategoryDatabase, code, wrappedError)
}

func NewNetworkError(code ErrorCode, wrappedError error) *GameError {
	return newGameErrorWithWrappedError(CategoryNetwork, code, wrappedError)
}

func NewThirdPartyError(code ErrorCode, wrappedError error) *GameError {
	return newGameErrorWithWrappedError(CategoryThirdParty, code, wrappedError)
}

func NewGameplayError(code ErrorCode) *GameError {
	return newGameError(CategoryGameplay, code, "")
}

func NewGameplayErrorWithMsg(code ErrorCode, clientMsg string) *GameError {
	return newGameError(CategoryGameplay, code, clientMsg)
}

func newGameError(category ErrorCategory, code ErrorCode, clientMsg string) *GameError {
	return &GameError{
		Category:  category,
		Code:      code,
		ClientMsg: clientMsg,
	}
}

func newGameErrorWithWrappedError(category ErrorCategory, code ErrorCode, wrappedError error) *GameError {
	return &GameError{
		Category:     category,
		Code:         code,
		ClientMsg:    "Server Internal Error.",
		InternalMsg:  wrappedError.Error(),
		WrappedError: wrappedError,
	}
}
