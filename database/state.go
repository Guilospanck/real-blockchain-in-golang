// The State Component is responsible for:
//  • Adding new transactions to Mempool
//  • Validating transactions against the current State (sufficient sender balance)
//  • Changing the state
//  • Persisting transactions to disk
//  • Calculating accounts balances by replaying all transactions since Genesis in a sequence

package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

/*
  State of the blockchain.
  - Balances of the accounts
  - Size of the mempool
  - where is the database metadata to save
*/
type State struct {
	Balances  map[Account]uint
	txMempool []Tx

	dbFile          *os.File
	latestBlockHash Hash
}

func NewStateFromDisk() (*State, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	genFilePath := filepath.Join(cwd, "database", "genesis.json")
	gen, err := loadGenesis(genFilePath)
	if err != nil {
		return nil, err
	}

	balances := make(map[Account]uint)
	for account, balance := range gen.Balances {
		balances[account] = balance
	}

	blockDbFilePath := filepath.Join(cwd, "database", "block.db")
	f, err := os.OpenFile(blockDbFilePath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)

	state := &State{
		Balances:        balances,
		txMempool:       make([]Tx, 0),
		dbFile:          f,
		latestBlockHash: Hash{},
	}

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		blockFsJson := scanner.Bytes()
		var blockFs BlockFS
		err = json.Unmarshal(blockFsJson, &blockFs)
		if err != nil {
			return nil, err
		}

		if err := state.applyBlock(blockFs.Value); err != nil {
			return nil, err
		}

		state.latestBlockHash = blockFs.Key
	}

	return state, nil
}

func (s *State) apply(tx Tx) error {
	if tx.IsReward() {
		s.Balances[tx.To] += tx.Value
		return nil
	}

	if s.Balances[tx.From] < tx.Value {
		return fmt.Errorf("insufficient balance")
	}

	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value

	return nil
}

func (s *State) applyBlock(block Block) error {
	for _, tx := range block.TXs {
		if err := s.apply(tx); err != nil {
			return err
		}
	}

	return nil
}

func (s *State) AddBlock(block Block) error {
	for _, tx := range block.TXs {
		if err := s.AddTx(tx); err != nil {
			return err
		}
	}

	return nil
}

func (s *State) AddTx(tx Tx) error {
	if err := s.apply(tx); err != nil {
		return err
	}

	s.txMempool = append(s.txMempool, tx)

	return nil
}

func (s *State) LatestBlockHash() Hash {
	return s.latestBlockHash
}

func (s *State) Persist() (Hash, error) {
	// create a new Block with only the new TXs
	block := NewBlock(
		s.latestBlockHash,
		uint64(time.Now().Unix()),
		s.txMempool,
	)

	blockHash, err := block.Hash()
	if err != nil {
		return Hash{}, nil
	}

	blockFS := BlockFS{Key: blockHash, Value: block}

	blockFsJson, err := json.Marshal(blockFS)
	if err != nil {
		return Hash{}, nil
	}

	fmt.Printf("Persisting new Block to disk: \n")
	fmt.Printf("\t%s\n", blockFsJson)

	if _, err := s.dbFile.Write(append(blockFsJson, '\n')); err != nil {
		return Hash{}, err
	}

	s.latestBlockHash = blockHash

	// Freeing mempool
	s.txMempool = []Tx{}

	return s.latestBlockHash, nil
}

func (s *State) Close() {
	s.dbFile.Close()
}
