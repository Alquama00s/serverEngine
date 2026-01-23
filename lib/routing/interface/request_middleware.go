package routingInterface

import (
	"regexp"

	routingModel "github.com/Alquama00s/serverEngine/lib/routing/model"
)

type RequestProcessor interface {
	GetProcessor() func(req *routingModel.Request) (*routingModel.Request, error, *routingModel.Response)
	Process(fun func(req *routingModel.Request) (*routingModel.Request, error, *routingModel.Response))
	GetRegex() *regexp.Regexp
	SetRegex(regex string)
	GetRegexString() string
	GetPriority() int
	SetPriority(priority int)
}
