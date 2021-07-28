/*
Copyright 2021.

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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	graphnodev1alpha1 "github.com/streamingfast/graphnode-operator/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var oneMibiByte = resource.MustParse("1Mi")

// TestRunReconciler reconciles a TestRun object
type TestRunReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=graphnode.streamingfast.io,resources=testruns,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=graphnode.streamingfast.io,resources=testruns/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=graphnode.streamingfast.io,resources=testruns/finalizers,verbs=update

var logger = log.Log.WithName("global")

// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *TestRunReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	testRun := &graphnodev1alpha1.TestRun{}
	err := r.Get(ctx, req.NamespacedName, testRun)
	if err != nil {
		if apierrors.IsNotFound(err) {
			logger.V(0).Info("no testrun found")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	logger.V(1).Info(fmt.Sprintf("Managing TestRun:%+v", testRun.Spec))
	if err := r.EnsurePVC(ctx, testRun); err != nil {
		logger.Error(err, "trying to ensurePVC")
		return ctrl.Result{}, err
	}

	if err := r.EnsureJob(ctx, testRun); err != nil {
		logger.Error(err, "trying to ensureJOB")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func deterministicPVCName(testRunName string) string {
	return fmt.Sprintf("sql-%s", testRunName)
}

func deterministicJobName(testRunName string) string {
	return fmt.Sprintf("job-%s", testRunName)
}

func (r *TestRunReconciler) EnsureJob(ctx context.Context, testRun *graphnodev1alpha1.TestRun) error {
	job := jobDef(testRun)

	controllerutil.SetControllerReference(testRun, job, r.Scheme) // force ownership
	err := r.Create(ctx, job)
	if err == nil {
		logger.Info("Created job", "namespace", job.Namespace, "name", job.Name)
	}
	if apierrors.IsAlreadyExists(err) { //FIXME check before ? what is optimal here ?
		return nil
	}
	return err
}

func (r *TestRunReconciler) EnsurePVC(ctx context.Context, testRun *graphnodev1alpha1.TestRun) error {
	pvcName := deterministicPVCName(testRun.Name)
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:        pvcName,
			Namespace:   testRun.Namespace,
			Annotations: nil,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					"storage": testRun.Spec.PGDataStorage,
				},
			},
			StorageClassName: &testRun.Spec.StorageClassName,
		},
	}

	controllerutil.SetControllerReference(testRun, pvc, r.Scheme) // force ownership
	err := r.Create(ctx, pvc)
	if err == nil {
		logger.Info("Created pvc", "namespace", pvc.Namespace, "name", pvc.Name)
	}
	if apierrors.IsAlreadyExists(err) { //FIXME check before ? what is optimal here ?
		return nil
	}
	return err
}

func jobDef(testRun *graphnodev1alpha1.TestRun) *batchv1.Job {
	//var tolerations []corev1.Toleration
	//var affinity corev1.Affinity
	//if batchConfig.NodePool != "" {
	//        tolerations = tolerationsForPool(batchConfig.NodePool)
	//        affinity = affinityForPool(batchConfig.NodePool)
	//}
	backoffLimit := int32(0)

	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:        deterministicJobName(testRun.Name),
			Namespace:   testRun.Namespace,
			Annotations: map[string]string{"cluster-autoscaler.kubernetes.io/safe-to-evict": "false"}, // no eviction permitted
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: &backoffLimit,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{"cluster-autoscaler.kubernetes.io/safe-to-evict": "false"}, // no eviction permitted
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:      "sql",
							Image:     testRun.Spec.PostgresTarballerDockerImage,
							Command:   []string{"/restorer.sh"},
							Resources: testRun.Spec.PostgresResources,
							Env: []corev1.EnvVar{
								{Name: "PGDATA", Value: "/data/db"},
								{Name: "POSTGRES_PASSWORD", Value: testRun.Spec.PostgresPassword},
								{Name: "POSTGRES_USER", Value: testRun.Spec.PostgresUser},
								{Name: "POSTGRES_DB", Value: testRun.Spec.PostgresDBName},
								{Name: "SRC_TARBALL_URL", Value: testRun.Spec.TarballsURL},
								{Name: "SIGNALING_FOLDER", Value: "/data/signal"},
								//"SRC_TARBALL_FILENAME" - latest alphabetically by default
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "dbdata",
									MountPath: "/data/db",
								},
								{
									Name:      "signal",
									MountPath: "/data/signal",
								},
								{
									Name:      "output",
									MountPath: "/data/output",
								},
							},
						},
						{
							Name:  "graph-node",
							Image: testRun.Spec.GraphnodeDockerImage,
							Command: []string{
								"sh",
								"-c",
								"git clone $GITREPO --branch $GITBRANCH --single-branch /data/graph-node && cd /data/graph-node;" +
									"cargo build -p graph-node;" +
									"echo Waiting for ${SIGNALING_FOLDER}/dbready...; while sleep 1; do test -e ${SIGNALING_FOLDER}/dbready && break; done;" +
									"./target/debug/graph-node --ethereum-rpc=${ETHEREUM_RPC} --postgres-url=\"postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${POSTGRES_DB}\" --ipfs=${IPFS_ADDR};" +
									"OUTPUT_TAG=$(date +%s)-$(git rev-parse --short HEAD);" +
									"echo \"cd /data/output && gsutil -m cp -r . ${OUTPUT_URL}/${OUTPUT_TAG}\" > /data/signal/complete.tmp && " + // the tarballer-postgresql will execute this command
									"chmod +x /data/signal/complete.tmp && mv /data/signal/complete.tmp /data/signal/complete",
							},
							Resources: testRun.Spec.GraphnodeResources,
							Env: []corev1.EnvVar{
								{Name: "GITREPO", Value: testRun.Spec.GitRepo},
								{Name: "GITBRANCH", Value: testRun.Spec.GitBranch},
								{Name: "ETHEREUM_RPC", Value: testRun.Spec.EthereumRPCAddress},
								{Name: "IPFS_ADDR", Value: testRun.Spec.IPFSAddr},
								{Name: "POSTGRES_PASSWORD", Value: testRun.Spec.PostgresPassword},
								{Name: "POSTGRES_USER", Value: testRun.Spec.PostgresUser},
								{Name: "POSTGRES_DB", Value: testRun.Spec.PostgresDBName},
								{Name: "GRAPH_STOP_BLOCK", Value: fmt.Sprintf("%d", testRun.Spec.StopBlock)},
								{Name: "SIGNALING_FOLDER", Value: "/data/signal"},
								{Name: "GRAPH_DEBUG_POI_FILE", Value: "/data/output/poi.csv"},
								{Name: "OUTPUT_URL", Value: testRun.Spec.TestOutputURL},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "signal",
									MountPath: "/data/signal",
								},
								{
									Name:      "output",
									MountPath: "/data/output",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name:         "dbdata",
							VolumeSource: corev1.VolumeSource{PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: deterministicPVCName(testRun.Name)}},
						},
						{
							Name:         "signal",
							VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{SizeLimit: &oneMibiByte}},
						},
						{
							Name:         "output",
							VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{SizeLimit: &testRun.Spec.OutputDirSize}},
						},
					},
					RestartPolicy:      corev1.RestartPolicyNever,
					ServiceAccountName: testRun.Spec.ServiceAccountName,
					//                                       Tolerations:        tolerations,
					//                                       Affinity:           &affinity,
				},
			},
		},
	}

}

// SetupWithManager sets up the controller with the Manager.
func (r *TestRunReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&graphnodev1alpha1.TestRun{}).
		Complete(r)
}
