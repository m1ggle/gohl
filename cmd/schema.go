package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/c-bata/go-prompt"
)

type SchemaCache struct {
	mu      sync.RWMutex
	Tables  []prompt.Suggest
	Columns map[string][]prompt.Suggest // table -> columns
	DBName  string
}

var globalCache = &SchemaCache{
	Columns: make(map[string][]prompt.Suggest),
}

func (c *SchemaCache) FetchMetadata(db *sql.DB, dbName string) error {
	c.mu.Lock()
	c.DBName = dbName
	c.mu.Unlock()

	// Fetch Tables
	rows, err := db.Query("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ?", dbName)
	if err != nil {
		return fmt.Errorf("failed to fetch tables: %w", err)
	}
	defer rows.Close()

	var tables []prompt.Suggest
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err == nil {
			tables = append(tables, prompt.Suggest{Text: name, Description: "Table"})
		}
	}

	c.mu.Lock()
	c.Tables = tables
	c.mu.Unlock()

	// Fetch Columns
	colRows, err := db.Query("SELECT TABLE_NAME, COLUMN_NAME, COLUMN_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ?", dbName)
	if err != nil {
		return fmt.Errorf("failed to fetch columns: %w", err)
	}
	defer colRows.Close()

	newCols := make(map[string][]prompt.Suggest)
	for colRows.Next() {
		var tableName, colName, colType string
		if err := colRows.Scan(&tableName, &colName, &colType); err == nil {
			newCols[tableName] = append(newCols[tableName], prompt.Suggest{
				Text:        colName,
				Description: fmt.Sprintf("Column (%s)", colType),
			})
		}
	}

	c.mu.Lock()
	c.Columns = newCols
	c.mu.Unlock()

	return nil
}

func StartAsyncRefresh(db *sql.DB, dbName string) {
	go func() {
		if err := globalCache.FetchMetadata(db, dbName); err != nil {
			log.Printf("Metadata refresh failed: %v", err)
		} else {
			// log.Println("Metadata cache populated successfully")
		}
	}()
}

func getTableSuggestions() []prompt.Suggest {
	globalCache.mu.RLock()
	defer globalCache.mu.RUnlock()
	return globalCache.Tables
}

func getColumnSuggestions(tableName string) []prompt.Suggest {
	globalCache.mu.RLock()
	defer globalCache.mu.RUnlock()
	if tableName == "" {
		// Return all columns if no table specified (less performance but more helpful)
		var all []prompt.Suggest
		for _, cols := range globalCache.Columns {
			all = append(all, cols...)
		}
		return all
	}
	return globalCache.Columns[tableName]
}

func getAllColumnSuggestions() []prompt.Suggest {
	globalCache.mu.RLock()
	defer globalCache.mu.RUnlock()
	var all []prompt.Suggest
	seen := make(map[string]bool)
	for _, cols := range globalCache.Columns {
		for _, c := range cols {
			if !seen[c.Text] {
				all = append(all, c)
				seen[c.Text] = true
			}
		}
	}
	return all
}
