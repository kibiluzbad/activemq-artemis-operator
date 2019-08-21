package activemqartemis

import (
	brokerv1alpha1 "github.com/rh-messaging/activemq-artemis-operator/pkg/apis/broker/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"strings"
)

const (
	statefulSetSizeUpdated          = 1 << 0
	statefulSetClusterConfigUpdated = 1 << 1
	statefulSetSSLConfigUpdated     = 1 << 2
	statefulSetImageUpdated         = 1 << 3
	statefulSetPersistentUpdated    = 1 << 4
	statefulSetAioUpdated           = 1 << 5
	statefulSetCommonConfigUpdated  = 1 << 6
)

type ActiveMQArtemisReconciler struct {
	statefulSetUpdates uint32
}

type ActiveMQArtemisIReconciler interface {
	Process(customResource *brokerv1alpha1.ActiveMQArtemis, currentStatefulSet *appsv1.StatefulSet) uint32
}

func (reconciler *ActiveMQArtemisReconciler) Process(customResource *brokerv1alpha1.ActiveMQArtemis, currentStatefulSet *appsv1.StatefulSet) uint32 {

	// Ensure the StatefulSet size is the same as the spec
	if *currentStatefulSet.Spec.Replicas != customResource.Spec.Size {
		currentStatefulSet.Spec.Replicas = &customResource.Spec.Size
		reconciler.statefulSetUpdates |= statefulSetSizeUpdated
	}

	if clusterConfigSyncCausedUpdateOn(customResource, currentStatefulSet) {
		reconciler.statefulSetUpdates |= statefulSetClusterConfigUpdated
	}

	if sslConfigSyncCausedUpdateOn(customResource, currentStatefulSet) {
		reconciler.statefulSetUpdates |= statefulSetSSLConfigUpdated
	}

	if imageSyncCausedUpdateOn(customResource, currentStatefulSet) {
		reconciler.statefulSetUpdates |= statefulSetImageUpdated
	}

	if aioSyncCausedUpdateOn(customResource, currentStatefulSet) {
		reconciler.statefulSetUpdates |= statefulSetAioUpdated
	}

	if persistentSyncCausedUpdateOn(customResource, currentStatefulSet) {
		reconciler.statefulSetUpdates |= statefulSetPersistentUpdated
	}

	if commonConfigSyncCausedUpdateOn(customResource, currentStatefulSet) {
		reconciler.statefulSetUpdates |= statefulSetCommonConfigUpdated
	}

	return reconciler.statefulSetUpdates
}

// https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang
func remove(s []corev1.EnvVar, i int) []corev1.EnvVar {
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1]
}

func clusterConfigSyncCausedUpdateOn(customResource *brokerv1alpha1.ActiveMQArtemis, currentStatefulSet *appsv1.StatefulSet) bool {

	foundClustered := false
	foundClusterUser := false
	foundClusterPassword := false

	clusteredNeedsUpdate := false
	clusterUserNeedsUpdate := false
	clusterPasswordNeedsUpdate := false

	statefulSetUpdated := false

	// TODO: Remove yuck
	// ensure password and username are valid if can't via openapi validation?
	if customResource.Spec.ClusterConfig.ClusterPassword != "" &&
		customResource.Spec.ClusterConfig.ClusterUserName != "" {

		envVarArray := []corev1.EnvVar{}
		// Find the existing values
		for _, v := range currentStatefulSet.Spec.Template.Spec.Containers[0].Env {
			if v.Name == "AMQ_CLUSTERED" {
				foundClustered = true
				if v.Value == "false" {
					clusteredNeedsUpdate = true
				}
			}
			if v.Name == "AMQ_CLUSTER_USER" {
				foundClusterUser = true
				if v.Value != customResource.Spec.ClusterConfig.ClusterUserName {
					clusterUserNeedsUpdate = true
				}
			}
			if v.Name == "AMQ_CLUSTER_PASSWORD" {
				foundClusterPassword = true
				if v.Value != customResource.Spec.ClusterConfig.ClusterPassword {
					clusterPasswordNeedsUpdate = true
				}
			}
		}

		if !foundClustered || clusteredNeedsUpdate {
			newClusteredValue := corev1.EnvVar{
				"AMQ_CLUSTERED",
				"true",
				nil,
			}
			envVarArray = append(envVarArray, newClusteredValue)
			statefulSetUpdated = true
		}

		if !foundClusterUser || clusterUserNeedsUpdate {
			newClusteredValue := corev1.EnvVar{
				"AMQ_CLUSTER_USER",
				customResource.Spec.ClusterConfig.ClusterUserName,
				nil,
			}
			envVarArray = append(envVarArray, newClusteredValue)
			statefulSetUpdated = true
		}

		if !foundClusterPassword || clusterPasswordNeedsUpdate {
			newClusteredValue := corev1.EnvVar{
				"AMQ_CLUSTER_PASSWORD",
				customResource.Spec.ClusterConfig.ClusterPassword,
				nil,
			}
			envVarArray = append(envVarArray, newClusteredValue)
			statefulSetUpdated = true
		}

		if statefulSetUpdated {
			envVarArrayLen := len(envVarArray)
			if envVarArrayLen > 0 {
				for i := 0; i < len(currentStatefulSet.Spec.Template.Spec.Containers); i++ {
					for j := len(currentStatefulSet.Spec.Template.Spec.Containers[i].Env) - 1; j >= 0; j-- {
						if ("AMQ_CLUSTERED" == currentStatefulSet.Spec.Template.Spec.Containers[i].Env[j].Name && clusteredNeedsUpdate) ||
							("AMQ_CLUSTER_USER" == currentStatefulSet.Spec.Template.Spec.Containers[i].Env[j].Name && clusterUserNeedsUpdate) ||
							("AMQ_CLUSTER_PASSWORD" == currentStatefulSet.Spec.Template.Spec.Containers[i].Env[j].Name && clusterPasswordNeedsUpdate) {
							currentStatefulSet.Spec.Template.Spec.Containers[i].Env = remove(currentStatefulSet.Spec.Template.Spec.Containers[i].Env, j)
						}
					}
				}

				containerArrayLen := len(currentStatefulSet.Spec.Template.Spec.Containers)
				for i := 0; i < containerArrayLen; i++ {
					for j := 0; j < envVarArrayLen; j++ {
						currentStatefulSet.Spec.Template.Spec.Containers[i].Env = append(currentStatefulSet.Spec.Template.Spec.Containers[i].Env, envVarArray[j])
					}
				}
			}
		}
	} else {
		for i := 0; i < len(currentStatefulSet.Spec.Template.Spec.Containers); i++ {
			for j := len(currentStatefulSet.Spec.Template.Spec.Containers[i].Env) - 1; j >= 0; j-- {
				if "AMQ_CLUSTERED" == currentStatefulSet.Spec.Template.Spec.Containers[i].Env[j].Name ||
					"AMQ_CLUSTER_USER" == currentStatefulSet.Spec.Template.Spec.Containers[i].Env[j].Name ||
					"AMQ_CLUSTER_PASSWORD" == currentStatefulSet.Spec.Template.Spec.Containers[i].Env[j].Name {
					currentStatefulSet.Spec.Template.Spec.Containers[i].Env = remove(currentStatefulSet.Spec.Template.Spec.Containers[i].Env, j)
					statefulSetUpdated = true
				}
			}
		}
	}

	return statefulSetUpdated
}

func sslConfigSyncCausedUpdateOn(customResource *brokerv1alpha1.ActiveMQArtemis, currentStatefulSet *appsv1.StatefulSet) bool {

	foundKeystore := false
	foundKeystorePassword := false
	foundKeystoreTruststoreDir := false
	foundTruststore := false
	foundTruststorePassword := false

	keystoreNeedsUpdate := false
	keystorePasswordNeedsUpdate := false
	keystoreTruststoreDirNeedsUpdate := false
	truststoreNeedsUpdate := false
	truststorePasswordNeedsUpdate := false

	statefulSetUpdated := false

	// TODO: Remove yuck
	// ensure password and username are valid if can't via openapi validation?
	if customResource.Spec.SSLConfig.KeyStorePassword != "" &&
		customResource.Spec.SSLConfig.KeystoreFilename != "" {

		envVarArray := []corev1.EnvVar{}
		// Find the existing values
		for _, v := range currentStatefulSet.Spec.Template.Spec.Containers[0].Env {
			if v.Name == "AMQ_KEYSTORE" {
				foundKeystore = true
				if v.Value != customResource.Spec.SSLConfig.KeystoreFilename {
					keystoreNeedsUpdate = true
				}
			}
			if v.Name == "AMQ_KEYSTORE_PASSWORD" {
				foundKeystorePassword = true
				if v.Value != customResource.Spec.SSLConfig.KeyStorePassword {
					keystorePasswordNeedsUpdate = true
				}
			}
			if v.Name == "AMQ_KEYSTORE_TRUSTSTORE_DIR" {
				foundKeystoreTruststoreDir = true
				if v.Value != "/etc/amq-secret-volume" {
					keystoreTruststoreDirNeedsUpdate = true
				}
			}
		}

		if !foundKeystore || keystoreNeedsUpdate {
			newClusteredValue := corev1.EnvVar{
				"AMQ_KEYSTORE",
				customResource.Spec.SSLConfig.KeystoreFilename,
				nil,
			}
			envVarArray = append(envVarArray, newClusteredValue)
			statefulSetUpdated = true
		}

		if !foundKeystorePassword || keystorePasswordNeedsUpdate {
			newClusteredValue := corev1.EnvVar{
				"AMQ_KEYSTORE_PASSWORD",
				customResource.Spec.SSLConfig.KeyStorePassword,
				nil,
			}
			envVarArray = append(envVarArray, newClusteredValue)
			statefulSetUpdated = true
		}

		if !foundKeystoreTruststoreDir || keystoreTruststoreDirNeedsUpdate {
			newClusteredValue := corev1.EnvVar{
				"AMQ_KEYSTORE_TRUSTSTORE_DIR",
				"/etc/amq-secret-volume",
				nil,
			}
			envVarArray = append(envVarArray, newClusteredValue)
			statefulSetUpdated = true
		}

		if statefulSetUpdated {
			envVarArrayLen := len(envVarArray)
			if envVarArrayLen > 0 {
				for i := 0; i < len(currentStatefulSet.Spec.Template.Spec.Containers); i++ {
					for j := len(currentStatefulSet.Spec.Template.Spec.Containers[i].Env) - 1; j >= 0; j-- {
						if ("AMQ_KEYSTORE" == currentStatefulSet.Spec.Template.Spec.Containers[i].Env[j].Name && keystoreNeedsUpdate) ||
							("AMQ_KEYSTORE_PASSWORD" == currentStatefulSet.Spec.Template.Spec.Containers[i].Env[j].Name && keystorePasswordNeedsUpdate) ||
							("AMQ_KEYSTORE_TRUSTSTORE_DIR" == currentStatefulSet.Spec.Template.Spec.Containers[i].Env[j].Name && keystoreTruststoreDirNeedsUpdate) {
							currentStatefulSet.Spec.Template.Spec.Containers[i].Env = remove(currentStatefulSet.Spec.Template.Spec.Containers[i].Env, j)
						}
					}
				}

				containerArrayLen := len(currentStatefulSet.Spec.Template.Spec.Containers)
				for i := 0; i < containerArrayLen; i++ {
					for j := 0; j < envVarArrayLen; j++ {
						currentStatefulSet.Spec.Template.Spec.Containers[i].Env = append(currentStatefulSet.Spec.Template.Spec.Containers[i].Env, envVarArray[j])
					}
				}
			}
		}
	} else {
		for i := 0; i < len(currentStatefulSet.Spec.Template.Spec.Containers); i++ {
			for j := len(currentStatefulSet.Spec.Template.Spec.Containers[i].Env) - 1; j >= 0; j-- {
				if "AMQ_KEYSTORE" == currentStatefulSet.Spec.Template.Spec.Containers[i].Env[j].Name ||
					"AMQ_KEYSTORE_PASSWORD" == currentStatefulSet.Spec.Template.Spec.Containers[i].Env[j].Name ||
					"AMQ_KEYSTORE_TRUSTSTORE_DIR" == currentStatefulSet.Spec.Template.Spec.Containers[i].Env[j].Name {
					currentStatefulSet.Spec.Template.Spec.Containers[i].Env = remove(currentStatefulSet.Spec.Template.Spec.Containers[i].Env, j)
					statefulSetUpdated = true
				}
			}
		}
	}

	if customResource.Spec.SSLConfig.TrustStorePassword != "" &&
		customResource.Spec.SSLConfig.TrustStoreFilename != "" {

		envVarArray := []corev1.EnvVar{}
		// Find the existing values
		for _, v := range currentStatefulSet.Spec.Template.Spec.Containers[0].Env {
			if v.Name == "AMQ_TRUSTSTORE" {
				foundTruststore = true
				if v.Value != customResource.Spec.SSLConfig.TrustStoreFilename {
					truststoreNeedsUpdate = true
				}
			}
			if v.Name == "AMQ_TRUSTSTORE_PASSWORD" {
				foundTruststorePassword = true
				if v.Value != customResource.Spec.SSLConfig.TrustStorePassword {
					truststorePasswordNeedsUpdate = true
				}
			}
		}

		if !foundTruststore || truststoreNeedsUpdate {
			newClusteredValue := corev1.EnvVar{
				"AMQ_TRUSTSTORE",
				customResource.Spec.SSLConfig.TrustStoreFilename,
				nil,
			}
			envVarArray = append(envVarArray, newClusteredValue)
			statefulSetUpdated = true
		}

		if !foundTruststorePassword || truststorePasswordNeedsUpdate {
			newClusteredValue := corev1.EnvVar{
				"AMQ_TRUSTSTORE_PASSWORD",
				customResource.Spec.SSLConfig.TrustStorePassword,
				nil,
			}
			envVarArray = append(envVarArray, newClusteredValue)
			statefulSetUpdated = true
		}

		if statefulSetUpdated {
			envVarArrayLen := len(envVarArray)
			if envVarArrayLen > 0 {
				for i := 0; i < len(currentStatefulSet.Spec.Template.Spec.Containers); i++ {
					for j := len(currentStatefulSet.Spec.Template.Spec.Containers[i].Env) - 1; j >= 0; j-- {
						if ("AMQ_TRUSTSTORE" == currentStatefulSet.Spec.Template.Spec.Containers[i].Env[j].Name && truststoreNeedsUpdate) ||
							("AMQ_TRUSTSTORE_PASSWORD" == currentStatefulSet.Spec.Template.Spec.Containers[i].Env[j].Name && truststorePasswordNeedsUpdate) {
							currentStatefulSet.Spec.Template.Spec.Containers[i].Env = remove(currentStatefulSet.Spec.Template.Spec.Containers[i].Env, j)
						}
					}
				}

				containerArrayLen := len(currentStatefulSet.Spec.Template.Spec.Containers)
				for i := 0; i < containerArrayLen; i++ {
					for j := 0; j < envVarArrayLen; j++ {
						currentStatefulSet.Spec.Template.Spec.Containers[i].Env = append(currentStatefulSet.Spec.Template.Spec.Containers[i].Env, envVarArray[j])
					}
				}
			}
		}
	} else {
		for i := 0; i < len(currentStatefulSet.Spec.Template.Spec.Containers); i++ {
			for j := len(currentStatefulSet.Spec.Template.Spec.Containers[i].Env) - 1; j >= 0; j-- {
				if "AMQ_TRUSTSTORE" == currentStatefulSet.Spec.Template.Spec.Containers[i].Env[j].Name ||
					"AMQ_TRUSTSTORE_PASSWORD" == currentStatefulSet.Spec.Template.Spec.Containers[i].Env[j].Name {
					currentStatefulSet.Spec.Template.Spec.Containers[i].Env = remove(currentStatefulSet.Spec.Template.Spec.Containers[i].Env, j)
					statefulSetUpdated = true
				}
			}
		}
	}

	if statefulSetUpdated {
		sslConfigSyncEnsureSecretVolumeMountExists(customResource, currentStatefulSet)
	}

	return statefulSetUpdated
}

func sslConfigSyncEnsureSecretVolumeMountExists(customResource *brokerv1alpha1.ActiveMQArtemis, currentStatefulSet *appsv1.StatefulSet) {

	secretVolumeExists := false
	secretVolumeMountExists := false

	for i := 0; i < len(currentStatefulSet.Spec.Template.Spec.Volumes); i++ {
		if currentStatefulSet.Spec.Template.Spec.Volumes[i].Name == customResource.Spec.SSLConfig.SecretName {
			secretVolumeExists = true
			break
		}
	}
	if !secretVolumeExists {
		volume := corev1.Volume{
			Name: "broker-secret-volume",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: customResource.Spec.SSLConfig.SecretName,
				},
			},
		}

		currentStatefulSet.Spec.Template.Spec.Volumes = append(currentStatefulSet.Spec.Template.Spec.Volumes, volume)
	}

	for i := 0; i < len(currentStatefulSet.Spec.Template.Spec.Containers); i++ {
		for j := 0; j < len(currentStatefulSet.Spec.Template.Spec.Containers[i].VolumeMounts); j++ {
			if currentStatefulSet.Spec.Template.Spec.Containers[i].VolumeMounts[j].Name == "broker-secret-volume" {
				secretVolumeMountExists = true
				break
			}
		}
		if !secretVolumeMountExists {
			volumeMount := corev1.VolumeMount{
				Name:      "broker-secret-volume",
				MountPath: "/etc/amq-secret-volume",
				ReadOnly:  true,
			}
			currentStatefulSet.Spec.Template.Spec.Containers[i].VolumeMounts = append(currentStatefulSet.Spec.Template.Spec.Containers[i].VolumeMounts, volumeMount)
		}
	}
}

func aioSyncCausedUpdateOn(customResource *brokerv1alpha1.ActiveMQArtemis, currentStatefulSet *appsv1.StatefulSet) bool {

	foundAio := false
	foundNio := false
	var extraArgs string = ""
	extraArgsNeedsUpdate := false

	// Find the existing values
	for _, v := range currentStatefulSet.Spec.Template.Spec.Containers[0].Env {
		if v.Name == "AMQ_EXTRA_ARGS" {
			if strings.Index(v.Value, "--aio") > -1 {
				foundAio = true
			}
			if strings.Index(v.Value, "--nio") > -1 {
				foundNio = true
			}
			extraArgs = v.Value
			break
		}
	}

	if customResource.Spec.Aio && foundNio {
		extraArgs = strings.Replace(extraArgs, "--nio", "--aio", 1)
		extraArgsNeedsUpdate = true
	}

	if !customResource.Spec.Aio && foundAio {
		extraArgs = strings.Replace(extraArgs, "--aio", "--nio", 1)
		extraArgsNeedsUpdate = true
	}

	newExtraArgsValue := corev1.EnvVar{}
	if extraArgsNeedsUpdate {
		newExtraArgsValue = corev1.EnvVar{
			"AMQ_EXTRA_ARGS",
			extraArgs,
			nil,
		}

		containerArrayLen := len(currentStatefulSet.Spec.Template.Spec.Containers)
		for i := 0; i < containerArrayLen; i++ {
			//for j := 0; j < envVarArrayLen; j++ {
			for j := 0; j < len(currentStatefulSet.Spec.Template.Spec.Containers[i].Env); j++ {
				if "AMQ_EXTRA_ARGS" == currentStatefulSet.Spec.Template.Spec.Containers[i].Env[j].Name {
					currentStatefulSet.Spec.Template.Spec.Containers[i].Env = remove(currentStatefulSet.Spec.Template.Spec.Containers[i].Env, j)
					currentStatefulSet.Spec.Template.Spec.Containers[i].Env = append(currentStatefulSet.Spec.Template.Spec.Containers[i].Env, newExtraArgsValue)
					break
				}
			}
		}
	}

	return extraArgsNeedsUpdate
}

func persistentSyncCausedUpdateOn(customResource *brokerv1alpha1.ActiveMQArtemis, currentStatefulSet *appsv1.StatefulSet) bool {

	foundDataDir := false
	foundDataDirLogging := false

	dataDirNeedsUpdate := false
	dataDirLoggingNeedsUpdate := false

	statefulSetUpdated := false

	// TODO: Remove yuck
	// ensure password and username are valid if can't via openapi validation?
	if customResource.Spec.Persistent {

		envVarArray := []corev1.EnvVar{}
		// Find the existing values
		for _, v := range currentStatefulSet.Spec.Template.Spec.Containers[0].Env {
			if v.Name == "AMQ_DATA_DIR" {
				foundDataDir = true
				if v.Value == "false" {
					dataDirNeedsUpdate = true
				}
			}
			if v.Name == "AMQ_DATA_DIR_LOGGING" {
				foundDataDirLogging = true
				if v.Value != "true" {
					dataDirLoggingNeedsUpdate = true
				}
			}
		}

		if !foundDataDir || dataDirNeedsUpdate {
			newDataDirValue := corev1.EnvVar{
				"AMQ_DATA_DIR",
				"/opt/" + customResource.Name + "/data",
				nil,
			}
			envVarArray = append(envVarArray, newDataDirValue)
			statefulSetUpdated = true
		}

		if !foundDataDirLogging || dataDirLoggingNeedsUpdate {
			newDataDirLoggingValue := corev1.EnvVar{
				"AMQ_DATA_DIR_LOGGING",
				"true",
				nil,
			}
			envVarArray = append(envVarArray, newDataDirLoggingValue)
			statefulSetUpdated = true
		}

		if statefulSetUpdated {
			envVarArrayLen := len(envVarArray)
			if envVarArrayLen > 0 {
				for i := 0; i < len(currentStatefulSet.Spec.Template.Spec.Containers); i++ {
					for j := len(currentStatefulSet.Spec.Template.Spec.Containers[i].Env) - 1; j >= 0; j-- {
						if ("AMQ_DATA_DIR" == currentStatefulSet.Spec.Template.Spec.Containers[i].Env[j].Name && dataDirNeedsUpdate) ||
							("AMQ_DATA_DIR_LOGGING" == currentStatefulSet.Spec.Template.Spec.Containers[i].Env[j].Name && dataDirLoggingNeedsUpdate) {
							currentStatefulSet.Spec.Template.Spec.Containers[i].Env = remove(currentStatefulSet.Spec.Template.Spec.Containers[i].Env, j)
						}
					}
				}

				containerArrayLen := len(currentStatefulSet.Spec.Template.Spec.Containers)
				for i := 0; i < containerArrayLen; i++ {
					for j := 0; j < envVarArrayLen; j++ {
						currentStatefulSet.Spec.Template.Spec.Containers[i].Env = append(currentStatefulSet.Spec.Template.Spec.Containers[i].Env, envVarArray[j])
					}
				}
			}
		}
	} else {

		for i := 0; i < len(currentStatefulSet.Spec.Template.Spec.Containers); i++ {
			for j := len(currentStatefulSet.Spec.Template.Spec.Containers[i].Env) - 1; j >= 0; j-- {
				if "AMQ_DATA_DIR" == currentStatefulSet.Spec.Template.Spec.Containers[i].Env[j].Name ||
					"AMQ_DATA_DIR_LOGGING" == currentStatefulSet.Spec.Template.Spec.Containers[i].Env[j].Name {
					currentStatefulSet.Spec.Template.Spec.Containers[i].Env = remove(currentStatefulSet.Spec.Template.Spec.Containers[i].Env, j)
					statefulSetUpdated = true
				}
			}
		}
	}

	return statefulSetUpdated
}

func imageSyncCausedUpdateOn(customResource *brokerv1alpha1.ActiveMQArtemis, currentStatefulSet *appsv1.StatefulSet) bool {

	// At implementation time only one container
	if strings.Compare(currentStatefulSet.Spec.Template.Spec.Containers[0].Image, customResource.Spec.Image) != 0 {
		containerArrayLen := len(currentStatefulSet.Spec.Template.Spec.Containers)
		for i := 0; i < containerArrayLen; i++ {
			currentStatefulSet.Spec.Template.Spec.Containers[i].Image = customResource.Spec.Image
		}
		return true
	}

	return false
}

func commonConfigSyncCausedUpdateOn(customResource *brokerv1alpha1.ActiveMQArtemis, currentStatefulSet *appsv1.StatefulSet) bool {

	foundCommonUser := false
	foundCommonPassword := false

	commonUserNeedsUpdate := false
	commonPasswordNeedsUpdate := false

	statefulSetUpdated := false

	// TODO: Remove yuck
	// ensure password and username are valid if can't via openapi validation?
	if customResource.Spec.CommonConfig.Password != "" &&
		customResource.Spec.CommonConfig.UserName != "" {

		envVarArray := []corev1.EnvVar{}
		// Find the existing values
		for _, v := range currentStatefulSet.Spec.Template.Spec.Containers[0].Env {
			if v.Name == "AMQ_USER" {
				foundCommonUser = true
				if v.Value != customResource.Spec.CommonConfig.UserName {
					commonUserNeedsUpdate = true
				}
			}
			if v.Name == "AMQ_PASSWORD" {
				foundCommonPassword = true
				if v.Value != customResource.Spec.CommonConfig.Password {
					commonPasswordNeedsUpdate = true
				}
			}
		}

		if !foundCommonUser || commonUserNeedsUpdate {
			newCommonedValue := corev1.EnvVar{
				"AMQ_USER",
				customResource.Spec.CommonConfig.UserName,
				nil,
			}
			envVarArray = append(envVarArray, newCommonedValue)
			statefulSetUpdated = true
		}

		if !foundCommonPassword || commonPasswordNeedsUpdate {
			newCommonedValue := corev1.EnvVar{
				"AMQ_PASSWORD",
				customResource.Spec.CommonConfig.Password,
				nil,
			}
			envVarArray = append(envVarArray, newCommonedValue)
			statefulSetUpdated = true
		}

		if statefulSetUpdated {
			envVarArrayLen := len(envVarArray)
			if envVarArrayLen > 0 {
				for i := 0; i < len(currentStatefulSet.Spec.Template.Spec.Containers); i++ {
					for j := len(currentStatefulSet.Spec.Template.Spec.Containers[i].Env) - 1; j >= 0; j-- {
						if ("AMQ_USER" == currentStatefulSet.Spec.Template.Spec.Containers[i].Env[j].Name && commonUserNeedsUpdate) ||
							("AMQ_PASSWORD" == currentStatefulSet.Spec.Template.Spec.Containers[i].Env[j].Name && commonPasswordNeedsUpdate) {
							currentStatefulSet.Spec.Template.Spec.Containers[i].Env = remove(currentStatefulSet.Spec.Template.Spec.Containers[i].Env, j)
						}
					}
				}

				containerArrayLen := len(currentStatefulSet.Spec.Template.Spec.Containers)
				for i := 0; i < containerArrayLen; i++ {
					for j := 0; j < envVarArrayLen; j++ {
						currentStatefulSet.Spec.Template.Spec.Containers[i].Env = append(currentStatefulSet.Spec.Template.Spec.Containers[i].Env, envVarArray[j])
					}
				}
			}
		}
	}

	return statefulSetUpdated
}
