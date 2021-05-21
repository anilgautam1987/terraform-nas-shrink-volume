package netappStorageProvApi

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitscm.cisco.com/sa/tf-provider-netapp-storageprovapi/netapp-storageprovapi/spApiClient"
)

func VolumeGet(c *spApiClient.Client, volName string) (*VolumeAsNas, error) {
	fmt.Println("----------------------------- Manage Storage Volume ----------------------------")

	volResp := VolumeResp{}
	reqData, err := c.Get(fmt.Sprintf(GetnasurlVol, volName))
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(reqData, &volResp); err != nil {
		fmt.Println("ERROR: Converting Request Data To JSON")
		fmt.Printf("    Raw Data: %v\n", string(reqData))
		return nil, err
	}

	if len(volResp.Results) != 1 {
		err := errors.New(fmt.Sprintf("Expected 1 Qtree, got %v", len(volResp.Results)))
		return nil, err
	}
	myVolume := &VolumeAsNas{}
	myVolume.Name = volResp.Results[0].Name
	myVolume.Datacenter = volResp.Results[0].Aggr.Node.Cluster.Datacenter
	myVolume.Lifecycle = volResp.Results[0].Svm.Lifecycle
	myVolume.NetworkType = determineNetworkType(volResp.Results[0].Vlg.FriendlyName)
	myVolume.Fabric = volResp.Results[0].Svm.FabricName
	myVolume.Purpose = determinePurpose(volResp.Results[0].Vlg.FriendlyName)
	myVolume.IsSelfManaged = determineSelfManaged(volResp.Results[0].Vlg.FriendlyName)
	myVolume.ProviderUuid = volResp.Results[0].ProviderUuid
	myVolume.NasStorageType = "Volume"
	myVolume.NasStoragePath = fmt.Sprintf("%v:/%v/%v/%v",
		volResp.Results[0].Vlg.LifName,
		volResp.Results[0].Vlg.FriendlyName,
		volResp.Results[0].FriendlyName,
		volResp.Results[0].FriendlyName)
	myVolume.UsableSizeGib = volResp.Results[0].UsableCapacityGb

	return myVolume, nil
}

func VolumePut(c *spApiClient.Client, input VolSize) string {
	fmt.Println("----------------------------- Put Volume ----------------------------")
	fmt.Printf("%+v\n", input)

	//   GetNasUrl_Vol_Resize string = "%s/nas/netapp/volume/%s"
	url := fmt.Sprintf(GetnasurlVolResize, input.VolumeName)
	jsonReq, err := json.Marshal(input)
	if err != nil {
		fmt.Printf("%v \n", err)
	}
	reqData, _ := c.Put(url, jsonReq)
	volResp := VolumeResp{}

	if err := json.Unmarshal(reqData, &volResp); err != nil {
		fmt.Println("ERROR: Converting Request Data To JSON")
		fmt.Printf("    Raw Data: %v\n", string(reqData))
	}

	if len(volResp.Results) != 1 {
		err := errors.New(fmt.Sprintf("Expected 1 Qtree, got %v \n", len(volResp.Results)))
		fmt.Printf("volResp result %v", err)
	}

	fmt.Printf("last blk volResp result %v \n", volResp.Results)
	return string(reqData)
}



