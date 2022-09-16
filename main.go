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

func init() {
	config := config.SetAddressPrefixes()
	config.Seal()
}

func main() {
	timeout := 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Read configuration from the config.toml file
	config, err := config.Read(config.DefaultConfigPath)
	if err != nil {
		panic(fmt.Errorf("failed to read config.toml file: %w", err))
	}

	// Connect the Tendermint RPC endpoint
	rpcClient, err := client.ConnectRPCWithTimeout(config.RPC.Address, timeout)
	if err != nil {
		panic(fmt.Errorf("failed to connect RPC client: %w", err))
	}

	// Connect the gRPC server
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
	address := wallet.Address(privKey)
	baseAccount, _ := gRPCConn.GetAccount(ctx, address.String())
	accNum := baseAccount.GetAccountNumber()
	accSeq := baseAccount.GetSequence()
	gasLimit := config.TxConfig.GasLimit
	fees, _ := sdk.ParseCoinsNormalized(config.TxConfig.Fees)

	// You can append more messages...
	msgs := []sdk.Msg{
		MsgMMOrder(address),
	}

	tx := client.NewTx(
		chainID,
		accNum,
		accSeq,
		gasLimit,
		fees,
		msgs...,
	)
	txBytes, err := client.SignTx(tx, chain.MakeEncodingConfig().TxConfig, privKey)
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
	// $ crescentd query liquidity pairs --node <YOUR_NODE> -o json | jq
	// $ crescentd query liquidity pools --node <YOUR_NODE> -o json | jq
	fmt.Println("Sent transaction successfully!")
	fmt.Println("TxHash: ", resp.TxResponse.TxHash)

	//
	// You can use gRPCConn to query pair(s) and pool(s)
	// See the implemented queries in client/grpc.go file
	// ...
	//
}

// Example of MMOrder Message
//
// MsgMMOrder create a transaction message to make market making order (MMOrder).
// An MMOrder is a group of multiple buy/sell limit orders which are
// distributed evenly based on its parameters.
func MsgMMOrder(addr sdk.AccAddress) sdk.Msg {
	// the order lifespan; it is the duration that the order lives until it is expired.
	// even if you have 0 for order life span, an order requires at least one batch to be executed
	// and maximum order life span is 24 hours for performance reason.
	var orderLifeSpan = 30 * time.Second

	// Adjust the following values for your needs
	msg := liquiditytypes.MsgMMOrder{
		Orderer:       addr.String(),
		PairId:        1,                            // the pair id that you would like to target
		MaxSellPrice:  sdk.MustNewDecFromStr("102"), // the maximum sell price
		MinSellPrice:  sdk.MustNewDecFromStr("101"), // the minimum sell price
		SellAmount:    sdk.NewInt(1000000),          // the total amount of base coin of sell orders
		MaxBuyPrice:   sdk.MustNewDecFromStr("100"), // the maximum buy price
		MinBuyPrice:   sdk.MustNewDecFromStr("99"),  // the minimum buy price
		BuyAmount:     sdk.NewInt(1000000),          // the total amount of base coin of buy orders
		OrderLifespan: orderLifeSpan,
	}
	return &msg
}

// Example of CreatePool Message
//
// MsgCreatePool creates a transaction message to create new pool.
func MsgCreatePool(addr sdk.AccAddress) sdk.Msg {
	msg := liquiditytypes.MsgCreatePool{
		Creator: addr.String(),
		PairId:  1,
		DepositCoins: sdk.NewCoins(
			sdk.NewInt64Coin("uusd", 100_000_000),
			sdk.NewInt64Coin("uatom", 100_000_000),
		),
	}
	return &msg
}

// Example of CreatePair Message
//
// MsgCreatePair creates a transaction message to create new pair.
func MsgCreatePair(addr sdk.AccAddress) sdk.Msg {
	msg := liquiditytypes.MsgCreatePair{
		Creator:        addr.String(),
		BaseCoinDenom:  "uusd",
		QuoteCoinDenom: "uatom",
	}
	return &msg
}
