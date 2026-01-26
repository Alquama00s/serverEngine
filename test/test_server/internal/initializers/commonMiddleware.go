package initializers

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Alquama00s/serverEngine"
	routingModel "github.com/Alquama00s/serverEngine/lib/routing/model"
)

type CommonMiddleware struct {
}

func (c *CommonMiddleware) Init() {
	//add required response headers
	serverEngine.Registrar().
		RegisterSimpleResProcessor("/*",
			func(res *routingModel.Response) (*routingModel.Response, error) {
				res.Logger.
					Debug().Str("pkg", "controller.GlobalConfigs").
					Msg("headers response processors")

				if res.Headers == nil {
					res.Headers = &http.Header{}
				}
				for name, values := range res.Request.RawRequest.Header {
					for _, value := range values {
						fmt.Printf("%s: %s\n", name, value)
					}
				}
				res.Headers.Add("Access-Control-Allow-Origin", "*")
				res.Headers.Add("Access-Control-Allow-Methods", "*")
				res.Headers.Add("Access-Control-Allow-Headers", "Content-Type")
				res.Headers.Add("Access-Control-Allow-Headers", "Authorization")
				res.Headers.Add("Content-type", "application/json")
				return res, nil
			})

	// //context injection for gorm
	// serverEngine.Registrar().
	// 	RegisterSimpleReqProcessor("/*", func(req *lib.Request) (*lib.Request, error, *lib.Response) {
	// 		if req.RequestPrincipal != nil {
	// 			req.
	// 		}
	// 		return req, nil, nil
	// 	})

	//http level pre processor to validate the request
	serverEngine.Registrar().
		RegisterSimpleReqProcessor("/*",
			func(req *routingModel.Request) (*routingModel.Request, error, *routingModel.Response) {
				req.Logger.
					Debug().Str("pkg", "controller.GlobalConfigs").
					Msg("General Request Validator")
				if req.RawRequest.Method == http.MethodPost {
					s, err := io.ReadAll(req.RawRequest.Body)
					if err != nil {
						return nil, err, nil
					}
					if len(s) == 0 {
						return nil, errors.New("Request body can't be null"), nil
					}
					req.RawBody = s
				}
				return req, nil, nil
			})

	// //performs error handling on all routes
	// serverEngine.Registrar().ErrorHandler(func(err error, req *lib.Request, res *lib.Response) error {
	// 	req.Logger.
	// 		Debug().Str("pkg", "controller.GlobalConfigs").Msg("handling error")
	// 	switch err.(type) {
	// 	case validator.FieldError:
	// 		return err
	// 	case validator.ValidationErrors:
	// 		verr := err.(validator.ValidationErrors)[0]
	// 		return fmt.Errorf("field %s: wanted %s %s, got `%s`", verr.Field(), verr.Tag(), verr.Param(), verr.Value())
	// 	case *lib.ErrorResponse:
	// 		return err
	// 	default:
	// 		req.Logger.Err(err).Msg("handling error")
	// 		return fmt.Errorf("unknown error occured")
	// 	}
	// })
}
