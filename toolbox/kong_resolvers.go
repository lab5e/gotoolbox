package toolbox

import (
	"os"
	"strings"

	"github.com/alecthomas/kong"
)

// EnvVarResolver is an environment variable resolver for Kong -- ie
// it will default to the environment variable for the parameters. The parameters
// are named similarly to the parameter names but in upper case and with underscores.
// ie "--some-parameter-name" will use the environment variable "SOME_PARAMETER_NAME"
func EnvVarResolver() kong.ResolverFunc {
	return func(context *kong.Context, parent *kong.Path, flag *kong.Flag) (interface{}, error) {
		envName := strings.ToUpper(strings.Replace(flag.Name, "-", "_", -1))
		val := os.Getenv(envName)
		if val != "" {
			return val, nil
		}
		return nil, nil
	}
}
