package geminicli

// LoadCodeAssist types for project discovery
type LoadCodeAssistRequest struct {
	CloudaicompanionProject string         `json:"cloudaicompanionProject,omitempty"`
	Metadata                ClientMetadata `json:"metadata"`
}

type ClientMetadata struct {
	IDEType     string `json:"ideType,omitempty"`
	Platform    string `json:"platform,omitempty"`
	PluginType  string `json:"pluginType,omitempty"`
	DuetProject string `json:"duetProject,omitempty"`
}

type LoadCodeAssistResponse struct {
	CurrentTier             *UserTier `json:"currentTier,omitempty"`
	CloudaicompanionProject string    `json:"cloudaicompanionProject,omitempty"`
}

type UserTier struct {
	ID string `json:"id"`
}
