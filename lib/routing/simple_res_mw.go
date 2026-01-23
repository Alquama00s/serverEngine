package routing

import (
	"regexp"

	routingModel "github.com/Alquama00s/serverEngine/lib/routing/model"
)

type SimpleResMiddleWare struct {
	regExString       string
	regex             *regexp.Regexp
	ResponseProcessor func(req *routingModel.Response) (*routingModel.Response, error)
	Priority          int
}

func (s *SimpleResMiddleWare) Process(fun func(req *routingModel.Response) (*routingModel.Response, error)) {
	s.ResponseProcessor = fun
}

func (s *SimpleResMiddleWare) GetProcessor() func(req *routingModel.Response) (*routingModel.Response, error) {
	return s.ResponseProcessor
}

func (s *SimpleResMiddleWare) GetRegex() *regexp.Regexp {
	if s.regex == nil {
		s.regex = regexp.MustCompile(s.regExString)
	}
	return s.regex
}
func (s *SimpleResMiddleWare) SetRegex(regex string) {
	s.regExString = regex
	s.regex = regexp.MustCompile(s.regExString)
}
func (s *SimpleResMiddleWare) GetRegexString() string {
	return s.regExString
}

func (s *SimpleResMiddleWare) GetPriority() int {
	return s.Priority
}
func (s *SimpleResMiddleWare) SetPriority(priority int) {
	s.Priority = priority
}
