package lightning

import (
	"context"
	"fmt"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/invoicesrpc"
	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	"github.com/lightningnetwork/lnd/macaroons"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/macaroon.v2"
	"io"
	"io/ioutil"
	"path/filepath"
)

type service struct {
	config        *Config
	conn          *grpc.ClientConn
	Client        lnrpc.LightningClient
	RouterClient  routerrpc.RouterClient
	InvoiceClient invoicesrpc.InvoicesClient
}

func NewService(config *Config) Service {
	c := &service{
		config: config,
	}
	c.conn = GetLNDClientConn(config)
	c.Client = lnrpc.NewLightningClient(c.conn)
	c.RouterClient = routerrpc.NewRouterClient(c.conn)
	c.InvoiceClient = invoicesrpc.NewInvoicesClient(c.conn)
	log.Infof("lightning  started")

	return c
}

func (s *service) DecodePayReq(request *lnrpc.PayReqString) (*lnrpc.PayReq, error) {
	return s.Client.DecodePayReq(context.Background(), request)
}

func (s *service) SettleInvoice(request *invoicesrpc.SettleInvoiceMsg) (*invoicesrpc.SettleInvoiceResp, error) {
	return s.InvoiceClient.SettleInvoice(context.Background(), request)
}

func (s *service) NewHoldInvoice(rebalanceHash []byte, amtMsat uint64, memo string, cb InvoiceCallBack) (*invoicesrpc.AddHoldInvoiceResp, error) {
	ctxt := context.Background()
	resp, err := s.InvoiceClient.AddHoldInvoice(ctxt, &invoicesrpc.AddHoldInvoiceRequest{
		Hash:      rebalanceHash,
		ValueMsat: int64(amtMsat),
		Memo:      memo,
	})
	if err != nil {
		log.WithField("error", err).Error("failed tp create rebalance hold invoice")
		return nil, err

	}
	//payReqBytes := []byte(resp.PaymentRequest)
	s.SubscribeSingleInvoice(&invoicesrpc.SubscribeSingleInvoiceRequest{
		RHash: rebalanceHash,
	}, cb)
	return resp, nil

}
func (s *service) Close() {
	s.conn.Close()
}

func (s *service) MakeHashPaymentAndMonitor(destination []byte, chanId uint64, hash []byte, addr []byte, amtMsat uint64,
	cb PaymentCallBack) error {
	ctxt := context.Background()
	paymentClient, err := s.RouterClient.SendPaymentV2(ctxt, &routerrpc.SendPaymentRequest{
		Dest: destination,
		// TODO - pass fee as parameter
		AmtMsat:           int64(amtMsat), //11000,
		PaymentHash:       hash,
		FeeLimitMsat:      0,
		OutgoingChanIds:   []uint64{chanId},
		MaxParts:          0,
		NoInflightUpdates: false,
		TimeoutSeconds:    10,
		FinalCltvDelta:    40 + 3,
		PaymentAddr:       addr,
	})
	if err != nil {
		return err
	}
	go func() {
		for {
			payment, err := paymentClient.Recv()
			if err != nil {
				if err == io.EOF {
					return
				}
				log.Println(err)
				return
			}
			cb(payment)
		}
	}()
	return nil
}

func GetLNDClientConn(config *Config) *grpc.ClientConn {

	defaultLndDir, _ := homedir.Expand(config.DataDir)
	RpcHostPort := fmt.Sprintf("%v:%v", config.Host, config.Port)

	defaultTLSCertPath := filepath.Join(defaultLndDir, defaultTLSCertFilename)
	defaultMacaroonPath := filepath.Join(defaultLndDir,
		"data/chain/bitcoin",
		config.Network,
		defaultMacaroonFilename,
	)

	creds, err := credentials.NewClientTLSFromFile(defaultTLSCertPath, "")
	if err != nil {
		log.Fatal(err)
	}

	macaroonBytes, err := ioutil.ReadFile(defaultMacaroonPath)
	if err != nil {
		log.Fatal("Cannot read macaroon file", err)
	}

	mac := &macaroon.Macaroon{}
	if err = mac.UnmarshalBinary(macaroonBytes); err != nil {
		log.Fatal("Cannot unmarshal macaroon", err)
	}

	macCredential, _ := macaroons.NewMacaroonCredential(mac)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(100 * 1024 * 1024)),
		grpc.WithPerRPCCredentials(macCredential),
	}

	conn, err := grpc.Dial(RpcHostPort, opts...)
	if err != nil {
		log.Fatal(err)
	}

	return conn
}

func (s *service) SubscribeSingleInvoice(req *invoicesrpc.SubscribeSingleInvoiceRequest, cb InvoiceCallBack) error {
	ctxt := context.Background()
	singleInvoiceClient, err := s.InvoiceClient.SubscribeSingleInvoice(ctxt, req)
	if err != nil {
		return fmt.Errorf("can't SubscribeSingleInvoice for request %v - %v", req, err)
	}
	go func() {
		for {
			singleInvoiceUpdate, err := singleInvoiceClient.Recv()
			if err != nil {
				if err == io.EOF {
					//log.Println("singleInvoiceClient.Recv() closed")
					return
				}
				log.Println("got error from singleInvoiceClient.Recv()", err)
				return
			}
			cb(singleInvoiceUpdate)
		}
	}()
	return nil
}

func (s *service) GetInfo(ctx context.Context, request *lnrpc.GetInfoRequest) (*lnrpc.GetInfoResponse, error) {
	return s.Client.GetInfo(ctx, request)
}

func (s *service) ListChannels(ctx context.Context, request *lnrpc.ListChannelsRequest) (*lnrpc.ListChannelsResponse, error) {
	return s.Client.ListChannels(ctx, request)
}

func (s *service) ChanInfo(ctx context.Context, request *lnrpc.ChanInfoRequest) (*lnrpc.ChannelEdge, error) {
	return s.Client.GetChanInfo(ctx, request)
}

func (s *service) FeeReport(ctx context.Context, request *lnrpc.FeeReportRequest, opts ...grpc.CallOption) (*lnrpc.FeeReportResponse, error) {
	return s.Client.FeeReport(ctx, request, opts...)
}

func (s *service) SignMessage(ctx context.Context, request *lnrpc.SignMessageRequest, opts ...grpc.CallOption) (*lnrpc.SignMessageResponse, error) {
	return s.Client.SignMessage(ctx, request, opts...)
}

func (s *service) DescribeGraph(ctx context.Context, request *lnrpc.ChannelGraphRequest, opts ...grpc.CallOption) (*lnrpc.ChannelGraph, error) {
	return s.Client.DescribeGraph(ctx, request, opts...)
}
