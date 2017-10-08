package core

import (
	"mvbc/core/utils"
	"reflect"
)

var (
	TRANSACTION_POW = utils.ArrayOfBytes(TRANSACTION_POW_COMPLEXITY, POW_PREFIX)
	BLOCK_POW       = utils.ArrayOfBytes(BLOCK_POW_COMPLEXITY, POW_PREFIX)

	TEST_TRANSACTION_POW = utils.ArrayOfBytes(TEST_TRANSACTION_POW_COMPLEXITY, POW_PREFIX)
	TEST_BLOCK_POW       = utils.ArrayOfBytes(TEST_BLOCK_POW_COMPLEXITY, POW_PREFIX)
)

func CheckProofOfWork(prefix []byte, hash []byte) bool {

	if len(prefix) > 0 {
		return reflect.DeepEqual(prefix, hash[:len(prefix)])
	}
	return true
}
