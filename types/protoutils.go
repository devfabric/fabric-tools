/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package types

import (
	"github.com/gogo/protobuf/proto"
	cb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
	pb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/pkg/errors"
)

// UnmarshalChannelHeader returns a ChannelHeader from bytes
func UnmarshalChannelHeader(bytes []byte) (*cb.ChannelHeader, error) {
	chdr := &cb.ChannelHeader{}
	err := proto.Unmarshal(bytes, chdr)
	return chdr, errors.Wrap(err, "error unmarshaling ChannelHeader")
}

// GetChaincodeActionPayload Get ChaincodeActionPayload from bytes
func GetChaincodeActionPayload(capBytes []byte) (*peer.ChaincodeActionPayload, error) {
	cap := &peer.ChaincodeActionPayload{}
	err := proto.Unmarshal(capBytes, cap)
	return cap, errors.Wrap(err, "error unmarshaling ChaincodeActionPayload")
}

// GetChaincodeInvocationSpec Get GetChaincodeInvocationSpec from bytes
func GetChaincodeInvocationSpec(chaincodeProposalPayload []byte) (*pb.ChaincodeInvocationSpec, error) {
	cpp := &pb.ChaincodeProposalPayload{}
	err := proto.Unmarshal(chaincodeProposalPayload, cpp)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshaling ChaincodeProposalPayload")
	}

	cis := &pb.ChaincodeInvocationSpec{}
	err = proto.Unmarshal(cpp.Input, cis)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshaling ChaincodeInvocationSpec")
	}

	return cis, nil
}

// GetChaincodeAction Get GetChaincodeAction from bytes
func GetChaincodeAction(ProposalResponsePayload []byte) (*pb.ChaincodeAction, error) {
	prp := &pb.ProposalResponsePayload{}
	err := proto.Unmarshal(ProposalResponsePayload, prp)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshaling ProposalResponsePayload")
	}
	chaincodeAction := &pb.ChaincodeAction{}
	err = proto.Unmarshal(prp.Extension, chaincodeAction)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshaling ChaincodeAction")
	}
	return chaincodeAction, nil
}

// GetEvents Get GetEvents from bytes
func GetEvents(eventByte []byte) (*Event, error) {
	events := Event{}
	if len(eventByte) > 0 {
		chaincodeEvent := &pb.ChaincodeEvent{}
		err := proto.Unmarshal(eventByte, chaincodeEvent)
		if err != nil {
			return nil, errors.Wrap(err, "error unmarshaling ChaincodeEvent")
		}
		events.EventName = chaincodeEvent.EventName
		events.Payload = chaincodeEvent.Payload
	}
	return &events, nil
}

// GetRWSet Get GetRWSet from bytes
func GetRWSet(setByte []byte) ([]NsRwSet, error) {
	mytxRWSet := []NsRwSet{}
	if len(setByte) > 0 {
		txRWSet := &rwsetutil.TxRwSet{}
		if err := txRWSet.FromProtoBytes(setByte); err != nil {
			return nil, errors.Wrap(err, "get txrwset error")
		}
		for _, nsRWSet := range txRWSet.NsRwSets {
			if nsRWSet.KvRwSet == nil {
				continue
			}
			reads := []KVRead{}
			for _, r := range nsRWSet.KvRwSet.Reads {
				kvRead := KVRead{
					Key:     r.Key,
					Version: Version{},
				}
				if r.Version != nil {
					kvRead.Version.BlockNum = r.Version.BlockNum
					kvRead.Version.TxNum = r.Version.TxNum
				}
				reads = append(reads, kvRead)
			}

			writes := []KVWrite{}
			for _, w := range nsRWSet.KvRwSet.Writes {
				kvwrite := KVWrite{
					Key:      w.Key,
					IsDelete: w.IsDelete,
					Value:    w.Value,
				}
				writes = append(writes, kvwrite)
			}
			nsRwSet := NsRwSet{
				NameSpace: nsRWSet.NameSpace,
				Reads:     reads,
				Writes:    writes,
			}
			mytxRWSet = append(mytxRWSet, nsRwSet)
		}
	}

	return mytxRWSet, nil
}
