package permission

// Permission : base Permission object model
type Permission struct {
	Reference    string `json:"_ref,omitempty"`
	Group        string `json:"group,omitempty"`
	Object       string `json:"object,omitempty"`
	Permission   string `json:"permission,omitempty"`
	ResourceType string `json:"resource_type,omitempty"`
	Role         string `json:"role,omitempty"`
}
