package tx

import (
	"bytes"
	"encoding/hex"
	"github.com/tybc/blockchain"
	"github.com/tybc/core/types"
	"github.com/tybc/errors"
	"github.com/tybc/log"
	"github.com/tybc/wallet"
)

var (
	ErrSubmitTx = errors.New("sumbit tx")
)

type SubmitTxRequest struct {
	Password  string      `json:"wallet_password"`
	TxInputs  []ReqInput  `json:"tx_inputs"`
	TxOutputs []ReqOutput `json:"tx_outputs"`
}

type ReqInput struct {
	SpendOutputId string `json:"spend_output_id"`
}

type ReqOutput struct {
	Address string `json:"address"`
	Amount  uint64 `json:"amount"`
}

type SumbitTxResponse struct {
	TxId string `json:"tx_id"`
}

func (req *SubmitTxRequest) MapTx() (*types.Tx, error) {

	ins := make([]types.TxInput, len(req.TxInputs))
	ios := make([]types.TxOutput, len(req.TxOutputs))

	for _, inp := range req.TxInputs {
		b, err := hex.DecodeString(inp.SpendOutputId)
		if err != nil {
			return nil, errors.WithDetail(ErrSubmitTx, "invalid spend_output_id format")
		}
		ins = append(ins, types.TxInput{
			SpendOutputId: types.BytesToHash(b),
		})
	}

	for _, iop := range req.TxOutputs {
		addr, err := hex.DecodeString(iop.Address)
		if err != nil {
			return nil, errors.WithDetail(ErrSubmitTx, "invalid output.address format")
		}

		hashId := bytes.Join([][]byte{

		}, []byte{})

		ios = append(ios, types.TxOutput{
			Id:         hashId,
			IsCoinBase: false,
			Address:    addr,
			Amount:     iop.Amount,
		})
	}

	tx := &types.Tx{
		TxInput:  ins,
		TxOutput: ios,
	}

	return tx, nil
}

func SubmitTx(chain *blockchain.Chain, reqTx *SubmitTxRequest) (*SumbitTxResponse, error) {

	if len(reqTx.TxInputs) == 0 {
		return nil, errors.WithDetail(ErrSubmitTx, "no input data")
	}

	if len(reqTx.TxOutputs) == 0 {
		return nil, errors.WithDetail(ErrSubmitTx, "no output data")
	}

	// request data map to tx
	tx, err := reqTx.MapTx()
	if err != nil {
		return nil, err
	}

	//check and set utxo
	for _, input := range tx.TxInput {
		input.SetUtxo(&chain.Store)
	}

	//TODO sign transaction
	if wt, err := wallet.My(reqTx.Password); err != nil {
		return nil, err
	} else {
		log.Logger.Infof("pub %s", wt.Pub)
	}

	//TODO check if exist on txpool

	return nil, nil
}
