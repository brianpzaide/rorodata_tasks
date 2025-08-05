package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const create_clusters_table = `CREATE TABLE IF NOT EXISTS clusters (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    region TEXT,
    createdAt INTEGER NOT NULL
);`

const create_vms_table = `CREATE TABLE IF NOT EXISTS vms (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    clusterId INTEGER REFERENCES clusters(id),
	name TEXT,
    instanceType TEXT,
    state TEXT CHECK (state IN ('running', 'stopped', 'rebooting')),
    ip TEXT,
	createdAt INTEGER NOT NULL
);`

const create_tags_table = `CREATE TABLE IF NOT EXISTS tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE
);`

const create_vms_tags_table = `CREATE TABLE IF NOT EXISTS vms_tags (
    tag_id INTEGER NOT NULL REFERENCES tags(id),
    vm_id INTEGER NOT NULL REFERENCES vms(id)
);`

const create_cluster = `INSERT INTO clusters (name, region, createdAt) VALUES (?, ?, ?);`

const delete_cluster = `DELETE FROM clusters WHERE id = ?;`

const create_vm = `INSERT INTO vms (clusterId, ip, name, instanceType, state, createdAt) VALUES (?, ?, ?, ?, ?, ?);`

const delete_vm = `DELETE FROM vms WHERE clusterId = ? AND id = ?;`

const delete_vm_tag_association = `DELETE FROM vms_tags WHERE vm_id = ?;`

const start_vm = `WITH aux_vmid AS (
  SELECT DISTINCT vm_id
  FROM vms_tags
  INNER JOIN tags ON vms_tags.tag_id = tags.id
  WHERE tags.name IN (%s)
)
UPDATE vms
SET state = 'running'
WHERE vms.id IN (SELECT vm_id FROM aux_vmid)
RETURNING *;`

const stop_vm = `WITH aux_vmid AS (
  SELECT DISTINCT vm_id
  FROM vms_tags
  INNER JOIN tags ON vms_tags.tag_id = tags.id
  WHERE tags.name IN (%s)
)
UPDATE vms
SET state = 'stopped'
WHERE vms.id IN (SELECT vm_id FROM aux_vmid)
RETURNING *;`

const reboot_vm = `WITH aux_vmid AS (
  SELECT DISTINCT vm_id
  FROM vms_tags
  INNER JOIN tags ON vms_tags.tag_id = tags.id
  WHERE tags.name IN (%s)
)
UPDATE vms
SET state = 'rebooting'
WHERE vms.id IN (SELECT vm_id FROM aux_vmid)
RETURNING *;`

const create_tag = `INSERT INTO tags (name) VALUES (?);`

const fetch_tag_id = `SELECT id from tags WHERE name = ?;`

const add_tags_to_vm = `INSERT INTO vms_tags (tag_id, vm_id) VALUES %s;`

type SqliteModel struct {
	dsn string
}

func getDBConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewSqliteModel(dsn string) (*SqliteModel, error) {
	db, err := getDBConnection(dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	_, err = db.Exec(create_clusters_table)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(create_vms_table)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(create_tags_table)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(create_vms_tags_table)
	if err != nil {
		return nil, err
	}

	return &SqliteModel{dsn: dsn}, nil
}

func (m *SqliteModel) Close() {

}
