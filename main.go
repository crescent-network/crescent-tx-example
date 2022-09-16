package main

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"

	chain "github.com/crescent-network/crescent/v3/app"
	liquiditytypes "github.com/crescent-network/crescent/v3/x/liquidity/types"

	"github.com/crescent-network/crescent-tx-example/client"
	"github.com/crescent-network/crescent-tx-example/config"
	"github.com/crescent-network/crescent-tx-example/wallet"
)

var (
	timeout = 5 * time.Second
)

func init() {
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

	var (
		baseCoinDenom  = "stake"
		quoteCoinDenom = "uatom"
	)

	msg1 := liquiditytypes.MsgCreatePair{
		Creator:        creator.String(),
		BaseCoinDenom:  baseCoinDenom,
		QuoteCoinDenom: quoteCoinDenom,
	}
	msg2 := liquiditytypes.MsgCreatePool{
		Creator: creator.String(),
		PairId:  1,
		DepositCoins: sdk.NewCoins(
			sdk.NewInt64Coin(baseCoinDenom, 100_000_000),
			sdk.NewInt64Coin(quoteCoinDenom, 100_000_000),
		),
	}
	msgs := []sdk.Msg{&msg1, &msg2}

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

	// Query to see if the pair and the pool is created
	// $ crescentd query liquidity pairs -o json | jq
	// $ crescentd query liquidity pools -o json | jq
	fmt.Println("Sent transaction successfully!")
	fmt.Println("TxHash: ", resp.TxResponse.TxHash)

	//
	// You can use gRPCConn to query pair(s) and pool(s)
	// See the implemented queries in client/grpc.go file
	// ...
	//
}
