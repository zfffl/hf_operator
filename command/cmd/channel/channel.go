package channel

import (
	"HFOperator/command/cmd/channel/consenter"
	"HFOperator/command/cmd/channel/ordorg"
	"io"

	"github.com/spf13/cobra"
)

func NewChannelCmd(stdOut io.Writer, stdErr io.Writer) *cobra.Command {
	channelCmd := &cobra.Command{
		Use: "channel",
	}
	channelCmd.AddCommand(
		newCreateChannelCMD(stdOut, stdErr),
		newJoinChannelCMD(stdOut, stdErr),
		newSignUpdateChannelCMD(stdOut, stdErr),
		newAddAnchorPeerCMD(stdOut, stdErr),
		newUpdateChannelCMD(stdOut, stdErr),
		newGenerateChannelCMD(stdOut, stdErr),
		newInspectChannelCMD(stdOut, stdErr),
		newTopChannelCMD(stdOut, stdErr),
		newAddOrgToChannelCMD(stdOut, stdErr),
		ordorg.NewOrdOrgCmd(stdOut, stdErr),
		consenter.NewConsenterCmd(stdOut, stdErr),
		newDelAnchorPeerCMD(stdOut, stdErr),
	)
	return channelCmd
}
