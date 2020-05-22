package manager

import (
	"fmt"
	"math/rand"
	"os"

	"code.uni-ledger.com/switch/fabric-tools/config"
	"code.uni-ledger.com/switch/fabric-tools/printer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	mspapi "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	contextImpl "github.com/hyperledger/fabric-sdk-go/pkg/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/orderer"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Manager 管理sdk
type Manager struct {
	Cfg    *config.Cfg
	Logger *logrus.Logger

	sdk            *fabsdk.FabricSDK
	endpointConfig fab.EndpointConfig

	orgIDByPeer map[string]string
	peersByOrg  map[string][]fab.Peer
	peers       []fab.Peer

	sessions map[string]context.ClientProvider

	Printer printer.Printer
}

// NewManager 返回配置实体
func NewManager(cfg *config.Cfg) (*Manager, error) {
	var options []fabsdk.Option
	options = append(options, fabsdk.WithCorePkg(&cryptoSuiteProviderFactory{}))

	sdk, err := fabsdk.New(cfg.Provider(), options...)
	if err != nil {
		return nil, errors.Errorf("Error initializing SDK: %s", err)
	}

	ctx, err := sdk.Context()()
	if err != nil {
		return nil, errors.WithMessage(err, "Error creating anonymous provider")
	}
	endpointConfig := ctx.EndpointConfig()
	networkConfig := endpointConfig.NetworkConfig()

	orgIDByPeer := make(map[string]string)

	var allPeers []fab.Peer
	allPeersByOrg := make(map[string][]fab.Peer)

	for orgID := range networkConfig.Organizations {
		peersConfig, ok := endpointConfig.PeersConfig(orgID)
		if !ok {
			return nil, errors.Errorf("failed to get peer configs for org [%s]", orgID)
		}

		// cfg.Logger.Debugf("Peers for org [%s]: %v\n", orgID, peersConfig)

		var peers []fab.Peer
		for _, p := range peersConfig {
			endorser, err := ctx.InfraProvider().CreatePeerFromConfig(&fab.NetworkPeer{PeerConfig: p})
			if err != nil {
				return nil, errors.Wrapf(err, "failed to create peer from config")
			}
			peers = append(peers, endorser)
			orgIDByPeer[endorser.URL()] = orgID
		}
		allPeersByOrg[orgID] = peers
		allPeers = append(allPeers, peers...)
	}

	m := Manager{
		Cfg:    cfg,
		Logger: cfg.Logger,

		sdk:            sdk,
		endpointConfig: endpointConfig,

		sessions: make(map[string]context.ClientProvider),
	}

	// 筛选出指定的节点
	peers, err := m.getPeers(allPeers, cfg.PeerURLs(), cfg.OrgIDs())
	if err != nil {
		return nil, err
	}

	// 通过组织id存放peer列表
	peersByOrg := make(map[string][]fab.Peer)
	cfg.Logger.Debugf("Selected Peers:\n")
	for _, peer := range peers {
		cfg.Logger.Debugf("- URL: %s\n", peer.URL())
		orgID := orgIDByPeer[peer.URL()]
		if orgID == "" {
			return nil, errors.Errorf("unable to find org for peer: %s", peer.URL())
		}
		peersByOrg[orgID] = append(peersByOrg[orgID], peer)
	}

	m.peers = peers
	m.peersByOrg = peersByOrg
	m.orgIDByPeer = orgIDByPeer

	m.CheckManager()

	m.Printer = printer.NewBlockPrinterWithOpts(
		printer.AsOutputFormat("DISPLAY"),
		printer.AsWriterType("STDOUT"),
		&printer.FormatterOpts{Base64Encode: true})

	return &m, nil
}

// NewManagerFromOpts 从自定义配置生成Manager
func NewManagerFromOpts(opts *config.Opts) (*Manager, error) {
	cfg, err := config.NewCfg(opts)
	if err != nil {
		return nil, err
	}
	return NewManager(cfg)
}

// LedgerClient 使用当前用户创建一个帐本客户端
func (M *Manager) LedgerClient(user string) (*ledger.Client, error) {
	channelProvider, err := M.ChannelProvider(user, M.OrgID())
	if err != nil {
		return nil, errors.Errorf("error creating channel provider: %v", err)
	}
	c, err := ledger.New(channelProvider)
	if err != nil {
		return nil, errors.Errorf("error creating new ledger client: %v", err)
	}
	return c, nil
}

// User 返回一个用户实体
func (M *Manager) User() (mspapi.SigningIdentity, error) {
	return M.OrgUser(M.OrgID(), M.Cfg.OperateUser)
}

// LocalContext 创建一个本地实体
func (M *Manager) LocalContext() (context.Local, error) {
	user, err := M.User()
	if err != nil {
		return nil, errors.Errorf("error getting user: %s", err)
	}
	contextProvider, err := M.context(user)
	if err != nil {
		return nil, errors.Errorf("error getting context for user [%s,%s]: %v", user.Identifier().MSPID, user.Identifier().ID, err)
	}
	return contextImpl.NewLocal(contextProvider)
}

// ChannelProvider 返回通道提供者实体
func (M *Manager) ChannelProvider(username string, orgID string) (context.ChannelProvider, error) {
	if M.IsExsitOrgID(orgID) != true {
		return nil, errors.Errorf("error origanization id (%s) is not exsit", orgID)
	}

	channelID := M.Cfg.ChannelName()
	user, err := M.OrgUser(orgID, username)
	if err != nil {
		return nil, err
	}
	M.Log().Debugf("creating channel provider for user [%s] in org [%s]...", user.Identifier().ID, user.Identifier().MSPID)
	clientContext, err := M.context(user)
	if err != nil {
		return nil, errors.Errorf("error getting client context for user [%s,%s]: %v", user.Identifier().MSPID, user.Identifier().ID, err)
	}
	channelProvider := func() (context.Channel, error) {
		return contextImpl.NewChannel(clientContext, channelID)
	}
	return channelProvider, nil
}

// ChannelClient 返回通道客户端
func (M *Manager) ChannelClient(username string, orgID string, opts ...channel.ClientOption) (*channel.Client, error) {
	channelProvider, err := M.ChannelProvider(username, orgID)
	if err != nil {
		return nil, err
	}
	return channel.New(channelProvider, opts...)
}

// PeerFromURL 返回指定peer的实体
func (M *Manager) PeerFromURL(url string) (fab.Peer, bool) {
	for _, peer := range M.peers {
		if url == peer.URL() {
			return peer, true
		}
	}
	return nil, false
}

// PeersByOrg 返回带peer的组织map
func (M *Manager) PeersByOrg() map[string][]fab.Peer {
	return M.peersByOrg
}

// CheckManager 检查运行时参数
func (M *Manager) CheckManager() {
	if M.Cfg.ChannelName() == "" {
		M.Log().Error("channel id is nil")
		os.Exit(1)
	}

	if len(M.Peers()) == 0 {
		M.Log().Error("no peers to choose from!")
		os.Exit(1)
	}

	peer := M.Peers()[0]
	_, err := M.OrgOfPeer(peer.URL())
	if err != nil {
		M.Log().Error(err)
		os.Exit(1)
	}
}

// RandomOrderer 随机一个从配置中orderer
func (M *Manager) RandomOrderer() (fab.Orderer, error) {
	orderers, err := M.Orderers()
	if err != nil {
		return nil, err
	}
	if len(orderers) == 0 {
		return nil, errors.New("No orders found")
	}
	return orderers[rand.Intn(len(orderers))], nil
}

// Orderers 返回配置的（或指定的）所以orderer
func (M *Manager) Orderers() ([]fab.Orderer, error) {
	ordererConfigs := M.endpointConfig.OrderersConfig()
	ordererURL := M.Cfg.OrdererURL()

	var orderers []fab.Orderer
	for _, ordererConfig := range ordererConfigs {
		if ordererURL == "" || ordererConfig.URL == ordererURL {
			newOrderer, err := orderer.New(M.endpointConfig, orderer.FromOrdererConfig(&ordererConfig))
			if err != nil {
				return nil, errors.WithMessage(err, "creating orderer failed")
			}
			orderers = append(orderers, newOrderer)
		}
	}

	return orderers, nil
}

// getPeers 筛选指定peer
func (M *Manager) getPeers(allPeers []fab.Peer, peerURLs []string, orgIDs []string) ([]fab.Peer, error) {
	selectAll := false
	if len(peerURLs) == 0 && len(orgIDs) == 0 {
		selectAll = true
	}
	var selectedPeers []fab.Peer
	var allPeerURLs []string
	for _, peer := range allPeers {
		allPeerURLs = append(allPeerURLs, peer.URL())
		orgID := M.orgIDByPeer[peer.URL()]
		if selectAll || containsString(peerURLs, peer.URL()) || len(peerURLs) == 0 && containsString(orgIDs, orgID) {
			selectedPeers = append(selectedPeers, peer)
		}
	}

	for _, url := range peerURLs {
		if !containsString(allPeerURLs, url) {
			return nil, fmt.Errorf("invalid peer URL: %s", url)
		}
	}
	return selectedPeers, nil
}

// Log 返回log实体
func (M *Manager) Log() *logrus.Logger {
	return M.Logger
}

// OrgOfPeer 返回指定peer的组织id
func (M *Manager) OrgOfPeer(peerURL string) (string, error) {
	orgID, ok := M.orgIDByPeer[peerURL]
	if !ok {
		return "", errors.Errorf("org not found for peer %s", peerURL)
	}
	return orgID, nil
}

// Peers 返回所有peer
func (M *Manager) Peers() []fab.Peer {
	peers := []fab.Peer{}
	peerURLs := M.Cfg.Opts.PeerURLs()

	for _, peer := range M.peers {
		for _, url := range peerURLs {
			if url == peer.URL() {
				peers = append(peers, peer)
			}
		}
	}

	if len(peers) == 0 {
		peers = M.peers
	}

	return peers
}

// OrgID 返回peer列表中第一个的组织id
func (M *Manager) OrgID() string {
	orgID, _ := M.OrgOfPeer(M.Peers()[0].URL())
	return orgID
}

// OrgUser 返回指定组织的用户实体
func (M *Manager) OrgUser(orgID, username string) (mspapi.SigningIdentity, error) {
	if username == "" {
		return nil, errors.Errorf("no username specified")
	}

	mspClient, err := msp.New(M.sdk.Context(), msp.WithOrg(orgID))
	if err != nil {
		return nil, errors.Errorf("error creating MSP client: %s", err)
	}

	user, err := mspClient.GetSigningIdentity(username)
	if err != nil {
		return nil, errors.Errorf("GetSigningIdentity returned error: %v", err)
	}

	return user, nil
}

// OrgAdminUser 返回给定组织的预注册管理用户
func (M *Manager) OrgAdminUser(orgID string) (mspapi.SigningIdentity, error) {
	userName := M.Cfg.Admin()
	return M.OrgUser(orgID, userName)
}

// ResourceMgmtClient 返回当前用户的资源管理客户端
func (M *Manager) ResourceMgmtClient() (*resmgmt.Client, error) {
	return M.ResourceMgmtClientForOrg(M.OrgID())
}

// IsExsitOrgID 判断指定组织id是否存在
func (M *Manager) IsExsitOrgID(orgID string) bool {
	for _, value := range M.orgIDByPeer {
		if orgID == value {
			return true
		}
	}

	return false
}

// ResourceMgmtClientForOrg 给定组织的资源管理客户端
func (M *Manager) ResourceMgmtClientForOrg(orgID string) (*resmgmt.Client, error) {
	user, err := M.OrgAdminUser(orgID)
	if err != nil {
		return nil, err
	}
	return M.ResourceMgmtClientForUser(user)
}

// ResourceMgmtClientForUser 返回给定用户的Fabric客户端
func (M *Manager) ResourceMgmtClientForUser(user mspapi.SigningIdentity) (*resmgmt.Client, error) {
	M.Log().Debugf("create resmgmt client for user [%s] in org [%s]...", user.Identifier().ID, user.Identifier().MSPID)
	session, err := M.context(user)
	if err != nil {
		return nil, errors.Errorf("error getting session for user [%s,%s]: %v", user.Identifier().MSPID, user.Identifier().ID, err)
	}
	c, err := resmgmt.New(session)
	if err != nil {
		return nil, errors.Errorf("error creating new resmgmt client for user [%s,%s]: %v", user.Identifier().MSPID, user.Identifier().ID, err)
	}
	return c, nil
}

func (M *Manager) context(user mspapi.SigningIdentity) (context.ClientProvider, error) {
	key := user.Identifier().MSPID + "_" + user.Identifier().ID
	session := M.sessions[key]
	if session == nil {
		session = M.sdk.Context(fabsdk.WithIdentity(user))
		M.Log().Debugf("Created session for user [%s] in org [%s]", user.Identifier().ID, user.Identifier().MSPID)
		M.sessions[key] = session
	}
	return session, nil
}

// EventClient 返回事件客户端
func (M *Manager) EventClient(user string, opts ...event.ClientOption) (*event.Client, error) {
	channelProvider, err := M.ChannelProvider(user, M.OrgID())
	if err != nil {
		return nil, errors.Errorf("error creating channel provider: %v", err)
	}
	c, err := event.New(channelProvider, opts...)
	if err != nil {
		return nil, errors.Errorf("error creating new event client: %v", err)
	}
	return c, nil
}
