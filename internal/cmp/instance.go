// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/models"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/logger"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
)

// instance implements functions related to cmp instances
type instance struct {
	// expose Instance API service to instances related operations
	iClient           *client.InstancesApiService
	serviceInstanceID string
}

func newInstance(iClient *client.InstancesApiService, serviceInstanceID string) *instance {
	return &instance{
		iClient:           iClient,
		serviceInstanceID: serviceInstanceID,
	}
}

// Create instance
func (i *instance) Create(ctx context.Context, d *utils.Data) error {
	logger.Debug("Creating new instance")

	req := &models.CreateInstanceBody{
		ZoneId: d.GetJSONNumber("cloud_id"),
		Instance: &models.CreateInstanceBodyInstance{
			Name: d.GetString("name"),
			InstanceType: &models.CreateInstanceBodyInstanceInstanceType{
				Code: d.GetString("instance_code"),
			},
			Plan: &models.CreateInstanceBodyInstancePlan{
				Id: d.GetJSONNumber("plan_id"),
			},
			Site: &models.CreateInstanceBodyInstanceSite{
				Id: int32(d.GetInt("group_id")),
			},
			Layout: &models.CreateInstanceBodyInstanceLayout{
				Id: d.GetJSONNumber("layout_id"),
			},
		},
		Volumes:           getVolume(d.GetListMap("volumes")),
		NetworkInterfaces: getNetwork(d.GetListMap("networks")),
		Config:            getConfig(d.GetSMap("config")),
		Tags:              getTags(d.GetMap("tags")),
	}

	// Pre check
	if err := d.Error(); err != nil {
		return err
	}

	resp, err := i.iClient.CreateAnInstance(ctx, i.serviceInstanceID, req)
	if err != nil {
		return err
	}
	d.SetID(strconv.Itoa(int(resp.Instance.Id)))

	// post check
	if err := d.Error(); err != nil {
		return err
	}

	return nil
}

// Update instance including poweroff, powerOn, restart, suspend
// changing network, volumes and instance properties such as labels
// groups and tags
func (i *instance) Update(ctx context.Context, d *utils.Data) error {
	logger.Debug("Updating the instance")

	return nil
}

// Delete instance and set ID as ""
func (i *instance) Delete(ctx context.Context, d *utils.Data) error {
	id := d.GetID()
	logger.Debugf("Deleting instance with ID : %d", id)

	// Precheck
	if err := d.Error(); err != nil {
		return err
	}

	res, err := i.iClient.DeleteAnInstance(ctx, i.serviceInstanceID, int32(id))
	if err != nil {
		return err
	}
	if !res.Success {
		return fmt.Errorf("%s", res.Message)
	}
	d.SetID("")

	// post check
	if err := d.Error(); err != nil {
		return err
	}

	return nil
}

// Read instance and set state values accordingly
func (i *instance) Read(ctx context.Context, d *utils.Data) error {
	id := d.GetID()

	logger.Debug("Get instance with ID %d", id)

	// Precheck
	if err := d.Error(); err != nil {
		return err
	}

	resp, err := i.iClient.GetASpecificInstance(ctx, i.serviceInstanceID, int32(id))
	if err != nil {
		return err
	}
	d.SetID(strconv.Itoa(int(resp.Instance.Id)))
	d.SetString("status", resp.Instance.Status)

	// post check
	if err := d.Error(); err != nil {
		return err
	}

	return nil
}

func getVolume(volumes []map[string]interface{}) []models.CreateInstanceBodyVolumes {
	volumesModel := make([]models.CreateInstanceBodyVolumes, 0, len(volumes))
	for i := range volumes {
		vID, _ := utils.ParseInt(volumes[i]["size"].(string))
		volumesModel = append(volumesModel, models.CreateInstanceBodyVolumes{
			Id:          -1,
			Name:        volumes[i]["name"].(string),
			Size:        int32(vID),
			DatastoreId: volumes[i]["datastore_id"],
			RootVolume:  true,
		})
	}

	return volumesModel
}

func getNetwork(networksMap []map[string]interface{}) []models.CreateInstanceBodyNetworkInterfaces {
	networks := make([]models.CreateInstanceBodyNetworkInterfaces, 0, len(networksMap))
	for _, n := range networksMap {
		networks = append(networks, models.CreateInstanceBodyNetworkInterfaces{
			Network: &models.CreateInstanceBodyNetwork{
				Id: int32(n["id"].(int)),
			},
		})
	}

	return networks
}

func getConfig(c map[string]interface{}) *models.CreateInstanceBodyConfig {
	config := &models.CreateInstanceBodyConfig{
		ResourcePoolId: utils.JSONNumber(c["resource_pool_id"]),
		Template:       int32(c["template_id"].(int)),
	}

	return config
}

func getTags(t map[string]interface{}) []models.CreateInstanceBodyTag {
	tags := make([]models.CreateInstanceBodyTag, 0, len(t))
	for k, v := range t {
		tags = append(tags, models.CreateInstanceBodyTag{
			Name:  k,
			Value: v.(string),
		})
	}

	return tags
}