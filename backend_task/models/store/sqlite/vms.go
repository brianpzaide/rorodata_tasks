package sqlite

import (
	"errors"
	"fmt"
	"rorodata_backend_task/models"
	"strings"
	"time"
)

const ()

func (m *SqliteModel) CreateVM(clusterId int64, name, instanceType string, tags []string) (models.VM, error) {
	db, err := getDBConnection(m.dsn)
	if err != nil {
		return models.VM{}, err
	}
	defer db.Close()
	ip := models.GenerateRandomIP()
	r, err := db.Exec(create_vm, clusterId, ip, name, instanceType, "running", time.Now().UTC().Unix())
	if err != nil {
		if err.Error() == "FOREIGN KEY constraint failed" {
			return models.VM{}, models.ErrNoClusterFound
		}
		return models.VM{}, err
	}
	id, _ := r.LastInsertId()

	insertTags(db, id, tags)

	return models.VM{
		ID:           id,
		IP:           ip,
		ClusterID:    clusterId,
		Name:         name,
		State:        "running",
		InstanceType: instanceType,
		Tags:         tags,
	}, nil
}

func (m *SqliteModel) DeleteVM(clusterId, vmId int64) error {
	db, err := getDBConnection(m.dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	db.Exec(delete_vm_tag_association, vmId)

	r, err := db.Exec(delete_vm, clusterId, vmId)
	if err != nil {
		fmt.Println(err.Error())
		if err.Error() == "FOREIGN KEY constraint failed" {
			return models.ErrNoClusterFound
		}
		return err
	}
	if n, _ := r.RowsAffected(); n == 0 {
		return models.ErrRecordNotFound
	}
	return nil
}

func (m *SqliteModel) Operate(operation string, tags []string) ([]models.VM, error) {
	switch operation {
	case "start":
		return m.startVM(tags)
	case "stop":
		return m.stopVM(tags)
	case "reboot":
		return m.rebootVM(tags)
	default:
		return nil, errors.ErrUnsupported
	}
}

func (m *SqliteModel) startVM(tags []string) ([]models.VM, error) {
	db, err := getDBConnection(m.dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	placeholder := genPlaceholder(tags)

	query := fmt.Sprintf(start_vm, placeholder)

	args := convertToInterfaceSliceVMS(tags)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	vms := make([]models.VM, 0)

	for rows.Next() {
		var vm models.VM
		var createdAt int64
		var err = rows.Scan(
			&vm.ID,
			&vm.ClusterID,
			&vm.Name,
			&vm.InstanceType,
			&vm.State,
			&vm.IP,
			&createdAt,
		)
		if err != nil {
			return nil, err
		}
		vms = append(vms, vm)
	}

	return vms, nil
}

func (m *SqliteModel) stopVM(tags []string) ([]models.VM, error) {
	db, err := getDBConnection(m.dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	placeholder := genPlaceholder(tags)

	query := fmt.Sprintf(stop_vm, placeholder)

	args := convertToInterfaceSliceVMS(tags)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	vms := make([]models.VM, 0)

	for rows.Next() {
		var vm models.VM
		var createdAt int64
		var err = rows.Scan(
			&vm.ID,
			&vm.ClusterID,
			&vm.Name,
			&vm.InstanceType,
			&vm.State,
			&vm.IP,
			&createdAt,
		)
		if err != nil {
			return nil, err
		}
		vms = append(vms, vm)
	}

	return vms, nil
}

func (m *SqliteModel) rebootVM(tags []string) ([]models.VM, error) {
	db, err := getDBConnection(m.dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	placeholder := genPlaceholder(tags)

	query := fmt.Sprintf(reboot_vm, placeholder)

	args := convertToInterfaceSliceVMS(tags)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	vms := make([]models.VM, 0)

	for rows.Next() {
		var vm models.VM
		var createdAt int64
		var err = rows.Scan(
			&vm.ID,
			&vm.ClusterID,
			&vm.Name,
			&vm.InstanceType,
			&vm.State,
			&vm.IP,
			&createdAt,
		)
		if err != nil {
			return nil, err
		}
		vms = append(vms, vm)
	}

	go func() {
		time.Sleep(5 * time.Second)
		m.startVM(tags)
	}()

	return vms, nil
}

func genPlaceholder(tags []string) string {
	placeHolders := make([]string, 0)

	for _, _ = range tags {
		placeHolders = append(placeHolders, "?")
	}

	return strings.Join(placeHolders, ", ")
}

func convertToInterfaceSliceVMS(vals []string) []interface{} {
	args := make([]interface{}, len(vals))
	for i, v := range vals {
		args[i] = v
	}
	return args
}
