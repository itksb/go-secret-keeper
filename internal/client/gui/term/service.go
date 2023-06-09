package term

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/itksb/go-secret-keeper/internal/client/command"
	secret2 "github.com/itksb/go-secret-keeper/internal/client/keeper/secret"
	"github.com/itksb/go-secret-keeper/internal/client/session"
	"github.com/itksb/go-secret-keeper/pkg/contract"
	"github.com/itksb/go-secret-keeper/pkg/keeper/secret"
	"io"
	"log"
	"strconv"
	"strings"
)

var completer = readline.NewPrefixCompleter(
	readline.PcItem("login"),
	readline.PcItem("register"),
	readline.PcItem("list"),
	readline.PcItem("delete"),
	readline.PcItem("savetext"),
	readline.PcItem("saveloginpass"),
	readline.PcItem("help"),
	readline.PcItem("exit"),
)

// TerminalService - terminal service
type TerminalService struct {
	l                      contract.IApplicationLogger
	session                session.ISession
	loginCmdFabric         loginCmdFabric
	registerCmdFabric      registerCmdFabric
	listSecretsCmdFabric   listSecretsCmdFabric
	deleteSecretCmdFabric  deleteSecretCmdFabric
	saveLoginPassCmdFabric saveLoginPassCmdFabric
	saveTextCmdFabric      saveTextCmdFabric
}

type loginCmdFabric func(
	session session.ISession,
	login string,
	password string,
) *command.LoginCommand

type registerCmdFabric func(
	session session.ISession,
	login string,
	password string,
) *command.RegisterCommand

type listSecretsCmdFabric func(
	userID string,
	processFunc command.SecretsProcessorFunc,
) *command.ListSecretsCommand

type deleteSecretCmdFabric func(
	userID string,
	secretID string,
) *command.DeleteSecretCommand

type saveLoginPassCmdFabric func(
	userID string,
	secretID string,
	login string,
	password string,
	meta string,
	secretItemProcessorFunc command.SecretItemProcessorFunc,
) *command.SaveLoginPassCommand

type saveTextCmdFabric func(
	userID string,
	secretID string,
	text string,
	meta string,
	secretItemProcessorFunc command.SecretItemProcessorFunc,
) *command.SaveTextCommand

// NewTerminalService - create new terminal service
func NewTerminalService(
	l contract.IApplicationLogger,
	session session.ISession,
	loginCmdFabric loginCmdFabric,
	registerCmdFabric registerCmdFabric,
	listSecretsCmdFabric listSecretsCmdFabric,
	deleteSecretCmdFabric deleteSecretCmdFabric,
	saveLoginPassCmdFabric saveLoginPassCmdFabric,
	saveTextCmdFabric saveTextCmdFabric,
) *TerminalService {
	return &TerminalService{
		l:                      l,
		session:                session,
		loginCmdFabric:         loginCmdFabric,
		registerCmdFabric:      registerCmdFabric,
		listSecretsCmdFabric:   listSecretsCmdFabric,
		deleteSecretCmdFabric:  deleteSecretCmdFabric,
		saveLoginPassCmdFabric: saveLoginPassCmdFabric,
		saveTextCmdFabric:      saveTextCmdFabric,
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

		var cmdErr error

		switch {
		case line == "login":
			cmdErr = nil
			rline.SetPrompt(createLoginPrompt("enter login fot the server"))
			var login string
			login, cmdErr = rline.Readline()
			if err != nil {
				rline.Write([]byte(fmt.Sprintf("error while reading login: %s \r\n", cmdErr.Error())))
				continue
			}
			rline.SetPrompt(nemoPrompt)
			rline.Write([]byte(fmt.Sprintf("your login: %s \r\n", login)))
			cmdErr = nil
			var pswd []byte
			pswd, cmdErr = rline.ReadPasswordWithConfig(setPasswordCfg)
			if cmdErr != nil {
				rline.Write([]byte(fmt.Sprintf("error while reading password: %s \r\n", cmdErr.Error())))
				continue
			}
			loginCmd := s.loginCmdFabric(s.session, login, string(pswd))
			cmdErr = loginCmd.Execute()
			if cmdErr != nil {
				rline.Write([]byte(fmt.Sprintf("error while executing login command: %s \r\n", cmdErr.Error())))
			}

		case line == "register":
			cmdErr = nil
			rline.SetPrompt(createLoginPrompt("enter login to register on the server"))
			var login string
			login, cmdErr = rline.Readline()
			if err != nil {
				rline.Write([]byte(fmt.Sprintf("error while reading login: %s \r\n", cmdErr.Error())))
				continue
			}
			rline.SetPrompt(nemoPrompt)
			rline.Write([]byte(fmt.Sprintf("your login: %s \r\n", login)))
			cmdErr = nil
			var pswd []byte
			pswd, cmdErr = rline.ReadPasswordWithConfig(setPasswordCfg)
			if cmdErr != nil {
				rline.Write([]byte(fmt.Sprintf("error while executing register command: %s \r\n", cmdErr.Error())))
				continue
			}
			registerCmd := s.registerCmdFabric(s.session, login, string(pswd))
			cmdErr = registerCmd.Execute()
			if cmdErr != nil {
				rline.Write([]byte(fmt.Sprintf("error while reading password: %s \r\n", cmdErr.Error())))
			}
		case line == "list":
			cmdErr = nil
			acc := s.GetAccount()
			if acc == nil {
				continue
			}
			listSecretsCmd := s.listSecretsCmdFabric(
				acc.GetID(),
				func(secrets []contract.IUserSecretItem) error {
					for _, sec := range secrets {
						switch sec.GetType() {
						case contract.UserSecretLoginPasswd:
							s.renderUserSecretLoginPasswd(sec)
						case contract.UserSecretTypeTextData:
						default:
							rline.Write([]byte(fmt.Sprintf("unknown secret type: %s \r\n", sec.GetType())))
							continue
						}
					}
					return nil
				})

			cmdErr = listSecretsCmd.Execute()
			if cmdErr != nil {
				rline.Write([]byte(fmt.Sprintf("error while executing list command: %s \r\n", cmdErr.Error())))
			}
		case line == "delete":
			cmdErr = nil
			acc := s.GetAccount()
			if acc == nil {
				continue
			}
			rline.SetPrompt(createLoginPrompt("enter secret id to delete"))
			var secretID string
			secretID, cmdErr = rline.Readline()
			if cmdErr != nil {
				rline.Write([]byte(fmt.Sprintf("error while reading secret id: %s \r\n", cmdErr.Error())))
				continue
			}
			rline.SetPrompt(createLoginPrompt(acc.GetLogin()))
			deleteCmd := s.deleteSecretCmdFabric(acc.GetID(), secretID)
			cmdErr = deleteCmd.Execute()
			if cmdErr != nil {
				rline.Write([]byte(fmt.Sprintf("error while executing delete command: %s \r\n", cmdErr.Error())))
			}

		case line == "savetext":
			cmdErr = nil
			acc := s.GetAccount()
			if acc == nil {
				continue

			}
			rline.SetPrompt(createLoginPrompt("enter secret id to update or leave empty to create"))
			var secretID string
			secretID, cmdErr = rline.Readline()
			if cmdErr != nil {
				rline.Write([]byte(fmt.Sprintf("error while reading secret id: %s \r\n", cmdErr.Error())))
				continue
			}

			rline.SetPrompt(createLoginPrompt("enter text to save or update"))
			var text string
			text, cmdErr = rline.Readline()
			if cmdErr != nil {
				rline.Write([]byte(fmt.Sprintf("error while reading text: %s \r\n", cmdErr.Error())))
				continue
			}

			cmdErr = nil

			rline.SetPrompt(createLoginPrompt("enter meta to save or update"))
			var meta string
			meta, cmdErr = rline.Readline()
			if err != nil {
				rline.Write([]byte(fmt.Sprintf("error while reading meta: %s \r\n", cmdErr.Error())))
				continue
			}
			rline.SetPrompt(acc.GetLogin())
			rline.Write([]byte(fmt.Sprintf("your meta: %s \r\n", meta)))

			saveCmd := s.saveTextCmdFabric(
				acc.GetID(),
				secretID,
				text,
				meta,
				func(item contract.IUserSecretItem) error {
					if string(item.GetID()) != secretID {
						rline.Write([]byte(fmt.Sprintf("secret was created: %s \r\n", item.GetID())))
						return nil
					} else {
						rline.Write([]byte(fmt.Sprintf("secret was updated: %s \r\n", item.GetID())))
						return nil
					}
				},
			)

			cmdErr = saveCmd.Execute()
			if cmdErr != nil {
				rline.Write([]byte(fmt.Sprintf("error while executing save command: %s \r\n", cmdErr.Error())))
			}

		case line == "saveloginpass":
			cmdErr = nil
			acc := s.GetAccount()
			if acc == nil {
				continue

			}
			rline.SetPrompt(createLoginPrompt("enter secret id to update or leave empty to create"))
			var secretID string
			secretID, cmdErr = rline.Readline()
			if cmdErr != nil {
				rline.Write([]byte(fmt.Sprintf("error while reading secret id: %s \r\n", cmdErr.Error())))
				continue
			}

			rline.SetPrompt(createLoginPrompt("enter login to save"))
			var login string
			login, cmdErr = rline.Readline()
			if err != nil {
				rline.Write([]byte(fmt.Sprintf("error while reading login: %s \r\n", cmdErr.Error())))
				continue
			}
			rline.SetPrompt(acc.GetLogin())
			rline.Write([]byte(fmt.Sprintf("your login: %s \r\n", login)))

			cmdErr = nil
			var pswd []byte
			pswd, cmdErr = rline.ReadPasswordWithConfig(setPasswordCfg)
			if cmdErr != nil {
				rline.Write([]byte(fmt.Sprintf("error while reading password: %s \r\n", cmdErr.Error())))
				continue
			}

			rline.SetPrompt(createLoginPrompt("enter meta to save or update"))
			var meta string
			meta, cmdErr = rline.Readline()
			if err != nil {
				rline.Write([]byte(fmt.Sprintf("error while reading meta: %s \r\n", cmdErr.Error())))
				continue
			}
			rline.SetPrompt(acc.GetLogin())
			rline.Write([]byte(fmt.Sprintf("your meta: %s \r\n", meta)))

			saveCmd := s.saveLoginPassCmdFabric(
				acc.GetID(),
				secretID,
				login,
				string(pswd),
				meta,
				func(item contract.IUserSecretItem) error {
					if string(item.GetID()) != secretID {
						rline.Write([]byte(fmt.Sprintf("secret was created: %s \r\n", item.GetID())))
						return nil
					} else {
						rline.Write([]byte(fmt.Sprintf("secret was updated: %s \r\n", item.GetID())))
						return nil
					}
				},
			)

			cmdErr = saveCmd.Execute()
			if cmdErr != nil {
				rline.Write([]byte(fmt.Sprintf("error while executing save command: %s \r\n", cmdErr.Error())))
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

func (s *TerminalService) GetAccount() contract.IAccount {
	accRaw, err := s.session.GetValue(session.AccountKey)
	if err != nil {
		s.l.Errorf("error while getting account from session: %s", err.Error())
		return nil
	}
	acc, ok := accRaw.(contract.IAccount)
	if !ok {
		s.l.Errorf("error while getting account from session: %s", err.Error())
		return nil
	}

	return acc
}

func (s *TerminalService) renderUserSecretLoginPasswd(sec contract.IUserSecretItem) {
	item, ok := sec.(*secret.LoginPasswdSecretItem)
	if !ok {
		s.l.Errorf("error while casting secret to LoginPasswdSecretItem")
		return
	}
	packer := secret2.LoginPasswdSecretItemPacker{Entity: *item}
	login, passwd, err := packer.Read()
	if err != nil {
		s.l.Errorf("error while reading secret: %s", err.Error())
		return
	}
	fmt.Printf(
		"id: %s \nlogin: %s \npassword: %s \nmeta: %s",
		sec.GetID(),
		login,
		passwd,
		sec.GetMeta(),
	)
}

func (s *TerminalService) renderUserSecretTypeTextData(sec contract.IUserSecretItem) {
	item, ok := sec.(*secret.TextDataSecretItem)
	if !ok {
		s.l.Errorf("error while casting secret to LoginPasswdSecretItem")
		return
	}
	packer := secret2.TextDataSecretItemPacker{Entity: *item}
	text, err := packer.Read()
	if err != nil {
		s.l.Errorf("error while reading secret: %s", err.Error())
		return
	}
	fmt.Printf(
		"id: %s \ntext: %s \nmeta: %s",
		sec.GetID(),
		text,
		sec.GetMeta(),
	)
}
