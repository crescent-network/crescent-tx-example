package client

import (
	"fmt"

	sdkclient "github.com/cosmos/cosmos-sdk/client"
	clienttx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
)

type Tx struct {
	ChainID        string
	AccNum, AccSeq uint64
	GasLimit       uint64
	Fees           sdk.Coins
	Msgs           []sdk.Msg
}

func NewTx(chainID string, accNum, accSeq, gasLimit uint64, fees sdk.Coins, msgs ...sdk.Msg) *Tx {
	return &Tx{
		ChainID:  chainID,
		AccNum:   accNum,
		AccSeq:   accSeq,
		GasLimit: gasLimit,
		Fees:     fees,
		Msgs:     msgs,
	}
}

func SignTx(tx *Tx, txCfg sdkclient.TxConfig, privKey *secp256k1.PrivKey) ([]byte, error) {
	b := txCfg.NewTxBuilder()
	if err := b.SetMsgs(tx.Msgs...); err != nil {
		return nil, fmt.Errorf("set msgs: %w", err)
	}
	b.SetGasLimit(tx.GasLimit)
	b.SetFeeAmount(tx.Fees)

	mode := txCfg.SignModeHandler().DefaultMode()
	sig := signing.SignatureV2{
		PubKey: privKey.PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode: mode,
		},
		Sequence: tx.AccSeq,
	}
	if err := b.SetSignatures(sig); err != nil {
		return nil, fmt.Errorf("set signatures: %w", err)
	}
	data := authsigning.SignerData{
		ChainID:       tx.ChainID,
		AccountNumber: tx.AccNum,
		Sequence:      tx.AccSeq,
	}
	sig, err := clienttx.SignWithPrivKey(mode, data, b, privKey, txCfg, tx.AccSeq)
	if err != nil {
		return nil, fmt.Errorf("sign with priv key: %w", err)
	}
	if err := b.SetSignatures(sig); err != nil {
		return nil, fmt.Errorf("set signatures: %w", err)
	}
	txBytes, err := txCfg.TxEncoder()(b.GetTx())
	if err != nil {
		return nil, fmt.Errorf("encode tx: %w", err)
	}
	return txBytes, nil
}
