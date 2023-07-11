package smilecdr

type ClientSecret struct {
	Pid         int    `json:"pid,omitempty"`
	Secret      string `json:"secret,omitempty"`
	Description string `json:"description,omitempty"`
	Expiration  string `json:"expiration,omitempty"`
	activation  string `json:"activation,omitempty"`
}

type UserPermission struct {
	Permission string `json:"permission,omitempty"`
	Argument   string `json:"argument,omitempty"`
}

type OpenIdClient struct {
	Pid                         int              `json:"pid,omitempty"`
	NodeId                      string           `json:"nodeId,omitempty"`
	ModuleId                    string           `json:"moduleId,omitempty"`
	AccessTokenValiditySeconds  int              `json:"accessTokenValiditySeconds,omitempty"`
	AllowedGrantTypes           []string         `json:"allowedGrantTypes,omitempty"`
	AutoApproveScopes           []string         `json:"autoApproveScopes,omitempty"`
	AutoGrantScopes             []string         `json:"autoGrantScopes,omitempty"`
	ClientId                    string           `json:"clientId,omitempty"`
	ClientName                  string           `json:"clientName,omitempty"`
	ClientSecrets               []ClientSecret   `json:"clientSecrets,omitempty"`
	FixedScope                  bool             `json:"fixedScope,omitempty"`
	RefreshTokenValiditySeconds int              `json:"refreshTokenValiditySeconds,omitempty"`
	RegisteredRedirectUris      []string         `json:"registeredRedirectUris,omitempty"`
	Scopes                      []string         `json:"scopes,omitempty"`
	SecretRequired              bool             `json:"secretRequired,omitempty"`
	SecretClientCanChange       bool             `json:"secretClientCanChange,omitempty"`
	Enabled                     bool             `json:"enabled,omitempty"`
	CanIntrospectAnyTokens      bool             `json:"canIntrospectAnyTokens,omitempty"`
	CanIntrospectOwnTokens      bool             `json:"canIntrospectOwnTokens,omitempty"`
	AlwaysRequireApproval       bool             `json:"alwaysRequireApproval,omitempty"`
	CanReissueTokens            bool             `json:"canReissueTokens,omitempty"`
	Permissions                 []UserPermission `json:"permissions,omitempty"`
	AttestationAccepted         bool             `json:"rememberedScopes,omitempty"`
	PublicJwksUri               string           `json:"publicJwksUri,omitempty"`
	ArchivedAt                  string           `json:"archivedAt,omitempty"`
	CreatedByAppSphere          string           `json:"createdByAppSphere,omitempty"`
}
