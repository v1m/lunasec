// Copyright 2021 by LunaSec (owned by Refinery Labs, Inc)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package command

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"lunasec/lunatrace/pkg/constants"
	"lunasec/lunatrace/pkg/util"
	"os"
	"strings"
)

func EnableGlobalFlags(globalBoolFlags map[string]bool) {
	verbose := globalBoolFlags["verbose"]
	debug := globalBoolFlags["debug"]
	jsonFlag := globalBoolFlags["json"]
	ignoreWarnings := globalBoolFlags["ignore-warnings"]
	logToStderr := globalBoolFlags["log-to-stderr"]

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if verbose || debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if ignoreWarnings {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	}

	if debug {
		// include file and line number when logging
		log.Logger = log.With().Caller().Logger()
	}

	outStream := os.Stdout
	if logToStderr {
		outStream = os.Stderr
	}

	if !jsonFlag {
		// pretty print output to the console if we are not interested in parsable output
		consoleOutput := zerolog.ConsoleWriter{Out: outStream}
		consoleOutput.FormatFieldName = func(i interface{}) string {
			return fmt.Sprintf("\n\t%s: ", util.Colorize(constants.ColorBlue, i))
		}
		consoleOutput.FormatFieldValue = func(i interface{}) string {
			switch t := i.(type) {
			case string:
				return fmt.Sprintf("%s", util.FixStringSlashes(t))
			case []string:
				return fmt.Sprintf("[%s]", strings.Join(util.FixStringSliceSlashes(t), ", "))
			default:
				return fmt.Sprintf("%s", i)
			}
		}

		consoleOutput.FormatLevel = func(i interface{}) string {
			if i == nil {
				return util.Colorize(constants.ColorBold, "Scan Result:")
			}

			level := i.(string)

			var formattedLevel string
			switch level {
			case "warn":
				formattedLevel = util.Colorize(constants.ColorYellow, level)
			case "error":
				formattedLevel = util.Colorize(constants.ColorRed, level)
			case "info":
				formattedLevel = util.Colorize(constants.ColorBlue, level)
			case "debug":
				formattedLevel = util.Colorize(constants.ColorGreen, level)
			default:
				formattedLevel = util.Colorize(constants.ColorWhite, level)
			}
			return fmt.Sprintf("| %s |", formattedLevel)
		}

		log.Logger = log.Output(consoleOutput)
	}
}
