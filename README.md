# go-pbl-test / offchain

Test for the offchain Go transactions in [Cardano Go PBL](https://github.com/gimbalabs/cardano-go-pbl) — Module 203.

## Lessons

| Lesson | What it does |
|--------|-------------|
| 203.1  | Build and print Datum / Redeemer structs |
| 203.2  | Mint or burn tokens via native script |
| 203.3  | Lock ADA at a Plutus contract |
| 203.4  | Unlock ADA from a Plutus contract |
| 203.6  | Mint or burn tokens via Plutus script |

## Setup
in the offchain directory
```bash
cp .env.example .env   # add your keys
go run .
```

**.env**
```
BLOCKFROST_PROJECT_ID=preprod...
WALLET_MNEMONIC="word1 word2 ..."
```

## Usage

Uncomment the lesson you want to run in `main.go`, supply any required tx hashes, then:

```bash
go run .
```
