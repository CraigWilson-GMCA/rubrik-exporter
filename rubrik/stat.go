//
// rubrik-exporter
//
// Exports metrics from rubrik backup for prometheus
//
// License: Apache License Version 2.0,
// Organization: Claranet GmbH
// Author: Martin Weber <martin.weber@de.clara.net>
//

package rubrik

import (
	// "fmt"
	"encoding/json"
	"io/ioutil"
	"net/url"
)

type VmStorageList struct {
	*ResultList
	Data []VmStorage `json:"data"`
}

type VmStorage struct {
	ID                     string
	Logicalbytes           float64 `json:"logicalBytes"`
	IngestedBytes          float64 `json:"ingestedBytes"`
	ExclusivePhysicalBytes float64 `json:"exclusivePhysicalBytes"`
	SharedPhysicalBytes    float64 `json:"sharedPhysicalBytes"`
	IndexStorageBytes      float64 `json:"indexStorageBytes"`
}

type SystemStorage struct {
	Total         int `json:"total"`
	Used          int `json:"used"`
	Available     int `json:"available"`
	Snapshot      int `json:"snapshot"`
	LiveMount     int `json:"liveMount"`
	Miscellaneous int `json:"miscellaneous"`
}

type DataLocationUsageList struct {
	*ResultList
	Data []DataLocationUsage `json:"data"`
}

type DataLocationUsage struct {
	LocationID                 string `json:"locationId"`
	DataDownloaded             int    `json:"dataDownloaded"`
	DataArchived               int    `json:"dataArchived"`
	NumVMsArchived             int    `json:"numVMsArchived"`
	NumFilesetsArchived        int    `json:"numFilesetsArchived"`
	NumLinuxFilesetsArchived   int    `json:"numLinuxFilesetsArchived"`
	NumWindowsFilesetsArchived int    `json:"numWindowsFilesetsArchived"`
	NumShareFilesetsArchived   int    `json:"numShareFilesetsArchived"`
	NumMssqlDbsArchived        int    `json:"numMssqlDbsArchived"`
	NumHypervVmsArchived       int    `json:"numHypervVmsArchived"`
	NumNutanixVmsArchived      int    `json:"numNutanixVmsArchived"`
	NumManagedVolumesArchived  int    `json:"numManagedVolumesArchived"`
}

// GetSystemStorage ...
func (r Rubrik) GetSystemStorage() SystemStorage {
	resp, err := r.makeRequest("GET", "/api/internal/stats/system_storage", RequestParams{})
	if err != nil || resp == nil {
		return SystemStorage{}
	}
	defer resp.Body.Close()

	data := json.NewDecoder(resp.Body)
	var d SystemStorage
	data.Decode(&d)

	return d
}

// GetPerVMStorage ...
func (r Rubrik) GetPerVMStorage() []VmStorage {
	resp, err := r.makeRequest("GET", "/api/internal/stats/per_vm_storage", RequestParams{})
	if err != nil || resp == nil {
		return []VmStorage{}
	}
	defer resp.Body.Close()

	data := json.NewDecoder(resp.Body)
	var d VmStorageList
	data.Decode(&d)

	return d.Data
}

// GetStreamCount ...
func (r Rubrik) GetStreamCount() int {
	resp, err := r.makeRequest("GET", "/api/internal/stats/streams/count", RequestParams{})
	if err != nil || resp == nil {
		return 0
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var data map[string]int
	json.Unmarshal(body, &data)

	return data["count"]
}

// GetDataLocationUsage ...
func (r Rubrik) GetDataLocationUsage() []DataLocationUsage {
	resp, err := r.makeRequest("GET", "/api/internal/stats/data_location/usage", RequestParams{})
	if err != nil || resp == nil {
		return []DataLocationUsage{}
	}
	defer resp.Body.Close()

	var data DataLocationUsageList
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&data)

	return data.Data
}

func (r Rubrik) GetPhysicalIngest() []TimeStat {
	resp, err := r.makeRequest("GET", "/api/internal/stats/physical_ingest/time_series", RequestParams{params: url.Values{"range": []string{"-10min"}}})
	if err != nil || resp == nil {
		return []TimeStat{}
	}
	defer resp.Body.Close()

	var data []TimeStat
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&data)

	return data
}

func (r Rubrik) GetArchivalBandwith(locationID string, timerange string) []TimeStat {
	if timerange == "" {
		timerange = "-1h"
	}

	resp, err := r.makeRequest("GET", "/api/internal/stats/archival/bandwidth/time_series",
		RequestParams{params: url.Values{"data_location_id": []string{locationID}, "range": []string{timerange}}})
	if err != nil || resp == nil {
		return []TimeStat{}
	}
	defer resp.Body.Close()

	var data []TimeStat
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&data)

	return data
}

// GetRunawayRemaining - Get the number of days remaining before the system fills up.
func (r Rubrik) GetRunawayRemaining() int {
	resp, err := r.makeRequest("GET", "/api/internal/stats/runway_remaining", RequestParams{})
	if err != nil || resp == nil {
		return 0
	}
	defer resp.Body.Close()

	var data map[string]int
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&data)

	return data["days"]
}

// GetAverageStorageGrowthPerDay - Get average storage growth per day.
func (r Rubrik) GetAverageStorageGrowthPerDay() int {
	resp, err := r.makeRequest("GET", "/api/internal/stats/average_storage_growth_per_day", RequestParams{})
	if err != nil || resp == nil {
		return 0
	}
	defer resp.Body.Close()

	var data map[string]int
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&data)

	return data["bytes"]
}
