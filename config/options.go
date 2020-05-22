package config

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

// Opts 自定义配置
type Opts struct {
	SdkCfgFile string
	LogLevel   string
	Logger     *logrus.Logger

	PeerAssignURLs   string //选的哪些可以连接的peer（用，分割节点）
	OrdererAssignURL string //选定连接的orderer（仅支持指定一个）

	Orgs        string //选定哪些可以连接的组织
	AdminUser   string //组织管理员名称
	OperateUser string //操作账户不配置将会初始为AdminUser名
	ChannelID   string //通道名

	// TxFilePath string //创世块交易
}

// NewOpts 返回默认的opts
func NewOpts() *Opts {
	return &Opts{
		SdkCfgFile: defaultCfgFile,
		AdminUser:  defaultAdminUser,
	}
}

// User 返回user名
func (O *Opts) User() string {
	if O.OperateUser == "" {
		O.OperateUser = O.AdminUser
	}
	return O.OperateUser
}

// OrdererURL 返回orderer的url
func (O *Opts) OrdererURL() string {
	return O.OrdererAssignURL
}

// TxFile  用于创建通道的.tx 文件路径
// func (O *Opts) TxFile() string {
// 	return O.TxFilePath
// }

// ChannelName 返回指定的通道
func (O *Opts) ChannelName() string {
	return O.ChannelID
}

// PeerURLs 返回指定的peer的url
func (O *Opts) PeerURLs() []string {
	var urls []string
	if len(strings.TrimSpace(O.PeerAssignURLs)) > 0 {
		s := strings.Split(O.PeerAssignURLs, ",")
		for _, orgID := range s {
			urls = append(urls, orgID)
		}
	}
	return urls
}

// OrgIDs 返回组织id的列表 将opts.Orgs中逗号分割的进行转换成列表
func (O *Opts) OrgIDs() []string {
	var orgIDs []string
	if len(strings.TrimSpace(O.Orgs)) > 0 {
		s := strings.Split(O.Orgs, ",")
		for _, orgID := range s {
			orgIDs = append(orgIDs, orgID)
		}
	}
	return orgIDs
}

// Admin 返回管理员
func (O *Opts) Admin() string {
	if O.AdminUser == "" {
		O.AdminUser = defaultAdminUser
	}
	return O.AdminUser
}

// InitCfgPath 设置配置文件路径
func (O *Opts) InitCfgPath(flags *pflag.FlagSet) {
	flags.StringVar(&O.SdkCfgFile, cfgFlag, defaultCfgFile, cfgFileDescription)
}

// InitChannel 设置通道
func (O *Opts) InitChannel(flags *pflag.FlagSet) {
	flags.StringVar(&O.ChannelID, channelFlag, defaultChannel, channelDescription)
}

// InitLogLevel 设置通道
func (O *Opts) InitLogLevel(flags *pflag.FlagSet) {
	flags.StringVar(&O.LogLevel, logLevelFlag, defaultlogLevel, logLevelDescription)
}

// InitUser 设置通道
func (O *Opts) InitUser(flags *pflag.FlagSet) {
	flags.StringVar(&O.OperateUser, userFlag, defaultUser, userDescription)
}

// InitOrgs 初始化通道
func (O *Opts) InitOrgs(flags *pflag.FlagSet) {
	flags.StringVar(&O.Orgs, orgsFlag, defaultOrgs, orgsDescription)
}

// InitPeers  指定节点
func (O *Opts) InitPeers(flags *pflag.FlagSet) {
	flags.StringVar(&O.PeerAssignURLs, peerFlag, defaultPeer, peerDescription)
}
