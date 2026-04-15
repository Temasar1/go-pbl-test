package lesson203_3

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/Salvionied/apollo"
	"github.com/Salvionied/apollo/constants"
	"github.com/Salvionied/apollo/serialization/PlutusData"
	"github.com/Salvionied/apollo/serialization/Redeemer"
	"github.com/Salvionied/apollo/txBuilding/Backend/BlockFrostChainContext"
	"github.com/fxamacker/cbor/v2"
	"github.com/joho/godotenv"
)

const compiledCode = "585401010029800aba2aba1aab9eaab9dab9a4888896600264653001300600198031803800cc0180092225980099b8748000c01cdd500144c9289bae30093008375400516401830060013003375400d149a26cac8009"

func RunLesson203_3Transaction() error {
	godotenv.Load()

	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	bfc, err := BlockFrostChainContext.NewBlockfrostChainContext(
		constants.BLOCKFROST_BASE_URL_PREPROD,
		int(constants.PREPROD),
		cfg.BlockfrostProjectID,
	)
	if err != nil {
		return fmt.Errorf("blockfrost init: %w", err)
	}

	scriptBytes, err := hex.DecodeString(compiledCode)
	if err != nil {
		return fmt.Errorf("decode script hex: %w", err)
	}
	mintingScript := PlutusData.PlutusV3Script(scriptBytes)

	scriptHash, err := mintingScript.Hash()
	if err != nil {
		return fmt.Errorf("script hash: %w", err)
	}
	policyID := hex.EncodeToString(scriptHash.Bytes())
	fmt.Println("Policy ID:", policyID)

	cc := apollo.NewEmptyBackend()
	builder := apollo.New(&cc)

	builder, err = builder.SetWalletFromMnemonic(cfg.WalletMnemonic, constants.PREPROD)
	if err != nil {
		return fmt.Errorf("set wallet: %w", err)
	}
	builder, err = builder.SetWalletAsChangeAddress()
	if err != nil {
		return fmt.Errorf("set change address: %w", err)
	}

	walletAddr := builder.GetWallet().GetAddress()
	utxos, err := bfc.Utxos(*walletAddr)
	if err != nil {
		return fmt.Errorf("fetch utxos: %w", err)
	}
	if len(utxos) == 0 {
		return fmt.Errorf("no UTxOs at %s", walletAddr.String())
	}
	fmt.Printf("Found %d UTxO(s)\n", len(utxos))

	redeemer := PlutusData.PlutusData{Value: 1}
	mintUnit := apollo.NewUnit(policyID, "gimbalabs_go_test_token", 1000)

	builder, _, err = builder.
		AddLoadedUTxOs(utxos...).
		AttachV3Script(mintingScript).
		MintAssetsWithRedeemer(mintUnit, redeemer).
		PayToAddressBech32(walletAddr.String(), 2_000_000, mintUnit).
		Complete()
	if err != nil {
		return fmt.Errorf("tx build failed: %w", err)
	}

	// --- Evaluate first ---
	tx := builder.GetTx()
	txCbor, err := cbor.Marshal(tx)
	if err != nil {
		return fmt.Errorf("marshal for eval: %w", err)
	}

	eval, err := bfc.EvaluateTx(txCbor)
	if err != nil {
		fmt.Printf("⚠️  EvaluateTx error: %v\n", err)
	} else {
		fmt.Printf("✅ EvaluateTx result: %+v\n", eval)
	}

	// --- Force-set ExUnits on ALL redeemers ---
	redeemers := builder.GetRedeemers()
	fmt.Printf("Redeemers count: %d\n", len(redeemers))

	for key, r := range redeemers {
		fmt.Printf("  Redeemer key=%s mem=%d steps=%d\n", key, r.ExUnits.Mem, r.ExUnits.Steps)

		// Try eval result first
		if eval != nil {
			if exUnits, ok := eval[key]; ok {
				r.ExUnits = exUnits
				redeemers[key] = r
				fmt.Printf("  → Applied eval ExUnits: mem=%d steps=%d\n", exUnits.Mem, exUnits.Steps)
				continue
			}
		}

		// Always override — never leave at 0
		r.ExUnits = Redeemer.ExecutionUnits{
			Mem:   14_000_000,
			Steps: 10_000_000_000,
		}
		redeemers[key] = r
		fmt.Println("  → Applied fallback ExUnits")
	}

	builder = builder.UpdateRedeemers(redeemers)

	// Verify ExUnits were actually set
	finalRedeemers := builder.GetRedeemers()
	for key, r := range finalRedeemers {
		fmt.Printf("Final redeemer key=%s mem=%d steps=%d\n", key, r.ExUnits.Mem, r.ExUnits.Steps)
	}

	builder = builder.Sign()
	signedTx := builder.GetTx()

	signedCbor, err := cbor.Marshal(signedTx)
	if err != nil {
		return fmt.Errorf("marshal signed tx: %w", err)
	}

	// Decode the CBOR to confirm ExUnits in the final tx before submitting
	fmt.Println("Signed Tx CBOR:", hex.EncodeToString(signedCbor))

	txID, err := bfc.SubmitTx(*signedTx)
	if err != nil {
		return fmt.Errorf("submit failed: %w", err)
	}

	fmt.Println("✅ Tx Hash:", hex.EncodeToString(txID.Payload))
	return nil
}

type AppConfig struct {
	BlockfrostProjectID string
	WalletMnemonic      string
}

func loadConfig() (AppConfig, error) {
	cfg := AppConfig{
		BlockfrostProjectID: os.Getenv("BLOCKFROST_PROJECT_ID"),
		WalletMnemonic:      os.Getenv("WALLET_MNEMONIC"),
	}
	if cfg.BlockfrostProjectID == "" {
		return cfg, fmt.Errorf("missing BLOCKFROST_PROJECT_ID")
	}
	if cfg.WalletMnemonic == "" {
		return cfg, fmt.Errorf("missing WALLET_MNEMONIC")
	}
	return cfg, nil
}