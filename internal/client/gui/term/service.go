package term

import (
	"github.com/chzyer/readline"
	"github.com/itksb/go-secret-keeper/pkg/contract"
)

type TerminalService struct {
	l contract.IApplicationLogger
}

// NewTerminalService - create new terminal service
func NewTerminalService(
	l contract.IApplicationLogger,
) *TerminalService {
	return &TerminalService{
		l: l,
	}
}

var completer = readline.NewPrefixCompleter(
	readline.PcItem("login"),
	readline.PcItem("help"),
)

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func (s *TerminalService) Start() error {

	rline, err := readline.NewEx(&readline.Config{
		Prompt:                 "\u001B[31m»\u001B[0m ",
		HistoryLimit:           0,
		DisableAutoSaveHistory: true,
		AutoComplete:           completer,
		InterruptPrompt:        "^C",
		EOFPrompt:              "exit",
		FuncFilterInputRune:    filterInput,
	})

	if err != nil {
		s.l.Errorf("error while starting readline %s", err.Error())
		return err
	}

	defer func() {
		errClose := rline.Close()
		if errClose != nil {
			s.l.Errorf("error while closing readline %s", errClose.Error())
		}
	}()

	return nil
}
