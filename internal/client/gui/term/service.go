package term

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/itksb/go-secret-keeper/internal/client/command"
	"github.com/itksb/go-secret-keeper/internal/client/session"
	"github.com/itksb/go-secret-keeper/pkg/contract"
	"io"
	"log"
	"strconv"
	"strings"
)

var completer = readline.NewPrefixCompleter(
	readline.PcItem("login"),
	readline.PcItem("help"),
	readline.PcItem("exit"),
)

// TerminalService - terminal service
type TerminalService struct {
	l              contract.IApplicationLogger
	session        session.ISession
	loginCmdFabric loginCmdFabric
}

type loginCmdFabric func(
	session session.ISession,
	login string,
	password string,
) *command.LoginCommand

// NewTerminalService - create new terminal service
func NewTerminalService(
	l contract.IApplicationLogger,
	session session.ISession,
	loginCmdFabric loginCmdFabric,
) *TerminalService {
	return &TerminalService{
		l:              l,
		session:        session,
		loginCmdFabric: loginCmdFabric,
	}
}

func usage(w io.Writer) {
	io.WriteString(w, "commands:\n")
	io.WriteString(w, completer.Tree("    "))
}

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func (s *TerminalService) Start() error {

	const nemoPrompt = "\u001B[31m(nemo)»\u001B[0m "
	createLoginPrompt := func(login string) string {
		return fmt.Sprintf("\u001B[31m%s»\u001B[0m ", login)
	}

	rline, err := readline.NewEx(&readline.Config{
		Prompt:                 nemoPrompt,
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

	rline.CaptureExitSignal()

	setPasswordCfg := rline.GenPasswordConfig()
	setPasswordCfg.SetListener(func(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
		rline.SetPrompt(fmt.Sprintf("Enter password(%v): ", len(line)))
		rline.Refresh()
		return nil, 0, false
	})

	for {
		line, rlErr := rline.Readline()
		if s.isAuthorized() {
			rline.SetPrompt(createLoginPrompt(s.GetLogin()))
		} else {
			rline.SetPrompt(nemoPrompt)
		}
		if rlErr == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if rlErr == io.EOF {
			break
		}

		line = strings.TrimSpace(line)

		switch {
		case line == "login":
			rline.SetPrompt(createLoginPrompt("enter login fot the server"))
			login, err := rline.Readline()
			if err != nil {
				rline.Write([]byte(fmt.Sprintf("error while reading login: %s \r\n", err.Error())))
			}
			rline.SetPrompt(nemoPrompt)
			rline.Write([]byte(fmt.Sprintf("your login: %s \r\n", login)))
			var pswd []byte
			pswd, err = rline.ReadPasswordWithConfig(setPasswordCfg)
			if err != nil {
				rline.Write([]byte(fmt.Sprintf("error while reading password: %s \r\n", err.Error())))
			}
			rline.Write([]byte(fmt.Sprintf("your password: %s \r\n", string(pswd))))
			loginCmd := s.loginCmdFabric(s.session, login, string(pswd))
			err = loginCmd.Execute()
			if err != nil {
				rline.Write([]byte(fmt.Sprintf("error while reading password: %s \r\n", err.Error())))
			}

		case line == "help":
			usage(rline.Stderr())
		case line == "exit":
			goto exit
		case line == "":
		default:
			log.Println("unknown command: ", strconv.Quote(line))
		}
	}

exit:
	return nil
}

// isAuthorized - check if user is authorized
func (s *TerminalService) isAuthorized() bool {
	_, err := s.session.GetValue(session.TokenKey)
	if err != nil {
		return false
	}
	return true
}

// GetLogin - get login from session
func (s *TerminalService) GetLogin() string {
	accRaw, err := s.session.GetValue(session.AccountKey)
	if err != nil {
		return ""
	}
	acc, ok := accRaw.(contract.IAccount)
	if !ok {
		return ""
	}

	return acc.GetLogin()
}
