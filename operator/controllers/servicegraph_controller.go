/*
Copyright 2021 Tutkovics András.

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

package controllers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "k8s.io/api/apps/v1"
	autoscalev1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"

	diptervv1beta1 "example.com/m/v2/api/v1beta1"
)

// ServicegraphReconciler reconciles a Servicegraph object
type ServicegraphReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=dipterv.my.domain,resources=servicegraphs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dipterv.my.domain,resources=servicegraphs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dipterv.my.domain,resources=servicegraphs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Servicegraph object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *ServicegraphReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	//ctx = context.Background()
	//log := r.Log.WithValues("servicegraph", req.NamespacedName)
	log := r.Log.WithValues("servicegraph", req.NamespacedName)

	log.Info("New Reconcile loop started.")

	servicegraph := &diptervv1beta1.Servicegraph{}
	err := r.Get(ctx, req.NamespacedName, servicegraph)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("Servicegraph resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get ServiceGraph resource")
		return ctrl.Result{}, err
	}

	for _, node := range servicegraph.Spec.Nodes {
		// Check if the deployment for the node already exists, if not create a new one
		found := &appsv1.Deployment{}

		err = r.Get(ctx, types.NamespacedName{Name: node.Name, Namespace: "default"}, found)
		if err != nil && errors.IsNotFound(err) {
			//fmt.Printf("######### CREATE: %d node type: %T\n", i, node)
			// Define a new deployment for the node
			dep := r.deploymentForNode(node, servicegraph)
			log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)

			err = r.Create(ctx, dep)
			if err != nil {
				log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
				return ctrl.Result{}, err
			}
			// Deployment created successfully - return and requeue
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			log.Error(err, "Failed to get Deployment")
			return ctrl.Result{}, err
		}

		foundHPA := &autoscalev1.HorizontalPodAutoscaler{}
		err = r.Get(ctx, types.NamespacedName{Name: node.Name, Namespace: "default"}, foundHPA)

		// Check maxReplicas --> if 0 HPA was not specified
		if node.HPA.MaxReplicas != 0 || node.HPA.MinReplicas != 0 {

			if err != nil && errors.IsNotFound(err) {
				hpa := r.hpaForNode(node, servicegraph)
				log.Info("Creating a new HPA", "HPA.Namespace", hpa.Namespace, "HPA.Name", hpa.Name)

				err = r.Create(ctx, hpa)
				if err != nil {
					log.Error(err, "Failed to create new HPA", "HPA.Namespace", hpa.Namespace, "HPA.Name", hpa.Name)
					return ctrl.Result{}, err
				}
				// HPA created successfully - return and requeue
				return ctrl.Result{Requeue: true}, nil
			} else if err != nil {
				log.Error(err, "Failed to get HorizontalPodAutoscaler")
				return ctrl.Result{}, err
			}
		} else {
			// If HPA was not specified --> Ensure the deployment size is the same as the spec
			size := int32(node.Replicas)
			if *found.Spec.Replicas != size {
				found.Spec.Replicas = &size
				err = r.Update(ctx, found)
				if err != nil {
					log.Error(err, "Failed to update Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
					return ctrl.Result{}, err
				}
				// Spec updated - return and requeue
				return ctrl.Result{Requeue: true}, nil
			}
		}

		foundSVC := &corev1.Service{}
		err = r.Get(ctx, types.NamespacedName{Name: node.Name, Namespace: "default"}, foundSVC)

		if err != nil && errors.IsNotFound(err) {
			svc := r.serviceForNode(node, servicegraph)
			log.Info("Creating a new SVC", "SVC.Namespace", svc.Namespace, "SVC.Name", svc.Name)

			err = r.Create(ctx, svc)
			if err != nil {
				log.Error(err, "Failed to create new SVC", "SVC.Namespace", svc.Namespace, "SVC.Name", svc.Name)
				return ctrl.Result{}, err
			}
			// SVC created successfully - return and requeue
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			log.Error(err, "Failed to get Service")
			return ctrl.Result{}, err
		}

	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ServicegraphReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&diptervv1beta1.Servicegraph{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&autoscalev1.HorizontalPodAutoscaler{}).
		Complete(r)
}

func (r *ServicegraphReconciler) deploymentForNode(node *diptervv1beta1.Node, sg *diptervv1beta1.Servicegraph) *appsv1.Deployment {
	ls := r.labelsForNode(node)
	replicas := int32(node.Replicas)
	args := r.createCommandForNode(node)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      node.Name,
			Namespace: "default", // TODO: use correct ns
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						// Image:   "tuti/service-graph-simulator:json2",
						Image:           "tuti/service-graph-simulator:latest",
						ImagePullPolicy: "Always",  //corev1.PullPolicy PullAlways,
						Name:            node.Name, // from: "servicenode" --> so we can get container resource usage
						Command:         args,
						Ports: []corev1.ContainerPort{{
							ContainerPort: int32(node.ContainerPort),
							Name:          "listen",
						}},
						// Set CPU and Memory limit and requests
						Resources: corev1.ResourceRequirements{
							Limits: map[corev1.ResourceName]resource.Quantity{
								corev1.ResourceCPU: resource.MustParse(
									strconv.FormatUint(uint64(node.Resources.CPU), 10) + "m"),
								corev1.ResourceMemory: resource.MustParse(
									strconv.FormatUint(uint64(node.Resources.Memory), 10) + "Mi"),
							},
							Requests: map[corev1.ResourceName]resource.Quantity{
								corev1.ResourceCPU: resource.MustParse(
									strconv.FormatUint(uint64(node.Resources.CPU), 10) + "m"),
								corev1.ResourceMemory: resource.MustParse(
									strconv.FormatUint(uint64(node.Resources.Memory), 10) + "Mi"),
							},
						},
					}},
				},
			},
		},
	}

	fmt.Printf("Created Deployment: %+v", dep)

	// Set deployment instance as the owner and controller
	ctrl.SetControllerReference(sg, dep, r.Scheme)
	return dep
}

func (r *ServicegraphReconciler) hpaForNode(node *diptervv1beta1.Node, sg *diptervv1beta1.Servicegraph) *autoscalev1.HorizontalPodAutoscaler {
	//ls := r.labelsForNode(node)
	// replicas := int32(node.Replicas)
	//args := r.createCommandForNode(node)
	utilization := int32(node.HPA.Utilization)
	min_replicas := int32(node.HPA.MinReplicas)
	max_replicas := int32(node.HPA.MaxReplicas)

	hpa := &autoscalev1.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      node.Name,
			Namespace: "default",
		},
		Spec: autoscalev1.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalev1.CrossVersionObjectReference{
				Kind:       "Deployment",
				Name:       node.Name,
				APIVersion: "apps/v1",
			},
			MinReplicas:                    &min_replicas,
			MaxReplicas:                    max_replicas,
			TargetCPUUtilizationPercentage: &utilization,
		},
	}

	fmt.Printf("Created HPA: %+v", hpa)

	// Set deployment instance as the owner and controller
	ctrl.SetControllerReference(sg, hpa, r.Scheme)
	return hpa
}

// labelsForNode returns the labels for selecting the resources
func (r *ServicegraphReconciler) labelsForNode(node *diptervv1beta1.Node) map[string]string {
	return map[string]string{"app": "servicegraph", "node": node.Name}
}

func (r *ServicegraphReconciler) serviceForNode(node *diptervv1beta1.Node, sg *diptervv1beta1.Servicegraph) *corev1.Service {
	ls := r.labelsForNode(node)
	cPort := int32(node.ContainerPort)
	nPort := int32(node.NodePort)
	sType := corev1.ServiceTypeClusterIP

	if nPort != 0 {
		sType = corev1.ServiceTypeNodePort
	}

	port := []corev1.ServicePort{
		corev1.ServicePort{
			Port:     cPort,
			NodePort: nPort,
		},
	}

	service := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      node.Name,
			Namespace: "default",
		},
		Spec: corev1.ServiceSpec{
			Selector: ls,
			Ports:    port,
			Type:     sType,
		},
	}

	//fmt.Printf("\nSERVICE #### Nodeport: %d ######\n%+v\n", node.NodePort, service)

	// Set Service instance as the owner and controller
	ctrl.SetControllerReference(sg, service, r.Scheme)

	return service
}

func (r *ServicegraphReconciler) createCommandForNode(node *diptervv1beta1.Node) []string {
	// eg: -name Backend -delay 90 -port 9999 -cpu 90 -memory 900 -endpoint-url /read -endpoint-cpu 99 -endpoint-delay 192 -endpoint-url /index -endpoint-cpu 22 -endpoint-delay 111 -endpoint-call='"back-end:9898/staus__front-end:9876/health"' -endpoint-call "database:1234/asd?q=user1"
	var cmd []string

	cmd = append(cmd, "/app/main") //start my application

	// application parameters
	cmd = append(cmd, fmt.Sprintf("-name=%s", node.Name))
	cmd = append(cmd, fmt.Sprintf("-port=%d", node.ContainerPort))
	cmd = append(cmd, fmt.Sprintf("-cpu=%d", node.Resources.CPU))
	cmd = append(cmd, fmt.Sprintf("-memory=%d", node.Resources.Memory))

	for _, ep := range node.Endpoints {
		// convert list of callouts to "call-out1__call-out2..." format
		var tmpArray []string
		for _, ca := range ep.CallOuts {
			tmpArray = append(tmpArray, string(ca.URL))
		}
		callOutParsed := strings.Join(tmpArray, "__")

		cmd = append(cmd, fmt.Sprintf("-endpoint-url=%s", ep.Path))
		cmd = append(cmd, fmt.Sprintf("-endpoint-delay=%d", ep.Delay))
		cmd = append(cmd, fmt.Sprintf("-endpoint-call='%s'", callOutParsed))
		cmd = append(cmd, fmt.Sprintf("-endpoint-cpu=%d", ep.CPULoad))
	}

	fmt.Printf("Command to run in container: %s", cmd)
	return cmd
}
