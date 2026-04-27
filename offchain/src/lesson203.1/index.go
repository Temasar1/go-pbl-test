package lesson203_1

import (
	"fmt"

	"github.com/Salvionied/apollo/serialization/PlutusData"
)

// BuildDatum constructs the Datum for hello_world.
// Blueprint: hello_world/Datum — constructor 0, fields: [owner: #bytes]
// PlutusArray = tagged constructor; PlutusBytes = raw byte field.
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

// BuildRedeemer constructs the Redeemer for hello_world.
// Blueprint: hello_world/Redeemer — constructor 0, fields: [msg: #bytes]
// msg must be []byte("HelloSpendRedeemer") to satisfy the spend handler.
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

// RunLesson203_1 prints the Datum and Redeemer built for the hello_world validator.
func RunLesson203_1() {
	ownerPkh := make([]byte, 28) // replace with real 28-byte payment key hash
	datum := BuildDatum(ownerPkh)
	redeemer := BuildRedeemer([]byte("HelloSpendRedeemer"))

	fmt.Printf("Datum : %v\n", datum)
	fmt.Printf("Redeemer : %v\n", redeemer)
}
