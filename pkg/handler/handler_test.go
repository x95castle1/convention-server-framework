package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/x95castle1/convention-server-framework/pkg/convention"

	webhookv1alpha1 "github.com/vmware-tanzu/cartographer-conventions/webhook/api/v1alpha1"

	"github.com/google/go-cmp/cmp"
	"github.com/vmware-tanzu/cartographer-conventions/webhook"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	imageDefault            = "sample-accelerators/tanzu-java-web-app"
	Prefix                  = "x95castle1.org"
	ReadinessId             = Prefix + "-readiness"
	ReadinessAnnotation     = Prefix + "/readinessProbe"
	LivenessId              = Prefix + "-liveness"
	LivenessAnnotation      = Prefix + "/livenessProbe"
	StartupId               = Prefix + "-startup"
	StartUpAnnotation       = Prefix + "/startupProbe"
	ArgsId                  = Prefix + "-args"
	ArgsAnnotation          = Prefix + "/args"
	StorageId               = Prefix + "-storage"
	StorageAnnotation       = Prefix + "/storage"
	WorkloadNameId          = Prefix + "-carto-run-workload-name"
	WorkloadNameLabel       = "carto.run/workload-name"
	WorkloadNameEnvVariable = "CARTO_RUN_WORKLOAD_NAME"
)

func Test_addConventions(t *testing.T) {
	testdataPath := "testdata"
	var emptyAppliedConventions []string
	l := zap.NewNop().Sugar()
	type args struct {
		logger   *zap.SugaredLogger
		template *corev1.PodTemplateSpec
		images   []webhook.ImageConfig
	}
	tests := []struct {
		skip               bool
		name               string
		args               args
		want               []string
		wantErr            bool
		validateTemplate   bool
		wantedTemplateFile string
	}{
		{
			name: "no images",
			args: args{
				logger:   l,
				template: getMockTemplate(""),
				images:   make([]webhook.ImageConfig, 0),
			},
			want:    emptyAppliedConventions,
			wantErr: false,
		},
		{
			name: "readinessProbe",
			args: args{
				logger:   l,
				template: getMockTemplateWithImageAndAnnotation("", fmt.Sprintf("%s/readinessProbe", Prefix), "{\"exec\":{\"command\":[\"cat\",\"/tmp/healthy\"]},\"initialDelaySeconds\":5,\"periodSeconds\":5}"),
				images: []webhook.ImageConfig{
					{
						Image: imageDefault,
						BOMs: []webhookv1alpha1.BOM{
							{
								Name: "cnb-app:dependencies",
								Raw:  getFileBytes(testdataPath + "/boms/bom.cdx.not_springboot.json"),
							},
						},
					},
				},
			},
			want:               []string{fmt.Sprintf("%s-readiness", Prefix)},
			wantErr:            false,
			validateTemplate:   true,
			wantedTemplateFile: "readinessProbe.json",
		},
		{
			name: "livenessProbe",
			args: args{
				logger:   l,
				template: getMockTemplateWithImageAndAnnotation("", fmt.Sprintf("%s/livenessProbe", Prefix), "{\"exec\":{\"command\":[\"cat\",\"/tmp/healthy\"]},\"initialDelaySeconds\":5,\"periodSeconds\":5}"),
				images: []webhook.ImageConfig{
					{
						Image: imageDefault,
						BOMs: []webhookv1alpha1.BOM{
							{
								Name: "cnb-app:dependencies",
								Raw:  getFileBytes(testdataPath + "/boms/bom.cdx.not_springboot.json"),
							},
						},
					},
				},
			},
			want:               []string{fmt.Sprintf("%s-liveness", Prefix)},
			wantErr:            false,
			validateTemplate:   true,
			wantedTemplateFile: "livenessProbe.json",
		},
		{
			name: "startupProbe",
			args: args{
				logger:   l,
				template: getMockTemplateWithImageAndAnnotation("", fmt.Sprintf("%s/startupProbe", Prefix), "{\"exec\":{\"command\":[\"cat\",\"/tmp/healthy\"]},\"initialDelaySeconds\":5,\"periodSeconds\":5}"),
				images: []webhook.ImageConfig{
					{
						Image: imageDefault,
						BOMs: []webhookv1alpha1.BOM{
							{
								Name: "cnb-app:dependencies",
								Raw:  getFileBytes(testdataPath + "/boms/bom.cdx.not_springboot.json"),
							},
						},
					},
				},
			},
			want:               []string{fmt.Sprintf("%s-startup", Prefix)},
			wantErr:            false,
			validateTemplate:   true,
			wantedTemplateFile: "startupProbe.json",
		},
		{
			name: "cartoRunWorkloadName",
			args: args{
				logger:   l,
				template: getMockTemplateWithImageAndLabel("", "carto.run/workload-name", "my-workload"),
				images: []webhook.ImageConfig{
					{
						Image: imageDefault,
						BOMs: []webhookv1alpha1.BOM{
							{
								Name: "cnb-app:dependencies",
								Raw:  getFileBytes(testdataPath + "/boms/bom.cdx.not_springboot.json"),
							},
						},
					},
				},
			},
			want:               []string{fmt.Sprintf("%s-carto-run-workload-name", Prefix)},
			wantErr:            false,
			validateTemplate:   true,
			wantedTemplateFile: "cartoRunWorkloadName.json",
		},
		{
			name: "args",
			args: args{
				logger:   l,
				template: getMockTemplateWithImageAndAnnotation("", fmt.Sprintf("%s/args", Prefix), "[\"one\",\"two\",\"three\"]"),
				images: []webhook.ImageConfig{
					{
						Image: imageDefault,
						BOMs: []webhookv1alpha1.BOM{
							{
								Name: "cnb-app:dependencies",
								Raw:  getFileBytes(testdataPath + "/boms/bom.cdx.not_springboot.json"),
							},
						},
					},
				},
			},
			want:               []string{fmt.Sprintf("%s-args", Prefix)},
			wantErr:            false,
			validateTemplate:   true,
			wantedTemplateFile: "args.json",
		},
		{
			name: "storage",
			args: args{
				logger:   l,
				template: getMockTemplateWithImageAndAnnotation("", fmt.Sprintf("%s/storage", Prefix), "{\"volumeMounts\":[{\"mountPath\":\"/test\",\"name\":\"test\"}],\"volumes\":[{\"name\":\"test\",\"emptyDir\":{}}]}"),
				images: []webhook.ImageConfig{
					{
						Image: imageDefault,
						BOMs: []webhookv1alpha1.BOM{
							{
								Name: "cnb-app:dependencies",
								Raw:  getFileBytes(testdataPath + "/boms/bom.cdx.not_springboot.json"),
							},
						},
					},
				},
			},
			want:               []string{fmt.Sprintf("%s-storage", Prefix)},
			wantErr:            false,
			validateTemplate:   true,
			wantedTemplateFile: "storage.json",
		},
	}

	for _, tt := range tests {
		if tt.skip {
			t.Logf("skipping test %s", tt.name)
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := AddConventions(tt.args.logger, tt.args.template, tt.args.images, Conventions)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddConventions() = unwanted error: %v", err)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("%s() = (-expected, +actual): %s", tt.name, diff)
			}

			if tt.validateTemplate {
				if tt.wantedTemplateFile == "" {
					t.Errorf("AddConventions(): [%s] Template JSON file not provided", tt.name)
				} else {
					pts := &corev1.PodTemplateSpec{}
					tf := getFileBytes(testdataPath + "/templates/" + tt.wantedTemplateFile)

					err = json.Unmarshal(tf, pts)
					if err != nil {
						t.Errorf("AddConventions(): [%s] %v", tt.name, err)
					}
					if diff := cmp.Diff(tt.args.template, pts); diff != "" {
						t.Errorf("%s() = (-expected, +actual): %s", tt.name, diff)
					}
				}
			}
		})
	}
}

func getFileBytes(file string) []byte {
	b, _ := os.ReadFile(file)
	return b
}

func getMockTemplate(img string) *corev1.PodTemplateSpec {
	if img == "" {
		img = imageDefault
	}
	return &corev1.PodTemplateSpec{
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "workload",
					Image: img,
					Ports: make([]corev1.ContainerPort, 0),
				},
			},
		},
	}
}

func getMockTemplateWithImageAndAnnotation(img string, key string, value string) *corev1.PodTemplateSpec {
	if img == "" {
		img = imageDefault
	}
	return &corev1.PodTemplateSpec{
		ObjectMeta: v1.ObjectMeta{Annotations: map[string]string{key: value}},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "workload",
					Image: img,
					Ports: make([]corev1.ContainerPort, 0),
				},
			},
		},
	}
}

func getMockTemplateWithImageAndLabel(img string, key string, value string) *corev1.PodTemplateSpec {
	if img == "" {
		img = imageDefault
	}
	return &corev1.PodTemplateSpec{
		ObjectMeta: v1.ObjectMeta{Labels: map[string]string{key: value}},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "workload",
					Image: img,
					Ports: make([]corev1.ContainerPort, 0),
				},
			},
		},
	}
}

var Conventions = []convention.Convention{
	&convention.BasicConvention{
		Id: ReadinessId,
		Applicable: func(ctx context.Context, target *corev1.PodTemplateSpec, metadata convention.ImageMetadata) bool {
			return getAnnotation(target, ReadinessAnnotation) != ""
		},
		Apply: func(ctx context.Context, target *corev1.PodTemplateSpec, containerIdx int, metadata convention.ImageMetadata, imageName string) error {
			readinessProbe := getAnnotation(target, ReadinessAnnotation)

			for i := range target.Spec.Containers {
				c := &target.Spec.Containers[i]

				if c.ReadinessProbe == nil {
					p, err := getProbe(readinessProbe)
					if err != nil {
						return err
					}
					log.Printf("Adding ReadinessProbe %+v", p)
					c.ReadinessProbe = p
				}
			}
			return nil
		},
	},

	&convention.BasicConvention{
		Id: LivenessId,
		Applicable: func(ctx context.Context, target *corev1.PodTemplateSpec, metadata convention.ImageMetadata) bool {
			return getAnnotation(target, LivenessAnnotation) != ""
		},
		Apply: func(ctx context.Context, target *corev1.PodTemplateSpec, containerIdx int, metadata convention.ImageMetadata, imageName string) error {
			livenessProbe := getAnnotation(target, LivenessAnnotation)

			for i := range target.Spec.Containers {
				c := &target.Spec.Containers[i]

				if c.LivenessProbe == nil {
					p, err := getProbe(livenessProbe)
					if err != nil {
						return err
					}
					log.Printf("Adding LivenessProbe %+v", p)
					c.LivenessProbe = p
				}
			}
			return nil
		},
	},

	&convention.BasicConvention{
		Id: StartupId,
		Applicable: func(ctx context.Context, target *corev1.PodTemplateSpec, metadata convention.ImageMetadata) bool {
			return getAnnotation(target, StartUpAnnotation) != ""
		},
		Apply: func(ctx context.Context, target *corev1.PodTemplateSpec, containerIdx int, metadata convention.ImageMetadata, imageName string) error {
			startupProbe := getAnnotation(target, StartUpAnnotation)

			for i := range target.Spec.Containers {
				c := &target.Spec.Containers[i]

				if c.StartupProbe == nil {
					p, err := getProbe(startupProbe)
					if err != nil {
						return err
					}
					log.Printf("Adding StartupProbe %+v", p)
					c.StartupProbe = p
				}
			}
			return nil
		},
	},

	&convention.BasicConvention{
		Id: WorkloadNameId,
		Applicable: func(ctx context.Context, target *corev1.PodTemplateSpec, metadata convention.ImageMetadata) bool {
			return getLabel(target, WorkloadNameLabel) != ""
		},
		Apply: func(ctx context.Context, target *corev1.PodTemplateSpec, containerIdx int, metadata convention.ImageMetadata, imageName string) error {
			value := getLabel(target, WorkloadNameLabel)

			for i := range target.Spec.Containers {
				c := &target.Spec.Containers[i]
				addEnvVar(c, corev1.EnvVar{
					Name:  WorkloadNameEnvVariable,
					Value: value,
				})
			}

			return nil
		},
	},

	&convention.BasicConvention{
		Id: ArgsId,
		Applicable: func(ctx context.Context, target *corev1.PodTemplateSpec, metadata convention.ImageMetadata) bool {
			return getAnnotation(target, ArgsAnnotation) != ""
		},
		Apply: func(ctx context.Context, target *corev1.PodTemplateSpec, containerIdx int, metadata convention.ImageMetadata, imageName string) error {
			arguments := getAnnotation(target, ArgsAnnotation)

			for i := range target.Spec.Containers {
				c := &target.Spec.Containers[i]

				a, err := getArguments(arguments)
				if err != nil {
					return err
				}

				c.Args = append(c.Args, a...)
			}
			return nil
		},
	},

	&convention.BasicConvention{
		Id: StorageId,
		Applicable: func(ctx context.Context, target *corev1.PodTemplateSpec, metadata convention.ImageMetadata) bool {
			return getAnnotation(target, StorageAnnotation) != ""
		},
		Apply: func(ctx context.Context, target *corev1.PodTemplateSpec, containerIdx int, metadata convention.ImageMetadata, imageName string) error {
			storage := getAnnotation(target, StorageAnnotation)

			for i := range target.Spec.Containers {
				c := &target.Spec.Containers[i]

				s, err := getStorage(storage)
				if err != nil {
					return err
				}

				c.VolumeMounts = append(c.VolumeMounts, s.VolumeMounts...)
				target.Spec.Volumes = append(target.Spec.Volumes, s.Volumes...)
			}
			return nil
		},
	},
}

// getArguments parse the arguments into a string array
func getArguments(arguments string) ([]string, error) {
	var a []string
	err := json.Unmarshal([]byte(arguments), &a)
	return a, err
}

// getStorage parse the arguments into a storage struct
func getStorage(arguments string) (Storage, error) {
	var s Storage
	err := json.Unmarshal([]byte(arguments), &s)
	return s, err
}

// getAnnotation gets the annotation on PodTemplateSpec
func getAnnotation(pts *corev1.PodTemplateSpec, key string) string {
	if pts.Annotations == nil || len(pts.Annotations[key]) == 0 {
		return ""
	}
	return pts.Annotations[key]
}

// getLabel gets the label on PodTemplateSpec
func getLabel(pts *corev1.PodTemplateSpec, key string) string {
	if pts.Labels == nil || len(pts.Labels[key]) == 0 {
		return ""
	}
	return pts.Labels[key]
}

func getProbe(config string) (*corev1.Probe, error) {
	probe := corev1.Probe{}
	err := json.Unmarshal([]byte(config), &probe)
	return &probe, err
}

func addEnvVar(container *corev1.Container, envvar corev1.EnvVar) bool {
	for _, e := range container.Env {
		if e.Name == envvar.Name {
			return false
		}
	}
	container.Env = append(container.Env, envvar)
	return true
}

type Storage struct {
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"`
	Volumes      []corev1.Volume      `json:"volumes,omitempty"`
}
