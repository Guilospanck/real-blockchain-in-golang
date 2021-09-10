package database

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// SHA-256 produces a 256-bit (32 bytes) hash value.
// It's usually represented as a hexadecimal number of 64 digits.
type Hash [32]byte

func (h Hash) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(h[:])), nil
}

func (h *Hash) UnmarshalText(data []byte) error {
	_, err := hex.Decode(h[:], data)
	return err
}

type Block struct {
	Header BlockHeader `json:"header"`
	TXs    []Tx        `json:"payload"`
}

type BlockHeader struct {
	Parent Hash   `json:"parent"`
	Time   uint64 `json:"time"`
}

type BlockFS struct {
	Key   Hash  `json:"hash"`
	Value Block `json:"block"`
}

func NewBlock(parent Hash, time uint64, txs []Tx) Block {
	return Block{
		BlockHeader{
			Parent: parent,
			Time:   time,
		},
		txs,
	}
}

func (b Block) Hash() (Hash, error) {
	blockJson, err := json.Marshal(b)
	if err != nil {
		return Hash{}, err
	}

	return sha256.Sum256(blockJson), nil
}
