package contract

type IApplication interface {
	Run() error
	Stop() error
}
