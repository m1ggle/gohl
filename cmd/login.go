package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Connect to a MySQL database and start interactive shell",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, dbName, err := GetDBConnection(cmd)
		if err != nil {
			return err
		}
		defer db.Close()

		StartAsyncRefresh(db, dbName)

		fmt.Println("Connected successfully")
		fmt.Println("Enter SQL statements; type 'exit' or 'quit' to disconnect.")

		runInteractiveShell(db)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func getSuggestions(doc prompt.Document) []prompt.Suggest {
	text := strings.ToUpper(strings.TrimSpace(doc.TextBeforeCursor()))
	if len(text) == 0 {
		// 如果没有输入任何字符，则不显示任何提示
		return []prompt.Suggest{}
	}

	// 提供基于当前输入的建议
	words := strings.Fields(text)
	if len(words) > 0 {
		lastWord := strings.ToUpper(words[len(words)-1])
		switch lastWord {
		case "SEL", "SELECT":
			return getAllColumnSuggestions()
		case "SH", "SHOW":
			return []prompt.Suggest{{Text: "TABLES", Description: "Show tables"}, {Text: "DATABASES", Description: "Show databases"}}
		case "DE", "DELETE":
			return []prompt.Suggest{{Text: "FROM", Description: "Delete from table"}}
		case "UP", "UPDATE":
			return getTableSuggestions()
		case "IN", "INSERT":
			return []prompt.Suggest{{Text: "INTO", Description: "Insert into table"}}
		case "CR", "CREATE":
			return []prompt.Suggest{{Text: "TABLE", Description: "Create table"}}
		case "DR", "DROP":
			return []prompt.Suggest{{Text: "TABLE", Description: "Drop table"}}
		case "AL", "ALTER":
			return []prompt.Suggest{{Text: "TABLE", Description: "Alter table"}}
		case "TAB", "TABLE":
			return getTableSuggestions()
		case "DAT", "DATABASE":
			return []prompt.Suggest{{Text: "DATABASE", Description: "Database keyword"}}
		case "FROM", "JOIN", "INTO":
			return getTableSuggestions()
		case "WHERE", "SET", "BY", "ON", "AND", "OR":
			return getAllColumnSuggestions()
		case "US", "USE":
			return []prompt.Suggest{{Text: "mysql", Description: "System database"}}
		case "EX", "EXIT":
			return []prompt.Suggest{{Text: "EXIT", Description: "Exit the shell"}}
		case "QUI", "QUIT":
			return []prompt.Suggest{{Text: "QUIT", Description: "Quit the shell"}}
		case "COM", "COMMIT":
			return []prompt.Suggest{{Text: "COMMIT", Description: "Commit transaction"}}
		case "ROL", "ROLLBACK":
			return []prompt.Suggest{{Text: "ROLLBACK", Description: "Rollback transaction"}}
		}
		// If we are after a keyword but not one of the above, maybe we are after a table name or something else.
		// For now, if we match part of a known table, suggest it.
		if len(words) >= 2 {
			prevWord := strings.ToUpper(words[len(words)-2])
			if prevWord == "FROM" || prevWord == "JOIN" || prevWord == "UPDATE" || prevWord == "INTO" {
				return getTableSuggestions()
			}
		}
	}

	return []prompt.Suggest{}
}

func runInteractiveShell(db *sql.DB) {
	p := prompt.New(
		func(in string) {
			q := strings.TrimSpace(in)
			if q == "" {
				return
			}
			up := strings.ToUpper(q)
			if up == "EXIT" || up == "QUIT" || up == "\\q" {
				fmt.Println("Goodbye")
				os.Exit(0)
			}

			if strings.HasPrefix(up, "SELECT") || strings.HasPrefix(up, "SHOW") || strings.HasPrefix(up, "DESCRIBE") || strings.HasPrefix(up, "EXPLAIN") {
				// Basic heuristic for queries that return rows
				runQuery(db, q)
			} else {
				// Non-query
				res, err := db.Exec(q)
				if err != nil {
					fmt.Fprintln(os.Stderr, "exec error:", err)
					return
				}
				if n, err := res.RowsAffected(); err == nil {
					fmt.Printf("%d rows affected\n", n)
				} else {
					fmt.Println("OK")
				}
			}
		},
		getSuggestions,
		prompt.OptionPrefix("> "),
		prompt.OptionTitle("MySQL CLI"),
		prompt.OptionDescriptionBGColor(prompt.Yellow),
		prompt.OptionDescriptionTextColor(prompt.Black),
		prompt.OptionSelectedDescriptionBGColor(prompt.Blue),
		prompt.OptionSelectedDescriptionTextColor(prompt.White),
		prompt.OptionSuggestionBGColor(prompt.Cyan),
		prompt.OptionSuggestionTextColor(prompt.Black),
		prompt.OptionSelectedSuggestionBGColor(prompt.Green),
		prompt.OptionSelectedSuggestionTextColor(prompt.White),
		prompt.OptionScrollbarBGColor(prompt.DefaultColor),
		prompt.OptionScrollbarThumbColor(prompt.DefaultColor),
		prompt.OptionMaxSuggestion(6),
		prompt.OptionHistory([]string{}),
	)
	p.Run()
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
