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

package chaincode

import (
	"HFOperator/controllers/certs"
	"HFOperator/controllers/utils"
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/operator-framework/operator-lib/status"
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	hlfv1alpha1 "HFOperator/api/hlf.zhangfuli.com/v1alpha1"
)

// FabricChaincodeReconciler reconciles a FabricChaincode object
type FabricChaincodeReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
	Config *rest.Config
}

const chaincodeFinalizer = "finalizer.chaincode.hlf.zhangfuli.com"

type SecretChaincodeData struct {
	Updated     bool
	Enabled     bool
	Certificate []byte
	PrivateKey  []byte
	RootCert    []byte
}

func CreateChaincodeCryptoMaterial(conf *hlfv1alpha1.FabricChaincode, caName string, caurl string, enrollID string, enrollSecret string, tlsCertString string, hosts []string) (*x509.Certificate, *ecdsa.PrivateKey, *x509.Certificate, error) {
	tlsCert, tlsKey, tlsRootCert, err := certs.EnrollUser(certs.EnrollUserRequest{
		TLSCert:    tlsCertString,
		URL:        caurl,
		Name:       caName,
		MSPID:      "Org1MSP",
		User:       enrollID,
		Secret:     enrollSecret,
		Hosts:      hosts,
		CN:         "",
		Profile:    "tls",
		Attributes: nil,
	})
	if err != nil {
		return nil, nil, nil, err
	}
	return tlsCert, tlsKey, tlsRootCert, nil
}
func (r *FabricChaincodeReconciler) getDeploymentName(fabricChaincode *hlfv1alpha1.FabricChaincode) string {
	return fmt.Sprintf("%s", fabricChaincode.Name)
}
func (r *FabricChaincodeReconciler) getServiceName(fabricChaincode *hlfv1alpha1.FabricChaincode) string {
	return fmt.Sprintf("%s", fabricChaincode.Name)
}
func (r *FabricChaincodeReconciler) getSecretName(fabricChaincode *hlfv1alpha1.FabricChaincode) string {
	return fmt.Sprintf("%s-certs", fabricChaincode.Name)
}

func (r *FabricChaincodeReconciler) finalizeChaincode(reqLogger logr.Logger, m *hlfv1alpha1.FabricChaincode) error {
	ns := m.Namespace
	if ns == "" {
		ns = "default"
	}
	//releaseName := m.Name
	reqLogger.Info("Successfully finalized chaincode")
	kubeClientset, err := kubernetes.NewForConfig(r.Config)
	if err != nil {
		return err
	}
	deploymentName := r.getDeploymentName(m)
	ctx := context.Background()
	err = kubeClientset.AppsV1().Deployments(ns).Delete(ctx, deploymentName, metav1.DeleteOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			reqLogger.Info(fmt.Sprintf("Deployment %s not found", deploymentName))
		} else {
			reqLogger.Error(err, "Failed to delete deployment")
			return err
		}
	}
	serviceName := r.getServiceName(m)
	err = kubeClientset.CoreV1().Services(ns).Delete(ctx, serviceName, metav1.DeleteOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			reqLogger.Info(fmt.Sprintf("Service %s not found", serviceName))
		} else {
			reqLogger.Error(err, "Failed to delete service")
			return err
		}
	}
	secretName := r.getSecretName(m)
	err = kubeClientset.CoreV1().Secrets(ns).Delete(ctx, secretName, metav1.DeleteOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			reqLogger.Info(fmt.Sprintf("Secret %s not found", secretName))
		} else {
			reqLogger.Error(err, "Failed to delete secret")
			return err
		}
	}
	return nil
}

func (r *FabricChaincodeReconciler) addFinalizer(reqLogger logr.Logger, m *hlfv1alpha1.FabricChaincode) error {
	reqLogger.Info("Adding Finalizer for the Chaincode")
	controllerutil.AddFinalizer(m, chaincodeFinalizer)

	// Update CR
	err := r.Update(context.TODO(), m)
	if err != nil {
		reqLogger.Error(err, "Failed to update Chaincode with finalizer")
		return err
	}
	return nil
}

const (
	CertificateSecretKey = "tls.crt"
	PrivateKeySecretKey  = "tls.key"
	RootCertSecretKey    = "tlsroot.crt"
)

func (r FabricChaincodeReconciler) getCryptoMaterial(ctx context.Context, labels map[string]string, ns string, fabricChaincode *hlfv1alpha1.FabricChaincode) (*SecretChaincodeData, error) {
	secretChaincodeData := &SecretChaincodeData{
		Enabled: true,
		Updated: false,
	}
	if fabricChaincode.Spec.Credentials == nil {
		secretChaincodeData.Enabled = false
		return secretChaincodeData, nil
	}
	secretName := r.getSecretName(fabricChaincode)
	tlsCAUrl := fmt.Sprintf("https://%s:%d", fabricChaincode.Spec.Credentials.Cahost, fabricChaincode.Spec.Credentials.Caport)

	kubeClientset, err := kubernetes.NewForConfig(r.Config)
	if err != nil {
		return nil, err
	}

	updateSecretData := false
	secret, err := kubeClientset.CoreV1().Secrets(ns).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		updateSecretData = true
	} else {
		x509Cert, err := utils.ParseX509Certificate(secret.Data[CertificateSecretKey])
		// renew certificates data if certificate is about to expire (7 days before expiration)
		if err != nil || x509Cert.NotAfter.Before(time.Now().Add(time.Hour*24*7)) {
			updateSecretData = true
		}
		if secret.Data[CertificateSecretKey] != nil &&
			len(secret.Data[CertificateSecretKey]) > 0 &&
			secret.Data[PrivateKeySecretKey] != nil &&
			len(secret.Data[PrivateKeySecretKey]) > 0 &&
			secret.Data[RootCertSecretKey] != nil &&
			len(secret.Data[RootCertSecretKey]) > 0 {
			updateSecretData = false
		} else {
			updateSecretData = true
		}
	}
	secretChaincodeData.Updated = updateSecretData
	if updateSecretData {
		cacert, err := base64.StdEncoding.DecodeString(fabricChaincode.Spec.Credentials.Catls.Cacert)
		if err != nil {
			return nil, err
		}
		tlsCert, tlsKey, tlsRootCert, err := CreateChaincodeCryptoMaterial(
			fabricChaincode,
			fabricChaincode.Spec.Credentials.Caname,
			tlsCAUrl,
			fabricChaincode.Spec.Credentials.Enrollid,
			fabricChaincode.Spec.Credentials.Enrollsecret,
			string(cacert),
			fabricChaincode.Spec.Credentials.Csr.Hosts,
		)
		if err != nil {
			err = errors.New("Failed to create chaincode crypto material")
			return nil, err
		}
		key, err := utils.EncodePrivateKey(tlsKey)
		if err != nil {
			return nil, err
		}
		secretChaincodeData.Certificate = utils.EncodeX509Certificate(tlsCert)
		secretChaincodeData.RootCert = utils.EncodeX509Certificate(tlsRootCert)
		secretChaincodeData.PrivateKey = key
	} else {
		secretChaincodeData.Certificate = secret.Data[CertificateSecretKey]
		secretChaincodeData.PrivateKey = secret.Data[PrivateKeySecretKey]
		secretChaincodeData.RootCert = secret.Data[RootCertSecretKey]
	}

	if err != nil {
		if apierrors.IsNotFound(err) {
			// creating secret
			secretData := map[string][]byte{
				"tls.crt":     secretChaincodeData.Certificate,
				"tlsroot.crt": secretChaincodeData.RootCert,
				"tls.key":     secretChaincodeData.PrivateKey,
			}
			secret, err = kubeClientset.CoreV1().Secrets(ns).Create(
				ctx,
				&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      secretName,
						Namespace: ns,
						Labels:    labels,
					},

					Data: secretData,
				},
				metav1.CreateOptions{},
			)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		secretData := map[string][]byte{
			"tls.crt":     secretChaincodeData.Certificate,
			"tlsroot.crt": secretChaincodeData.RootCert,
			"tls.key":     secretChaincodeData.PrivateKey,
		}
		secret.Data = secretData
		secret, err = kubeClientset.CoreV1().Secrets(ns).Update(
			ctx,
			secret,
			metav1.UpdateOptions{},
		)
		if err != nil {
			return nil, err
		}
	}
	return secretChaincodeData, nil
}

// +kubebuilder:rbac:groups=hlf.zhangfuli.com,resources=fabricchaincodes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=hlf.zhangfuli.com,resources=fabricchaincodes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=hlf.zhangfuli.com,resources=fabricchaincodes/finalizers,verbs=get;update;patch
func (r *FabricChaincodeReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	reqLogger := r.Log.WithValues("hlf", req.NamespacedName)
	fabricChaincode := &hlfv1alpha1.FabricChaincode{}
	//releaseName := req.Name

	err := r.Get(ctx, req.NamespacedName, fabricChaincode)
	if err != nil {
		log.Debugf("Error getting the object %s error=%v", req.NamespacedName, err)
		if apierrors.IsNotFound(err) {
			reqLogger.Info("Chaincode resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		reqLogger.Error(err, "Failed to get Chaincode.")
		return ctrl.Result{}, err
	}
	isChaincodeMarkedToBeDeleted := fabricChaincode.GetDeletionTimestamp() != nil
	if isChaincodeMarkedToBeDeleted {
		if utils.Contains(fabricChaincode.GetFinalizers(), chaincodeFinalizer) {
			if err := r.finalizeChaincode(reqLogger, fabricChaincode); err != nil {
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(fabricChaincode, chaincodeFinalizer)
			err := r.Update(ctx, fabricChaincode)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}
	if !utils.Contains(fabricChaincode.GetFinalizers(), chaincodeFinalizer) {
		if err := r.addFinalizer(reqLogger, fabricChaincode); err != nil {
			return ctrl.Result{}, err
		}
	}
	log.Infof("Chaincode %s reconciled", req.NamespacedName)
	ns := req.Namespace
	if ns == "" {
		ns = "default"
	}
	labels := map[string]string{
		"app":       "fabric-chaincode",
		"chaincode": fabricChaincode.Name,
	}
	kubeClientset, err := kubernetes.NewForConfig(r.Config)
	if err != nil {
		r.setConditionStatus(ctx, fabricChaincode, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricChaincode)
	}

	cryptoData, err := r.getCryptoMaterial(ctx, labels, ns, fabricChaincode)
	if err != nil {
		r.setConditionStatus(ctx, fabricChaincode, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricChaincode)
	}
	deploymentName := fmt.Sprintf("%s", fabricChaincode.Name)
	serviceName := fmt.Sprintf("%s", fabricChaincode.Name)

	chaincodeAddress := "0.0.0.0:7052"
	envVars := []corev1.EnvVar{
		{
			Name:  "CHAINCODE_ID",
			Value: fabricChaincode.Spec.PackageID,
		},
		{
			Name:  "CORE_CHAINCODE_ID",
			Value: fabricChaincode.Spec.PackageID,
		},
		{
			Name:  "CHAINCODE_SERVER_ADDRESS",
			Value: chaincodeAddress,
		},
		{
			Name:  "CORE_CHAINCODE_ADDRESS",
			Value: chaincodeAddress,
		},
	}
	var volumes []corev1.Volume
	secretName := r.getSecretName(fabricChaincode)
	var volumeMounts []corev1.VolumeMount
	if cryptoData.Enabled {
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      secretName,
			ReadOnly:  true,
			MountPath: "/config/certs",
		})
		volumes = append(volumes, corev1.Volume{
			Name: secretName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: secretName,
				},
			},
		})
		envVars = append(envVars, []corev1.EnvVar{
			{
				Name:  "CHAINCODE_TLS_DISABLED",
				Value: "false",
			},
			{
				Name:  "CORE_PEER_TLS_ENABLED",
				Value: "true",
			},
			{
				Name:  "CHAINCODE_TLS_KEY",
				Value: "/config/certs/tls.key",
			},
			{
				Name:  "CORE_CHAINCODE_TLS_KEY_FILE",
				Value: "/config/certs/tls.key",
			},
			{
				Name:  "CHAINCODE_TLS_CERT",
				Value: "/config/certs/tls.crt",
			},
			{
				Name:  "CORE_CHAINCODE_TLS_CERT_FILE",
				Value: "/config/certs/tls.crt",
			},
			{
				Name:  "CHAINCODE_CLIENT_CA_CERT",
				Value: "/config/certs/tlsroot.crt",
			},
			{
				Name:  "CORE_CHAINCODE_TLS_CLIENT_CACERT_FILE",
				Value: "/config/certs/tlsroot.crt",
			},
		}...)
	} else {
		envVars = append(envVars, []corev1.EnvVar{
			{
				Name:  "CHAINCODE_TLS_DISABLED",
				Value: "true",
			},
			{
				Name:  "CORE_PEER_TLS_ENABLED",
				Value: "false",
			},
		}...)
	}
	if len(fabricChaincode.Spec.Env) > 0 {
		envVars = append(envVars, fabricChaincode.Spec.Env...)
	}

	podSpec := corev1.PodSpec{
		Volumes:        volumes,
		InitContainers: nil,
		Containers: []corev1.Container{
			{
				Env:             envVars,
				Name:            "chaincode",
				Image:           fabricChaincode.Spec.Image,
				ImagePullPolicy: fabricChaincode.Spec.ImagePullPolicy,
				VolumeMounts:    volumeMounts,
			},
		},
		EphemeralContainers: nil,
		RestartPolicy:       corev1.RestartPolicyAlways,
		ImagePullSecrets:    fabricChaincode.Spec.ImagePullSecrets,
		Affinity:            fabricChaincode.Spec.Affinity,
		Tolerations:         fabricChaincode.Spec.Tolerations,
	}
	replicas := fabricChaincode.Spec.Replicas
	appv1Deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: ns,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: func(i int32) *int32 { return &i }(int32(replicas)),
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},

				Spec: podSpec,
			},
			Strategy: appsv1.DeploymentStrategy{
				Type:          appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: nil,
			},
		},
		Status: appsv1.DeploymentStatus{},
	}

	deployment, err := kubeClientset.AppsV1().Deployments(ns).Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			// creating deployment
			deployment, err = kubeClientset.AppsV1().Deployments(ns).Create(
				ctx,
				appv1Deployment,
				metav1.CreateOptions{},
			)
			if err != nil {
				r.setConditionStatus(ctx, fabricChaincode, hlfv1alpha1.FailedStatus, false, err, false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricChaincode)
			}
		}
		r.setConditionStatus(ctx, fabricChaincode, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricChaincode)
	} else {
		deployment.Spec = appv1Deployment.Spec
		if cryptoData.Updated {
			if deployment.Spec.Template.ObjectMeta.Annotations == nil {
				deployment.Spec.Template.ObjectMeta.Annotations = make(map[string]string)
			}
			deployment.Spec.Template.ObjectMeta.Annotations["hlf.zhangfuli.com/updatedsecrettime"] = time.Now().UTC().Format(time.RFC3339)
		}
		deployment, err = kubeClientset.AppsV1().Deployments(ns).Update(
			ctx,
			deployment,
			metav1.UpdateOptions{},
		)
		if err != nil {
			err = errors.New("failed to update the deployment")
			r.setConditionStatus(ctx, fabricChaincode, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricChaincode)
		}
	}
	serviceObjectMeta := metav1.ObjectMeta{
		Name:      serviceName,
		Namespace: ns,
	}
	defaultServiceIPFamily := corev1.IPv4Protocol
	serviceSpec := corev1.ServiceSpec{
		IPFamily: &defaultServiceIPFamily,
		Ports: []corev1.ServicePort{
			{
				Name:       "chaincode",
				Protocol:   "TCP",
				Port:       7052,
				TargetPort: intstr.FromInt(7052),
			},
		},

		Selector: labels,
		Type:     "ClusterIP",
	}

	service, err := kubeClientset.CoreV1().Services(ns).Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			// creating service
			service, err = kubeClientset.CoreV1().Services(ns).Create(
				ctx,
				&corev1.Service{
					ObjectMeta: serviceObjectMeta,
					Spec:       serviceSpec,
				},
				metav1.CreateOptions{},
			)
			if err != nil {
				r.setConditionStatus(ctx, fabricChaincode, hlfv1alpha1.FailedStatus, false, err, false)
				return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricChaincode)
			}
		}
		r.setConditionStatus(ctx, fabricChaincode, hlfv1alpha1.FailedStatus, false, err, false)
		return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricChaincode)
	} else {
		service.Spec.Ports = serviceSpec.Ports
		service.Spec.Type = serviceSpec.Type
		service.Spec.Selector = serviceSpec.Selector

		service, err = kubeClientset.CoreV1().Services(ns).Update(
			ctx,
			service,
			metav1.UpdateOptions{},
		)
		if err != nil {
			r.setConditionStatus(ctx, fabricChaincode, hlfv1alpha1.FailedStatus, false, err, false)
			return r.updateCRStatusOrFailReconcile(ctx, r.Log, fabricChaincode)
		}
	}
	r.setConditionStatus(ctx, fabricChaincode, hlfv1alpha1.RunningStatus, true, nil, false)
	return ctrl.Result{}, nil
}

var (
	ErrClientK8s = errors.New("k8sAPIClientError")
)

func (r *FabricChaincodeReconciler) updateCRStatusOrFailReconcile(ctx context.Context, log logr.Logger, p *hlfv1alpha1.FabricChaincode) (
	ctrl.Result, error) {
	if err := r.Status().Update(ctx, p); err != nil {
		log.Error(err, fmt.Sprintf("%v failed to update the application status", ErrClientK8s))
		return ctrl.Result{Requeue: false, RequeueAfter: 0}, err
	}
	return ctrl.Result{Requeue: false, RequeueAfter: 0}, nil
}

func (r *FabricChaincodeReconciler) setConditionStatus(ctx context.Context, p *hlfv1alpha1.FabricChaincode, conditionType hlfv1alpha1.DeploymentStatus, statusFlag bool, err error, statusUnknown bool) (update bool) {
	statusStr := func() corev1.ConditionStatus {
		if statusUnknown {
			return corev1.ConditionUnknown
		}
		if statusFlag {
			return corev1.ConditionTrue
		} else {
			return corev1.ConditionFalse
		}
	}
	if p.Status.Status != conditionType {
		depCopy := client.MergeFrom(p.DeepCopy())
		p.Status.Status = conditionType
		err = r.Status().Patch(ctx, p, depCopy)
		if err != nil {
			log.Warnf("Failed to update status to %s: %v", conditionType, err)
		}
	}
	if err != nil {
		p.Status.Message = err.Error()
	} else {
		p.Status.Message = ""
	}
	condition := func() status.Condition {
		if err != nil {
			return status.Condition{
				Type:    status.ConditionType(conditionType),
				Status:  statusStr(),
				Reason:  status.ConditionReason(err.Error()),
				Message: err.Error(),
			}
		}
		return status.Condition{
			Type:   status.ConditionType(conditionType),
			Status: statusStr(),
		}
	}

	return p.Status.Conditions.SetCondition(condition())
}

func (r *FabricChaincodeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	managedBy := ctrl.NewControllerManagedBy(mgr)
	return managedBy.
		For(&hlfv1alpha1.FabricChaincode{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}

func actionLogger(logger logr.Logger) func(format string, v ...interface{}) {
	return func(format string, v ...interface{}) {
		logger.Info(fmt.Sprintf(format, v...))
	}
}
