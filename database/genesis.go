// Used only to read the genesis.json file
package database

import (
	"encoding/json"
	"io/ioutil"
)

var genesisJson = `
{
  "genesis_time": "2021-09-09T19:40:22.415Z",
  "chain_id": "the-blockchain-bar-ledger_TBB",
  "balances": {
    "guilospanck": 1000000
  }
}
`

type genesis struct {
	Balances map[Account]uint `json:"balances"`
}

func loadGenesis(path string) (genesis, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return genesis{}, err
	}

	var loadedGenesis genesis
	err = json.Unmarshal(content, &loadedGenesis)
	if err != nil {
		return genesis{}, err
	}

	return loadedGenesis, nil
}

func writeGenesisToDisk(path string) error {
	return ioutil.WriteFile(path, []byte(genesisJson), 0644)
}
