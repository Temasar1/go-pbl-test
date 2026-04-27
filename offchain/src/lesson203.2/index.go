package lesson203_2

import (
	"encoding/hex"
	"fmt"

	"github.com/Salvionied/apollo"
	"github.com/Salvionied/apollo/constants"
	"github.com/Salvionied/apollo/serialization/NativeScript"
	"github.com/Salvionied/apollo/txBuilding/Backend/BlockFrostChainContext"

	"go-pbl-test/config"
)

func RunLesson203_2(cfg config.AppConfig) error {
	bfc, err := BlockFrostChainContext.NewBlockfrostChainContext(
		constants.BLOCKFROST_BASE_URL_PREPROD,
		int(constants.PREPROD),
		cfg.BlockfrostProjectID,
	)
	if err != nil {
		return err
	}

	apollob := apollo.New(&bfc)
	apollob, err = apollob.SetWalletFromMnemonic(cfg.WalletMnemonic, constants.PREPROD)
	if err != nil {
		return err
	}
	apollob, err = apollob.SetWalletAsChangeAddress()
	if err != nil {
		return err
	}

	// GetAddress().PaymentPart is []byte — the 28-byte payment key hash.
	pkh := apollob.GetWallet().GetAddress().PaymentPart
	script := NativeScript.NativeScript{
		Tag:     NativeScript.ScriptPubKey,
		KeyHash: pkh,
	}

	// script.Hash() returns (serialization.ScriptHash, error).
	// ScriptHash is [28]byte — use [:] to get []byte for hex encoding.
	scriptHash, err := script.Hash()
	if err != nil {
		return err
	}
	policyId := hex.EncodeToString(scriptHash[:])
	fmt.Println("Policy ID:", policyId)

	mintUnit := apollo.NewUnit(policyId, "GimbalToken", 1_000_000)

	utxos, err := bfc.Utxos(*apollob.GetWallet().GetAddress())
	if err != nil {
		return err
	}

	apollob, _, err = apollob.
		AddLoadedUTxOs(utxos...).
		MintAssetsWithNativeScript(mintUnit, script).
		AddRequiredSignerFromBech32(
			apollob.GetWallet().GetAddress().String(),
			true, false,
		).
		PayToAddressBech32(
			apollob.GetWallet().GetAddress().String(),
			2_000_000,
			mintUnit,
		).
		Complete()
	if err != nil {
		return err
	}

	apollob = apollob.Sign()
	txId, err := apollob.Submit()
	if err != nil {
		return err
	}

	fmt.Println("Tx hash:", hex.EncodeToString(txId.Payload))
	return nil
}
