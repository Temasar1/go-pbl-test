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

// RunLesson203_2 mints or burns tokenName via a native-script policy tied to the wallet key.
// amount > 0 mints, amount < 0 burns, amount == 0 is a no-op.
func RunLesson203_2(cfg config.AppConfig, amount int, tokenName string) error {
	if amount == 0 {
		return nil
	}

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

	pkh := apollob.GetWallet().GetAddress().PaymentPart
	script := NativeScript.NativeScript{
		Tag:     NativeScript.ScriptPubKey,
		KeyHash: pkh,
	}

	scriptHash, err := script.Hash()
	if err != nil {
		return err
	}
	policyId := hex.EncodeToString(scriptHash[:])
	fmt.Println("Policy ID:", policyId)

	unit := apollo.NewUnit(policyId, tokenName, amount)

	utxos, err := bfc.Utxos(*apollob.GetWallet().GetAddress())
	if err != nil {
		return err
	}

	apollob = apollob.
		AddLoadedUTxOs(utxos...).
		MintAssetsWithNativeScript(unit, script).
		AddRequiredSignerFromBech32(apollob.GetWallet().GetAddress().String(), true, false)

	if amount > 0 {
		apollob = apollob.PayToAddressBech32(apollob.GetWallet().GetAddress().String(), 2_000_000, unit)
	}

	apollob, _, err = apollob.Complete()
	if err != nil {
		return err
	}

	apollob = apollob.Sign()
	txId, err := apollob.Submit()
	if err != nil {
		return err
	}

	action := "Mint"
	if amount < 0 {
		action = "Burn"
	}
	fmt.Printf("%s tx hash: %s\n", action, hex.EncodeToString(txId.Payload))
	return nil
}
