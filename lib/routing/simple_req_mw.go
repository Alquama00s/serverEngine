package routing

import (
	"regexp"

	routingModel "github.com/Alquama00s/serverEngine/lib/routing/model"
)

type SimpleReqMiddleWare struct {
	regExString      string
	regex            *regexp.Regexp
	RequestProcessor func(req *routingModel.Request) (*routingModel.Request, error, *routingModel.Response)
	Priority         int
}

func (s *SimpleReqMiddleWare) Process(fun func(req *routingModel.Request) (*routingModel.Request, error, *routingModel.Response)) {
	s.RequestProcessor = fun
}

func (s *SimpleReqMiddleWare) GetProcessor() func(req *routingModel.Request) (*routingModel.Request, error, *routingModel.Response) {
	return s.RequestProcessor
}

func (s *SimpleReqMiddleWare) GetRegex() *regexp.Regexp {
	if s.regex == nil {
		s.regex = regexp.MustCompile(s.regExString)
	}
	return s.regex
}

func (s *SimpleReqMiddleWare) SetRegex(regex string) {
	s.regExString = regex
	s.regex = regexp.MustCompile(s.regExString)
}

func (s *SimpleReqMiddleWare) GetRegexString() string {
	return s.regExString
}

func (s *SimpleReqMiddleWare) GetPriority() int {
	return s.Priority
}
func (s *SimpleReqMiddleWare) SetPriority(priority int) {
	s.Priority = priority
}
