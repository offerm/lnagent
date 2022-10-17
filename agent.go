package lnagent

import (
	"context"
	"encoding/json"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/offerm/lnagent/lightning"
	"github.com/offerm/lnagent/protobuf"
	"github.com/offerm/lnagent/rebalancer"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
	"reflect"
	"sync"
	"time"
)

type emptyObject struct{}

// TODO: handle cleanup on startup
type agent struct {
	lnagentConfig *Config
	lnService     lightning.Service
	rebalancer    *rebalancer.Rebalancer

	coordinatorClient protobuf.CoordinatorClient

	events chan *protobuf.TaskResponse

	cancelCtx  context.Context
	cancelFunc context.CancelFunc
	// todo: change to empty object
	done chan emptyObject
}

func (agent *agent) Stop() {
	agent.cancelFunc()
	<-agent.done
}

func NewAgent(lnagentConfig *Config, lnService lightning.Service) *agent {
	agent := &agent{
		lnagentConfig: lnagentConfig,
		events:        make(chan *protobuf.TaskResponse, 100),
	}

	agent.lnService = lnService

	agent.rebalancer = rebalancer.NewRebalancer(agent.events, agent.lnService)

	agent.cancelCtx, agent.cancelFunc = context.WithCancel(context.Background())

	agent.done = make(chan emptyObject)

	return agent
}

// Run execute the agent
func (agent *agent) Run() error {
	log.Info("agent is starting")

	conn := GetLNCClientConn(agent.lnagentConfig)
	defer conn.Close()
	agent.coordinatorClient = protobuf.NewCoordinatorClient(conn)
	//agent.cancelOneSidePendingHoldInvoices()
	agent.loop()

	agent.lnService.Close()

	agent.done <- emptyObject{}

	//agent.Cleanup()
	return nil
}
func (agent *agent) getPubKey() string {
	for {
		info, err := agent.lnService.GetInfo(context.Background(), &lnrpc.GetInfoRequest{})
		if err == nil {
			return info.IdentityPubkey
		}
		log.Errorf("failed to query pubkey - %v", err)
		time.Sleep(5 * time.Second)
	}
}

func (agent *agent) loop() {
	wg := sync.WaitGroup{}
	ctxt := context.Background()
	infoTicker := time.Tick(10 * time.Second)
	defer func() { agent.coordinatorClient = nil }()

	go func() {
		wg.Add(1)
		for {
			select {
			case <-agent.cancelCtx.Done():
				wg.Done()
				return
			default:
				md := metadata.Pairs("pubkey", agent.getPubKey())
				ctx := metadata.NewOutgoingContext(agent.cancelCtx, md)
				taskClient, err := agent.coordinatorClient.Tasks(ctx)
				if err != nil {
					log.Errorf("failed to open Tasks feed with lncoordinator - %v", err)
					time.Sleep(5 * time.Second)
					continue
				}
				log.Info("Tasks feed established")
				ctx, cancel := context.WithCancel(agent.cancelCtx)
				go func() {
					wg.Add(1)
					defer log.Info("stop sending events to coordinator")
					log.Infof("start sending events to coordinator")
					if agent.lnagentConfig == nil {
						return
					}
					for {
						select {
						case event, ok := <-agent.events:
							if !ok {
								return
							}
							log.Tracef("sending event %v", reflect.TypeOf(event.Response))
							err := taskClient.Send(event)
							if err != nil {
								log.Errorf("failed to send to %v swayID %v - %v", event.Pubkey, event.Swap_ID, err)
							}
						case <-ctx.Done():
							wg.Done()
							return
						}

					}
				}()
				for {
					task, err := taskClient.Recv()
					select {
					case <-agent.cancelCtx.Done():
						cancel()
						wg.Done()
						return
					default:
						if err != nil {
							cancel()
							log.Errorf("got error while waiting for Task - %v", err)
							time.Sleep(5 * time.Second)
							break
						}
						agent.executeTask(task)

					}
				}

			}

		}

	}()

	for {
		select {
		case <-agent.cancelCtx.Done():
			wg.Wait()
			return

		case <-infoTicker:
			info, err := agent.lnService.GetInfo(ctxt, &lnrpc.GetInfoRequest{})
			if err != nil {
				log.Errorf("can't GetInfo - %v", err)
				continue
			}

			channels, err := agent.lnService.ListChannels(ctxt, &lnrpc.ListChannelsRequest{
				ActiveOnly: true,
			})
			if err != nil {
				log.Errorf("can't ListChannels - %v", err)
				continue
			}
			config := agent.lnagentConfig
			statusUpdate := &protobuf.StatusUpdateRequest{
				Pubkey:                       info.IdentityPubkey,
				RebalanceBudgetPpmPercent:    config.RebalanceBudget,
				ChannelRebalanbceLowTrigger:  config.ChannelLowTrigger,
				ChannelRebalanbceHighTrigger: config.ChannelHighTrigger,
				ChannelRebalanceTarget:       config.ChannelTarget,
			}

			for _, channel := range channels.Channels {
				channelInfo := &protobuf.Channel{
					Active:           channel.Active,
					RemotePubkey:     channel.RemotePubkey,
					ChannelPoint:     channel.ChannelPoint,
					ChanId:           channel.ChanId,
					Capacity:         channel.Capacity,
					LocalBalance:     channel.LocalBalance,
					RemoteBalance:    channel.RemoteBalance,
					CommitFee:        channel.CommitFee,
					CommitWeight:     channel.CommitWeight,
					FeePerKw:         channel.FeePerKw,
					UnsettledBalance: channel.UnsettledBalance,
					Private:          channel.Private,
					ChanStatusFlags:  channel.ChanStatusFlags,
					LocalConstraints: &protobuf.ChannelConstraints{
						CsvDelay:          channel.LocalConstraints.CsvDelay,
						ChanReserveSat:    channel.LocalConstraints.ChanReserveSat,
						DustLimitSat:      channel.LocalConstraints.DustLimitSat,
						MaxPendingAmtMsat: channel.LocalConstraints.MaxPendingAmtMsat,
						MinHtlcMsat:       channel.LocalConstraints.MinHtlcMsat,
						MaxAcceptedHtlcs:  channel.LocalConstraints.MaxAcceptedHtlcs,
					},
					RemoteConstraints: &protobuf.ChannelConstraints{
						CsvDelay:          channel.RemoteConstraints.CsvDelay,
						ChanReserveSat:    channel.RemoteConstraints.ChanReserveSat,
						DustLimitSat:      channel.RemoteConstraints.DustLimitSat,
						MaxPendingAmtMsat: channel.RemoteConstraints.MaxPendingAmtMsat,
						MinHtlcMsat:       channel.RemoteConstraints.MinHtlcMsat,
						MaxAcceptedHtlcs:  channel.RemoteConstraints.MaxAcceptedHtlcs,
					},
				}
				statusUpdate.Channels = append(statusUpdate.Channels, channelInfo)
				//log.WithField("channel", channel).Info("")
			}

			feeReport, err := agent.lnService.FeeReport(ctxt, &lnrpc.FeeReportRequest{})
			if err != nil {
				log.Errorf("can't FeeReport - %v", err)
				continue
			}
			for _, channelFees := range feeReport.ChannelFees {
				for _, channel := range statusUpdate.Channels {
					if channel.ChanId == channelFees.ChanId {
						channel.BaseFeeMsat = channelFees.BaseFeeMsat
						channel.FeePerMil = channelFees.BaseFeeMsat
						channel.FeeRate = channelFees.FeeRate
					}
				}

			}
			//log.WithField("statusUpdate", statusUpdate).Trace("")
			agent.sendUpdate(statusUpdate)
		}
	}
}

func (agent *agent) executeTask(task *protobuf.Task) {

	swapID := rebalancer.SwapID(task.Swap_ID)
	switch taskType := task.Type.(type) {
	case *protobuf.Task_InitType:
		agent.rebalancer.TaskInit(swapID, taskType.InitType)
	case *protobuf.Task_SwapType:
		agent.rebalancer.TaskSwap(swapID, taskType.SwapType)
	case *protobuf.Task_CancelType:
		agent.rebalancer.TaskCancel(swapID, taskType.CancelType)
	default:
		agent.rebalancer.TaskUnKnow(swapID)
	}
}

// sendUpdate is called to report the most up-to-date status of the channels
func (agent *agent) sendUpdate(statusUpdate *protobuf.StatusUpdateRequest) {
	ctxt := context.Background()
	b, err := json.Marshal(statusUpdate.Channels)
	if err != nil {
		log.Panicf("can't marshal statusUpdate.Channels")
	}
	signMessageResp, err := agent.lnService.SignMessage(ctxt, &lnrpc.SignMessageRequest{Msg: b})
	if err != nil {
		log.Error("can't sign StatusUpdateRequest message, skipping this message")
		return
	}
	statusUpdate.StatusSignature = signMessageResp.Signature
	// todo: add metadata
	_, err = agent.coordinatorClient.StatusUpdate(ctxt, statusUpdate)
	if err != nil {
		log.Errorf("failed to send StatusUpdate to coordinator - %v", err)
	}
}
