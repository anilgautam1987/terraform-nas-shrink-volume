package netappStorageProvApi

import (
  "fmt"
  "errors"
  "encoding/json"
  "gitscm.cisco.com/sa/tf-provider-netapp-storageprovapi/netapp-storageprovapi/spApiClient"
)

func QtreeGet(c *spApiClient.Client, qtreeName string) (*QtreeAsNas, error) {

  qtreeResp := QtreeResp{}

  reqData, err := c.Get(fmt.Sprintf(GeturlQtree, qtreeName))
  if err != nil {
    return nil, err
  }

  if err := json.Unmarshal(reqData, &qtreeResp); err != nil {
    fmt.Println("ERROR: Converting Request Data To JSON")
    fmt.Printf("    Raw Data: %v\n", string(reqData))
    return nil, err
  }

  if len(qtreeResp.Results) != 1 {
    err := errors.New(fmt.Sprintf("Expected 1 Qtree, got %v", len(qtreeResp.Results)))
    return nil, err
  }

  myQtree := &QtreeAsNas{}
  myQtree.Datacenter     = qtreeResp.Results[0].Volume.Aggr.Node.Cluster.Datacenter
  myQtree.Lifecycle      = qtreeResp.Results[0].Lifecycle
  myQtree.NetworkType    = determineNetworkType(qtreeResp.Results[0].Volume.Vlg.FriendlyName)
  myQtree.Fabric         = qtreeResp.Results[0].Volume.Svm.FabricName
  myQtree.Purpose        = determinePurpose(qtreeResp.Results[0].Volume.Vlg.FriendlyName)
  myQtree.IsSelfManaged  = determineSelfManaged(qtreeResp.Results[0].Volume.Vlg.FriendlyName)
  myQtree.ProviderUuid   = qtreeResp.Results[0].ProviderUuid
  myQtree.NasStorageType = "Qtree"
  myQtree.NasStoragePath = fmt.Sprintf("%v:/%v/%v/%v",
                                       qtreeResp.Results[0].Volume.Vlg.LifName,
                                       qtreeResp.Results[0].Volume.Vlg.FriendlyName,
                                       qtreeResp.Results[0].Volume.FriendlyName,
                                       qtreeResp.Results[0].FriendlyName)
  myQtree.UsableSizeGib  = qtreeResp.Results[0].TotalCapacityGb

  return myQtree, nil
}

