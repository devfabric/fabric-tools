package chaincode

import (
	"fmt"
	"strings"
)

// IsChainCodeNotFoundError 判断链码是否正确安装
func IsChainCodeNotFoundError(err error) bool {
	msg := fmt.Sprintf("%s", err)
	return strings.Contains(msg, "CHAINCODE_NAME_NOT_FOUND")
}
