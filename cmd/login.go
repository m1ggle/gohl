package cmd

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Connect to a MySQL database and start interactive shell",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := GetDBConnection(cmd)
		if err != nil {
			return err
		}
		defer db.Close()

		fmt.Println("Connected successfully")
		fmt.Println("Enter SQL statements; type 'exit' or 'quit' to disconnect.")

		runInteractiveShell(db)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func runInteractiveShell(db *sql.DB) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "input error:", err)
			break
		}
		q := strings.TrimSpace(line)
		if q == "" {
			continue
		}
		up := strings.ToUpper(q)
		if up == "EXIT" || up == "QUIT" || up == "\\q" {
			fmt.Println("Goodbye")
			break
		}

		if strings.HasPrefix(up, "SELECT") || strings.HasPrefix(up, "SHOW") || strings.HasPrefix(up, "DESCRIBE") || strings.HasPrefix(up, "EXPLAIN") {
			// Basic heuristic for queries that return rows
			runQuery(db, q)
		} else {
			// Non-query
			res, err := db.Exec(q)
			if err != nil {
				fmt.Fprintln(os.Stderr, "exec error:", err)
				continue
			}
			if n, err := res.RowsAffected(); err == nil {
				fmt.Printf("%d rows affected\n", n)
			} else {
				fmt.Println("OK")
			}
		}
	}
}

func runQuery(db *sql.DB, q string) {
	rows, err := db.Query(q)
	if err != nil {
		fmt.Fprintln(os.Stderr, "query error:", err)
		return
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		fmt.Fprintln(os.Stderr, "columns error:", err)
		return
	}
	vals := make([]sql.RawBytes, len(cols))
	scanArgs := make([]interface{}, len(vals))
	for i := range vals {
		scanArgs[i] = &vals[i]
	}

	table := tablewriter.NewTable(os.Stdout,
		tablewriter.WithHeaderAlignment(tw.AlignLeft),
		tablewriter.WithRowAlignment(tw.AlignLeft))
	table.Header(cols)

	// table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	// table.SetAlignment(tablewriter.ALIGN_LEFT)
	// table.SetCenterSeparator("|")
	// table.SetColumnSeparator("|")
	// table.SetRowSeparator("-")
	// table.SetHeaderLine(true)
	// table.Border(true)
	// table.SetTablePadding("\t") // pad with tabs for consistency if needed, but space is better for tables

	count := 0
	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			fmt.Fprintln(os.Stderr, "scan error:", err)
			return
		}
		row := make([]string, len(vals))
		for i, col := range vals {
			if col == nil {
				row[i] = "NULL"
			} else {
				row[i] = string(col)
			}
		}
		table.Append(row)
		count++
	}
	table.Render()
	fmt.Printf("%d rows in set\n", count)
}
