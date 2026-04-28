package main

import (
	"fmt"
	"os"

	"go-pbl-test/config"
	lesson203_1 "go-pbl-test/src/lesson203.1"
	lesson203_2 "go-pbl-test/src/lesson203.2"
	lesson203_3 "go-pbl-test/src/lesson203.3"
	lesson203_4 "go-pbl-test/src/lesson203.4"
	lesson203_6 "go-pbl-test/src/lesson203.6"
)

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// 203.1 — print Datum / Redeemer structs (no network)
	// lesson203_1.RunLesson203_1()

	// 203.2 — mint or burn via native script
	// lesson203_2.RunLesson203_2(cfg, 1_000_000, "GimbalToken")
	// lesson203_2.RunLesson203_2(cfg, -50_000, "GimbalToken")

	// 203.3 — lock 5 ADA at the hello_world contract
	// lesson203_3.RunLesson203_3Transaction(cfg)

	// 203.4 — unlock the UTxO locked by 203.3
	// lesson203_4.RunLesson203_4Transaction(cfg, "LOCK_TX_HASH", 0)

	// 203.6 — mint or burn via Plutus script
	// lesson203_6.RunLesson203_6Transaction(cfg, 1_000, "HelloToken", "", 0)
	// lesson203_6.RunLesson203_6Transaction(cfg, -500, "HelloToken", "MINT_TX_HASH", 0)

	_ = lesson203_1.RunLesson203_1
	_ = lesson203_2.RunLesson203_2
	_ = lesson203_3.RunLesson203_3Transaction
	_ = lesson203_4.RunLesson203_4Transaction
	_ = lesson203_6.RunLesson203_6Transaction
	_ = cfg
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
