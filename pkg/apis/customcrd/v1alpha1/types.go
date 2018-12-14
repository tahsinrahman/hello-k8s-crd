package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
import corev1 "k8s.io/api/core/v1"

//+genclient
//+genclient:noStatus
//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CustomDeployment is the definition of a CustomDeployment Resource
type CustomDeployment struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Spec CustomDeploymentSpec
	//Status CustomDeploymentStatus
}

// CustomDeploymentSpec defines the specs of CustomDeploymentSpec
type CustomDeploymentSpec struct {
	Replicas int
	Selector *metav1.LabelSelector
	Template corev1.PodTemplateSpec
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type CustomDeploymentList struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Item []CustomDeployment
}
