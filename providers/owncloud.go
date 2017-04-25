package providers

import (
	"net/url"
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
	return &OwncloudProvider{ProviderData: p}
}

func (p *OwncloudProvider) GetEmailAddress(s *SessionState) (string, error) {
	return s.User + "@" + p.LoginURL.Host, nil
}
