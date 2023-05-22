package secret

type UserSecretItemTypeID int

// взять любые два типа.
const (
	CredentialsType UserSecretItemTypeID = iota
	Binary
)

// тесты на все контроллеры сервера

type IUserSecretItem interface {
	GetID() string
	GetType() string
}

// graceful shutdown on the server and on the client
// дождаться пока запишется пароль \выполнится операция и позволить клиенту выйти
// context with timeout

// пришел сигнал от ОС, убеждаемся что никакия команда не выполняется.

// ./client - secretKey=file.key -config=config.json -user -password
// ->  list
// get
// set
// delete
