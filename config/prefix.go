package config

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	chain "github.com/crescent-network/crescent/v4/app"
)

// SetAddressPrefixes sets chain specific bech32 prefixes.
func SetAddressPrefixes() *sdk.Config {
	config := sdk.GetConfig()
	config.SetPurpose(chain.Purpose)
	config.SetCoinType(chain.CoinType)
	config.SetBech32PrefixForAccount(chain.Bech32PrefixAccAddr, chain.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(chain.Bech32PrefixValAddr, chain.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(chain.Bech32PrefixConsAddr, chain.Bech32PrefixConsPub)
	return config
}
