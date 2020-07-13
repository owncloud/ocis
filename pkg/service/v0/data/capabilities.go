package data

import (
	"encoding/xml"
)

// ocsBool implements the xml/json Marshaler interface. The OCS API inconsistency require us to parse boolean values
// as native booleans for json requests but "truthy" 0/1 values for xml requests.
type ocsBool bool

func (c *ocsBool) MarshalJSON() ([]byte, error) {
	if *c {
		return []byte("true"), nil
	}

	return []byte("false"), nil
}

func (c ocsBool) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if c {
		return e.EncodeElement("1", start)
	}

	return e.EncodeElement("0", start)
}

// CapabilitiesData TODO document
type CapabilitiesData struct {
	Capabilities *Capabilities `json:"capabilities" xml:"capabilities"`
	Version      *Version      `json:"version" xml:"version"`
}

// Capabilities groups several capability aspects
type Capabilities struct {
	Core          *CapabilitiesCore          `json:"core" xml:"core"`
	Checksums     *CapabilitiesChecksums     `json:"checksums" xml:"checksums"`
	Files         *CapabilitiesFiles         `json:"files" xml:"files" mapstructure:"files"`
	Dav           *CapabilitiesDav           `json:"dav" xml:"dav"`
	FilesSharing  *CapabilitiesFilesSharing  `json:"files_sharing" xml:"files_sharing" mapstructure:"files_sharing"`
	Notifications *CapabilitiesNotifications `json:"notifications" xml:"notifications"`
}

// CapabilitiesCore holds webdav config
type CapabilitiesCore struct {
	PollInterval int     `json:"pollinterval" xml:"pollinterval" mapstructure:"poll_interval"`
	WebdavRoot   string  `json:"webdav-root,omitempty" xml:"webdav-root,omitempty" mapstructure:"webdav_root"`
	Status       *Status `json:"status" xml:"status"`
}

// Status holds basic status information
type Status struct {
	Installed      ocsBool `json:"installed" xml:"installed"`
	Maintenance    ocsBool `json:"maintenance" xml:"maintenance"`
	NeedsDBUpgrade ocsBool `json:"needsDbUpgrade" xml:"needsDbUpgrade"`
	Version        string  `json:"version" xml:"version"`
	VersionString  string  `json:"versionstring" xml:"versionstring"`
	Edition        string  `json:"edition" xml:"edition"`
	ProductName    string  `json:"productname" xml:"productname"`
	Hostname       string  `json:"hostname,omitempty" xml:"hostname,omitempty"`
}

// CapabilitiesChecksums holds available hashes
type CapabilitiesChecksums struct {
	SupportedTypes      []string `json:"supportedTypes" xml:"supportedTypes>element" mapstructure:"supported_types"`
	PreferredUploadType string   `json:"preferredUploadType" xml:"preferredUploadType" mapstructure:"preferred_upload_type"`
}

// CapabilitiesFilesTusSupport TODO this must be a summary of storages
type CapabilitiesFilesTusSupport struct {
	Version            string `json:"version" xml:"version"`
	Resumable          string `json:"resumable" xml:"resumable"`
	Extension          string `json:"extension" xml:"extension"`
	MaxChunkSize       int    `json:"max_chunk_size" xml:"max_chunk_size" mapstructure:"max_chunk_size"`
	HTTPMethodOverride string `json:"http_method_override" xml:"http_method_override" mapstructure:"http_method_override"`
}

// CapabilitiesFiles TODO this is storage specific, not global. What effect do these options have on the clients?
type CapabilitiesFiles struct {
	PrivateLinks     ocsBool                      `json:"privateLinks" xml:"privateLinks" mapstructure:"private_links"`
	BigFileChunking  ocsBool                      `json:"bigfilechunking" xml:"bigfilechunking"`
	Undelete         ocsBool                      `json:"undelete" xml:"undelete"`
	Versioning       ocsBool                      `json:"versioning" xml:"versioning"`
	BlacklistedFiles []string                     `json:"blacklisted_files" xml:"blacklisted_files>element" mapstructure:"blacklisted_files"`
	TusSupport       *CapabilitiesFilesTusSupport `json:"tus_support" xml:"tus_support" mapstructure:"tus_support"`
}

// CapabilitiesDav holds dav endpoint config
type CapabilitiesDav struct {
	Chunking                       string   `json:"chunking" xml:"chunking"`
	Trashbin                       string   `json:"trashbin" xml:"trashbin"`
	Reports                        []string `json:"reports" xml:"reports>element" mapstructure:"reports"`
	ChunkingParallelUploadDisabled bool     `json:"chunkingParallelUploadDisabled" xml:"chunkingParallelUploadDisabled"`
}

// CapabilitiesFilesSharing TODO document
type CapabilitiesFilesSharing struct {
	APIEnabled                    ocsBool                                  `json:"api_enabled" xml:"api_enabled" mapstructure:"api_enabled"`
	Resharing                     ocsBool                                  `json:"resharing" xml:"resharing"`
	GroupSharing                  ocsBool                                  `json:"group_sharing" xml:"group_sharing" mapstructure:"group_sharing"`
	AutoAcceptShare               ocsBool                                  `json:"auto_accept_share" xml:"auto_accept_share" mapstructure:"auto_accept_share"`
	ShareWithGroupMembersOnly     ocsBool                                  `json:"share_with_group_members_only" xml:"share_with_group_members_only" mapstructure:"share_with_group_members_only"`
	ShareWithMembershipGroupsOnly ocsBool                                  `json:"share_with_membership_groups_only" xml:"share_with_membership_groups_only" mapstructure:"share_with_membership_groups_only"`
	SearchMinLength               int                                      `json:"search_min_length" xml:"search_min_length" mapstructure:"search_min_length"`
	DefaultPermissions            int                                      `json:"default_permissions" xml:"default_permissions" mapstructure:"default_permissions"`
	UserEnumeration               *CapabilitiesFilesSharingUserEnumeration `json:"user_enumeration" xml:"user_enumeration" mapstructure:"user_enumeration"`
	Federation                    *CapabilitiesFilesSharingFederation      `json:"federation" xml:"federation"`
	Public                        *CapabilitiesFilesSharingPublic          `json:"public" xml:"public"`
	User                          *CapabilitiesFilesSharingUser            `json:"user" xml:"user"`
}

// CapabilitiesFilesSharingPublic TODO document
type CapabilitiesFilesSharingPublic struct {
	Enabled            ocsBool                                   `json:"enabled" xml:"enabled"`
	SendMail           ocsBool                                   `json:"send_mail" xml:"send_mail" mapstructure:"send_mail"`
	SocialShare        ocsBool                                   `json:"social_share" xml:"social_share" mapstructure:"social_share"`
	Upload             ocsBool                                   `json:"upload" xml:"upload"`
	Multiple           ocsBool                                   `json:"multiple" xml:"multiple"`
	SupportsUploadOnly ocsBool                                   `json:"supports_upload_only" xml:"supports_upload_only" mapstructure:"supports_upload_only"`
	Password           *CapabilitiesFilesSharingPublicPassword   `json:"password" xml:"password"`
	ExpireDate         *CapabilitiesFilesSharingPublicExpireDate `json:"expire_date" xml:"expire_date" mapstructure:"expire_date"`
}

// CapabilitiesFilesSharingPublicPassword TODO document
type CapabilitiesFilesSharingPublicPassword struct {
	EnforcedFor *CapabilitiesFilesSharingPublicPasswordEnforcedFor `json:"enforced_for" xml:"enforced_for" mapstructure:"enforced_for"`
	Enforced    ocsBool                                            `json:"enforced" xml:"enforced"`
}

// CapabilitiesFilesSharingPublicPasswordEnforcedFor TODO document
type CapabilitiesFilesSharingPublicPasswordEnforcedFor struct {
	ReadOnly   ocsBool `json:"read_only" xml:"read_only,omitempty" mapstructure:"read_only"`
	ReadWrite  ocsBool `json:"read_write" xml:"read_write,omitempty" mapstructure:"read_write"`
	UploadOnly ocsBool `json:"upload_only" xml:"upload_only,omitempty" mapstructure:"upload_only"`
}

// CapabilitiesFilesSharingPublicExpireDate TODO document
type CapabilitiesFilesSharingPublicExpireDate struct {
	Enabled ocsBool `json:"enabled" xml:"enabled"`
}

// CapabilitiesFilesSharingUser TODO document
type CapabilitiesFilesSharingUser struct {
	SendMail ocsBool `json:"send_mail" xml:"send_mail" mapstructure:"send_mail"`
}

// CapabilitiesFilesSharingUserEnumeration TODO document
type CapabilitiesFilesSharingUserEnumeration struct {
	Enabled          ocsBool `json:"enabled" xml:"enabled"`
	GroupMembersOnly ocsBool `json:"group_members_only" xml:"group_members_only" mapstructure:"group_members_only"`
}

// CapabilitiesFilesSharingFederation holds outgoing and incoming flags
type CapabilitiesFilesSharingFederation struct {
	Outgoing ocsBool `json:"outgoing" xml:"outgoing"`
	Incoming ocsBool `json:"incoming" xml:"incoming"`
}

// CapabilitiesNotifications holds a list of notification endpoints
type CapabilitiesNotifications struct {
	Endpoints []string `json:"ocs-endpoints" xml:"ocs-endpoints>element" mapstructure:"endpoints"`
}

// Version holds version information
type Version struct {
	Major   int    `json:"major" xml:"major"`
	Minor   int    `json:"minor" xml:"minor"`
	Micro   int    `json:"micro" xml:"micro"` // = patch level
	String  string `json:"string" xml:"string"`
	Edition string `json:"edition" xml:"edition"`
}
