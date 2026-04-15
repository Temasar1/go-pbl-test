package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
)

type Blueprint struct {
	Validators []struct {
		Title        string `json:"title"`
		CompiledCode string `json:"compiledCode"`
		Hash         string `json:"hash"`
	} `json:"validators"`
}

func loadScript(path, title string) (compiledCode []byte, scriptHash string, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, "", err
	}

	var bp Blueprint
	if err := json.Unmarshal(data, &bp); err != nil {
		return nil, "", err
	}

	for _, v := range bp.Validators {
		if v.Title == title {
			code, err := hex.DecodeString(v.CompiledCode)
			if err != nil {
				return nil, "", err
			}
			return code, v.Hash, nil
		}
	}

	return nil, "", fmt.Errorf("validator %q not found", title)
}
