package cb

import (
	"github.com/afex/hystrix-go/hystrix"
)

func StartHystrix(timeoutInSecond int, serviceName string) error {


	setting := hystrix.CommandConfig{
		Timeout: 1000 * timeoutInSecond,
	}

	hystrix.ConfigureCommand(serviceName, setting)
	return nil
}
