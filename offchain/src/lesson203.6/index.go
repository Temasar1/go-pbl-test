package lesson203_6

import (
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/Salvionied/apollo"
	"github.com/Salvionied/apollo/constants"
	"github.com/Salvionied/apollo/serialization/PlutusData"
	"github.com/Salvionied/apollo/txBuilding/Backend/BlockFrostChainContext"

	"go-pbl-test/config"
)

//go:embed plutus.json
var plutusJSON []byte

type blueprint struct {
	Validators []struct {
		Title        string `json:"title"`
		CompiledCode string `json:"compiledCode"`
		Hash         string `json:"hash"`
	} `json:"validators"`
}

func loadScript(title string) (PlutusData.PlutusV3Script, string, error) {
	var bp blueprint
	if err := json.Unmarshal(plutusJSON, &bp); err != nil {
		return nil, "", fmt.Errorf("parsing plutus.json: %w", err)
	}
	for _, v := range bp.Validators {
		if v.Title == title {
			code, err := hex.DecodeString(v.CompiledCode)
			if err != nil {
				return nil, "", fmt.Errorf("decoding compiledCode for %q: %w", title, err)
			}
			return PlutusData.PlutusV3Script(code), v.Hash, nil
		}
	}
	return nil, "", fmt.Errorf("validator %q not found in plutus.json", title)
}

func buildMintRedeemer() PlutusData.PlutusData {
	return PlutusData.PlutusData{
		PlutusDataType: PlutusData.PlutusArray,
		TagNr:          121,
		Value: PlutusData.PlutusIndefArray{
			PlutusData.PlutusData{
				PlutusDataType: PlutusData.PlutusBytes,
				Value:          []byte("HelloMintRedeemer"),
			},
		},
	}
}

// RunLesson203_6Transaction mints or burns tokenName via the hello_world Plutus mint validator.
// amount > 0 mints, amount < 0 burns, amount == 0 is a no-op.
// When burning, pass burnTxHash and burnTxIndex to pin the UTxO holding the tokens.
// Leave burnTxHash as "" when minting.
func RunLesson203_6Transaction(cfg config.AppConfig, amount int, tokenName string, burnTxHash string, burnTxIndex int) error {
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

	mintScript, policyId, err := loadScript("lesson203_1.hello_world.mint")
	if err != nil {
		return err
	}
	fmt.Println("Policy ID:", policyId)

	walletUtxos, err := bfc.Utxos(*apollob.GetWallet().GetAddress())
	if err != nil {
		return err
	}

	unit := apollo.NewUnit(policyId, tokenName, amount)

	apollob = apollob.AddLoadedUTxOs(walletUtxos...)

	if amount < 0 && burnTxHash != "" {
		tokenUtxo, err := apollob.UtxoFromRef(burnTxHash, burnTxIndex)
		if err != nil {
			return fmt.Errorf("fetching token UTxO %s#%d: %w", burnTxHash, burnTxIndex, err)
		}
		apollob = apollob.AddLoadedUTxOs(*tokenUtxo)
	}

	apollob = apollob.
		AttachV3Script(mintScript).
		MintAssetsWithRedeemer(unit, buildMintRedeemer())

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
