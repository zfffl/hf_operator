package cmd

import (
	"HFOperator/command/cmd/ca"
	"HFOperator/command/cmd/chaincode"
	"HFOperator/command/cmd/channel"
	"HFOperator/command/cmd/externalchaincode"
	"HFOperator/command/cmd/inspect"
	"HFOperator/command/cmd/orderer"
	"HFOperator/command/cmd/peer"
	"HFOperator/command/cmd/utils"

	"github.com/spf13/cobra"
	// Workaround for authentication plugins https://krew.sigs.k8s.io/docs/developer-guide/develop/best-practices/#auth-plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

const (
	hlfDesc = `
kubectl plugin to manage HLF operator CRDs.`
)

// NewCmdHLF creates a new root command for command
func NewCmdHLF() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "hlf",
		Short:        "manage HLF operator CRDs",
		Long:         hlfDesc,
		SilenceUsage: true,
	}

	cmd.AddCommand(ca.NewCACmd(cmd.OutOrStdout(), cmd.ErrOrStderr()))
	cmd.AddCommand(peer.NewPeerCmd(cmd.OutOrStdout(), cmd.ErrOrStderr()))
	cmd.AddCommand(orderer.NewOrdNodeCmd(cmd.OutOrStdout(), cmd.ErrOrStderr()))
	cmd.AddCommand(inspect.NewInspectHLFConfig(cmd.OutOrStdout()))
	cmd.AddCommand(utils.NewUtilsCMD(cmd.OutOrStdout(), cmd.ErrOrStderr()))
	cmd.AddCommand(chaincode.NewChaincodeCmd(cmd.OutOrStdout(), cmd.ErrOrStderr()))
	cmd.AddCommand(channel.NewChannelCmd(cmd.OutOrStdout(), cmd.ErrOrStderr()))
	cmd.AddCommand(externalchaincode.NewExternalChaincodeCmd(cmd.OutOrStdout(), cmd.ErrOrStderr()))
	//cmd.AddCommand(org.NewOrgCmd(cmd.OutOrStdout(), cmd.ErrOrStderr()))
	//cmd.AddCommand(consortium.NewConsortiumCmd(cmd.OutOrStdout(), cmd.ErrOrStderr()))
	//cmd.AddCommand(ordservice.NewOrdServiceCmd(cmd.OutOrStdout(), cmd.ErrOrStderr()))
	return cmd
}
