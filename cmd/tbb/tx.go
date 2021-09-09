package main

import (
	"fmt"
	"os"

	"github.com/Guilospanck/Real-Blockchain-In-Golang/database"
	"github.com/spf13/cobra"
)

const flagFrom = "from"
const flagTo = "to"
const flagValue = "value"
const flagData = "data"

func txCmd() *cobra.Command {
	var txsCmd = &cobra.Command{
		Use:   "tx",
		Short: "Interact with txs (Add...).",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return incorrectUsageErr()
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	txsCmd.AddCommand(txAddCmd())

	return txsCmd
}

func txAddCmd() *cobra.Command {
	var txAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Adds a new TX to database.",
		Run: func(cmd *cobra.Command, args []string) {
			from, _ := cmd.Flags().GetString(flagFrom)
			to, _ := cmd.Flags().GetString(flagTo)
			value, _ := cmd.Flags().GetUint(flagValue)
			data, _ := cmd.Flags().GetString(flagData)

			// creating accounts
			fromAcc := database.NewAccount(from)
			toAcc := database.NewAccount(to)

			tx := database.NewTx(fromAcc, toAcc, value, data)

			state, err := database.NewStateFromDisk()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			defer state.Close() // close db file when everything is executed.

			error := state.Add(tx)
			if error != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			_, error = state.Persist()
			if error != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			fmt.Println("TX succesfully added to the ledger.")

		},
	}

	// info about the flags needed to make a transaction
	txAddCmd.Flags().String(flagFrom, "", "From what account to send tokens")
	txAddCmd.MarkFlagRequired(flagFrom)

	txAddCmd.Flags().String(flagTo, "", "To what account to send tokens")
	txAddCmd.MarkFlagRequired(flagTo)

	txAddCmd.Flags().Uint(flagValue, 0, "How many tokens to send")
	txAddCmd.MarkFlagRequired(flagValue)

	txAddCmd.Flags().String(flagData, "", "Possible values: 'reward'")

	return txAddCmd

}
