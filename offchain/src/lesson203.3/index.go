package lesson203_3

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Salvionied/apollo"
	"github.com/Salvionied/apollo/constants"
	"github.com/Salvionied/apollo/serialization/Address"
	"github.com/Salvionied/apollo/serialization/PlutusData"
	"github.com/Salvionied/apollo/txBuilding/Backend/BlockFrostChainContext"
)

const (
	BLOCKFROST_KEY = "preprodYOUR_KEY_HERE"
	MNEMONIC       = "word1 word2 word3 word4 word5 word6 word7 word8 word9 word10 word11 word12"
)

type Blueprint struct {
	Validators []struct {
		Title string `json:"title"`
		Hash  string `json:"hash"`
	} `json:"validators"`
}

// scriptAddress derives a preprod enterprise address from a validator hash.
// AddressFromBytes(payment, paymentIsScript, staking, stakingIsScript, network)
func scriptAddress(scriptHash string) Address.Address {
	hashBytes, _ := hex.DecodeString(scriptHash)
	return *Address.AddressFromBytes(hashBytes, true, nil, false, constants.PREPROD)
}

// buildDatum constructs the hello_world Datum.
// Blueprint: constructor 0, fields: [owner: #bytes] → tag 121
// PlutusArray = tagged constructor (Constr); PlutusBytes = raw byte field.
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

func RunLesson203_3Transaction() {
	bfc, err := BlockFrostChainContext.NewBlockfrostChainContext(
		constants.BLOCKFROST_BASE_URL_PREPROD,
		int(constants.PREPROD),
		BLOCKFROST_KEY,
	)
	if err != nil {
		panic(err)
	}

	apollob := apollo.New(&bfc)
	apollob, err = apollob.SetWalletFromMnemonic(MNEMONIC, constants.PREPROD)
	if err != nil {
		panic(err)
	}
	apollob, err = apollob.SetWalletAsChangeAddress()
	if err != nil {
		panic(err)
	}

	// Read the spend validator hash from plutus.json.
	data, err := os.ReadFile("plutus.json")
	if err != nil {
		panic(err)
	}
	var bp Blueprint
	json.Unmarshal(data, &bp)
	var spendHash string
	for _, v := range bp.Validators {
		if v.Title == "lesson203_1.hello_world.spend" {
			spendHash = v.Hash
		}
	}
	if spendHash == "" {
		panic("hello_world.spend not found in plutus.json")
	}

	contractAddress := scriptAddress(spendHash)
	fmt.Println("Script address:", contractAddress.String())

	// GetAddress().PaymentPart is []byte — the wallet's 28-byte payment key hash.
	// Set as datum.owner so only this wallet can satisfy the spend validator.
	walletPkh := apollob.GetWallet().GetAddress().PaymentPart
	datum := buildDatum(walletPkh)

	utxos, err := bfc.Utxos(*apollob.GetWallet().GetAddress())
	if err != nil {
		panic(err)
	}

	// isInline=true stores the datum bytes in the UTxO output itself.
	// The unlock transaction in 203.4 reads it directly from the UTxO.
	apollob, _, err = apollob.
		AddLoadedUTxOs(utxos...).
		PayToContract(contractAddress, &datum, 5_000_000, true).
		Complete()
	if err != nil {
		panic(err)
	}

	apollob = apollob.Sign()
	txId, err := apollob.Submit()
	if err != nil {
		panic(err)
	}

	fmt.Println("Lock tx hash:", hex.EncodeToString(txId.Payload))
	fmt.Println("Save this hash — you need it for 203.4.")
}