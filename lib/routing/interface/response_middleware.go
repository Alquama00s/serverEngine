package routingInterface

import (
	"regexp"

	routingModel "github.com/Alquama00s/serverEngine/lib/routing/model"
)

type ResponseProcessor interface {
	GetProcessor() func(req *routingModel.Response) (*routingModel.Response, error)
	Process(fun func(req *routingModel.Response) (*routingModel.Response, error))
	GetRegex() *regexp.Regexp
	SetRegex(regex string)
	GetRegexString() string
	GetPriority() int
	SetPriority(priority int)
}
