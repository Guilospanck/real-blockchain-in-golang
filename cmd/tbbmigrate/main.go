// usage: tbbmigrate
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Guilospanck/Real-Blockchain-In-Golang/database"
	"github.com/spf13/cobra"
)

func main() {
	var tbbCmd = &cobra.Command{
		Use:   "tbbmigrate",
		Short: "Run the migration from tx.db to blocks.db",
		Run: func(cmd *cobra.Command, args []string) {
			migrate()
		},
	}

	err := tbbCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func migrate() {
	cwd, _ := os.Getwd()
	state, err := database.NewStateFromDisk(cwd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer state.Close()

	block0 := database.NewBlock(
		database.Hash{},
		uint64(time.Now().Unix()),
		[]database.Tx{
			database.NewTx("guilospanck", "guilospanck", 3, ""),
			database.NewTx("guilospanck", "guilospanck", 700, "reward"),
		},
	)

	state.AddBlock(block0) // adds transactions to the mempool..
	block0Hash, _ := state.Persist()

	block1 := database.NewBlock(
		block0Hash,
		uint64(time.Now().Unix()),
		[]database.Tx{
			database.NewTx("guilospanck", "babayaga", 2000, ""),
			database.NewTx("guilospanck", "guilospanck", 100, "reward"),
			database.NewTx("babayaga", "guilospanck", 1, ""),
			database.NewTx("babayaga", "caesar", 1000, ""),
			database.NewTx("babayaga", "guilospanck", 50, ""),
			database.NewTx("guilospanck", "guilospanck", 600, "reward"),
		},
	)

	state.AddBlock(block1)
	state.Persist()
}
