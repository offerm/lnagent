package lightning

import (
	"context"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/invoicesrpc"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

type InvoiceCallBack func(*lnrpc.Invoice) error
type PaymentCallBack func(*lnrpc.Payment) error

// TODO: add context  to all
type Service interface {
	DecodePayReq(*lnrpc.PayReqString) (*lnrpc.PayReq, error)
	NewHoldInvoice([]byte, uint64, string, InvoiceCallBack) (*invoicesrpc.AddHoldInvoiceResp, error)
	MakeHashPaymentAndMonitor([]byte, uint64, []byte, []byte, uint64, PaymentCallBack) error
	SettleInvoice(*invoicesrpc.SettleInvoiceMsg) (*invoicesrpc.SettleInvoiceResp, error)
	GetInfo(context.Context, *lnrpc.GetInfoRequest) (*lnrpc.GetInfoResponse, error)
	ListChannels(context.Context, *lnrpc.ListChannelsRequest) (*lnrpc.ListChannelsResponse, error)
	ChanInfo(ctx context.Context, request *lnrpc.ChanInfoRequest) (*lnrpc.ChannelEdge, error)
	FeeReport(context.Context, *lnrpc.FeeReportRequest, ...grpc.CallOption) (*lnrpc.FeeReportResponse, error)
	SignMessage(context.Context, *lnrpc.SignMessageRequest, ...grpc.CallOption) (*lnrpc.SignMessageResponse, error)
	DescribeGraph(context.Context, *lnrpc.ChannelGraphRequest, ...grpc.CallOption) (*lnrpc.ChannelGraph, error)
	Close()
}

type Config struct {
	Host           string
	Port           int
	Network        string
	Implementation string
	DataDir        string
}

const (
	defaultTLSCertFilename  = "tls.cert"
	defaultMacaroonFilename = "admin.macaroon"
	defaultRpcPort          = "1231"
	defaultRpcHostPort      = "localhost:" + defaultRpcPort

	LND     = "lnd"
	MainNet = "mainnet"
	SimNet  = "simnet"
)

var (
	LnHostFlag = &cli.StringFlag{
		Name:    "ln-host",
		Usage:   "host name/ip of the lightning node",
		Aliases: []string{"lnh"},
		Value:   "localhost",
	}
	LnPortFlag = &cli.IntFlag{
		Name:    "ln-port",
		Usage:   "port of the lightning node",
		Aliases: []string{"lnp"},
		Value:   10009,
	}
	NetworkFlag = &cli.StringFlag{
		Name:    "network",
		Usage:   "lightning environment (mainnet, testnet, simnet)",
		Aliases: []string{"n"},
		Value:   "mainnet",
	}
	LndDirFlag = &cli.StringFlag{
		Name:     "lnd-dir",
		Usage:    "path to lnd directory",
		Required: true,
	}
	ImplementationFlag = &cli.StringFlag{
		Name:    "implementation",
		Usage:   "specify the lightning node implementation (lnd, c-lightning)",
		Aliases: []string{"impl"},
		Value:   "lnd",
	}
)
