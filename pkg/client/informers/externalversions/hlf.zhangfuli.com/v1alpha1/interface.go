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
// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	internalinterfaces "HFOperator/pkg/client/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// FabricCAs returns a FabricCAInformer.
	FabricCAs() FabricCAInformer
	// FabricChaincodes returns a FabricChaincodeInformer.
	FabricChaincodes() FabricChaincodeInformer
	// FabricOrdererNodes returns a FabricOrdererNodeInformer.
	FabricOrdererNodes() FabricOrdererNodeInformer
	// FabricOrderingServices returns a FabricOrderingServiceInformer.
	FabricOrderingServices() FabricOrderingServiceInformer
	// FabricPeers returns a FabricPeerInformer.
	FabricPeers() FabricPeerInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// FabricCAs returns a FabricCAInformer.
func (v *version) FabricCAs() FabricCAInformer {
	return &fabricCAInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// FabricChaincodes returns a FabricChaincodeInformer.
func (v *version) FabricChaincodes() FabricChaincodeInformer {
	return &fabricChaincodeInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// FabricOrdererNodes returns a FabricOrdererNodeInformer.
func (v *version) FabricOrdererNodes() FabricOrdererNodeInformer {
	return &fabricOrdererNodeInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// FabricOrderingServices returns a FabricOrderingServiceInformer.
func (v *version) FabricOrderingServices() FabricOrderingServiceInformer {
	return &fabricOrderingServiceInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// FabricPeers returns a FabricPeerInformer.
func (v *version) FabricPeers() FabricPeerInformer {
	return &fabricPeerInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}
