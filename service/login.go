package service

import (
	"github.com/io24m/hammer/shared"
	"github.com/io24m/hammer/util"
)

func Login(query *shared.Query) (string, error) {
	var options *shared.Options
	var data interface{}
	request := util.CreatRequest(shared.POST, shared.LOGIN, data, options)
	return request()
}
