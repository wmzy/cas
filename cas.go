package cas

import (
	"errors"
	"time"

	"github.com/parnurzeal/gorequest"
)

// CAS
type CAS struct {
	BaseUrl string
	client  *gorequest.SuperAgent
	version int
}

// NewClient creates a CAS with the provided baseURL.
func NewV1(baseURL string) *CAS {
	return &CAS{
		BaseUrl: baseURL,
		client:  gorequest.New().Timeout(10 * time.Second).Type("text"),
		version: 1,
	}
}

// NewClient creates a CAS with the provided baseURL.
func NewV2(baseURL string) *CAS {
	return &CAS{
		BaseUrl: baseURL,
		client:  gorequest.New().Timeout(10 * time.Second).Type("xml"),
		version: 2,
	}
}

func (self *CAS) GetLoginURL(callbackURL string) string {
	return self.BaseUrl + "/login?service=" + callbackURL
}

func (self *CAS) GetLogoutURL(callbackURL string) string {
	return self.BaseUrl + "/logout?service=" + callbackURL
}

func (self *CAS) getValidateURL() string {
	if self.version == 1 {
		return self.BaseUrl + "/validate"
	}

	return self.BaseUrl + "/serviceValidate"
}

func (self *CAS) getUserV1(ticket, callbackURL string) (*Authentication, error) {
	_, body, errs := self.client.Get(self.BaseUrl + "/validate").Query("ticket=" + ticket).Query("service=" + callbackURL).End()
	if errs != nil {
		return nil, errs[0]
	}
	if body == "no\n\n" {
		return nil, errors.New("not login")
	}

	return &Authentication{User: body[4 : len(body)-1]}, nil
}

func (self *CAS) GetUser(ticket, callbackURL string) (*Authentication, error) {
	if self.version == 1 {
		return self.getUserV1(ticket, callbackURL)
	}

	_, body, errs := self.client.Get(self.BaseUrl + "/serviceValidate").Query("ticket=" + ticket).Query("service=" + callbackURL).EndBytes()
	if errs != nil {
		return nil, errs[0]
	}
	return ParseServiceResponse(body)
}
