package command

// ICommand - command interface
type ICommand interface {
	Execute() error
}
