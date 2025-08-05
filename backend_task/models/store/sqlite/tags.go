package sqlite

import (
	"database/sql"
	"fmt"
	"strings"
)

func insertTags(db *sql.DB, vmId int64, tags []string) error {

	placeHolders := make([]string, 0)
	values := make([]int64, 0)

	for _, tag := range tags {
		r, err := db.Exec(create_tag, tag)
		if err != nil {
			if err.Error() == "UNIQUE constraint failed: tags.name" {
				row := db.QueryRow(fetch_tag_id, tag)
				var id int64
				if err := row.Scan(&id); err != nil {
					return fmt.Errorf("failed to fetch tag id for '%s': %w", tag, err)
				}
				placeHolders = append(placeHolders, "(?, ?)")
				values = append(values, id)
				values = append(values, vmId)
			} else {
				return err
			}
		} else {
			id, err := r.LastInsertId()
			if err != nil {
				return fmt.Errorf("could not get last insert id: %w", err)
			}
			placeHolders = append(placeHolders, "(?, ?)")
			values = append(values, id)
			values = append(values, vmId)
		}
	}

	query := fmt.Sprintf(add_tags_to_vm, strings.Join(placeHolders, ", "))

	args := convertToInterfaceSlice(values)
	db.Exec(query, args...)

	return nil
}

func convertToInterfaceSlice(vals []int64) []interface{} {
	args := make([]interface{}, len(vals))
	for i, v := range vals {
		args[i] = v
	}
	return args
}
