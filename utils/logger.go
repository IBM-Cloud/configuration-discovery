package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
)

// ContextKeyLogger -
var ContextKeyLogger = contextKey("logger")

// Mock up an implementation for terminal.UI and use log package implementation instead

// GetUI : Returns the cli ui logger from context
func GetLogger(ctx context.Context) terminal.UI {
	logger := ctx.Value(ContextKeyLogger)

	if cliLogger, ok := logger.(terminal.UI); ok {
		return cliLogger
	}
	return newlogMocker()
}

type logMocker struct {
	In     io.Reader
	Out    io.Writer
	ErrOut io.Writer
}

// NewStdUI initialize a terminal UI with os.Stdin and os.Stdout
func newlogMocker() terminal.UI {
	return &logMocker{
		In:     os.Stdin,
		Out:    terminal.Output,
		ErrOut: terminal.ErrOutput,
	}
}

func (l *logMocker) Say(format string, args ...interface{}) {
	if args != nil {
		log.Printf(format+"\n", args...)
	} else {
		log.Print(format + "\n")
	}
}

func (l *logMocker) Warn(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.Say(terminal.WarningColor(message))
}

func (l *logMocker) Ok() {
	l.Say(terminal.SuccessColor("OK"))
}

// todo: @srikar - change this
func (l *logMocker) Error(format string, args ...interface{}) {
	if args != nil {
		log.Printf(format+"\n", args...)
	} else {
		log.Print(format + "\n")
	}
}

func (l *logMocker) Failed(format string, args ...interface{}) {
	l.Error(terminal.FailureColor("FAILED"))
	l.Error(format, args...)
	l.Error("")
}

func (l *logMocker) Prompt(message string, options *terminal.PromptOptions) *terminal.Prompt {
	p := terminal.NewPrompt(message, options)
	p.Reader = l.In
	p.Writer = l.Out
	return p
}

func (l *logMocker) ChoicesPrompt(message string, choices []string, options *terminal.PromptOptions) *terminal.Prompt {
	p := terminal.NewChoicesPrompt(message, choices, options)
	p.Reader = l.In
	p.Writer = l.Out
	return p
}

func (l *logMocker) Ask(format string, args ...interface{}) (answer string, err error) {
	message := fmt.Sprintf(format, args...)
	err = l.Prompt(message, &terminal.PromptOptions{HideDefault: true, NoLoop: true}).Resolve(&answer)
	return
}

func (l *logMocker) AskForPassword(format string, args ...interface{}) (passwd string, err error) {
	message := fmt.Sprintf(format, args...)
	err = l.Prompt(message, &terminal.PromptOptions{HideInput: true, HideDefault: true, NoLoop: true}).Resolve(&passwd)
	return
}

func (l *logMocker) Confirm(format string, args ...interface{}) (yn bool, err error) {
	message := fmt.Sprintf(format, args...)
	err = l.Prompt(message, &terminal.PromptOptions{HideDefault: true, NoLoop: true}).Resolve(&yn)
	return
}

func (l *logMocker) ConfirmWithDefault(defaultBool bool, format string, args ...interface{}) (yn bool, err error) {
	yn = defaultBool
	message := fmt.Sprintf(format, args...)
	err = l.Prompt(message, &terminal.PromptOptions{HideDefault: true, NoLoop: true}).Resolve(&yn)
	return
}

func (l *logMocker) SelectOne(choices []string, format string, args ...interface{}) (int, error) {
	var selected string
	message := fmt.Sprintf(format, args...)

	err := l.ChoicesPrompt(message, choices, &terminal.PromptOptions{HideDefault: true}).Resolve(&selected)
	if err != nil {
		return -1, err
	}

	for i, c := range choices {
		if selected == c {
			return i, nil
		}
	}

	return -1, nil
}

func (l *logMocker) Table(headers []string) terminal.Table {
	return terminal.NewTable(l.Out, headers)
}

func (l *logMocker) Writer() io.Writer {
	return l.Out
}
