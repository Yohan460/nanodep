package godep

import (
	"context"
	"net/http"
)

// Profile corresponds to the Apple DEP API "Profile" structure.
// See https://developer.apple.com/documentation/devicemanagement/profile
type Profile struct {
	ProfileName           string   `json:"profile_name"`
	URL                   string   `json:"url"`
	AllowPairing          bool     `json:"allow_pairing,omitempty"`
	IsSupervised          bool     `json:"is_supervised,omitempty"`
	IsMultiUser           bool     `json:"is_multi_user,omitempty"`
	IsMandatory           bool     `json:"is_mandatory,omitempty"`
	AwaitDeviceConfigured bool     `json:"await_device_configured,omitempty"`
	IsMDMRemovable        bool     `json:"is_mdm_removable"` // default true
	SupportPhoneNumber    string   `json:"support_phone_number,omitempty"`
	AutoAdvanceSetup      bool     `json:"auto_advance_setup,omitempty"`
	SupportEmailAddress   string   `json:"support_email_address,omitempty"`
	OrgMagic              string   `json:"org_magic"`
	AnchorCerts           []string `json:"anchor_certs,omitempty"`
	SupervisingHostCerts  []string `json:"supervising_host_certs,omitempty"`
	Department            string   `json:"department,omitempty"`
	Devices               []string `json:"devices,omitempty"`
	Language              string   `json:"language,omitempty"`
	Region                string   `json:"region,omitempty"`
	ConfigurationWebURL   string   `json:"configuration_web_url,omitempty"`

	// See https://developer.apple.com/documentation/devicemanagement/skipkeys
	SkipSetupItems []string `json:"skip_setup_items,omitempty"`

	// additional undocumented key only returned when requesting a profile from Apple.
	ProfileUUID string `json:"profile_uuid,omitempty"`
}

// ProfileResponse corresponds to the Apple DEP API "AssignProfileResponse" structure.
// See https://developer.apple.com/documentation/devicemanagement/assignprofileresponse
type ProfileResponse struct {
	ProfileUUID string            `json:"profile_uuid"`
	Devices     map[string]string `json:"devices"`
}

// AssignProfiles uses the Apple "Assign a profile to a list of devices" API
// endpoint to assign a DEP profile UUID to a list of serial numbers.
// The name parameter specifies the configured DEP name to use.
// Note we use HTTP PUT for compatibility despite modern documentation
// listing HTTP POST for this endpoint.
// See https://developer.apple.com/documentation/devicemanagement/assign_a_profile
func (c *Client) AssignProfile(ctx context.Context, name, uuid string, serials ...string) (*ProfileResponse, error) {
	req := &struct {
		ProfileUUID string   `json:"profile_uuid"`
		Devices     []string `json:"devices"`
	}{
		ProfileUUID: uuid,
		Devices:     serials,
	}
	resp := new(ProfileResponse)
	// historically this has been an HTTP PUT and the DEP simulator depsim
	// requires this. however modern Apple documentation says this is a POST
	// now. we still use PUT here for compatibility.
	return resp, c.do(ctx, name, http.MethodPut, "/profile/devices", req, resp)
}

// DefineProfileResponse corresponds to the Apple DEP API "DefineProfileResponse" structure.
// See https://developer.apple.com/documentation/devicemanagement/defineprofileresponse
type DefineProfileResponse struct {
	ProfileUUID string   `json:"profile_uuid"`
	Devices     []string `json:"devices"`
}

// DefineProfile uses the Apple "Define a Profile" command to attempt to create a profile.
// This service defines a profile with Apple's servers that can then be assigned to specific devices.
// This command provides information about the MDM server that is assigned to manage one or more devices,
// information about the host that the managed devices can pair with, and various attributes that control
// the MDM association behavior of the device.
// See https://developer.apple.com/documentation/devicemanagement/define_a_profile
func (c *Client) DefineProfile(ctx context.Context, name string, profile *Profile) (*DefineProfileResponse, error) {
	resp := new(DefineProfileResponse)
	return resp, c.do(ctx, name, http.MethodPost, "/profile", profile, resp)
}

// ClearProfileResponse corresponds to the Apple DEP API "ClearProfileResponse" structure.
// See https://developer.apple.com/documentation/devicemanagement/clearprofileresponse
type ClearProfileResponse struct {
	Devices map[string]string `json:"devices"`
}

// RemoveProfile uses the Apple "Remove a Profile" API endpoint to "unassign"
// any DEP profile UUID from a list of serial numbers.
// A `profile_uuid` API paramater is listed in the documentation but we do not
// support it (nor does it appear to be used on the server-side).
// The name parameter specifies the configured DEP name to use.
// See https://developer.apple.com/documentation/devicemanagement/remove_a_profile-c2c
func (c *Client) RemoveProfile(ctx context.Context, name string, devices []string) (*ClearProfileResponse, error) {
	req := &struct {
		// ProfileUUID string `json:"profile_uuid,omitempty"`
		Devices []string `json:"devices"`
	}{
		Devices: devices,
	}
	resp := new(ClearProfileResponse)
	return resp, c.do(ctx, name, http.MethodDelete, "/profile/devices", req, resp)
}
