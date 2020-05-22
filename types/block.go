package types

import (
	"github.com/gogo/protobuf/proto"
	ledgerUtil "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/core/ledger/util"
	cb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	fabriccmn "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"

	"github.com/pkg/errors"
)

// BlockChainInfo 区块信息
type BlockChainInfo struct {
	Height            uint64 `json:"height,omitempty"`
	CurrentBlockHash  []byte `json:"current_block_hash,omitempty"`
	PreviousBlockHash []byte `json:"previous_block_hash,omitempty"`
}

//Block 解释后的区块文件
type Block struct {
	Height       uint64 `json:"height,omitempty"`
	PreviousHash []byte `json:"current_block_hash,omitempty"`
	DataHash     []byte `json:"previous_block_hash,omitempty"`

	Data []Transaction
}

func GetBlock(fabrBlock *fabriccmn.Block) (*Block, error) {
	ret := Block{
		Height:       fabrBlock.Header.Number,
		PreviousHash: fabrBlock.Header.PreviousHash,
		DataHash:     fabrBlock.Header.DataHash,
	}

	validationCode := []int32{}
	txValidationFlags := ledgerUtil.TxValidationFlags(fabrBlock.Metadata.Metadata[fabriccmn.BlockMetadataIndex_TRANSACTIONS_FILTER])
	for i := 0; i < len(txValidationFlags); i++ {
		validationCode = append(validationCode, int32(txValidationFlags[i]))
	}
	validationCodeLen := len(validationCode)

	datas := []Transaction{}
	for i, v := range fabrBlock.Data.Data {
		data, err := getBlockData(v)
		if err != nil {
			return nil, err
		}
		if i < validationCodeLen {
			data.ValidationCode = validationCode[i]
		}
		datas = append(datas, *data)
	}

	ret.Data = datas

	return &ret, nil
}

func getBlockData(dataBytes []byte) (*Transaction, error) {
	env := &cb.Envelope{}
	err := proto.Unmarshal(dataBytes, env)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshaling Envelope")
	}

	payload := &cb.Payload{}
	err = proto.Unmarshal(env.Payload, payload)
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
		// ValidationCode:  validationCode,
		TxID:            chdr.TxId,
		TransactionType: chdr.Type,
		ChannelID:       chdr.ChannelId,
		Timestamp:       chdr.Timestamp,
		Payload:         envelope,
	}

	return &ret, nil
}
