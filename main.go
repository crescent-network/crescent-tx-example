package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	chain "github.com/crescent-network/crescent/v4/app"

	"github.com/crescent-network/crescent-tx-example/client"
	"github.com/crescent-network/crescent-tx-example/config"
	"github.com/crescent-network/crescent-tx-example/wallet"
)

var (
	timeout = 5 * time.Second
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}) // human-friendly output

	config := config.SetAddressPrefixes()
	config.Seal()
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	config, err := config.Read(config.DefaultConfigPath)
	if err != nil {
		panic(fmt.Errorf("failed to read config.toml file: %w", err))
	}

	// Connect Tendermint RPC client
	rpcClient, err := client.ConnectRPCWithTimeout(config.RPC.Address, timeout)
	if err != nil {
		panic(fmt.Errorf("failed to connect RPC client: %w", err))
	}

	// Connect gRPC client
	gRPCConn, err := client.ConnectGRPCWithTimeout(ctx, config.GRPC.Address, config.GRPC.UseTLS, timeout)
	if err != nil {
		panic(fmt.Errorf("failed to connect gRPC client: %w", err))
	}
	defer gRPCConn.Close()

	// Recover private key from mnemonic phrases
	privKey, err := wallet.RecoverPrivKeyFromMnemonic(config.WalletConfig.Mnemonic, config.WalletConfig.Password)
	if err != nil {
		panic(fmt.Errorf("recovering private key: %w", err))
	}

	chainID, _ := rpcClient.NetworkChainID(ctx)
	creator := wallet.Address(privKey)
	baseAccount, _ := gRPCConn.GetAccount(ctx, creator.String())
	accNum := baseAccount.GetAccountNumber()
	accSeq := baseAccount.GetSequence()
	gasLimit := config.TxConfig.GasLimit
	fees, err := sdk.ParseCoinsNormalized(config.TxConfig.Fees)
	if err != nil {
		panic(fmt.Errorf("failed to parse coins %w", err))
	}

	// Create new MsgSend for test
	msg := banktypes.MsgSend{
		FromAddress: creator.String(),
		ToAddress:   "cre15rz2rwnlgr7nf6eauz52usezffwrxc0mxajpmw",
		Amount:      sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
	}
	msgs := []sdk.Msg{&msg}

	tx := client.NewTx(
		chainID,
		accNum,
		accSeq,
		gasLimit,
		fees,
		msgs...,
	)
	txCfg := chain.MakeEncodingConfig().TxConfig
	txBytes, err := client.SignTx(tx, txCfg, privKey)
	if err != nil {
		fmt.Printf("failed to sign transaction: %v", err)
		return
	}

	resp, err := gRPCConn.BroadcastTx(ctx, txBytes, sdktx.BroadcastMode_BROADCAST_MODE_BLOCK)
	if err != nil {
		fmt.Printf("failed to broadcast transaction: %v", err)
		return
	}

	log.Info().Msg("Go to the following link to see if transaction is successfully included in a block")
	log.Info().Msg("http://localhost:1317/cosmos/tx/v1beta1/txs/" + resp.TxResponse.TxHash)
}
