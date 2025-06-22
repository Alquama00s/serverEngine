package lib

import "regexp"

type ResponseProcessor interface {
	GetProcessor() func(req *Response) (*Response, error)
	Process(fun func(req *Response) (*Response, error))
	GetRegex() *regexp.Regexp
	SetRegex(regex string)
	GetRegexString() string
	GetPriority() int
	SetPriority(priority int)
}
