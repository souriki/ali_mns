package ali_mns

import (
	"time"

	"github.com/gogap/errors"
	"github.com/valyala/fasthttp"
)

type optionType string

const (
	clientOption  optionType = "clientOption"
	requestOption optionType = "requestOption"
)

const (
	optTimeout    = "Timeout"
	optReqTimeout = "ReqTimeout"
)

type optionValue struct {
	value interface{}
	typ   optionType
}

type optionParams map[string]optionValue

// Option http option
type Option func(optionParams) error

// Timeout ...
func Timeout(timeoutInSec int64) Option {
	return func(params optionParams) error {
		params[optTimeout] = optionValue{
			value: timeoutInSec,
			typ:   clientOption,
		}
		return nil
	}
}

// RequestTimeout ...
func RequestTimeout(d time.Duration) Option {
	return func(params optionParams) error {
		params[optReqTimeout] = optionValue{
			value: d,
			typ:   requestOption,
		}
		return nil
	}
}

func initMNSClientOption(cli *aliMNSClient, opts ...Option) error {
	params := optionParams{}
	for _, opt := range opts {
		err := opt(params)
		if err != nil {
			return err
		}
	}
	if optValue, ok := params[optTimeout]; ok && optValue.typ == clientOption {
		cli.Timeout = optValue.value.(int64)
	}
	return nil
}

func doRequestWithOption(p *aliMNSClient, req *fasthttp.Request, resp *fasthttp.Response, opts ...Option) error {
	params := optionParams{}
	for _, opt := range opts {
		err := opt(params)
		if err != nil {
			return err
		}
	}
	doRequest := func() error {
		return p.client.Do(req, resp)
	}
	if optValue, ok := params[optReqTimeout]; ok && optValue.typ == requestOption {
		doRequest = func() error {
			return p.client.DoTimeout(req, resp, optValue.value.(time.Duration))
		}
	}

	if err := doRequest(); err != nil {
		return ERR_SEND_REQUEST_FAILED.New(errors.Params{"err": err})
	}
	return nil
}
