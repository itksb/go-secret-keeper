package session

// ISession - session interface
type ISession interface {
	GetValue(key string) (interface{}, error)
	SetValue(key string, value interface{}) error
}

// TokenKey - token key
const TokenKey = "token"

// AccoutKey - account key
const AccountKey = "account"
