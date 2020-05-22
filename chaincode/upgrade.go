package chaincode

import (
	"strings"

	"fabric-tools/types"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"

	"github.com/pkg/errors"
)

// UpgradeCC 实例化链码
func (M *ChainCodeManager) UpgradeCC(cc *types.ChainCode, args [][]byte) error {
	resMgmtClient, err := M.ResourceMgmtClient()
	if err != nil {
		return err
	}

	chaincodePolicy, err := M.newChaincodePolicy(cc)
	if err != nil {
		return err
	}

	req := resmgmt.UpgradeCCRequest{
		Name:    cc.Name,
		Path:    cc.Path,
		Version: cc.Version,
		Args:    args,
		Policy:  chaincodePolicy,
		// CollConfig: collConfig,
	}

	_, err = resMgmtClient.UpgradeCC(M.Cfg.ChannelName(), req, resmgmt.WithTargets(M.Peers()...), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		if strings.Contains(err.Error(), "version already exists for chaincode with name '"+cc.Name+"'") {
			// Ignore
			M.Log().Warnf("Chaincode %s already instantiated.", cc.Name)
			M.Log().Warnf("...chaincode %s already instantiated.\n", cc.Name)
			return nil
		}
		return errors.Errorf("error instantiating chaincode: %v", err)
	}

	M.Log().Infof("...successfuly instantiated chaincode %s on channel %s.\n", cc.Name, M.Cfg.ChannelName())

	return nil
}
