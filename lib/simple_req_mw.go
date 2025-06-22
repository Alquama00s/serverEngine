package lib

import "regexp"

type SimpleReqMiddleWare struct {
	regExString      string
	regex            *regexp.Regexp
	RequestProcessor func(req *Request) (*Request, error, *Response)
	Priority         int
}

func (s *SimpleReqMiddleWare) Process(fun func(req *Request) (*Request, error, *Response)) {
	s.RequestProcessor = fun
}

func (s *SimpleReqMiddleWare) GetProcessor() func(req *Request) (*Request, error, *Response) {
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
