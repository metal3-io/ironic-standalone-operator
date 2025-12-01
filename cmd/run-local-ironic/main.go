package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/yaml"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
	"github.com/metal3-io/ironic-standalone-operator/pkg/ironic"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("run-local-ironic")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(metal3api.AddToScheme(scheme))
}

func main() {
	var inputFile string
	var outputFile string
	var verbose bool
	var generateOnly bool
	var tearDown bool

	flag.StringVar(&inputFile, "input", "", "YAML file containing Ironic resource and optional secrets (required)")
	flag.StringVar(&outputFile, "output", "", "Output file for generated manifests (required when using --generate-only)")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	flag.BoolVar(&generateOnly, "generate-only", false, "Only generate manifests, skip running podman kube play")
	flag.BoolVar(&tearDown, "down", false, "Tear down instead of creating")

	opts := zap.Options{
		Development: verbose,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	if inputFile == "" {
		setupLog.Error(errors.New("input file is required"), "missing required flag")
		flag.Usage()
		os.Exit(1)
	}

	// Validate output requirements
	if generateOnly && outputFile == "" {
		setupLog.Error(errors.New("output file is required when using --generate-only"), "missing required flag")
		flag.Usage()
		os.Exit(1)
	}

	if tearDown && generateOnly {
		setupLog.Error(errors.New("--down is not compatible with --generate-only"), "incompatible flags")
		flag.Usage()
		os.Exit(1)
	}

	// Generate temporary file if output not provided and not in generate-only mode
	if outputFile == "" {
		tempFile, err := os.CreateTemp("", "ironic-manifests-*.yaml")
		if err != nil {
			setupLog.Error(err, "failed to create temporary file")
			os.Exit(1)
		}
		outputFile = tempFile.Name()
		tempFile.Close() // Close immediately, we just need the name
	}

	versionInfo, err := ironic.NewVersionInfo(metal3api.Images{}, "", "")
	if err != nil {
		// This cannot happen in reality
		setupLog.Error(err, "invalid ironic version")
		os.Exit(1)
	}

	setupLog.V(1).Info("starting run-local-ironic", "input", inputFile, "output", outputFile)

	if err := runLocalIronic(inputFile, outputFile, versionInfo, generateOnly, tearDown); err != nil {
		setupLog.Error(err, "failed to run local ironic")
		os.Exit(1)
	}
}

func runLocalIronic(inputFile, outputFile string, versionInfo ironic.VersionInfo, generateOnly, tearDown bool) error {
	resources, err := ironic.ParseLocalIronic(inputFile, scheme)
	if err != nil {
		return fmt.Errorf("failed to parse input file: %w", err)
	}

	outputDir := filepath.Dir(outputFile)
	if err = os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	cctx := ironic.ControllerContext{
		VersionInfo: versionInfo,
		Logger:      setupLog,
	}

	manifests, err := ironic.GenerateLocalManifests(cctx, resources)
	if err != nil {
		return fmt.Errorf("failed to generate manifests: %w", err)
	}

	if err := writeManifestsToYAML(manifests, outputFile, scheme); err != nil {
		return fmt.Errorf("failed to write manifests: %w", err)
	}

	setupLog.V(1).Info("manifests generated", "file", outputFile)

	if !generateOnly {
		// Run podman kube play
		if err := runPodmanKubePlay(outputFile, tearDown); err != nil {
			return fmt.Errorf("failed to run podman kube play: %w", err)
		}
		if tearDown {
			setupLog.Info("ironic deployment stopped successfully")
		} else {
			setupLog.Info("ironic deployment started successfully")
		}
	} else {
		setupLog.Info("manifest generation complete, skipping podman kube play")
	}

	return nil
}

func writeManifestsToYAML(manifests []runtime.Object, filename string, scheme *runtime.Scheme) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	codecFactory := serializer.NewCodecFactory(scheme)
	// FIXME(dtantsur): I wonder what I should use instead of the legacy one
	encoder := codecFactory.LegacyCodec(corev1.SchemeGroupVersion, appsv1.SchemeGroupVersion)

	for i, manifest := range manifests {
		if i > 0 {
			if _, err := file.WriteString("---\n"); err != nil {
				return err
			}
		}

		// First encode to JSON
		jsonData, err := runtime.Encode(encoder, manifest)
		if err != nil {
			return err
		}

		// Then convert JSON to YAML
		yamlData, err := yaml.JSONToYAML(jsonData)
		if err != nil {
			return err
		}

		if _, err := file.Write(yamlData); err != nil {
			return err
		}
	}

	return nil
}

func runPodmanKubePlay(manifestFile string, tearDown bool) error {
	var cmd *exec.Cmd
	if tearDown {
		cmd = exec.CommandContext(context.TODO(), "podman", "kube", "play", "--down", manifestFile)
	} else {
		cmd = exec.CommandContext(context.TODO(), "podman", "kube", "play", "--network=host", manifestFile)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	setupLog.V(1).Info("running podman kube play", "args", cmd.Args)
	return cmd.Run()
}
