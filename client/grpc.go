package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	liquiditytypes "github.com/crescent-network/crescent/v2/x/liquidity/types"
)

type GRPCClient struct {
	conn *grpc.ClientConn
}

func ConnectGRPC(ctx context.Context, addr string, useTLS bool, opts ...grpc.DialOption) (*GRPCClient, error) {
	if !useTLS {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	}
	conn, err := grpc.DialContext(ctx, addr, append(opts, grpc.WithBlock())...)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}
	return &GRPCClient{conn: conn}, nil
}

func ConnectGRPCWithTimeout(ctx context.Context, addr string, useTLS bool, timeout time.Duration, opts ...grpc.DialOption) (*GRPCClient, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return ConnectGRPC(ctx, addr, useTLS, opts...)
}

func (c *GRPCClient) Close() error {
	return c.conn.Close()
}

func (c *GRPCClient) QueryPairs(ctx context.Context) (*liquiditytypes.QueryPairsResponse, error) {
	return liquiditytypes.NewQueryClient(c.conn).Pairs(ctx, &liquiditytypes.QueryPairsRequest{})
}

func (c *GRPCClient) QueryPair(ctx context.Context, id uint64) (*liquiditytypes.QueryPairResponse, error) {
	return liquiditytypes.NewQueryClient(c.conn).Pair(ctx, &liquiditytypes.QueryPairRequest{PairId: id})
}

func (c *GRPCClient) QueryPools(ctx context.Context) (*liquiditytypes.QueryPoolsResponse, error) {
	return liquiditytypes.NewQueryClient(c.conn).Pools(ctx, &liquiditytypes.QueryPoolsRequest{})
}

func (c *GRPCClient) QueryPool(ctx context.Context, id uint64) (*liquiditytypes.QueryPoolResponse, error) {
	return liquiditytypes.NewQueryClient(c.conn).Pool(ctx, &liquiditytypes.QueryPoolRequest{PoolId: id})
}

func (c *GRPCClient) GetAccount(ctx context.Context, addr string) (authtypes.BaseAccount, error) {
	resp, err := authtypes.NewQueryClient(c.conn).Account(ctx, &authtypes.QueryAccountRequest{Address: addr})
	if err != nil {
		return authtypes.BaseAccount{}, err
	}
	var acc authtypes.BaseAccount
	if err := acc.Unmarshal(resp.Account.Value); err != nil {
		return authtypes.BaseAccount{}, fmt.Errorf("unmarshal account: %w", err)
	}
	return acc, nil
}

func (c *GRPCClient) BroadcastTx(ctx context.Context, txBytes []byte, mode sdktx.BroadcastMode) (*sdktx.BroadcastTxResponse, error) {
	txClient := sdktx.NewServiceClient(c.conn)
	return txClient.BroadcastTx(ctx, &sdktx.BroadcastTxRequest{
		TxBytes: txBytes,
		Mode:    mode,
	})
}
