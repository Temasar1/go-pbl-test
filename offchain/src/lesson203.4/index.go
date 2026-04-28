package lesson203_4

import (
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/Salvionied/apollo"
	"github.com/Salvionied/apollo/constants"
	"github.com/Salvionied/apollo/serialization/PlutusData"
	"github.com/Salvionied/apollo/serialization/Redeemer"
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

func loadScript(title string) (PlutusData.PlutusV3Script, error) {
	var bp blueprint
	if err := json.Unmarshal(plutusJSON, &bp); err != nil {
		return nil, fmt.Errorf("parsing plutus.json: %w", err)
	}
	for _, v := range bp.Validators {
		if v.Title == title {
			code, err := hex.DecodeString(v.CompiledCode)
			if err != nil {
				return nil, fmt.Errorf("decoding compiledCode for %q: %w", title, err)
			}
			return PlutusData.PlutusV3Script(code), nil
		}
	}
	return nil, fmt.Errorf("validator %q not found in plutus.json", title)
}

func buildSpendRedeemer() Redeemer.Redeemer {
	return Redeemer.Redeemer{
		Tag: Redeemer.SPEND,
		Data: PlutusData.PlutusData{
			PlutusDataType: PlutusData.PlutusArray,
			TagNr:          121,
			Value: PlutusData.PlutusIndefArray{
				PlutusData.PlutusData{
					PlutusDataType: PlutusData.PlutusBytes,
					Value:          []byte("HelloSpendRedeemer"),
				},
			},
		},
		ExUnits: Redeemer.ExecutionUnits{Mem: 0, Steps: 0},
	}
}

// RunLesson203_4Transaction unlocks the UTxO locked by RunLesson203_3Transaction.
// lockTxHash is the tx hash printed by 203.3; lockTxIndex is almost always 0.
func RunLesson203_4Transaction(cfg config.AppConfig, lockTxHash string, lockTxIndex int) error {
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

	spendScript, err := loadScript("lesson203_1.hello_world.spend")
	if err != nil {
		return err
	}

	scriptUtxo, err := apollob.UtxoFromRef(lockTxHash, lockTxIndex)
	if err != nil {
		return fmt.Errorf("fetching script UTxO %s#%d: %w", lockTxHash, lockTxIndex, err)
	}

	walletUtxos, err := bfc.Utxos(*apollob.GetWallet().GetAddress())
	if err != nil {
		return err
	}

	redeemer := buildSpendRedeemer()

	apollob, _, err = apollob.
		AddLoadedUTxOs(walletUtxos...).
		CollectFrom(*scriptUtxo, redeemer).
		AddRequiredSignerFromBech32(apollob.GetWallet().GetAddress().String(), true, false).
		AttachV3Script(spendScript).
		PayToAddressBech32(apollob.GetWallet().GetAddress().String(), 4_500_000).
		Complete()
	if err != nil {
		return err
	}

	apollob = apollob.Sign()
	txId, err := apollob.Submit()
	if err != nil {
		return err
	}

	fmt.Println("Unlock tx hash:", hex.EncodeToString(txId.Payload))
	return nil
}
