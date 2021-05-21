package spApiClient

import (
  "bytes"
  "crypto/tls"
  "errors"
  "fmt"
  "io"
  "io/ioutil"
  "net/http"
)

const (
  BaseURL string = "https://"
)

type SpaConfig struct {
  SpaHostName     string
  OauthUrl        string
  ClientId        string
  ClientSecret    string
  UserName        string
  Password        string
}

type Client struct {
  BaseURL       string
  OauthToken    *Oauth2Token
  HTTPClient    *http.Client
}

func NewClient(spac SpaConfig, verify bool) *Client {
  tr := &http.Transport{
    TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
  }
  if verify == false {
    tr = &http.Transport{
      TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
  }

  oat := NewToken(spac.OauthUrl,
                  spac.ClientId,
                  spac.ClientSecret,
                  spac.UserName,
                  spac.Password,
  )

  return &Client{
    BaseURL: BaseURL + spac.SpaHostName + "%s",
    OauthToken: oat,
    HTTPClient: &http.Client{Transport: tr},
  }
}

func (c *Client) Get(url string) ([]byte, error) {
  var method string = "GET"
  var payload []byte = nil
  requestData, err := c.makeRequest(method, url, payload)
  return requestData, err
}

func (c *Client) Post(url string, payload []byte) ([]byte, error) {
  var method string = "POST"
  requestData, err := c.makeRequest(method, url, payload)
  return requestData, err
}

func (c *Client) Put(url string, payload []byte) ([]byte, error) {
  var method string = "PUT"
  requestData, err := c.makeRequest(method, url, payload)
  return requestData, err
}

func (c *Client) makeRequest(method string, url string, payload []byte) ([]byte, error) {
  fmt.Println("Getting Token...")
  err := c.OauthToken.GetToken()
  if err != nil {
    fmt.Println("ERROR: Getting Oauth Token For SPA Request")
    return nil, err
  }

  fmt.Printf("\n******************************************* NAS REST URLs **************************************\n")
  fmt.Printf("Getting SPA Method %v \n", method)
  fmt.Printf("Getting SPA REST API URL: %v \n", fmt.Sprintf(c.BaseURL, url) )
  fmt.Printf("\n****************************************************** END BLOCK *********************************************\n")

  var reader io.Reader
  if payload != nil{
    reader = bytes.NewReader(payload)
  }
  req, err := http.NewRequest(method, fmt.Sprintf(c.BaseURL, url), reader)
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.OauthToken.AccessToken))

  resp, err := c.HTTPClient.Do(req)

  if err != nil {
    fmt.Printf("ERROR: Executing SPA %v HTTP Request\n", method)
    return nil, err
  }

  defer resp.Body.Close()
  respData, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    fmt.Printf("ERROR: Reading data from SPA %v response\n", method)
    return nil, err
  }

  if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
    fmt.Printf("ERROR: SPA %v Response Code: %+v\n", method, string(resp.StatusCode))
    return nil, errors.New(string(respData))
  }
  return respData, nil
}


