package chaincode

import (
	"HFOperator/command/cmd/helpers"
	"encoding/json"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
)

type queryInstalledCmd struct {
	configPath string
	peer       string
	userName   string
}

func (c *queryInstalledCmd) validate() error {
	return nil
}
func (c *queryInstalledCmd) run(out io.Writer) error {
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
	peerName := peer.Name
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

	chaincodes, err := resClient.LifecycleQueryInstalledCC(
		resmgmt.WithTargetEndpoints(peerName),
	)
	if err != nil {
		return err
	}
	if len(chaincodes) == 0 {
		log.Infof("No chaincodes installed")
		return nil
	}
	var data [][]string
	for _, chaincode := range chaincodes {
		referencesJson, err := json.Marshal(chaincode.References)
		if err != nil {
			return err
		}
		data = append(data, []string{
			chaincode.PackageID, chaincode.Label, string(referencesJson),
		})
	}
	table := tablewriter.NewWriter(out)
	table.SetHeader([]string{"Package ID", "Label", "References"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)
	table.AppendBulk(data) // Add Bulk Data
	table.Render()
	return nil
}
func newChaincodeQueryInstalledCMD(out io.Writer, errOut io.Writer) *cobra.Command {
	c := &queryInstalledCmd{}
	cmd := &cobra.Command{
		Use: "queryinstalled",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run(out)
		},
	}
	persistentFlags := cmd.PersistentFlags()
	persistentFlags.StringVarP(&c.peer, "peer", "p", "", "Peer org to invoke the updates")
	persistentFlags.StringVarP(&c.userName, "user", "", "", "User")
	persistentFlags.StringVarP(&c.configPath, "config", "", "", "Configuration file for the SDK")
	cmd.MarkPersistentFlagRequired("user")
	cmd.MarkPersistentFlagRequired("peer")
	cmd.MarkPersistentFlagRequired("config")
	return cmd
}
