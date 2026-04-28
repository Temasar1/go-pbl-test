package lesson203_3

import (
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Salvionied/apollo"
	"github.com/Salvionied/apollo/constants"
	"github.com/Salvionied/apollo/serialization/Address"
	"github.com/Salvionied/apollo/serialization/PlutusData"
	"github.com/Salvionied/apollo/txBuilding/Backend/BlockFrostChainContext"

	"go-pbl-test/config"
)

//go:embed plutus.json
var plutusJSON []byte

type blueprint struct {
	Validators []struct {
		Title string `json:"title"`
		Hash  string `json:"hash"`
	} `json:"validators"`
}

func scriptAddress(scriptHash string) Address.Address {
	hashBytes, _ := hex.DecodeString(scriptHash)
	return *Address.AddressFromBytes(hashBytes, true, nil, false, constants.PREPROD)
}

func buildDatum(owner []byte) PlutusData.PlutusData {
	return PlutusData.PlutusData{
		PlutusDataType: PlutusData.PlutusArray,
		TagNr:          121,
		Value: PlutusData.PlutusIndefArray{
			PlutusData.PlutusData{
				PlutusDataType: PlutusData.PlutusBytes,
				Value:          owner,
			},
		},
	}
}

func RunLesson203_3Transaction(cfg config.AppConfig) error {
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

	var bp blueprint
	if err = json.Unmarshal(plutusJSON, &bp); err != nil {
		return fmt.Errorf("parsing plutus.json: %w", err)
	}

	var spendHash string
	for _, v := range bp.Validators {
		if v.Title == "lesson203_1.hello_world.spend" {
			spendHash = v.Hash
			break
		}
	}
	if spendHash == "" {
		return errors.New("hello_world.spend not found in plutus.json")
	}

	contractAddress := scriptAddress(spendHash)
	fmt.Println("Script address:", contractAddress.String())

	walletPkh := apollob.GetWallet().GetAddress().PaymentPart
	datum := buildDatum(walletPkh)

	utxos, err := bfc.Utxos(*apollob.GetWallet().GetAddress())
	if err != nil {
		return err
	}

	apollob, _, err = apollob.
		AddLoadedUTxOs(utxos...).
		PayToContract(contractAddress, &datum, 5_000_000, true).
		Complete()
	if err != nil {
		return err
	}

	apollob = apollob.Sign()
	txId, err := apollob.Submit()
	if err != nil {
		return err
	}

	fmt.Println("Lock tx hash:", hex.EncodeToString(txId.Payload))
	return nil
}
