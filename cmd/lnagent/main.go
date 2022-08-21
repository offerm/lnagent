package main

import (
	"github.com/offerm/lnagent"
	"github.com/offerm/lnagent/cmd/utils"
	"github.com/offerm/lnagent/lightning"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := &cli.App{
		UseShortOptionHandling: true,
		Name:                   "lnagent",
		Usage:                  "client side of lncoordinator. Re-balancing and more",
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "start the agent",
				Flags: []cli.Flag{
					lightning.LnHostFlag,
					lightning.LnPortFlag,
					lightning.NetworkFlag,
					utils.ImplementationFlag,
					lightning.LndDirFlag,
					utils.LNCHostFlag,
					utils.LNCPortFlag,
					utils.ChannelLowTriggerFlag,
					utils.ChannelHighTriggerFlag,
					utils.ChannelTargetFlag,
					utils.RebalanceBudgetFlag,
					utils.TorActiveFlag,
					utils.TorSocksFlag,
				},
				Action: runCmd,
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func runCmd(ctx *cli.Context) error {
	if err := utils.ValidateSetting(ctx.Float64(utils.RebalanceBudgetFlag.Name), ctx.Float64(utils.ChannelLowTriggerFlag.Name),
		ctx.Float64(utils.ChannelHighTriggerFlag.Name), ctx.Float64(utils.ChannelTargetFlag.Name)); err != nil {
		return err
	}

	lnConfig := &lightning.Config{
		Host:           ctx.String(lightning.LnHostFlag.Name),
		Port:           ctx.Int(lightning.LnPortFlag.Name),
		Network:        ctx.String(lightning.NetworkFlag.Name),
		Implementation: ctx.String(lightning.ImplementationFlag.Name),
		DataDir:        ctx.String(lightning.LndDirFlag.Name),
	}

	lnagentConfig := &lnagent.Config{
		Host:               ctx.String(utils.LNCHostFlag.Name),
		Port:               ctx.Int(utils.LNCPortFlag.Name),
		ChannelLowTrigger:  ctx.Float64(utils.ChannelLowTriggerFlag.Name),
		ChannelHighTrigger: ctx.Float64(utils.ChannelHighTriggerFlag.Name),
		ChannelTarget:      ctx.Float64(utils.ChannelTargetFlag.Name),
		RebalanceBudget:    ctx.Float64(utils.RebalanceBudgetFlag.Name),
		TorActive:          ctx.Bool(utils.TorActiveFlag.Name),
		TorSocks:           ctx.String(utils.TorSocksFlag.Name),
	}

	agent := lnagent.NewAgent(lnagentConfig, lnConfig)
	return agent.Run()
}
