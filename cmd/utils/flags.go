package utils

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

const LND = "lnd"
const MainNet = "mainnet"
const SimNet = "simnet"

var (
	LNCHostFlag = &cli.StringFlag{
		Name:     "lnc-host",
		Usage:    "lightning coordinator host",
		Aliases:  []string{"lnch"},
		Required: true,
	}
	LNCPortFlag = &cli.IntFlag{
		Name:    "lnc-port",
		Usage:   "lightning coordinator port",
		Aliases: []string{"lncp"},
		Value:   2222,
	}

	// agent only

	ImplementationFlag = &cli.StringFlag{
		Name:    "implementation",
		Usage:   "specify the lightning node implementation (lnd, c-lightning)",
		Aliases: []string{"impl"},

		Value: "lnd",
	}
	TorActiveFlag = &cli.BoolFlag{
		Name:  "tor.active",
		Usage: "tcp connections to use Tor",
		Value: false,
	}
	TorSocksFlag = &cli.StringFlag{
		Name:  "tor.socks",
		Usage: "The host:port that Tor's exposed SOCKS5 proxy is listening on (default: localhost:9050)",
		Value: "localhost:9050",
	}

	// coordinator only

	RebalanceBudgetFlag = &cli.Float64Flag{
		Name:    "rebalance-budget-ppm-percent",
		Usage:   "rebalance budget expressed as percentage of the channel's ppm (fee_per_mil). Budget may exceeds 100%",
		Aliases: []string{"budget"},
		Value:   0.0,
	}
	ChannelLowTriggerFlag = &cli.Float64Flag{
		Name:    "channel-low-trigger",
		Usage:   "when local amount (as percentage of capacity) fails below the channel-low-trigger, the channel becomes rebalance candidate. Expressed as percentage [0-100]",
		Aliases: []string{"low"},
		Value:   25.0,
	}
	ChannelHighTriggerFlag = &cli.Float64Flag{
		Name:    "channel-high-trigger",
		Usage:   "when local amount (as percentage of capacity) exceeds the channel-high-trigger, the channel becomes rebalance candidate. Expressed as percentage [0-100]",
		Aliases: []string{"high"},
		Value:   75.0,
	}
	ChannelTargetFlag = &cli.Float64Flag{
		Name:    "channel-target",
		Usage:   "the target of rebalance expressed as percentage [0-100] of the channel's capacity",
		Aliases: []string{"target"},
		Value:   50.0,
	}
)

func ValidateSetting(budget, low, high, target float64) error {
	if budget < 0.0 {
		return fmt.Errorf("rebalance budget %v can't be negative", budget)
	}

	if low < 0.0 || low > 100.0 {
		return fmt.Errorf("%v is an invalid value for %v. Please provide a value between 0 and 100", low, ChannelLowTriggerFlag.Name)
	}

	if high < 0.0 || high > 100.0 {
		return fmt.Errorf("%v is an invalid value for %v. Please provide a value between 0 and 100", high, ChannelHighTriggerFlag.Name)
	}

	if target < 0.0 || target > 100.0 {
		return fmt.Errorf("%v is an invalid value for %v. Please provide a value between 0 and 100", target, ChannelTargetFlag.Name)
	}
	if low >= target {
		return fmt.Errorf("%v value %v must be lower than %v value %v", ChannelLowTriggerFlag.Name, low, ChannelTargetFlag.Name, target)
	}
	if high <= target {
		return fmt.Errorf("%v value %v must be greater than %v value %v", ChannelHighTriggerFlag.Name, high, ChannelTargetFlag.Name, target)
	}
	return nil
}
