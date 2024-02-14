package utilities

import (
	"database/sql"
	"log"
)

// GetLanguage is used to retrieve languages from DB
func GetLanguage(psqldb *sql.DB) ([]string, error) {
	var languageCodes []string
	query := "SELECT code FROM public.language"
	row, err := psqldb.Query(query)
	if err != nil {
		return nil, err
	}
	// Iterate through the results and add codes to the slice
	for row.Next() {
		var code string
		if err := row.Scan(&code); err != nil {
			log.Fatal(err)
		}
		languageCodes = append(languageCodes, code)
	}

	// Check for errors from iterating over rows
	if err := row.Err(); err != nil {
		return []string{}, err
	}
	return languageCodes, nil

}
