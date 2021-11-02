package model

import (
	"context"
	"encoding/json"
	"errors"

	elastic "github.com/olivere/elastic/v7"
)

type Device struct {
	Name       string                   `json:"name"`
	ID         string                   `json:"id"`
	SensorData []map[string]interface{} `json:"sensor_data"`
}

func (d *Device) IndexName() string {
	return "devices"
}

func InsertNewDevice(dataJson string) error {
	var d Device
	_, err := DB.Index().Index(d.IndexName()).BodyJson(dataJson).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func GetAllDevices(devices *[]Device) error {
	var d Device

	searchSource := elastic.NewSearchSource()
	searchSource.Query(elastic.NewMatchAllQuery())

	searchResult, err := DB.Search().Index(d.IndexName()).SearchSource(searchSource).Do(context.Background())
	if err != nil {
		return err
	}

	for _, hit := range searchResult.Hits.Hits {
		var device Device

		err = json.Unmarshal(hit.Source, &device)
		if err != nil {
			return err
		}

		*devices = append(*devices, device)
	}

	return nil
}

func DeleteDeviceByID(deviceID string) (error, string) {
	var d Device
	searchSource := elastic.NewSearchSource()
	searchSource.Query(elastic.NewMatchQuery("id", deviceID))
	searchResult, err := DB.Search().Index(d.IndexName()).SearchSource(searchSource).Do(context.Background())
	if err != nil {
		return err, "Error occured while deleting the device"
	}
	if len(searchResult.Hits.Hits) == 0 {
		return errors.New("DEVICE_NOT_FOUND"), "Device not found"
	}
	for _, hit := range searchResult.Hits.Hits {
		_, err = DB.Delete().Index(d.IndexName()).Id(hit.Id).Do(context.Background())
		if err != nil {
			return err, "Error occured while deleting the device"
		}
	}
	return err, "Device deleted successfully!"
}

func UpdateDeviceByID(deviceID string, updateData map[string]interface{}) (error, string) {
	var d Device
	searchSource := elastic.NewSearchSource()
	searchSource.Query(elastic.NewMatchQuery("id", deviceID))
	searchResult, err := DB.Search().Index(d.IndexName()).SearchSource(searchSource).Do(context.Background())
	if err != nil {
		return err, "Error occured while updating the device"
	}
	if len(searchResult.Hits.Hits) == 0 {
		return errors.New("DEVICE_NOT_FOUND"), "Device not found"
	}
	for _, hit := range searchResult.Hits.Hits {
		_, err = DB.Update().Index(d.IndexName()).Id(hit.Id).Doc(updateData).Do(context.Background())
		if err != nil {
			return err, "Error occured while updating the device"
		}
	}
	return err, "Device updated successfully!"
}
