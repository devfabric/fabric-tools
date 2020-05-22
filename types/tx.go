package types

import (
	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	cb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
	pb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
)

const (
	HeaderType_MESSAGE = iota
	HeaderType_CONFIG
	HeaderType_CONFIG_UPDATE
	HeaderType_ENDORSER_TRANSACTION
	HeaderType_ORDERER_TRANSACTION
	HeaderType_DELIVER_SEEK_INFO
	HeaderType_CHAINCODE_PACKAGE
	HeaderType_PEER_ADMIN_OPERATION
	HeaderType_TOKEN_TRANSACTION
)

// ProcessedTransaction 交易查询
type Transaction struct {
	TxID            string
	TransactionType int32
	ValidationCode  int32 //交易确认码
	ChannelID       string
	Timestamp       *timestamp.Timestamp

	Payload interface{}
}

// GetProcessedTransaction 将pb.ProcessedTransaction转化为可用的ProcessedTransaction
func GetProcessedTransaction(pbpt *pb.ProcessedTransaction) (*Transaction, error) {
	if pbpt == nil {
		return nil, errors.New("pt is nil")
	}

	payload := &cb.Payload{}
	err := proto.Unmarshal(pbpt.TransactionEnvelope.Payload, payload)
	if err != nil {
		return nil, errors.Wrap(err, "no payload in envelope")
	}

	chdr, err := UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	envelope, err := getPayload(int(chdr.Type), payload.Data)
	if err != nil {
		return nil, err
	}

	ret := Transaction{
		ValidationCode:  pbpt.ValidationCode,
		TxID:            chdr.TxId,
		TransactionType: chdr.Type,
		ChannelID:       chdr.ChannelId,
		Timestamp:       chdr.Timestamp,
		Payload:         envelope,
	}

	return &ret, nil
}

func getPayload(headerType int, data []byte) (interface{}, error) {
	if headerType == HeaderType_CONFIG {
		return nil, nil
	} else if headerType == HeaderType_CONFIG_UPDATE {
		return nil, nil
	} else if headerType == HeaderType_ENDORSER_TRANSACTION {
		return getTransaction(data)
	}
	return nil, errors.New("headerType err")
}

type TransactionActions []Action

type Action struct {
	ChaincodeID   string
	ChainCodeArgs [][]byte

	Endorsements [][]byte
	Response     ChaincodeAction
	TxRwSet      string
	Events       Event
	NsRwSets     []NsRwSet
}

type Event struct {
	EventName string
	Payload   []byte
}
type ChaincodeAction struct {
	Status  int32
	Message string
	Payload []byte
}

type NsRwSet struct {
	NameSpace string
	Reads     []KVRead
	Writes    []KVWrite
}
type KVWrite struct {
	Key      string
	IsDelete bool
	Value    []byte
}
type KVRead struct {
	Key string
	Version
}
type Version struct {
	BlockNum uint64
	TxNum    uint64
}

func getTransaction(txBytes []byte) (TransactionActions, error) {
	ret := TransactionActions{}

	tx := &peer.Transaction{}
	err := proto.Unmarshal(txBytes, tx)
	if err != nil {
		return nil, errors.Wrap(err, "Bad envelope:error unmarshaling Transaction")
	}

	for _, action := range tx.Actions {
		cap, err := GetChaincodeActionPayload(action.Payload)
		if err != nil {
			return nil, errors.New(err.Error())
		}

		cis, err := GetChaincodeInvocationSpec(cap.ChaincodeProposalPayload)
		if err != nil {
			return nil, errors.Wrap(err, "error unmarshaling ChaincodeInvocationSpec")
		}

		endorsements := [][]byte{}
		for _, endorsement := range cap.Action.Endorsements {
			endorsements = append(endorsements, endorsement.Endorser)
		}

		chaincodeAction, err := GetChaincodeAction(cap.Action.ProposalResponsePayload)
		if err != nil {
			return nil, errors.Wrap(err, "error unmarshaling GetChaincodeAction")
		}
		ccac := ChaincodeAction{
			Status:  chaincodeAction.Response.Status,
			Message: chaincodeAction.Response.Message,
			Payload: chaincodeAction.Response.Payload,
		}

		events, err := GetEvents(chaincodeAction.Events)
		if err != nil {
			return nil, errors.Wrap(err, "error unmarshaling ChaincodeEvent")
		}

		mytxRWSet, err := GetRWSet(chaincodeAction.Results)
		if err != nil {
			return nil, errors.Wrap(err, "error unmarshaling ChaincodeEvent")
		}

		ac := Action{
			ChaincodeID:   cis.ChaincodeSpec.ChaincodeId.Name,
			ChainCodeArgs: cis.ChaincodeSpec.Input.Args,
			// Endorsements:  endorsements,
			Response: ccac,
			Events:   *events,
			NsRwSets: mytxRWSet,
		}
		ret = append(ret, ac)
	}

	return ret, nil
}
