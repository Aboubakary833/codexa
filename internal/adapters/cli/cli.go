package cli

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/aboubakary833/codexa/internal/ports"
	"github.com/spf13/cobra"
)

type CommandWrapper struct {
	rootCmd    *cobra.Command
	controller controller
}

func NewCommandWrapper(application ports.Application, errLogger *slog.Logger) CommandWrapper {
	rootCmd := &cobra.Command{
		Use:   "codexa",
		Short: "Codexa is a concise and descriptive snippets app.",
		Long:  "A simple terminal-based app designed to help devs quickly access pratical snippets.",
	}

	cw := CommandWrapper{
		controller: controller{
			app: application,
		},
		rootCmd: rootCmd,
	}

	// Register available commands
	cw.registerRunCmd()
	cw.registerOpenCmd()

	return cw
}

func (cw CommandWrapper) Execute() error {
	return cw.rootCmd.Execute()
}

func (cw CommandWrapper) registerRunCmd() {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Launch Codexa Terminal User Interface",
		Long:  "Open Codexa Terminal User Interface in browse mode",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cw.controller.render(nil)
		},
	}

	cw.rootCmd.AddCommand(cmd)
}

// registerOpenCmd register the command for opening a snippet.
func (cw CommandWrapper) registerOpenCmd() {
	var category, topic string

	cmd := &cobra.Command{
		Use:   "open",
		Short: "Open a snippet or a category list of snippets",
		Long:  "Open a given snippet content or list snippets for a given category",
		Example: `	codexa open php
	codexa open js array
	codexa open -c=typescript
	codexa open -c go -t slices
	codexa open --category javascript
	codexa open --category=go --topic=maps`,
		Args: cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if category == "" && len(args) >= 1 {
				category = args[0]
			}

			if category == "" {
				cmd.Usage()
				return
			}

			if topic == "" && len(args) == 2 {
				topic = args[1]
			}

			if topic == "" {
				cw.Exec(func() error {
					return cw.controller.
						renderTechSnippetsList(category)
				})
				return
			}

			cw.Exec(func() error {
				return cw.controller.
					renderSnippetContent(category, topic)
			})
		},
	}
	cmd.Flags().StringVarP(&category, "category", "c", "", "specify category")
	cmd.Flags().StringVarP(&topic, "topic", "t", "", "specify topic")

	cw.rootCmd.AddCommand(cmd)
}

func (cw CommandWrapper) Exec(fn func() error) {
	err := fn()

	if err == nil {
		return
	}

	if errors.Is(err, internalError) || errors.Is(err, timeoutError) {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(err.Error())
	os.Exit(0)
}
