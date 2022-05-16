package channel

import (
	"HFOperator/command/cmd/helpers"
	"bytes"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric/common/tools/protolator"

	"github.com/spf13/cobra"
	"io"
)

type inspectChannelCmd struct {
	configPath  string
	peer        string
	channelName string
	userName    string
}

func (c *inspectChannelCmd) validate() error {
	return nil
}
func (c *inspectChannelCmd) run(out io.Writer) error {
	oclient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	clientSet, err := helpers.GetKubeClient()
	if err != nil {
		return err
	}
	peer, err := helpers.GetPeerByFullName(clientSet, oclient, c.peer)
	if err != nil {
		return err
	}
	mspID := peer.Spec.MspID
	configBackend := config.FromFile(c.configPath)
	sdk, err := fabsdk.New(configBackend)
	if err != nil {
		return err
	}
	org1AdminClientContext := sdk.Context(
		fabsdk.WithUser(c.userName),
		fabsdk.WithOrg(mspID),
	)
	resClient, err := resmgmt.New(org1AdminClientContext)
	if err != nil {
		return err
	}
	block, err := resClient.QueryConfigBlockFromOrderer(
		c.channelName,
	)
	if err != nil {
		return err
	}
	cmnConfig, err := resource.ExtractConfigFromBlock(block)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	err = protolator.DeepMarshalJSON(&buf, cmnConfig)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(out, buf.String())
	if err != nil {
		return err
	}
	return nil
}
func newInspectChannelCMD(out io.Writer, errOut io.Writer) *cobra.Command {
	c := &inspectChannelCmd{}
	cmd := &cobra.Command{
		Use: "inspect",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run(out)
		},
	}
	persistentFlags := cmd.PersistentFlags()
	persistentFlags.StringVarP(&c.peer, "peer", "p", "", "Peer org to invoke the updates")
	persistentFlags.StringVarP(&c.userName, "user", "u", "", "User")
	persistentFlags.StringVarP(&c.channelName, "channel", "c", "", "Channel name")
	persistentFlags.StringVarP(&c.configPath, "config", "", "", "Configuration file for the SDK")
	cmd.MarkPersistentFlagRequired("channel")
	cmd.MarkPersistentFlagRequired("user")
	cmd.MarkPersistentFlagRequired("peer")
	cmd.MarkPersistentFlagRequired("config")
	return cmd
}
