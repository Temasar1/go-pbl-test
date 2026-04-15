package main

import (
	"encoding/hex"

	"github.com/Salvionied/apollo/serialization"
	"github.com/Salvionied/apollo/serialization/PlutusData"
)

func buildMintRedeemer(action []byte) PlutusData.PlutusData {
	return PlutusData.PlutusData{
		TagNr: 121, // constructor index 0
		Value: PlutusData.PlutusIndefArray{
			PlutusData.PlutusData{
				Value: serialization.NewCustomBytes(hex.EncodeToString(action)),
			},
		},
	}
}
