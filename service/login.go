package service

import (
	"github.com/io24m/hammer/shared"
	"github.com/io24m/hammer/util"
)

func Login(query *shared.Query) string {
	request := util.CreatRequest(shared.POST, shared.LOGIN, nil, nil)
	return request()
}
