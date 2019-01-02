package transfer

import (
	"context"
	"fmt"
	"github.com/GXK666/eosTransfer/log"
	"github.com/GXK666/eosTransfer/service/general"
	eos "github.com/eoscanada/eos-go"
	system "github.com/eoscanada/eos-go/system"
	token "github.com/eoscanada/eos-go/token"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"strings"
)

var Server general.ServiceServer

type Service struct {
	general.BaseService
	*eos.API
}

func Setup() {
	nodeUrl := viper.GetString("eosNodeUrl")
	if len(nodeUrl) == 0 {
		panic("node Url is null")
	}
	pivKey := viper.GetString("eosPrivateKeys")
	pivKeys := strings.Split(pivKey, ",")
	if len(pivKeys) == 0 {
		panic("privateKey is null")
	}

	signer := eos.NewKeyBag()

	for _, k := range pivKeys {
		err := signer.ImportPrivateKey(k)
		if nil != err {
			panic("ImportPrivateKey error")
		}
	}

	api := eos.New(nodeUrl)
	api.SetSigner(signer)
	Server = &Service{
		API: api,
	}
}

func (s *Service) TransferOut(ctx context.Context, req *general.TransferOutRequest) (*general.TransferOutResponse, error) {
	amount, err := eos.NewAsset(req.Amount)
	if err != err {
		return nil, fmt.Errorf("asset error : %#v", err)
	}
	nonce, _ := uuid.NewV4()
	rsp, err := s.SignPushActions(&eos.Action{
		Account: eos.AN(req.Contract),
		Name:    eos.ActN("transfer"),
		Authorization: []eos.PermissionLevel{
			{Actor: eos.AN(req.From), Permission: eos.PN("active")},
		},
		ActionData: eos.NewActionData(token.Transfer{
			From:     eos.AN(req.From),
			To:       eos.AN(req.To),
			Quantity: amount,
			Memo:     req.Memo,
		}),
	}, system.NewNonce(nonce.String()))

	if nil != err {
		log.Error("transfer %#v, rsp error :%v", req, err)
		return nil, fmt.Errorf("rsp error :%v", err)
	}
	if rsp.StatusCode != "" {
		log.Error("transfer %#v, rsp : %#v", req, rsp)
		return nil, fmt.Errorf("rsp : %#v", rsp)
	}

	log.Info("transfer %#v, success , %s", req, rsp.TransactionID)
	return &general.TransferOutResponse{Txid: rsp.TransactionID}, nil
}
