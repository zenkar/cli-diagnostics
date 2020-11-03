// Copyright 2020. Akamai Technologies, Inc

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// errorStatsCmd represents the errorStats command
var errorStatsCmd = &cobra.Command{
	Use:   estatsUse,
	Args:  cobra.ExactArgs(1),
	Short: estatsShortDescription,
	Long:  estatsLongDescription,
	Run: func(cmd *cobra.Command, args []string) {
		var url string
		if _, err := strconv.Atoi(args[0]); err == nil {
			url = fmt.Sprintf("/diagnostic-tools/v2/cpcodes/%s/estats", args[0])
		} else if checkAbsoluteURL(args[0]) {
			url = "/diagnostic-tools/v2/estats?url=" + args[0]
		} else {
			printWarning("URL or CP code is invalid, e.g., http://www.example.com or 123456")
			os.Exit(1)
		}
		resp, byt := doHTTPRequest("GET", url, nil)

		if resp.StatusCode == 200 {
			var responseStruct Wrapper
			var responseStructJson EstatsJson
			err := json.Unmarshal(*byt, &responseStruct)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if jsonString {
				responseStructJson.UrlorCpCode = args[0]
				responseStructJson.ReportedTime = getReportedTime()
				responseStructJson.Estats = responseStruct.Estats
				resJson, _ := json.MarshalIndent(responseStructJson, "", "  ")
				resJson = getDecodedResponse(resJson)
				fmt.Println(string(resJson))
				return
			}

			printErrorStats(responseStruct.Estats)
		} else {
			printResponseError(byt)
		}
	},
}

func init() {
	rootCmd.AddCommand(errorStatsCmd)

}
