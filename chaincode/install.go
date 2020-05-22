package chaincode

import (
	"net/http"

	"code.uni-ledger.com/switch/fabric-tools/types"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/pkg/errors"
)

// Install 安装指定链码
func (M *ChainCodeManager) Install(cc types.ChainCode, orgID string, targets []fab.Peer) error {
	resMgmtClient, err := M.ResourceMgmtClientForOrg(orgID)
	if err != nil {
		return err
	}

	ccPkg, err := gopackager.NewCCPackage(cc.Path, cc.Gopath)
	if err != nil {
		return err
	}

	req := resmgmt.InstallCCRequest{
		Name:    cc.Name,
		Path:    cc.Path,
		Version: cc.Version,
		Package: ccPkg,
	}

	responses, err := resMgmtClient.InstallCC(req, resmgmt.WithTargets(targets...))
	if err != nil {
		return errors.Errorf("InstallChaincode returned error: %v", err)
	}

	ccIDVersion := cc.Name + "." + cc.Version

	var errs []error
	for _, resp := range responses {
		if resp.Info == "already installed" {
			M.Log().Warnf("Chaincode %s already installed on peer: %s.\n", ccIDVersion, resp.Target)
		} else if resp.Status != http.StatusOK {
			errs = append(errs, errors.Errorf("installCC returned error from peer %s: %s", resp.Target, resp.Info))
		} else {
			M.Log().Infof("...successfuly installed chaincode %s on peer %s.\n", ccIDVersion, resp.Target)
		}
	}

	if len(errs) > 0 {
		M.Log().Warnf("Errors returned from InstallCC: %v\n", errs)
		return errs[0]
	}

	return nil
}
