package netappStorageProvApi

import (
  "regexp"
)

const (
  GeturlQtree  string = "/nas/netapp/qtree/?search=%s"
  GetnasurlVol string = "/nas/netapp/volume/?search=%s"
  GetnasurlVolResize string = "/nas/netapp/volume/%s?providername=IT"
)

// all constant as struct can be moved to constant.go;
// only generic methods are allowed

type Cluster struct {
  Name       string `json:"name"`
  Datacenter string `json:"datacenter"`
}

type Node struct {
  Name    string  `json:"name"`
  Cluster Cluster `json:"clustername"`
}

type Vlg struct {
  FriendlyName string `json:"friendlyname"`
  LifName      string `json:"lifname"`
}

type Svm struct {
  FriendlyName string `json:"friendlyname"`
  FabricName   string `json:"fabricname"`
  Lifecycle    string `json:"lifecycle"`
}

type Aggr struct {
  FriendlyName string `json:"friendlyname"`
  Node         Node   `json:"nodename"`
}

type Volume struct {
  Name                string  `json:"name"`
  FriendlyName        string  `json:"friendlyname"`
  Status              string  `json:"status"`
  SyncStatus          string  `json:"syncstatus"`
  ProviderUuid        string  `json:"provideruuid"`
  ProviderName        string  `json:"providername"`
  PhysicalCapacityGb  float32 `json:"physicalcapacitygb"`
  UsableCapacityGb    float32 `json:"usablecapacitygb"`
  AvailableCapacityGb float32 `json:"availablecapacitygb"`
  UsedCapacityGb      float32 `json:"usedcapacitygb"`
  Description         string  `json:"description"`
  IsGeneralPurpose    string  `json:"isgeneralpurpose"`
  MountPoint          string  `json:"mountpoint"`
  DbName              string  `json:"dbname"`
  SnapMirrored        string  `json:"snapmirrored"`
  SnapReservePct      int32   `json:"snapreservepct"`
  Aggr                Aggr    `json:"aggrname"`
  ClusterName         string  `json:"clustername"`
  Svm                 Svm     `json:"svmname"`
  Vlg                 Vlg     `json:"vlgname"`
}

type Qtree struct {
  Name               string  `json:"name"`
  FriendlyName       string  `json:"friendlyname"`
  Status             string  `json:"status"`
  SyncStatus         string  `json:"syncstatus"`
  ProviderUuid       string  `json:"provideruuid"`
  ProviderName       string  `json:"providername"`
  TotalCapacityGb    float32 `json:"totalcapacitygb"`
  AvailableCapaityGb float32 `json:"availablecapacitygb"`
  UsedCapaityGb      float32 `json:"usedcapacitygb"`
  Description        string  `json:"description"`
  Lifecycle          string  `json:"lifecycle"`
  Volume             Volume  `json:"volumename"`
}

type QtreeResp struct {
  Count    int     `json:"count"`
  Next     string  `json:"next"`
  Previous string  `json:"previous"`
  Results  []Qtree `json:"results"`
}

type QtreeAsNas struct {
  Datacenter     string
  Lifecycle      string
  NetworkType    string
  Fabric         string
  Purpose        string
  IsSelfManaged  string
  ProviderUuid   string
  NasStorageType string
  NasStoragePath string
  UsableSizeGib  float32
}

type VolumeResp struct {
  Count    int      `json:"count"`
  Next     string   `json:"next"`
  Previous string   `json:"previous"`
  Results  []Volume `json:"results"`
}

type VolumeAsNas struct {
  Name           string
  Datacenter     string
  Lifecycle      string
  NetworkType    string
  Fabric         string
  Purpose        string
  IsSelfManaged  string
  ProviderUuid   string
  NasStorageType string
  NasStoragePath string
  UsableSizeGib  float32
}

type Response struct {
  Message    string
  status      int
}

type VolSize struct {
  Workflow            string  `json:"workflow"`
  Requestor            string  `json:"requestor"`
  VolumeName          string  `json:"volumename"`
  ProviderUuid        string  `json:"provideruuid"`
  ProviderName        string  `json:"providername"`
  StorageType         string  `json:"storagetype"`
  ProvisioningType    string  `json:"provisioningtype"`
  NodeType            string  `json:"nodetype"`
  UsableSizeGiB       float32 `json:"usablesizegb"`
  OldUsableCapacityGB float32 `json:"oldusablecapacitygb"`
}

// constructor to set the volSize attributes

func NewVolSize(Workflow string, VolumeName string,
  ProviderUuid string, ProviderName string, StorageType string,
  ProvisioningType string, NodeType string, UsableSizeGiB float32,
  OldUsableCapacityGB float32, Requestor string) *VolSize {
  nas_attr := VolSize{}
  nas_attr.Workflow = Workflow
  nas_attr.Workflow = Requestor
  nas_attr.VolumeName = VolumeName
  nas_attr.ProviderUuid = ProviderUuid
  nas_attr.ProviderName = ProviderName
  nas_attr.StorageType = StorageType
  nas_attr.ProvisioningType = ProvisioningType
  nas_attr.NodeType = NodeType
  nas_attr.UsableSizeGiB = UsableSizeGiB
  nas_attr.OldUsableCapacityGB = OldUsableCapacityGB
  return &nas_attr
}

// generic methods only
func determineNetworkType(vlgFriendlyName string) string {
  reFindNetType := regexp.MustCompile(`(m|u)(i|p|sp|z|sz)(.*?)(\d{3})`)
  switch reFindNetType.FindStringSubmatch(vlgFriendlyName)[2] {
  case "i":
    return "Internal"
  case "p":
    return "Protected"
  case "sp":
    return "Simulated Protected"
  case "z":
    return "DMZ"
  case "sz":
    return "Simulated DMZ"
  default:
    return "Unknown"
  }
}

func determinePurpose(vlgFriendlyName string) string {
  reFindPurpose := regexp.MustCompile(`(m|u)(i|p|sp|z|sz)(.*?)(\d{3})`)
  return reFindPurpose.FindStringSubmatch(vlgFriendlyName)[3]
}

func determineSelfManaged(vlgFriendlyName string) string {
  reFindSelfManaged := regexp.MustCompile(`(m|u)(i|p|sp|z|sz)(.*?)(\d{3})`)
  switch reFindSelfManaged.FindStringSubmatch(vlgFriendlyName)[1] {
  case "m":
    return "No"
  case "u":
    return "Yes"
  default:
    return "Unknown"
  }
}

