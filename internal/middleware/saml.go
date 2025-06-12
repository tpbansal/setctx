package middleware

import (
	"encoding/base64"
	"encoding/xml"
	"net/http"
	"strings"

	"github.com/crewjam/saml/samlsp"
)

// SAMLAuthenticator implements SAML authentication
type SAMLAuthenticator struct {
	sp *samlsp.Middleware
}

// NewSAMLAuthenticator creates a new SAML authenticator
func NewSAMLAuthenticator(sp *samlsp.Middleware) *SAMLAuthenticator {
	return &SAMLAuthenticator{
		sp: sp,
	}
}

// Authenticate implements the Authenticator interface for SAML
func (a *SAMLAuthenticator) Authenticate(r *http.Request) (*User, error) {
	// Get the SAML assertion from the request
	session, err := a.sp.Session.GetSession(r)
	if err != nil {
		return nil, err
	}

	// Extract user information from the SAML assertion
	user := &User{
		ID:         session.(samlsp.SessionWithAttributes).GetAttributes().Get("uid"),
		Email:      session.(samlsp.SessionWithAttributes).GetAttributes().Get("email"),
		Name:       session.(samlsp.SessionWithAttributes).GetAttributes().Get("displayName"),
		Roles:      strings.Split(session.(samlsp.SessionWithAttributes).GetAttributes().Get("groups"), ","),
		AuthMethod: "saml",
		Claims:     make(map[string]interface{}),
	}

	// Add all attributes to the user's claims map
	for _, attr := range session.(samlsp.SessionWithAttributes).GetAttributes().All() {
		user.Claims[attr.Name] = attr.Values[0]
	}

	return user, nil
}

// HandleSAMLResponse handles the SAML response from the IdP
func (a *SAMLAuthenticator) HandleSAMLResponse(w http.ResponseWriter, r *http.Request) {
	// Parse the SAML response
	var response struct {
		XMLName      xml.Name `xml:"Response"`
		Assertion    string   `xml:"Assertion"`
		RelayState   string   `xml:"RelayState"`
		Destination  string   `xml:"Destination,attr"`
		InResponseTo string   `xml:"InResponseTo,attr"`
	}

	if err := xml.NewDecoder(r.Body).Decode(&response); err != nil {
		http.Error(w, "Invalid SAML response", http.StatusBadRequest)
		return
	}

	// Decode the assertion
	assertion, err := base64.StdEncoding.DecodeString(response.Assertion)
	if err != nil {
		http.Error(w, "Invalid SAML assertion", http.StatusBadRequest)
		return
	}

	// TODO: Validate the assertion and create a session
	// This would typically involve:
	// 1. Verifying the signature
	// 2. Checking the conditions (not before, not on or after)
	// 3. Creating a session with the user's attributes
	// 4. Redirecting to the original URL

	http.Redirect(w, r, response.RelayState, http.StatusFound)
}
