/*
Copyright 2021 zhangfuli.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"HFOperator/controllers/ca"
	"HFOperator/controllers/chaincode"
	"HFOperator/controllers/ordnode"
	"HFOperator/controllers/ordservice"
	"HFOperator/controllers/peer"
	"HFOperator/controllers/utils"
	"flag"
	"os"
	"path/filepath"

	hlfv1alpha1 "HFOperator/api/hlf.zhangfuli.com/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = hlfv1alpha1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	ctrl.SetLogger(zap.New(func(o *zap.Options) {
		o.Development = true
	}))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		LeaderElection:     enableLeaderElection,
		Port:               9443,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	caChartPath, err := filepath.Abs("./charts/hlf-ca")
	if err != nil {
		setupLog.Error(err, "unable to find the ca chart")
		os.Exit(1)
	}
	clientSet, err := utils.GetClientKubeWithConf(mgr.GetConfig())
	if err != nil {
		setupLog.Error(err, "unable to create client set", "controller", "FabricPeer")
		os.Exit(1)
	}
	if err = (&ca.FabricCAReconciler{
		Client:    mgr.GetClient(),
		Log:       ctrl.Log.WithName("controllers").WithName("FabricCA"),
		Scheme:    mgr.GetScheme(),
		Config:    mgr.GetConfig(),
		ClientSet: clientSet,
		ChartPath: caChartPath,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "FabricCA")
		os.Exit(1)
	}

	peerChartPath, err := filepath.Abs("./charts/hlf-peer")
	if err != nil {
		setupLog.Error(err, "unable to find the peer chart")
		os.Exit(1)
	}
	if err = (&peer.FabricPeerReconciler{
		Client:    mgr.GetClient(),
		Log:       ctrl.Log.WithName("controllers").WithName("FabricPeer"),
		Scheme:    mgr.GetScheme(),
		Config:    mgr.GetConfig(),
		ChartPath: peerChartPath,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "FabricPeer")
		os.Exit(1)
	}

	ordNodeChartPath, err := filepath.Abs("./charts/hlf-orderer")
	if err != nil {
		setupLog.Error(err, "unable to find the orderer node chart")
		os.Exit(1)
	}
	if err = (&ordnode.FabricOrdererNodeReconciler{
		Client:    mgr.GetClient(),
		Log:       ctrl.Log.WithName("controllers").WithName("FabricOrdererNode"),
		Scheme:    mgr.GetScheme(),
		Config:    mgr.GetConfig(),
		ChartPath: ordNodeChartPath,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "FabricOrdererNode")
		os.Exit(1)
	}

	ordServiceChartPath, err := filepath.Abs("./charts/hlf-ord")
	if err != nil {
		setupLog.Error(err, "unable to find the orderer chart")
		os.Exit(1)
	}
	if err = (&ordservice.FabricOrderingServiceReconciler{
		Client:    mgr.GetClient(),
		Log:       ctrl.Log.WithName("controllers").WithName("FabricOrderingService"),
		Scheme:    mgr.GetScheme(),
		Config:    mgr.GetConfig(),
		ChartPath: ordServiceChartPath,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "FabricOrderingService")
		os.Exit(1)
	}

	if err = (&chaincode.FabricChaincodeReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("FabricChaincode"),
		Scheme: mgr.GetScheme(),
		Config: mgr.GetConfig(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "FabricChaincode")
		os.Exit(1)
	}

	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
