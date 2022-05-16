package inspect

import (
	"HFOperator/command/cmd/helpers"
	"HFOperator/controllers/utils"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"sigs.k8s.io/yaml"
	"text/template"

	"github.com/Masterminds/sprig/v3"

	"github.com/spf13/cobra"
)

const (
	createDesc = `
'inspect' command creates creates a configuration file ready to use for the go sdk`
	createExample = `  kubectl hlf inspect --output hlf-cfg.yaml`
	yamlFormat    = "yaml"
	jsonFormat    = "json"
)

type inspectCmd struct {
	fileOutput    string
	organizations []string
	internal      bool
	format        string
}

func (c *inspectCmd) validate() error {
	return nil
}

type OrderingService struct {
	MSPID        string
	OrderingNode string
}
type OrderingNode struct {
	TLSCert string
	URL     string
}
type Peer struct {
	URL     string
	TLSCert string
}

const tmplGoConfig = `
name: hlf-network
version: 1.0.0
client:
  organization: "{{ .Organization }}"
{{- if not .Organizations }}
organizations: {}
{{- else }}
organizations:
  {{ range $mspID, $org := .Organizations }}
  {{$mspID}}:
    mspid: {{$mspID}}
    cryptoPath: /tmp/cryptopath
    users: {}
{{- if not $org.Peers }}
    peers: []
{{- else }}
    peers:
      {{- range $peer := $org.Peers }}
      - {{ $peer.Name }}
 	  {{- end }}
{{- end }}
{{- if not $org.OrderingServices }}
    orderers: []
{{- else }}
    orderers:
      {{- range $ordService := $org.OrderingServices }}
      {{- range $orderer := $ordService.Orderers }}
      - {{ $orderer.Name }}
 	  {{- end }}
 	  {{- end }}

    {{- end }}
{{- end }}
{{- end }}

{{- if not .Orderers }}
orderers: []
{{- else }}
orderers:
{{- range $ordService := .Orderers }}
{{- range $orderer := $ordService.Orderers }}
  {{$orderer.Name}}:
{{if $.Internal }}
    url: grpcs://{{ $orderer.PrivateURL }}
{{ else }}
    url: grpcs://{{ $orderer.PublicURL }}
{{ end }}
    grpcOptions:
      allow-insecure: false
    tlsCACerts:
      pem: |
{{ or $orderer.Status.TlsCACert $orderer.Status.TlsCert | indent 8 }}
{{- end }}
{{- end }}
{{- end }}

{{- if not .Peers }}
peers: []
{{- else }}
peers:
  {{- range $peer := .Peers }}
  {{$peer.Name}}:
{{if $.Internal }}
    url: grpcs://{{ $peer.PrivateURL }}
{{ else }}
    url: grpcs://{{ $peer.PublicURL }}
{{ end }}
    grpcOptions:
      hostnameOverride: ""
      ssl-target-name-override: ""
      allow-insecure: false
    tlsCACerts:
      pem: |
{{ $peer.Status.TlsCACert | indent 8 }}
{{- end }}
{{- end }}

{{- if not .CertAuths }}
certificateAuthorities: []
{{- else }}
certificateAuthorities:
{{- range $ca := .CertAuths }}
  
  {{ $ca.Name }}:
{{if $.Internal }}
    url: https://{{ $ca.PrivateURL }}
{{ else }}
    url: https://{{ $ca.PublicURL }}
{{ end }}
{{if $ca.EnrollID }}
    registrar:
        enrollId: {{ $ca.EnrollID }}
        enrollSecret: {{ $ca.EnrollPWD }}
{{ end }}
    caName: ca
    tlsCACerts:
      pem: 
       - |
{{ $ca.Status.TlsCert | indent 12 }}

{{- end }}
{{- end }}

channels:
  _default:
{{- if not .Orderers }}
    orderers: []
{{- else }}
    orderers:
{{- range $ordService := .Orderers }}
{{- range $orderer := $ordService.Orderers }}
      - {{$orderer.Name}}
{{- end }}
{{- end }}
{{- end }}
{{- if not .Peers }}
    peers: {}
{{- else }}
    peers:
{{- range $peer := .Peers }}
       {{$peer.Name}}:
        discover: true
        endorsingPeer: true
        chaincodeQuery: true
        ledgerQuery: true
        eventSource: true
{{- end }}
{{- end }}

`

func (c *inspectCmd) run(out io.Writer) error {
	oclient, err := helpers.GetKubeOperatorClient()
	if err != nil {
		return err
	}
	clientSet, err := helpers.GetKubeClient()
	if err != nil {
		return err
	}
	ns := ""
	certAuths, err := helpers.GetClusterCAs(clientSet, oclient, ns)
	if err != nil {
		return err
	}
	ordOrgs, orderers, err := helpers.GetClusterOrderers(clientSet, oclient, ns)
	if err != nil {
		return err
	}
	peerOrgs, clusterPeers, err := helpers.GetClusterPeers(clientSet, oclient, ns)
	if err != nil {
		return err
	}
	filterByOrgs := len(c.organizations) > 0
	orgMap := map[string]*helpers.Organization{}
	for _, v := range ordOrgs {
		if filterByOrgs && utils.Contains(c.organizations, v.MspID) {
			orgMap[v.MspID] = v
		}
	}
	for _, v := range peerOrgs {
		if filterByOrgs && utils.Contains(c.organizations, v.MspID) {
			orgMap[v.MspID] = v
		}
	}
	var peers []*helpers.ClusterPeer
	for _, peer := range clusterPeers {
		if filterByOrgs && utils.Contains(c.organizations, peer.MSPID) {
			peers = append(peers, peer)
		}
	}
	tmpl, err := template.New("test").Funcs(sprig.HermeticTxtFuncMap()).Parse(tmplGoConfig)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"Peers":         peers,
		"Orderers":      orderers,
		"Organizations": orgMap,
		"CertAuths":     certAuths,
		"Internal":      c.internal,
	})
	if err != nil {
		return err
	}

	var data []byte
	if c.format != yamlFormat && c.format != jsonFormat {
		fmt.Fprint(out, "Invalid output format... Default to yaml")
		c.format = yamlFormat
	}

	if c.format == jsonFormat {
		data, err = yaml.YAMLToJSON(buf.Bytes())
		if err != nil {
			return err
		}
	} else {
		data = buf.Bytes()
	}

	if c.fileOutput != "" {
		err = ioutil.WriteFile(c.fileOutput, data, 0644)
		if err != nil {
			return err
		}
	} else {
		_, err = fmt.Fprint(out, string(data))
		if err != nil {
			return err
		}
	}

	return nil
}

func NewInspectHLFConfig(out io.Writer) *cobra.Command {
	c := &inspectCmd{}
	cmd := &cobra.Command{
		Use:     "inspect",
		Short:   "Inspect the components deployed",
		Long:    createDesc,
		Example: createExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.validate(); err != nil {
				return err
			}
			return c.run(out)
		},
	}

	f := cmd.Flags()
	f.StringVar(&c.fileOutput, "output", "", "output file")
	f.BoolVar(&c.internal, "internal", false, "use kubernetes service names")
	f.StringArrayVarP(&c.organizations, "organizations", "o", []string{}, "organizations to export")
	f.StringVar(&c.format, "format", yamlFormat, "connection profile output format (yaml/json)")

	return cmd
}
