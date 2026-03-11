package cli

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/aboubakary833/codexa/internal/ports"
	"github.com/spf13/cobra"
)

type CommandWrapper struct {
	rootCmd    *cobra.Command
	controller controller
	appVersion string
}

func NewCommandWrapper(
	application ports.Application,
	appVersion	string,
	errLogger *slog.Logger,
	) CommandWrapper {
	rootCmd := &cobra.Command{
		Use:   "codexa",
		Short: "Codexa is a concise and descriptive snippets app.",
		Long:  "A simple terminal-based app designed to help devs quickly access pratical snippets.",
	}

	rootCmd.AddCommand(&cobra.Command{
        Use:   "version",
        Short: "Print the version number of Codexa",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Codexa version:", appVersion)
        },
    })

	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	cw := CommandWrapper{
		controller: controller{
			app: application,
			logger: errLogger,
		},
		rootCmd: rootCmd,
	}

	// Register available commands
	cw.registerRunCmd()
	cw.registerOpenCmd()
	cw.registerSyncCmd()

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
		Short: "Open a snippet or a list of snippets of category x",
		Long:  "Open a given snippet content or list snippets for a given tech category",
		Example: `	codexa open php
	codexa open js debounce
	codexa open -c=javascript
	codexa open -c go -t context-timeout
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
	cmd.Flags().StringVarP(&category, "category", "c", "", "specify tech category")
	cmd.Flags().StringVarP(&topic, "topic", "t", "", "specify snippet topic")

	cw.rootCmd.AddCommand(cmd)
}

// registerSyncCmd register the command for downloading/updating snippets
func (cw CommandWrapper) registerSyncCmd() {
	var category, id string

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync add or update a snippet or a set of snippets",
		Long:  "Sync create or update a snippet or a given tech category snippets",
		Example: `	codexa sync javascript
	codexa sync -c go
	codexa sync -c=js -id=debounce
	codexa sync --category go
	codexa sync --category=js --identifier=deep-clone`,
		Args: cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if category == "" && len(args) >= 1 {
				category = args[0]
			}

			if category == "" {
				cmd.Usage()
				return
			}

			if id == "" && len(args) == 2 {
				id = args[1]
			}

			if id == "" {
				cw.Exec(func() error {
					return cw.controller.syncTechCategory(category)
				})
				return
			}

			cw.Exec(func() error {
				return cw.controller.syncSnippet(category, id)
			})
		},
	}
	cmd.Flags().StringVarP(&category, "category", "c", "", "specify the tech category")
	cmd.Flags().StringVarP(&id, "identifier", "i", "", "specify the snippet identifier")

	cw.rootCmd.AddCommand(cmd)
}

func (cw CommandWrapper) Exec(fn func() error) {
	err := fn()

	if err == nil {
		return
	}

	fmt.Println(err.Error())
	os.Exit(1)
}
