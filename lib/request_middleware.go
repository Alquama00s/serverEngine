package lib

import "regexp"

type RequestProcessor interface {
	GetProcessor() func(req *Request) (*Request, error, *Response)
	Process(fun func(req *Request) (*Request, error, *Response))
	GetRegex() *regexp.Regexp
	SetRegex(regex string)
	GetRegexString() string
	GetPriority() int
	SetPriority(priority int)
}
