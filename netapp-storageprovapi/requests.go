package netappStorageProvApi 

import (
  "errors"
  "fmt"
  "time"
  "encoding/json"
  "gitscm.cisco.com/sa/tf-provider-netapp-storageprovapi/netapp-storageprovapi/spApiClient"
)

const (
  TimeOutSeconds  int64  =  3600
  RequestUrl      string =  "/nas/execution/%s"
)

type St2Output struct {
  Status            string        `json:"status"`
  StartTs           string        `json:"start_timestamp"`
  EndTs             string        `json:"end_timestamp"`
  Log               interface{}   `json:"log"`
  Parameters        interface{}   `json:"parameters"`
  Runner            interface{}   `json:"runner"`
  Children          interface{}   `json:"children"`
  ElapsedSeconds    float64       `json:"elapsed_seconds"`
  WebUrl            interface{}   `json:"web_url"`
  Result            interface{}   `json:"result"`
  Context           interface{}   `json:"context"`
  Action            interface{}   `json:"action"`
  LiveAction        interface{}   `json:"liveaction"`
  Id                string        `json:"id"`
}

type ExecutionResp struct {
  StatusCode  int         `json:"status_code"`
  Message     string      `json:"message"`
  Data        St2Output   `json:"data"`
}

func PollRequest(c *spApiClient.Client, reqId string) (string, error) {
  var endTime int64 = time.Now().Unix() + TimeOutSeconds

  execResp := ExecutionResp{}

  forLoop:for time.Now().Unix() <= endTime {
    reqData, err := c.Get(fmt.Sprintf(RequestUrl, reqId))

    if err != nil {
      return "", err
    }

    if err := json.Unmarshal(reqData, &execResp); err != nil {
      fmt.Println("ERROR: Converting Request Data To JSON")
      fmt.Printf("    Raw Data: %v\n", string(reqData))
      return "", err
    }

    switch execResp.Data.Status {
    case "succeeded", "failed":
      break forLoop
    }

    fmt.Printf("    Status=%v, Sleeping\n", execResp.Data.Status)
    time.Sleep(30 * time.Second)

  }

  switch execResp.Data.Status {
  case "succeeded", "failed":
    return execResp.Data.Status, nil
  default:
    var errMsg string = fmt.Sprintf("ERROR:  The request did not complete within %v seconds.  Last Status = %v", TimeOutSeconds, execResp.Data.Status)
    return "", errors.New(errMsg)
  }

}

