// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/models"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/logger"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
)

// snapshot implements functions related to cmp instances
type snapshot struct {
	// expose Instance API service to instances related operations
	iClient *client.InstancesApiService
}

func newSnapshot(iClient *client.InstancesApiService) *snapshot {
	return &snapshot{
		iClient: iClient,
	}
}

// Create snapshot
func (s *snapshot) Create(ctx context.Context, d *utils.Data) error {
	logger.Debug("Creating VMware snapshot of instance")
	instanceID := d.GetInt("instance_id")
	req := &models.SnapshotBody{
		Snapshot: &models.SnapshotBodySnapshot{
			Name:        d.GetString("name"),
			Description: d.GetString("description"),
		},
	}

	// Pre check
	if err := d.Error(); err != nil {
		return err
	}
	// create snapshot
	resp, err := utils.Retry(func() (interface{}, error) {
		return s.iClient.SnapshotAnInstance(ctx, instanceID, req)
	})
	if err != nil {
		return err
	}
	snapshotResp := resp.(models.Instances)
	if !snapshotResp.Success {
		return fmt.Errorf("%t", snapshotResp.Success)
	}

	// post check
	return d.Error()
}

// Read snapshot and set state values accordingly
func (s *snapshot) Read(ctx context.Context, d *utils.Data) error {
	instanceID := d.GetInt("instance_id")

	// Precheck
	if err := d.Error(); err != nil {
		return err
	}

	resp, err := utils.Retry(func() (interface{}, error) {
		return s.iClient.GetListOfSnapshotsForAnInstance(ctx, instanceID)
	})
	if err != nil {
		return err
	}
	snapshots := resp.(models.ListSnapshotResponse)
	d.SetID(strconv.Itoa(snapshots.Snapshots[0].ID))
	d.SetString("status", snapshots.Snapshots[0].Status)
	d.SetString("timestamp", time.Time.String(snapshots.Snapshots[0].DateCreated))

	// post check
	return d.Error()
}

func (s *snapshot) Delete(ctx context.Context, d *utils.Data) error {
	return nil
}

func (s *snapshot) Update(ctx context.Context, d *utils.Data) error {
	return nil
}
