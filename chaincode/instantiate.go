package chaincode

import (
	"strings"

	"fabric-tools/types"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	"github.com/pkg/errors"
)

// InstantiateCC 实例化链码
func (M *ChainCodeManager) InstantiateCC(cc *types.ChainCode, args [][]byte) error {
	resMgmtClient, err := M.ResourceMgmtClient()
	if err != nil {
		return err
	}

	chaincodePolicy, err := M.newChaincodePolicy(cc)
	if err != nil {
		return err
	}

	req := resmgmt.InstantiateCCRequest{
		Name:    cc.Name,
		Path:    cc.Path,
		Version: cc.Version,
		Args:    args,
		Policy:  chaincodePolicy,
		// CollConfig: collConfig,
	}

	_, err = resMgmtClient.InstantiateCC(M.Cfg.ChannelName(), req, resmgmt.WithTargets(M.Peers()...), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {

		if strings.Contains(err.Error(), "already exists") &&
			strings.Contains(err.Error(), cc.Name) {
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

func (M *ChainCodeManager) newChaincodePolicy(cc *types.ChainCode) (*common.SignaturePolicyEnvelope, error) {
	if cc.ChaincodePolicy != "" {
		// Create a signature policy from the policy expression passed in
		return newChaincodePolicy(cc.ChaincodePolicy)
	}

	// Default policy is 'signed by any member' for all known orgs
	return cauthdsl.AcceptAllPolicy, nil
}

func newChaincodePolicy(policyString string) (*common.SignaturePolicyEnvelope, error) {
	ccPolicy, err := cauthdsl.FromString(policyString)
	if err != nil {
		return nil, errors.Errorf("invalid chaincode policy [%s]: %s", policyString, err)
	}
	return ccPolicy, nil
}
