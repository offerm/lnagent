package lnagent

type Config struct {
	Host string
	Port int

	ChannelLowTrigger  float64
	ChannelHighTrigger float64
	ChannelTarget      float64
	RebalanceBudget    float64
	TorActive          bool
	TorSocks           string
}
