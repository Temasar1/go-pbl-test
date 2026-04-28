package lesson203_1

import (
	"fmt"

	"github.com/Salvionied/apollo/serialization/PlutusData"
)

func BuildDatum(owner []byte) PlutusData.PlutusData {
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

func BuildRedeemer(msg []byte) PlutusData.PlutusData {
	return PlutusData.PlutusData{
		PlutusDataType: PlutusData.PlutusArray,
		TagNr:          121,
		Value: PlutusData.PlutusIndefArray{
			PlutusData.PlutusData{
				PlutusDataType: PlutusData.PlutusBytes,
				Value:          msg,
			},
		},
	}
}

func RunLesson203_1() {
	ownerPkh := make([]byte, 28)
	datum := BuildDatum(ownerPkh)
	redeemer := BuildRedeemer([]byte("HelloSpendRedeemer"))

	fmt.Printf("Datum:    %v\n", datum)
	fmt.Printf("Redeemer: %v\n", redeemer)
}
