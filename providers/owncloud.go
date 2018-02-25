package providers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/bitly/oauth2_proxy/api"
)

type OwncloudProvider struct {
	*ProviderData
}

func NewOwncloudProvider(p *ProviderData) *OwncloudProvider {
	p.ProviderName = "Owncloud"
	if p.LoginURL == nil || p.LoginURL.String() == "" {
		p.LoginURL = &url.URL{
			Scheme: "http",
			Host:   "localhost",
			Path:   "/index.php/apps/oauth2/authorize",
		}
	}
	if p.RedeemURL == nil || p.RedeemURL.String() == "" {
		p.RedeemURL = &url.URL{
			Scheme: "http",
			Host:   "localhost",
			Path:   "/index.php/apps/oauth2/api/v1/token",
		}
	}
	if p.ValidateURL == nil || p.ValidateURL.String() == "" {
		p.RedeemURL = &url.URL{
			Scheme: "http",
			Host:   "localhost",
			Path:   "/ocs/v1.php/cloud/user",
		}
	}
	return &OwncloudProvider{ProviderData: p}
}

func (p *OwncloudProvider) Redeem(redirectURL, code string) (s *SessionState, err error) {
	if code == "" {
		err = errors.New("missing code")
		return
	}

	params := url.Values{}
	params.Add("redirect_uri", redirectURL)
	params.Add("client_id", p.ClientID)
	params.Add("client_secret", p.ClientSecret)
	params.Add("code", code)
	params.Add("grant_type", "authorization_code")
	if p.ProtectedResource != nil && p.ProtectedResource.String() != "" {
		params.Add("resource", p.ProtectedResource.String())
	}

	var req *http.Request
	req, err = http.NewRequest("POST", p.RedeemURL.String(), bytes.NewBufferString(params.Encode()))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", p.ClientID, p.ClientSecret)))))

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("got %d from %q %s", resp.StatusCode, p.RedeemURL.String(), body)
		return
	}

	// blindly try json and x-www-form-urlencoded
	var jsonResponse struct {
		AccessToken string `json:"access_token"`
	}
	err = json.Unmarshal(body, &jsonResponse)
	if err == nil {
		s = &SessionState{
			AccessToken: jsonResponse.AccessToken,
		}
		return
	}

	var v url.Values
	v, err = url.ParseQuery(string(body))
	if err != nil {
		return
	}
	if a := v.Get("access_token"); a != "" {
		s = &SessionState{AccessToken: a}
	} else {
		err = fmt.Errorf("no access token found %s", body)
	}
	return
}

func getOwncloudHeader(access_token string) http.Header {
	header := make(http.Header)
	header.Set("Authorization", fmt.Sprintf("Bearer %s", access_token))
	return header
}

func (p *OwncloudProvider) GetEmailAddress(s *SessionState) (string, error) {
	if s.AccessToken == "" {
		return "", errors.New("missing access token")
	}
	req, err := http.NewRequest("GET", p.ValidateURL.String()+"?format=json", nil)
	if err != nil {
		return "", err
	}
	req.Header = getOwncloudHeader(s.AccessToken)

	type result struct {
		Ocs struct {
			Meta struct {
				Status     string `json:"status"`
				StatusCode int    `json:"statuscode"`
				Message    string `json:"message"`
			} `json:"meta"`
			Data struct {
				Id          string `json:"id"`
				DisplayName string `json:"display-name"`
				Email       string `json:"email"`
			} `json:"data"`
		} `json:"ocs"`
	}
	var r result
	err = api.RequestJson(req, &r)
	if err != nil {
		return "", err
	}
	if r.Ocs.Data.Id == "" {
		return "", errors.New("no id")
	}
	return r.Ocs.Data.Id + "@" + p.LoginURL.Host, nil
}
