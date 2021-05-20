// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import "github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/models"

// CreateInstance CRUD operations for instances
type Instance interface {
	// CreateInstance will create instance and return nil if no error
	CreateInstance(instanceBody models.CreateInstanceBody) error
	// GetInstance will fetch instance details as per ID
	GetInstance(id int) (models.GetInstanceResponse, error)
}
