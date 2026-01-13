package config

import "github.com/moby/moby/api/types/container"

type ConfigType struct {
	Server   ServerConfig   `json:"server"`
	Network  NetworkConfig  `json:"network"`
	Security SecurityConfig `json:"security"`
	Docker   DockerConfig   `json:"docker"`
	Git      GitConfig      `json:"git"`
}

type ServerConfig struct {
	Port                 *int `json:"port,omitempty"`
	DefaultAppPort       *int `json:"defaultAppPort,omitempty"`
	APIReadHeaderTimeout *int `json:"apiReadHeaderTimeout,omitempty"`
	MaxAvatarSize        *int `json:"maxAvatarSize,omitempty"`
}

type NetworkConfig struct {
	WildcardDomain       *string `json:"wildCardDomain,omitempty"`
	MistAppName          *string `json:"mistAppName,omitempty"`
	DNSValidationTimeout *int    `json:"dnsValidationTimeout,omitempty"`
}

type SecurityConfig struct {
	JWTExpiry         *int  `json:"jwtExpiry,omitempty"`
	SecureCookies     *bool `json:"secureCookies,omitempty"`
	PasswordMinLength *int  `json:"passwordMinLength,omitempty"`
}

type DockerConfig struct {
	AutoCleanupContainers   *bool                        `json:"autoCleanupContainers,omitempty"`
	DefaultRestartPolicy    *container.RestartPolicyMode `json:"defaultRestartPolicy,omitempty"`
	BuildImageTimeout       *int                         `json:"buildImageTimeout,omitempty"`
	PullImageTimeout        *int                         `json:"pullImageTimeout,omitempty"`
	StartContainerTimeout   *int                         `json:"startContainerTimeout,omitempty"`
	StopContainerTimeout    *int                         `json:"stopContainerTimeout,omitempty"`
	RestartContainerTimeout *int                         `json:"restartContainerTimeout,omitempty"`
}

type GitConfig struct {
	GitCloneTimeout         *int  `json:"gitCloneTimeout,omitempty"`
	RemoveGitRepoAfterBuild *bool `json:"removeGitRepoAfterBuild,omitempty"`
}

type ResolvedServerConfig struct {
	Port                 int
	DefaultAppPort       int
	APIReadHeaderTimeout int
	MaxAvatarSize        int
}

type ResolvedNetworkConfig struct {
	WildcardDomain       string
	MistAppName          string
	DNSValidationTimeout int
}

type ResolvedSecurityConfig struct {
	JWTExpiry         int
	SecureCookies     bool
	PasswordMinLength int
}

type ResolvedDockerConfig struct {
	AutoCleanupContainers   bool
	DefaultRestartPolicy    container.RestartPolicyMode
	BuildImageTimeout       int
	PullImageTimeout        int
	StartContainerTimeout   int
	StopContainerTimeout    int
	RestartContainerTimeout int
}

type ResolvedGitConfig struct {
	GitCloneTimeout         int
	RemoveGitRepoAfterBuild bool
}
