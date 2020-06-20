/*
Copyright © 2020 Riku Lindblad <riku.lindblad@iki.fi>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/lepinkainen/gobot/botcore"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to an IRC network",
	Long: `Connects to a single network and starts feeding commands to 
the command server`,
	Run: command,
}

// Connect to the defined network
func command(cmd *cobra.Command, args []string) {

	nick := viper.GetString("nick")
	server := viper.GetString("server")
	verboseFlag := viper.GetBool("verbose")
	debugFlag := viper.GetBool("debug")
	tlsFlag := viper.GetBool("tls")
	channels := viper.GetStringSlice("channels")

	var channelList []string

	for _, channel := range channels {
		if strings.HasPrefix(channel, "#") {
			channelList = append(channelList, channel)
		} else {
			channelList = append(channelList, fmt.Sprintf("#%s", channel))
		}
	}

	config := botcore.IRCConfig{
		Nick:     nick,
		Server:   server,
		Channels: channelList,
		Verbose:  verboseFlag,
		Debug:    debugFlag,
		TLS:      tlsFlag,
	}
	botcore.Connect(config)
}

func init() {
	rootCmd.AddCommand(connectCmd)

	connectCmd.Flags().String("nick", "", "Nick to use when connecting")
	connectCmd.Flags().String("server", "", "Server:port combination for connecting")
	connectCmd.Flags().BoolP("verbose", "v", false, "Enable verbose callbacks")
	connectCmd.Flags().BoolP("debug", "d", false, "Enable protocol debugging")
	connectCmd.Flags().BoolP("tls", "", false, "Use TLS to connect")
}
