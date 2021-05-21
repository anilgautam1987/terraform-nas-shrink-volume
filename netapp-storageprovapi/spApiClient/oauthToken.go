package spApiClient

import (
  "errors"
  "fmt"
  "strings"
  "time"
  "io/ioutil"
  "net/http"
  "net/url"
  "encoding/json"
)

const (
  OauthRefreshTriggerSeconds int64 = 600
)

type Oauth2Return struct {
  AccessToken     string `json:"access_token"`
  RefreshToken    string `json:"refresh_token"`
  TokenType       string `json:"token_type"`
  ExpiresIn       int64  `json:"expires_in"`
}

type Oauth2Token struct {
  HTTPClient      *http.Client
  Url             string
  ClientId        string
  ClientSecret    string
  UserName        string
  Password        string
  AccessToken     string
  RefreshToken    string
  ExpirationDate  time.Time
}

func NewToken(url string, clientId string, clientSecret string, userName string, password string) *Oauth2Token {
  return &Oauth2Token{
    HTTPClient: &http.Client{},
    Url: url,
    ClientId: clientId,
    ClientSecret: clientSecret,
    UserName: userName,
    Password: password,
    AccessToken: "",
    RefreshToken: "",
    ExpirationDate: time.Time{},
  }
}

func (oat *Oauth2Token) GetToken() (error) {

  var payload url.Values = nil

  if oat.AccessToken != "" {
    expiresIn := oat.ExpirationDate.Sub(time.Now())
    if int64(expiresIn / time.Second) < OauthRefreshTriggerSeconds {
      payload = url.Values{
        "grant_type": {"refresh_token"},
        "refresh_token": {oat.RefreshToken},
      }
    } else {
      return nil
    }
  } else {
    payload = url.Values{
      "grant_type": {"password"},
      "username": {oat.UserName},
      "password": {oat.Password},
    }
  }

  oaReq, err := http.NewRequest("POST", oat.Url, strings.NewReader(payload.Encode()))
  oaReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  oaReq.SetBasicAuth(oat.ClientId, oat.ClientSecret)
  oaResp, err := oat.HTTPClient.Do(oaReq)

  if err != nil {
    fmt.Println("ERROR: Executing OAuth Token HTTP Request")
    return err
  }

  defer oaResp.Body.Close()
  oaRespData, err := ioutil.ReadAll(oaResp.Body)
  if err != nil {
    fmt.Printf("ERROR: Reading data from Oauth response\n")
    return err
  }

  if oaResp.StatusCode < http.StatusOK || oaResp.StatusCode >= http.StatusBadRequest {
    fmt.Printf("ERROR: Oauth Response Code: %v\n", string(oaResp.StatusCode))
    return errors.New(string(oaRespData))
  }

  oaRet := Oauth2Return{}

  if err := json.Unmarshal(oaRespData, &oaRet); err != nil {
    fmt.Printf("ERROR: Unmarshalling Oauth JSON Response\n")
    return err
  }

  oat.AccessToken = oaRet.AccessToken
  oat.RefreshToken = oaRet.RefreshToken
  oat.ExpirationDate = time.Now().Add(time.Second * time.Duration(oaRet.ExpiresIn))

  return nil
}
