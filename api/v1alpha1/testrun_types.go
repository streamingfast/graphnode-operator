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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TestRunSpec defines the desired state of TestRun
type TestRunSpec struct {
	GitBranch          string `json:"git_branch"`           //ex: master
	GitRepo            string `json:"git_repo"`             //ex: https://github.com/graph-protocol/graph-node
	StopBlock          uint64 `json:"stop_block"`           //ex: 123456
	IPFSAddr           string `json:"ipfs_addr"`            // ex: ipfs-v3.mydomain.com
	EthereumRPCAddress string `json:"ethereum_rpc_address"` // ex: mainnet:https://something/someKey

	GraphnodeDockerImage string `json:"graphnode_docker_image"` // ex: rust:1.52
	StorageClassName     string `json:"storage_class_name"`     // ex: standard
	ServiceAccountName   string `json:"service_account_name"`   // ex: some-data-reader
	TarballsURL          string `json:"tarballs_url"`           // ex: gs://mybucket/tarballs
	TestOutputURL        string `json:"output_url"`             // ex: gs://mybucket/tests_results

	PostgresTarballerDockerImage string `json:"postgres_tarballer_docker_image"` // ex: dfuse/tarballer:0.0.6
	PostgresDBName               string `json:"postgres_db_name"`                // ex: graph
	PostgresUser                 string `json:"postgres_user"`                   // ex: graph
	PostgresPassword             string `json:"postgres_password"`               // ex: changeme

	PGDataStorage resource.Quantity `json:"pgdata_storage"`  // ex: "20Gi"
	OutputDirSize resource.Quantity `json:"output_dir_size"` // ex: "1Gi"

	PostgresResources  corev1.ResourceRequirements `json:"postgres_resources"`  // ex: {"limit": {"memory":"10Gi", cpu: "2"  }, "request": {...}}
	GraphnodeResources corev1.ResourceRequirements `json:"graphnode_resources"` // see above
}

// TestRunStatus defines the observed state of TestRun
type TestRunStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// TestRun is the Schema for the testruns API
type TestRun struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TestRunSpec   `json:"spec,omitempty"`
	Status TestRunStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TestRunList contains a list of TestRun
type TestRunList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TestRun `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TestRun{}, &TestRunList{})
}
