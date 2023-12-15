package command

import (
	"context"
	"log"
	"os"

	"github.com/spf13/cobra"
	"requester/internal/model"
	"requester/internal/service"
)

const (
	flagURL       = "url"
	flagConfig    = "cfg"
	flagAmount    = "amount"
	flagPerSecond = "per-second"
)

func New(ctx context.Context, r *service.Requester) *cmd {
	cmd := &cmd{
		ctx:       ctx,
		requester: r,
		rootCmd:   &cobra.Command{},
	}

	cmd.init()

	return cmd
}

type cmd struct {
	ctx context.Context

	requester *service.Requester
	rootCmd   *cobra.Command
}

func (c *cmd) init() {
	c.addRunCmd()
}

func (c *cmd) Execute() {
	err := c.rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func (c *cmd) runCmd(cmd *cobra.Command, _ []string) error {
	var (
		reqData *model.RequestData
		err     error
	)

	if f := cmd.Flag(flagURL); f != nil {
		reqData, err = requestDataFromCmd(cmd)
	} else {
		reqData, err = requestDataFromFile(cmd.Flag(flagConfig).Value.String())
	}

	if err != nil {
		return err
	}

	return c.requester.Run(c.ctx, reqData)
}

func (c *cmd) addRunCmd() {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run requester",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if err := c.runCmd(cmd, args); err != nil {
				log.Fatalln(err)
			}
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			// нужно чтобы дефолтное значение считалось как заполненное
			if f := cmd.Flag(flagURL); f != nil && !f.Changed {
				cmd.Flag(flagConfig).Changed = true
			}
		},
	}

	c.rootCmd.AddCommand(runCmd)

	runCmd.Flags().String(flagConfig, "cfg.yaml", "")

	runCmd.Flags().String(flagURL, "", "")
	runCmd.Flags().Int(flagAmount, 0, "")
	runCmd.Flags().Int(flagPerSecond, 0, "")

	runCmd.MarkFlagsOneRequired(flagConfig, flagURL)

	runCmd.MarkFlagsRequiredTogether(flagURL, flagAmount, flagPerSecond)
}
