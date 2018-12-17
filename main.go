package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	crd_api "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	customResource "github.com/tahsinrahman/hello-k8s-crd/pkg/apis/customcrd.com/v1alpha1"

	myclient "github.com/tahsinrahman/hello-k8s-crd/pkg/client/clientset/versioned"
	crd "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/tools/clientcmd"

	kutil_crd "github.com/appscode/kutil/apiextensions/v1beta1"
)

func main() {
	kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		log.Fatal(err)
	}

	myClientset, err := myclient.NewForConfig(config)
	if err != nil {
		log.Fatal("error creating client for custom crd: ", err)
	}

	//kubeClientset, err := kubernetes.NewForConfig(config)
	//if err != nil {
	//log.Fatal("error creating client: ", err)
	//}

	crdClientset, err := crd.NewForConfig(config)
	if err != nil {
		log.Fatal("error creating client for apiextensions-apiserver: ", err)
	}

	mycrd := &crd_api.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CustomResourceDefinition",
			APIVersion: "apiextensions.k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "customdeployments.customcrd.com",
		},
		Spec: crd_api.CustomResourceDefinitionSpec{
			Group:   "customcrd.com",
			Version: "v1alpha1",
			Versions: []crd_api.CustomResourceDefinitionVersion{
				{
					Name:    "v1alpha1",
					Served:  true,
					Storage: true,
				},
			},
			Scope: crd_api.NamespaceScoped,
			Names: crd_api.CustomResourceDefinitionNames{
				Plural:     "customdeployments",
				Singular:   "customdeployment",
				ShortNames: []string{"cd"},
				Kind:       "CustomDeployment",
			},
		},
	}

	log.Println("creating crd")
	if err = kutil_crd.RegisterCRDs(crdClientset.ApiextensionsV1beta1(), []*crd_api.CustomResourceDefinition{mycrd}); err != nil {
		log.Fatal("error creating client for apiextensions-apiserver: ", err)
	}
	log.Println("successfully created CRD")

	log.Println(mycrd)

	customDeploymentResource := myClientset.CustomcrdV1alpha1().CustomDeployments("default")
	myDeploy, err := customDeploymentResource.Create(&customResource.CustomDeployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CustomDeployment",
			APIVersion: "customcrd.com/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-custom-deployment-resource",
			Labels: map[string]string{
				"app": "bookserver",
			},
		},
		Spec: customResource.CustomDeploymentSpec{
			Replicas: 5,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: "docker-pod-name",
					Labels: map[string]string{
						"name": "docker-pod-label",
					},
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						corev1.Volume{
							Name: "dind-storage",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
					Containers: []corev1.Container{
						corev1.Container{
							Name:  "docker-dind",
							Image: "docker:18.09-dind",
							VolumeMounts: []corev1.VolumeMount{
								corev1.VolumeMount{
									Name:      "dind-storage",
									MountPath: "/var/lib/docker",
								},
							},
							SecurityContext: &corev1.SecurityContext{
								Privileged: boolToP(true),
							},
						},
						corev1.Container{
							Name:    "docker-container",
							Image:   "docker:18.09",
							Command: []string{"/bin/sh", "-c"},
							Args:    []string{"docker run -p 4000:$PORT tahsin/booklist-api:0.0.1 --port=$PORT"},
							Ports: []corev1.ContainerPort{
								corev1.ContainerPort{
									Name:          "booklist-port",
									ContainerPort: 4000,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							Env: []corev1.EnvVar{
								corev1.EnvVar{
									Name:  "DOCKER_HOST",
									Value: "tcp://localhost:2375",
								},
								corev1.EnvVar{
									Name: "PORT",
									ValueFrom: &corev1.EnvVarSource{
										ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "api-configs",
											},
											Key: "port",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Println("error creating custom deployment: ", err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch

	log.Println("Shutting Down")

	if err := crdClientset.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(mycrd.ObjectMeta.Name, &metav1.DeleteOptions{}); err != nil {
		log.Println("error deleting crd: ", err)
	}
	if err := customDeploymentResource.Delete(myDeploy.ObjectMeta.Name, &metav1.DeleteOptions{}); err != nil {
		log.Println("error resource: ", err)
	}
}

func boolToP(value bool) *bool {
	return &value
}
