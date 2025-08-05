package store

import (
	"rorodata_backend_task/models"
	"rorodata_backend_task/models/store/sqlite"
)

type ClusterModelInterface interface {
	CreateCluster(name, region string) (models.Cluster, error)
	DeleteCluster(id int64) error
}

type VMModelInterface interface {
	CreateVM(clusterId int64, name, instanceType string, tags []string) (models.VM, error)
	DeleteVM(clusterId, vmId int64) error
	Operate(operation string, tags []string) ([]models.VM, error)
}

type Models interface {
	ClusterModelInterface
	VMModelInterface
	Close()
}

func New(dsn string) (Models, error) {
	m, err := sqlite.NewSqliteModel(dsn)
	if err != nil {
		return nil, err
	}
	return m, nil
}
