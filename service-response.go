package cas

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

// AuthenticationError Code values
const (
	INVALID_REQUEST            = "INVALID_REQUEST"
	INVALID_TICKET_SPEC        = "INVALID_TICKET_SPEC"
	UNAUTHORIZED_SERVICE       = "UNAUTHORIZED_SERVICE"
	UNAUTHORIZED_SERVICE_PROXY = "UNAUTHORIZED_SERVICE_PROXY"
	INVALID_PROXY_CALLBACK     = "INVALID_PROXY_CALLBACK"
	INVALID_TICKET             = "INVALID_TICKET"
	INVALID_SERVICE            = "INVALID_SERVICE"
	INTERNAL_ERROR             = "INTERNAL_ERROR"
)

// AuthenticationError represents a CAS AuthenticationFailure response
type AuthenticationError struct {
	Code    string
	Message string
}

// AuthenticationError provides a differentiator for casting.
func (e AuthenticationError) AuthenticationError() bool {
	return true
}

// Error returns the AuthenticationError as a string
func (e AuthenticationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Authentication captures authenticated user information
type Authentication struct {
	User                string         // Users login name
	ProxyGrantingTicket string         // Proxy Granting Ticket
	Proxies             []string       // List of proxies
	AuthenticationDate  time.Time      // Time at which authentication was performed
	IsNewLogin          bool           // Whether new authentication was used to grant the service ticket
	IsRememberedLogin   bool           // Whether a long term token was used to grant the service ticket
	MemberOf            []string       // List of groups which the user is a member of
	Attributes          UserAttributes // Additional information about the user
}

// UserAttributes represents additional data about the user
type UserAttributes map[string][]string

// Get retrieves an attribute by name.
//
// Attributes are stored in arrays. Get will only return the first element.
func (a UserAttributes) Get(name string) string {
	if v, ok := a[name]; ok {
		return v[0]
	}

	return ""
}

// Add appends a new attribute.
func (a UserAttributes) Add(name, value string) {
	a[name] = append(a[name], value)
}

// ParseServiceResponse returns a successful response or an error
func ParseServiceResponse(data []byte) (*Authentication, error) {
	var x xmlServiceResponse

	if err := xml.Unmarshal(data, &x); err != nil {
		return nil, err
	}

	if x.Failure != nil {
		msg := strings.TrimSpace(x.Failure.Message)
		err := &AuthenticationError{Code: x.Failure.Code, Message: msg}
		return nil, err
	}

	r := &Authentication{
		User:                x.Success.User,
		ProxyGrantingTicket: x.Success.ProxyGrantingTicket,
		Attributes:          make(UserAttributes),
	}

	if p := x.Success.Proxies; p != nil {
		r.Proxies = p.Proxies
	}

	if a := x.Success.Attributes; a != nil {
		r.AuthenticationDate = a.AuthenticationDate
		r.IsRememberedLogin = a.LongTermAuthenticationRequestTokenUsed
		r.IsNewLogin = a.IsFromNewLogin
		r.MemberOf = a.MemberOf

		if a.UserAttributes != nil {
			for _, ua := range a.UserAttributes.Attributes {
				if ua.Name == "" {
					continue
				}

				r.Attributes.Add(ua.Name, strings.TrimSpace(ua.Value))
			}

			for _, ea := range a.UserAttributes.AnyAttributes {
				r.Attributes.Add(ea.XMLName.Local, strings.TrimSpace(ea.Value))
			}
		}

		if a.ExtraAttributes != nil {
			for _, ea := range a.ExtraAttributes {
				r.Attributes.Add(ea.XMLName.Local, strings.TrimSpace(ea.Value))
			}
		}
	}

	return r, nil
}
