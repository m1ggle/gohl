package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:   "query [sql]",
	Short: "Execute a single SQL query and exit",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := GetDBConnection(cmd)
		if err != nil {
			return err
		}
		defer db.Close()

		q := args[0]
		up := strings.ToUpper(strings.TrimSpace(q))

		if strings.HasPrefix(up, "SELECT") || strings.HasPrefix(up, "SHOW") || strings.HasPrefix(up, "DESCRIBE") || strings.HasPrefix(up, "EXPLAIN") {
			runQuery(db, q) // reusing the helper from login.go if exported, or duplicate.
			// Since they are in the same package 'cmd', I can share functions if they are lowercase? Yes.
			// But I extracted 'runQuery' in login.go. I should make sure it's accessible.
			// In Go, functions in the same package are visible to each other.
		} else {
			res, err := db.Exec(q)
			if err != nil {
				return fmt.Errorf("exec error: %w", err)
			}
			if n, err := res.RowsAffected(); err == nil {
				fmt.Printf("%d rows affected\n", n)
			} else {
				fmt.Println("OK")
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)
}
