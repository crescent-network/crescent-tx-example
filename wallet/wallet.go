package wallet

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
)

func RecoverPrivKeyFromMnemonic(mnemonic, password string) (*secp256k1.PrivKey, error) {
	seed := bip39.NewSeed(mnemonic, password)
	master, ch := hd.ComputeMastersFromSeed(seed)
	priv, err := hd.DerivePrivateKeyForPath(master, ch, sdk.GetConfig().GetFullBIP44Path())
	if err != nil {
		return nil, fmt.Errorf("failed to derive private key: %w", err)
	}
	return &secp256k1.PrivKey{Key: priv}, nil
}

func Address(privKey *secp256k1.PrivKey) sdk.AccAddress {
	return sdk.AccAddress(privKey.PubKey().Address())
}
