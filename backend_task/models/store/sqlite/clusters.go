package sqlite

import (
	"rorodata_backend_task/models"
	"time"
)

func (m *SqliteModel) CreateCluster(name, region string) (models.Cluster, error) {
	db, err := getDBConnection(m.dsn)
	if err != nil {
		return models.Cluster{}, err
	}
	defer db.Close()
	r, err := db.Exec(create_cluster, name, region, time.Now().UTC().Unix())
	if err != nil {
		return models.Cluster{}, err
	}
	id, _ := r.LastInsertId()

	return models.Cluster{
		ID:     id,
		Name:   name,
		Region: region,
	}, nil
}

func (m *SqliteModel) DeleteCluster(id int64) error {

	db, err := getDBConnection(m.dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	r, err := db.Exec(delete_cluster, id)
	if err != nil {
		if err.Error() == "FOREIGN KEY constraint failed" {
			return models.ErrVMSExists
		}
		return err
	}
	if n, _ := r.RowsAffected(); n == 0 {
		return models.ErrRecordNotFound
	}
	return nil
}
