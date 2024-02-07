/*
Copyright 2024.

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

package boanlab

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	minicloudControllerType "github.com/boanlab/mini-cloud-server/controller/api/rest/types"
	"github.com/go-logr/logr"
	boanlabv1 "github.com/kubernetes-incubator/apiserver-builder/pkg/pkg/apis/boanlab/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// InstanceReconciler reconciles a Instance object
type InstanceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=boanlab,resources=instances,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=boanlab,resources=instances/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=boanlab,resources=instances/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *InstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Load instance
	var instance boanlabv1.Instance
	err := r.Get(ctx, req.NamespacedName, &instance)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Log.Error(err, "unable to fetch Instance")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Log.Error(err, "Failed to get Instance")
		return reconcile.Result{}, err
		//return ctrl.Result{}, client.IgnoreNotFound(err)
		//return reconcile.Result{}, nil
	}

	// Find helper-pod
	helperPod := &corev1.Pod{}
	var owner = instance.Spec.Environment.Owner
	var name = instance.Name
	// err = a.List(ctx, pods, client.InNamespace(req.Namespace), client.MatchingLabels(rs.Spec.Template.Labels))
	helperPodName := fmt.Sprintf("%s-%s", owner, name)

	err = r.Client.Get(ctx, types.NamespacedName{Name: helperPodName, Namespace: instance.Namespace}, helperPod)
	if err != nil && errors.IsNotFound(err) {
		// Create helper-pod : API call to mini-cloud-server Controller
		// TODO: migrate CR controller to mini-cloud-server controller

		log.Log.Info("Creating a new Instance(helper p od)", "helperPod.name", helperPodName)
		newInstance, err := r.PostInstance(instance)
		if err != nil {
			log.Log.Error(err, "Failed to create new Instance", "instance.Name", instance.Name, "helperPod.name", helperPodName)
			return reconcile.Result{}, err
		}

		// Helper-pod created successfully - return and requeue
		log.Log.Info("Success to create new Instance", "instance UUID", newInstance.UUID)

		// Update Status Subresource (InstanceID)
		instance.Status.InstanceID = newInstance.UUID

		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		log.Log.Error(err, "Failed to get helper-pod")
		return reconcile.Result{}, err
	}

	// Stop reconcile
	return ctrl.Result{}, nil
}

// PostInstance create a request(create new instance) to mini-cloud-server Controller.
func (r *InstanceReconciler) PostInstance(instance boanlabv1.Instance) (*minicloudControllerType.InstanceCreationResponse, error) {
	postRequest := minicloudControllerType.InstanceCreationRequest{
		Name:        instance.Name,
		Owner:       instance.Spec.Environment.Owner,
		Description: "post request from k8s",
		Os:          instance.Spec.Environment.Os,
		CpuSize:     instance.Spec.Resource.CpuLimit,
		RamSize:     instance.Spec.Resource.RamLimit,
		DiskSize:    instance.Spec.Resource.DiskLimit,
	}

	// load minicloud-Controller server info from env
	// TODO: add to k8s yaml file
	// TODO: migrate to main.go or misc/env.go
	hostIP := os.Getenv("NEBULA_REST_API_HOST_IP")
	hostPort := os.Getenv("NEBULA_REST_API_HOST_PORT")

	url := "http://" + hostIP + ":" + hostPort + "/note"

	// post request to minicloud-Controller serve
	buf, _ := json.Marshal(postRequest)
	response, err := http.Post(url, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		//log.Log.Error(err, "Cannot send [post] instance creation to Server ...")
		return nil, err
	}

	var createdInstance minicloudControllerType.InstanceCreationResponse
	err = json.NewDecoder(response.Body).Decode(&createdInstance)
	if err != nil {
		return nil, err
		//log.Log.Error(err, "Error decoding JSON response")
	}

	return &createdInstance, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *InstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&boanlabv1.Instance{}).
		Complete(r)
}
