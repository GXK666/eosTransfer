package transfer

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	pb "git.cochain.io/cochain/chains/service/eos"
	"github.com/GXK666/eosTransfer/log"
	"github.com/GXK666/eosTransfer/service/general"
	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/eoscanada/eos-go/system"
	"github.com/eoscanada/eos-go/token"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var Server general.ServiceServer

type Service struct {
	general.BaseService
	client pb.EosServiceClient
	keyBag *eos.KeyBag
}

func Setup() {
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

	cfg := viper.Sub("chainsRpc")
	creds, err := credentials.NewClientTLSFromFile(cfg.GetString("pemfile"), cfg.GetString("endpoint"))
	if err != nil {
		log.Errorw("credentials.NewClientTLSFromFile")
		panic("chains credentials.NewClientTLSFromFile")
	}
	conn, err := grpc.Dial(cfg.GetString("endpoint"), grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Errorw("HotWalletTransferOut()  grpc.Dial", "err", err)
		panic(err)
	}
	//defer conn.Close()
	client := pb.NewEosServiceClient(conn)

	Server = &Service{
		keyBag: signer,
		client: client,
	}
}

func (s *Service) SignPushActions(ctx context.Context, actions ...*eos.Action) (txid *string, err error) {
	ctxGetInfo, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	opts := &eos.TxOptions{}
	if rsp, err := s.client.GetInfo(ctxGetInfo, &pb.GetInfoRequest{}); err != nil {
		log.Errorw("rpc GetInfo ", "err", err)
		return nil, err
	} else {
		opts.ChainID, _ = hex.DecodeString(rsp.ChainId)
		opts.HeadBlockID, _ = hex.DecodeString(rsp.HeadBlockId)
	}

	tx := eos.NewTransaction(actions, opts)
	stx := eos.NewSignedTransaction(tx)

	pubKeys, err := s.keyBag.AvailableKeys()
	if nil != err {
		log.Errorw("keyBag.AvailableKeys ", "err", err)
		return nil, fmt.Errorf("%#v", err)
	}
	var pack *eos.PackedTransaction
	if signedTx, err := s.keyBag.Sign(stx, opts.ChainID, pubKeys...); nil != err {
		log.Errorw("keyBag.Sign", "err", err)
		return nil, fmt.Errorf("%#v", err)
	} else {
		stx = signedTx
	}

	var reqBody struct {
		Transaction     *eos.Transaction `json:"transaction"`
		ContextFreeData []eos.HexBytes   `json:"context_free_data"`

		Signatures []ecc.Signature `json:"signatures"`
	}
	reqBody.Transaction = tx
	reqBody.Signatures = stx.Signatures
	reqBody.ContextFreeData = stx.ContextFreeData

	txByte, err := json.Marshal(reqBody)
	if err != nil {
		log.Errorw("json.Marshal", "data", pack)
		return nil, err
	}
	if rsp, err := s.client.SendTransaction(ctx, &pb.SendTransactionRequest{Tx: string(txByte)}); err != nil {
		log.Errorw("rpc SendTransaction ", "err", err)
		return nil, err
	} else {
		return &rsp.TxId, nil
	}

	return nil, nil
}

func (s *Service) TransferOut(ctx context.Context, req *general.TransferOutRequest) (*general.TransferOutResponse, error) {
	fromAccount := viper.GetString("eosFromAccount")
	log.Infof("TransferOut start fromAccount: %s, req: %#v", fromAccount, req)
	amount, err := eos.NewAsset(req.Amount)
	if err != err {
		log.Errorf("TransferOut err, req: %#v  ,asset error : %#v", req, err)
		return nil, fmt.Errorf("asset error : %#v", err)
	}
	nonce, _ := uuid.NewV4()
	txid, err := s.SignPushActions(ctx, &eos.Action{
		Account: eos.AN(req.Contract),
		Name:    eos.ActN("transfer"),
		Authorization: []eos.PermissionLevel{
			{Actor: eos.AN(fromAccount), Permission: eos.PN("active")},
		},
		ActionData: eos.NewActionData(token.Transfer{
			From:     eos.AN(fromAccount),
			To:       eos.AN(req.To),
			Quantity: amount,
			Memo:     req.Memo,
		}),
	}, system.NewNonce(nonce.String()))

	if nil != err || txid == nil || len(*txid) != 64 {
		log.Errorf("TransferOut req: %#v, txid %#v,  error :%v", req, txid, err)
		return nil, fmt.Errorf("txid %#v, rsp error :%v", txid, err)
	}

	log.Infof("TransferOut end %#v, success , %s", req, *txid)
	return &general.TransferOutResponse{Txid: *txid}, nil
}

func (s *Service) GetTransferStatus(ctx context.Context, req *general.GetTransferStatusRequest) (*general.GetTransferStatusResponse, error) {
	rsp, err := s.client.GetTransactions(ctx, &pb.GetTransactionsRequest{Id: req.Txid})
	if nil != err {
		log.Errorf("GetTransactions txid %s, err %v", req.Txid, err)
		return nil, fmt.Errorf("GetTransactions txid %s, err %v", req.Txid, err)
	}
	if len(rsp.Transactions) == 0 {
		log.Errorf("not find transaction")
		return nil, fmt.Errorf("not find transaction")
	} else if len(rsp.Transactions) > 1 {
		log.Errorf("find transactions >  1")
		return nil, fmt.Errorf("find transactions >  1")
	}
	tx := rsp.Transactions[0]

	info, err := s.client.GetInfo(ctx, &pb.GetInfoRequest{})
	if nil != err {
		log.Errorf("get chain info %#v", err)
		return nil, fmt.Errorf("get chain info %#v", err)
	}

	blocks, err := s.client.GetBlocks(ctx, &pb.GetBlocksRequest{Id: tx.BlockId})
	if nil != err {
		log.Errorf("GetBlocks, err: %#v", err)
		return nil, fmt.Errorf("GetBlocks, err: %#v", err)
	}
	if len(blocks.Blocks) == 1 && blocks.Blocks[0].Num <= info.LastIrreversibleBlockNum {
		return &general.GetTransferStatusResponse{
			Status:   "irreversible",
			Txid:     tx.Id,
			BlockNum: tx.BlockNum,
			Blockid:  tx.BlockId,
		}, nil
	}

	return &general.GetTransferStatusResponse{
		Status:   tx.Status,
		Txid:     tx.Id,
		BlockNum: tx.BlockNum,
		Blockid:  tx.BlockId,
	}, nil
}

func (s *Service) GetSupportPubKey(ctx context.Context, req *general.GetSupportPubKeyRequest) (*general.GetSupportPubKeyResponse, error) {
	ret := []string{}
	if keys, err := s.keyBag.AvailableKeys(); nil != err {
		log.Errorf("keyBag.AvailableKeys, error : %#v", err)
		return nil, err
	} else {
		for _, pub := range keys {
			ret = append(ret, pub.String())
		}
	}

	return &general.GetSupportPubKeyResponse{
		PubKeys: ret,
	}, nil
}

func (s *Service) CheckAccountExist(ctx context.Context, req *general.CheckAccountRequest) (*general.CheckAccountResponse, error) {
	rsp, err := s.client.GetAccount(ctx, &pb.GetAccountRequest{Account: req.Account})
	if nil != err {
		log.Errorf("CheckAccountExist txid %s, err %v", req.Account, err)
		return nil, fmt.Errorf("CheckAccountExist txid %s, err %v", req.Account, err)
	}

	exits := rsp.AccountName == req.Account
	return &general.CheckAccountResponse{Exist: exits}, nil
}
