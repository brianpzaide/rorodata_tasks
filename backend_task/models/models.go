package models

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrVMSExists      = errors.New("virtual machines exist in this cluster. Please delete the VMS before deleting the cluster")
	ErrNoClusterFound = errors.New("No such cluster exists")
)

type Cluster struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Region    string    `json:"region"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type VM struct {
	ID           int64     `json:"id"`
	ClusterID    int64     `json:"cluster_id"`
	Name         string    `json:"name"`
	State        string    `json:"state"`
	Tags         []string  `json:"tags,omitempty"`
	InstanceType string    `json:"instance_type"`
	IP           string    `json:"ip"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
}

type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func GenerateRandomIP() string {
	ip := fmt.Sprintf("%d%d%d%d.%d%d%d%d.%d%d%d%d.%d%d%d%d",
		rand.IntN(10),
		rand.IntN(10),
		rand.IntN(10),
		rand.IntN(10),
		rand.IntN(10),
		rand.IntN(10),
		rand.IntN(10),
		rand.IntN(10),
		rand.IntN(10),
		rand.IntN(10),
		rand.IntN(10),
		rand.IntN(10),
		rand.IntN(10),
		rand.IntN(10),
		rand.IntN(10),
		rand.IntN(10),
	)
	return ip
}
