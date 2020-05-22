package config

const (
	// ProgramName 程序名
	ProgramName        = "tools"
	cfgFlag            = "cfg"
	defaultCfgFile     = "./sdk_config.yaml"
	cfgFileDescription = "The path of the config file"

	channelFlag        = "ch"
	defaultChannel     = "mychannel"
	channelDescription = "The channel Name"

	defaultAdminUser = "Admin"

	logLevelFlag        = "log"
	defaultlogLevel     = "INFO"
	logLevelDescription = "The Log level - ERROR, WARN, INFO, DEBUG"

	userFlag        = "user"
	defaultUser     = ""
	userDescription = "The user of client"

	orgsFlag        = "org"
	defaultOrgs     = ""
	orgsDescription = "The orgs of client"

	peerFlag        = "pr"
	defaultPeer     = ""
	peerDescription = "The peers of client"
)
