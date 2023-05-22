package contract

type IApplicationLogger interface {
	Infof(template string, args ...interface{})
	Errorf(template string, args ...interface{})
}
