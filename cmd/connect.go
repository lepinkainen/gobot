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
	"github.com/lepinkainen/gobot/botcore"
	"github.com/spf13/cobra"
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

	nick, _ := cmd.Flags().GetString("nick")
	server, _ := cmd.Flags().GetString("server")
	verboseFlag, _ := cmd.Flags().GetBool("verbose")
	debugFlag, _ := cmd.Flags().GetBool("debug")
	tlsFlag, _ := cmd.Flags().GetBool("tls")

	config := botcore.IRCConfig{
		Nick:     nick,
		Server:   server,
		Channels: []string{"#pyfibot.test2"},
		Verbose:  verboseFlag,
		Debug:    debugFlag,
		TLS:      tlsFlag,
	}
	botcore.Connect(config)
}

func init() {
	rootCmd.AddCommand(connectCmd)

	// TODO channels, multiple flags? config file?
	connectCmd.Flags().String("nick", "", "Nick to use when connecting")
	connectCmd.Flags().String("server", "", "Server:port combination for connecting")
	connectCmd.Flags().BoolP("verbose", "v", false, "Enable verbose callbacks")
	connectCmd.Flags().BoolP("debug", "d", false, "Enable protocol debugging")
	connectCmd.Flags().BoolP("tls", "", false, "Use TLS to connect")
}
