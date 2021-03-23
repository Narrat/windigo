package err

const (
	E_UNEXPECTED           ERROR = 0x8000ffff
	E_NOTIMPL              ERROR = 0x80004001
	E_OUTOFMEMORY          ERROR = 0x8007000e
	E_INVALIDARG           ERROR = 0x80070057
	E_NOINTERFACE          ERROR = 0x80004002
	E_POINTER              ERROR = 0x80004003
	E_HANDLE               ERROR = 0x80070006
	E_ABORT                ERROR = 0x80004004
	E_FAIL                 ERROR = 0x80004005
	E_ACCESSDENIED         ERROR = 0x80070005
	E_PENDING              ERROR = 0x8000000a
	E_BOUNDS               ERROR = 0x8000000b
	E_CHANGED_STATE        ERROR = 0x8000000c
	E_ILLEGAL_STATE_CHANGE ERROR = 0x8000000d
	E_ILLEGAL_METHOD_CALL  ERROR = 0x8000000e
)