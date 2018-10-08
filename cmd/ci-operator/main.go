package main

import (
	"context"
	"crypto/sha256"
	"encoding/base32"
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/openshift/ci-operator/pkg/defaults"

	coreapi "k8s.io/api/core/v1"
	rbacapi "k8s.io/api/rbac/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	coreclientset "k8s.io/client-go/kubernetes/typed/core/v1"
	rbacclientset "k8s.io/client-go/kubernetes/typed/rbac/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/retry"

	"github.com/ghodss/yaml"

	imageapi "github.com/openshift/api/image/v1"
	projectapi "github.com/openshift/api/project/v1"
	templateapi "github.com/openshift/api/template/v1"
	imageclientset "github.com/openshift/client-go/image/clientset/versioned/typed/image/v1"
	"github.com/openshift/client-go/project/clientset/versioned"
	templatescheme "github.com/openshift/client-go/template/clientset/versioned/scheme"

	"github.com/openshift/ci-operator/pkg/api"
	"github.com/openshift/ci-operator/pkg/interrupt"
	"github.com/openshift/ci-operator/pkg/junit"
	"github.com/openshift/ci-operator/pkg/steps"
)

const usage = `Orchestrate multi-stage image-based builds

The ci-operator reads a declarative configuration YAML file and executes a set of build
steps on an OpenShift cluster for image-based components. By default, all steps are run,
but a caller may select one or more targets (image names or test names) to limit to only
steps that those targets depend on. The build creates a new project to run the builds in
and can automatically clean up the project when the build completes.

ci-operator leverages declarative OpenShift builds and images to reuse previously compiled
artifacts. It makes building multiple images that share one or more common base layers
simple as well as running tests that depend on those images.

Since the command is intended for use in CI environments it requires an input environment
variable called the JOB_SPEC that defines the GitHub project to execute and the commit,
branch, and any PRs to merge onto the branch. See the kubernetes/test-infra project for
a description of JOB_SPEC.

The inputs of the build (source code, tagged images, configuration) are combined to form
a consistent name for the target namespace that will change if any of the inputs change.
This allows multiple test jobs to share common artifacts and still perform retries.

The standard build steps are designed for simple command-line actions (like invoking
"make test") but can be extended by passing one or more templates via the --template flag.
The name of the template defines the stage and the template must contain at least one
pod. The parameters passed to the template are the current process environment and a set
of dynamic parameters that are inferred from previous steps. These parameters are:

  NAMESPACE
    The namespace generated by the operator for the given inputs or the value of
    --namespace.

  IMAGE_FORMAT
    A string that points to the public image repository URL of the image stream(s)
    created by the tag step. Example:

      registry.svc.ci.openshift.org/ci-op-9o8bacu/stable:${component}

    Will cause the template to depend on all image builds.

  IMAGE_<component>
    The public image repository URL for an output image. If specified the template
    will depend on the image being built.

  LOCAL_IMAGE_<component>
    The public image repository URL for an image that was built during this run but
    was not part of the output (such as pipeline cache images). If specified the
    template will depend on the image being built.

  JOB_NAME
    The job name from the JOB_SPEC

  JOB_NAME_SAFE
    The job name in a form safe for use as a Kubernetes resource name.

  JOB_NAME_HASH
    A short hash of the job name for making tasks unique.

  RPM_REPO_<org>_<repo>
    If the job creates RPMs this will be the public URL that can be used as the
		baseurl= value of an RPM repository. The value of org and repo are uppercased
		and dashes are replaced with underscores.

Dynamic environment variables are overriden by process environment variables.

Both test and template jobs can gather artifacts created by pods. Set
--artifact-dir to define the top level artifact directory, and any test task
that defines artifact_dir or template that has an "artifacts" volume mounted
into a container will have artifacts extracted after the container has completed.
Errors in artifact extraction will not cause build failures.

In CI environments the inputs to a job may be different than what a normal
development workflow would use. The --override file will override fields
defined in the config file, such as base images and the release tag configuration.

After a successful build the --promote will tag each built image (in "images")
to the image stream(s) identified by the "promotion" config, which defaults to
the same image stream as the release configuration. You may add additional
images to promote and their target names via the "additional_images" map.
`

func main() {
	flagSet := flag.NewFlagSet("", flag.ExitOnError)
	opt := bindOptions(flagSet)
	flagSet.Parse(os.Args[1:])
	if opt.verbose {
		flag.CommandLine.Set("alsologtostderr", "true")
		flag.CommandLine.Set("v", "10")
	}
	if opt.help {
		fmt.Printf(usage)
		flagSet.SetOutput(os.Stdout)
		flagSet.Usage()
		os.Exit(0)
	}

	if err := opt.Validate(); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	if err := opt.Complete(); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	if err := opt.Run(); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

type stringSlice struct {
	values []string
}

func (s *stringSlice) String() string {
	return strings.Join(s.values, string(filepath.Separator))
}

func (s *stringSlice) Set(value string) error {
	s.values = append(s.values, value)
	return nil
}

type options struct {
	configSpecPath    string
	templatePaths     stringSlice
	secretDirectories stringSlice

	targets stringSlice
	promote bool

	verbose bool
	help    bool
	dry     bool

	writeParams string
	artifactDir string

	gitRef              string
	namespace           string
	baseNamespace       string
	idleCleanupDuration time.Duration
	cleanupDuration     time.Duration

	inputHash     string
	secrets       []*coreapi.Secret
	templates     []*templateapi.Template
	configSpec    *api.ReleaseBuildConfiguration
	jobSpec       *api.JobSpec
	clusterConfig *rest.Config

	givePrAuthorAccessToNamespace bool
}

func bindOptions(flag *flag.FlagSet) *options {
	opt := &options{
		idleCleanupDuration: time.Duration(1 * time.Hour),
		cleanupDuration:     time.Duration(12 * time.Hour),
	}

	// command specific options
	flag.BoolVar(&opt.help, "h", false, "short for --help")
	flag.BoolVar(&opt.help, "help", false, "See help for this command.")
	flag.BoolVar(&opt.verbose, "v", false, "Show verbose output.")

	// what we will run
	flag.StringVar(&opt.configSpecPath, "config", "", "The configuration file. If not specified the CONFIG_SPEC environment variable will be used.")
	flag.Var(&opt.targets, "target", "One or more targets in the configuration to build. Only steps that are required for this target will be run.")
	flag.BoolVar(&opt.dry, "dry-run", opt.dry, "Print the steps that would be run and the objects that would be created without executing any steps")

	// add to the graph of things we run or create
	flag.Var(&opt.templatePaths, "template", "A set of paths to optional templates to add as stages to this job. Each template is expected to contain at least one restart=Never pod. Parameters are filled from environment or from the automatic parameters generated by the operator.")
	flag.Var(&opt.secretDirectories, "secret-dir", "One or more directories that should converted into secrets in the test namespace.")

	// the target namespace and cleanup behavior
	flag.StringVar(&opt.namespace, "namespace", "", "Namespace to create builds into, defaults to build_id from JOB_SPEC. If the string '{id}' is in this value it will be replaced with the build input hash.")
	flag.StringVar(&opt.baseNamespace, "base-namespace", "stable", "Namespace to read builds from, defaults to stable.")
	flag.DurationVar(&opt.idleCleanupDuration, "delete-when-idle", opt.idleCleanupDuration, "If no pod is running for longer than this interval, delete the namespace. Set to zero to retain the contents. Requires the namespace TTL controller to be deployed.")
	flag.DurationVar(&opt.cleanupDuration, "delete-after", opt.cleanupDuration, "If namespace exists for longer than this interval, delete the namespace. Set to zero to retain the contents. Requires the namespace TTL controller to be deployed.")

	// actions to add to the graph
	flag.BoolVar(&opt.promote, "promote", false, "When all other targets complete, publish the set of images built by this job into the release configuration.")

	// output control
	flag.StringVar(&opt.artifactDir, "artifact-dir", "", "If set grab artifacts from test and template jobs.")
	flag.StringVar(&opt.writeParams, "write-params", "", "If set write an env-compatible file with the output of the job.")

	// experimental flags
	flag.StringVar(&opt.gitRef, "git-ref", "", "Populate the job spec from this local Git reference. If JOB_SPEC is set, the refs field will be overwritten.")
	flag.BoolVar(&opt.givePrAuthorAccessToNamespace, "give-pr-author-access-to-namespace", false, "Give view access to the temporarily created namespace to the PR author.")

	return opt
}

func (o *options) Validate() error {
	return nil
}

func (o *options) Complete() error {
	// Load the standard configuration from the path or env
	var configSpec string
	if len(o.configSpecPath) > 0 {
		data, err := ioutil.ReadFile(o.configSpecPath)
		if err != nil {
			return fmt.Errorf("--config error: %v", err)
		}
		configSpec = string(data)
	} else {
		var ok bool
		configSpec, ok = os.LookupEnv("CONFIG_SPEC")
		if !ok || len(configSpec) == 0 {
			return fmt.Errorf("CONFIG_SPEC environment variable is not set or empty and no --config file was set")
		}
	}
	if err := yaml.Unmarshal([]byte(configSpec), &o.configSpec); err != nil {
		return fmt.Errorf("invalid configuration: %v\nvalue:\n%s", err, string(configSpec))
	}

	if err := o.configSpec.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %v", err)
	}

	jobSpec, err := api.ResolveSpecFromEnv()
	if err != nil {
		if len(o.gitRef) == 0 {
			return fmt.Errorf("failed to determine job spec: no --git-ref passed and failed to resolve job spec from env: %v", err)
		}
		// Failed to read $JOB_SPEC but --git-ref was passed, so try that instead
		spec, refErr := jobSpecFromGitRef(o.gitRef)
		if refErr != nil {
			return fmt.Errorf("failed to determine job spec: failed to resolve --git-ref: %v", refErr)
		}
		jobSpec = spec
	} else if len(o.gitRef) > 0 {
		// Read from $JOB_SPEC but --git-ref was also passed, so merge them
		spec, err := jobSpecFromGitRef(o.gitRef)
		if err != nil {
			return fmt.Errorf("failed to determine job spec: failed to resolve --git-ref: %v", err)
		}
		jobSpec.Refs = spec.Refs
	}
	jobSpec.BaseNamespace = o.baseNamespace
	o.jobSpec = jobSpec

	if o.dry && o.verbose {
		config, _ := yaml.Marshal(o.configSpec)
		log.Printf("Resolved configuration:\n%s", string(config))
		job, _ := json.Marshal(o.jobSpec)
		log.Printf("Resolved job spec:\n%s", string(job))
	}
	refs := o.jobSpec.Refs
	if len(refs.Pulls) > 0 {
		var pulls []string
		for _, pull := range refs.Pulls {
			pulls = append(pulls, fmt.Sprintf("#%d %s @%s", pull.Number, shorten(pull.SHA, 8), pull.Author))
		}
		log.Printf("Resolved source https://github.com/%s/%s to %s@%s, merging: %s", refs.Org, refs.Repo, refs.BaseRef, shorten(refs.BaseSHA, 8), strings.Join(pulls, ", "))
	} else {
		log.Printf("Resolved source https://github.com/%s/%s to %s@%s", refs.Org, refs.Repo, refs.BaseRef, shorten(refs.BaseSHA, 8))
	}

	for _, path := range o.secretDirectories.values {
		secret := &coreapi.Secret{Data: make(map[string][]byte)}
		secret.Type = coreapi.SecretTypeOpaque
		secret.Name = filepath.Base(path)
		files, err := ioutil.ReadDir(path)
		if err != nil {
			return fmt.Errorf("could not read dir %s for secret: %v", path, err)
		}
		for _, f := range files {
			if f.IsDir() {
				continue
			}
			path := filepath.Join(path, f.Name())
			// if the file is a broken symlink or a symlink to a dir, skip it
			if fi, err := os.Stat(path); err != nil || fi.IsDir() {
				continue
			}
			secret.Data[f.Name()], err = ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("could not read file %s for secret: %v", path, err)
			}
		}
		o.secrets = append(o.secrets, secret)
	}

	for _, path := range o.templatePaths.values {
		contents, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("could not read dir %s for template: %v", path, err)
		}
		obj, gvk, err := templatescheme.Codecs.UniversalDeserializer().Decode(contents, nil, nil)
		if err != nil {
			return fmt.Errorf("unable to parse template %s: %v", path, err)
		}
		template, ok := obj.(*templateapi.Template)
		if !ok {
			return fmt.Errorf("%s is not a template: %v", path, gvk)
		}
		if len(template.Name) == 0 {
			template.Name = filepath.Base(path)
			template.Name = strings.TrimSuffix(template.Name, filepath.Ext(template.Name))
		}
		o.templates = append(o.templates, template)
	}

	clusterConfig, err := loadClusterConfig()
	if err != nil {
		return fmt.Errorf("failed to load cluster config: %v", err)
	}
	o.clusterConfig = clusterConfig

	return nil
}

func (o *options) Run() error {
	start := time.Now()
	defer func() {
		log.Printf("Ran for %s", time.Now().Sub(start).Truncate(time.Second))
	}()

	// load the graph from the configuration
	buildSteps, postSteps, err := defaults.FromConfig(o.configSpec, o.jobSpec, o.templates, o.writeParams, o.artifactDir, o.promote, o.clusterConfig, o.targets.values)
	if err != nil {
		return fmt.Errorf("failed to generate steps from config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	handler := func(s os.Signal) {
		if o.dry {
			os.Exit(0)
		}
		log.Printf("error: Process interrupted with signal %s, exiting in 2s ...", s)
		cancel()
		time.Sleep(2 * time.Second)
		os.Exit(1)
	}

	return interrupt.New(handler).Run(func() error {
		// Before we create the namespace, we need to ensure all inputs to the graph
		// have been resolved. We must run this step before we resolve the partial
		// graph or otherwise two jobs with different targets would create different
		// artifact caches.
		if err := o.resolveInputs(ctx, buildSteps); err != nil {
			return fmt.Errorf("could not resolve inputs: %v", err)
		}

		// convert the full graph into the subset we must run
		nodes, err := api.BuildPartialGraph(buildSteps, o.targets.values)
		if err != nil {
			return fmt.Errorf("could not build execution graph: %v", err)
		}

		if err := printExecutionOrder(nodes); err != nil {
			return fmt.Errorf("could not print execution order: %v", err)
		}

		// initialize the namespace if necessary and create any resources that must
		// exist prior to execution
		if err := o.initializeNamespace(); err != nil {
			return fmt.Errorf("could not initialize namespace: %v", err)
		}

		client, err := coreclientset.NewForConfig(o.clusterConfig)
		if err != nil {
			return fmt.Errorf("could not get core client for cluster config: %v", err)
		}
		eventRecorder := eventRecorder(client, o.namespace)
		runtimeObject := &coreapi.ObjectReference{Namespace: o.namespace}

		if !o.dry {
			eventRecorder.Event(runtimeObject, coreapi.EventTypeNormal, "CiJobStarted", eventJobDescription(o.jobSpec, o.namespace))
		}
		// execute the graph
		suites, err := steps.Run(ctx, nodes, o.dry)
		if err := o.writeJUnit(suites, "operator"); err != nil {
			log.Printf("warning: Unable to write JUnit result: %v", err)
		}
		if err != nil {
			if !o.dry {
				eventRecorder.Event(runtimeObject, coreapi.EventTypeWarning, "CiJobFailed", eventJobDescription(o.jobSpec, o.namespace))
				time.Sleep(time.Second)
			}
			return fmt.Errorf("could not run steps: %v", err)
		}

		for _, step := range postSteps {
			if err := step.Run(ctx, o.dry); err != nil {
				if !o.dry {
					eventRecorder.Event(runtimeObject, coreapi.EventTypeWarning, "PostStepFailed",
						fmt.Sprintf("Post step %s failed while %s", step.Name(), eventJobDescription(o.jobSpec, o.namespace)))
					time.Sleep(time.Second)
				}
				return fmt.Errorf("could not run post step %s: %v", step.Name(), err)
			}
		}

		if !o.dry {
			eventRecorder.Event(runtimeObject, coreapi.EventTypeNormal, "CiJobSucceeded", eventJobDescription(o.jobSpec, o.namespace))
			time.Sleep(time.Second)
		}

		return nil
	})
}

// loadClusterConfig loads connection configuration
// for the cluster we're deploying to. We prefer to
// use in-cluster configuration if possible, but will
// fall back to using default rules otherwise.
func loadClusterConfig() (*rest.Config, error) {
	clusterConfig, err := rest.InClusterConfig()
	if err == nil {
		return clusterConfig, nil
	}

	credentials, err := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	if err != nil {
		return nil, fmt.Errorf("could not load credentials from config: %v", err)
	}

	clusterConfig, err = clientcmd.NewDefaultClientConfig(*credentials, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("could not load client configuration: %v", err)
	}
	return clusterConfig, nil
}

func (o *options) resolveInputs(ctx context.Context, steps []api.Step) error {
	var inputs api.InputDefinition
	for _, step := range steps {
		definition, err := step.Inputs(ctx, o.dry)
		if err != nil {
			return fmt.Errorf("could not determine inputs for step %s: %v", step.Name(), err)
		}
		inputs = append(inputs, definition...)
	}

	// a change in the config for the build changes the output
	configSpec, err := yaml.Marshal(o.configSpec)
	if err != nil {
		panic(err)
	}
	inputs = append(inputs, string(configSpec))

	o.inputHash = inputHash(inputs)

	// input hash is unique for a given job definition and input refs
	if len(o.namespace) == 0 {
		o.namespace = "ci-op-{id}"
	}
	o.namespace = strings.Replace(o.namespace, "{id}", o.inputHash, -1)
	// TODO: instead of mutating this here, we should pass the parts of graph execution that are resolved
	// after the graph is created but before it is run down into the run step.
	o.jobSpec.Namespace = o.namespace
	log.Printf("Using namespace %s", o.namespace)

	return nil
}

func (o *options) initializeNamespace() error {
	if o.dry {
		return nil
	}
	projectGetter, err := versioned.NewForConfig(o.clusterConfig)
	if err != nil {
		return fmt.Errorf("could not get project client for cluster config: %v", err)
	}

	log.Printf("Creating namespace %s", o.namespace)
	retries := 5
	for {
		project, err := projectGetter.ProjectV1().ProjectRequests().Create(&projectapi.ProjectRequest{
			ObjectMeta: meta.ObjectMeta{
				Name: o.namespace,
			},
			DisplayName: fmt.Sprintf("%s - %s", o.namespace, o.jobSpec.Job),
			Description: jobDescription(o.jobSpec, o.configSpec),
		})
		if err != nil && !kerrors.IsAlreadyExists(err) {
			return fmt.Errorf("could not set up namespace for test: %v", err)
		}
		if err != nil {
			project, err = projectGetter.ProjectV1().Projects().Get(o.namespace, meta.GetOptions{})
			if err != nil {
				if kerrors.IsNotFound(err) {
					continue
				}
				// wait a few seconds for auth caches to catch up
				if kerrors.IsForbidden(err) && retries > 0 {
					retries--
					time.Sleep(time.Second)
					continue
				}
				return fmt.Errorf("cannot retrieve test namespace: %v", err)
			}
		}
		if project.Status.Phase == coreapi.NamespaceTerminating {
			log.Println("Waiting for namespace to finish terminating before creating another")
			time.Sleep(3 * time.Second)
			continue
		}
		break
	}

	if o.givePrAuthorAccessToNamespace {
		// Generate rolebinding for all the PR Authors.
		rbacClient, err := rbacclientset.NewForConfig(o.clusterConfig)
		if err != nil {
			return fmt.Errorf("could not get RBAC client for cluster config: %v", err)
		}
		for _, pull := range o.jobSpec.Refs.Pulls {
			log.Printf("Creating rolebinding for user %s in namespace %s", pull.Author, o.namespace)
			if _, err := rbacClient.RoleBindings(o.namespace).Create(&rbacapi.RoleBinding{
				ObjectMeta: meta.ObjectMeta{
					Name:      "ci-op-author-access",
					Namespace: o.namespace,
				},
				Subjects: []rbacapi.Subject{{Kind: "User", Name: pull.Author}},
				RoleRef: rbacapi.RoleRef{
					Kind: "ClusterRole",
					Name: "admin",
				},
			}); err != nil && !kerrors.IsAlreadyExists(err) {
				return fmt.Errorf("could not create role binding for: %v", err)
			}
		}
	}

	client, err := coreclientset.NewForConfig(o.clusterConfig)
	if err != nil {
		return fmt.Errorf("could not get core client for cluster config: %v", err)
	}

	updates := map[string]string{}
	if o.idleCleanupDuration > 0 {
		log.Printf("Setting a soft TTL of %s for the namespace\n", o.idleCleanupDuration.String())
		updates["ci.openshift.io/ttl.soft"] = o.idleCleanupDuration.String()
	}

	if o.cleanupDuration > 0 {
		log.Printf("Setting a hard TTL of %s for the namespace\n", o.cleanupDuration.String())
		updates["ci.openshift.io/ttl.hard"] = o.cleanupDuration.String()
	}

	if len(updates) > 0 {
		if err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			ns, err := client.Namespaces().Get(o.namespace, meta.GetOptions{})
			if err != nil {
				return err
			}

			if ns.Annotations == nil {
				ns.Annotations = make(map[string]string)
			}
			for key, value := range updates {
				ns.ObjectMeta.Annotations[key] = value
			}

			_, updateErr := client.Namespaces().Update(ns)
			if kerrors.IsForbidden(updateErr) {
				log.Printf("warning: Could not mark the namespace to be deleted later because you do not have permission to update the namespace (details: %v)", updateErr)
				return nil
			}
			return updateErr
		}); err != nil {
			return fmt.Errorf("could not update namespace to add TTLs: %v", err)
		}
	}

	log.Printf("Setting up pipeline imagestream for the test")
	imageGetter, err := imageclientset.NewForConfig(o.clusterConfig)
	if err != nil {
		return fmt.Errorf("could not get image client for cluster config: %v", err)
	}

	// create the image stream or read it to get its uid
	is, err := imageGetter.ImageStreams(o.jobSpec.Namespace).Create(&imageapi.ImageStream{
		ObjectMeta: meta.ObjectMeta{
			Namespace: o.jobSpec.Namespace,
			Name:      api.PipelineImageStream,
		},
		Spec: imageapi.ImageStreamSpec{
			// pipeline:* will now be directly referenceable
			LookupPolicy: imageapi.ImageLookupPolicy{Local: true},
		},
	})
	if err != nil {
		if !kerrors.IsAlreadyExists(err) {
			return fmt.Errorf("could not set up pipeline imagestream for test: %v", err)
		}
		is, _ = imageGetter.ImageStreams(o.jobSpec.Namespace).Get(api.PipelineImageStream, meta.GetOptions{})
	}
	if is != nil {
		isTrue := true
		o.jobSpec.SetOwner(&meta.OwnerReference{
			APIVersion: "image.openshift.io/v1",
			Kind:       "ImageStream",
			Name:       api.PipelineImageStream,
			UID:        is.UID,
			Controller: &isTrue,
		})
	}

	if len(o.secrets) > 0 {
		log.Printf("Populating secrets for test")
	}
	for _, secret := range o.secrets {
		_, err := client.Secrets(o.namespace).Create(secret)
		if kerrors.IsAlreadyExists(err) {
			existing, err := client.Secrets(o.namespace).Get(secret.Name, meta.GetOptions{})
			if err != nil {
				return fmt.Errorf("could not retrieve secret %s: %v", secret.Name, err)
			}
			for k, v := range secret.Data {
				existing.Data[k] = v
			}
			if _, err := client.Secrets(o.namespace).Update(existing); err != nil {
				return fmt.Errorf("could not update secret %s: %v", secret.Name, err)
			}
			log.Printf("Updated secret %s", secret.Name)
			continue
		}
		if err != nil {
			return fmt.Errorf("could not create secret %s: %v", secret.Name, err)
		}
		log.Printf("Created secret %s", secret.Name)
	}
	return nil
}

func (o *options) writeJUnit(suites *junit.TestSuites, name string) error {
	if len(o.artifactDir) == 0 || suites == nil {
		return nil
	}
	suites.Suites[0].Name = name
	out, err := xml.MarshalIndent(suites, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal jUnit XML: %v", err)
	}
	return ioutil.WriteFile(filepath.Join(o.artifactDir, fmt.Sprintf("junit_%s.xml", name)), out, 0640)
}

// oneWayEncoding can be used to encode hex to a 62-character set (0 and 1 are duplicates) for use in
// short display names that are safe for use in kubernetes as resource names.
var oneWayNameEncoding = base32.NewEncoding("bcdfghijklmnpqrstvwxyz0123456789").WithPadding(base32.NoPadding)

// inputHash returns a string that hashes the unique parts of the input to avoid collisions.
func inputHash(inputs api.InputDefinition) string {
	hash := sha256.New()

	// the inputs form a part of the hash
	for _, s := range inputs {
		hash.Write([]byte(s))
	}

	// Object names can't be too long so we truncate
	// the hash. This increases chances of collision
	// but we can tolerate it as our input space is
	// tiny.
	return oneWayNameEncoding.EncodeToString(hash.Sum(nil)[:5])
}

// eventJobDescription returns a string representing the pull requests and authors description, to be used in events.
func eventJobDescription(jobSpec *api.JobSpec, namespace string) string {
	pulls := []string{}
	authors := []string{}

	if len(jobSpec.Refs.Pulls) == 1 {
		pull := jobSpec.Refs.Pulls[0]
		return fmt.Sprintf("Running job %s for PR https://github.com/%s/%s/pull/%d in namespace %s from author %s",
			jobSpec.Job, jobSpec.Refs.Org, jobSpec.Refs.Repo, pull.Number, namespace, pull.Author)
	}
	for _, pull := range jobSpec.Refs.Pulls {
		pulls = append(pulls, fmt.Sprintf("https://github.com/%s/%s/pull/%d", jobSpec.Refs.Org, jobSpec.Refs.Repo, pull.Number))
		authors = append(authors, pull.Author)
	}
	return fmt.Sprintf("Running job %s for PRs (%s) in namespace %s from authors (%s)",
		jobSpec.Job, strings.Join(pulls, ", "), namespace, strings.Join(authors, ", "))
}

// jobDescription returns a string representing the job's description.
func jobDescription(job *api.JobSpec, config *api.ReleaseBuildConfiguration) string {
	var links []string
	for _, pull := range job.Refs.Pulls {
		links = append(links, fmt.Sprintf("https://github.com/%s/%s/pull/%d - %s", job.Refs.Org, job.Refs.Repo, pull.Number, pull.Author))
	}
	if len(links) > 0 {
		return fmt.Sprintf("%s on https://github.com/%s/%s\n\n%s", job.Job, job.Refs.Org, job.Refs.Repo, strings.Join(links, "\n"))
	}
	return fmt.Sprintf("%s on https://github.com/%s/%s ref=%s commit=%s", job.Job, job.Refs.Org, job.Refs.Repo, job.Refs.BaseRef, job.Refs.BaseSHA)
}

func jobSpecFromGitRef(ref string) (*api.JobSpec, error) {
	parts := strings.Split(ref, "@")
	if len(parts) != 2 {
		return nil, fmt.Errorf("must be ORG/NAME@REF")
	}
	prefix := strings.Split(parts[0], "/")
	if len(prefix) != 2 {
		return nil, fmt.Errorf("must be ORG/NAME@REF")
	}
	repo := fmt.Sprintf("https://github.com/%s/%s.git", prefix[0], prefix[1])
	out, err := exec.Command("git", "ls-remote", repo, parts[1]).Output()
	if err != nil {
		return nil, fmt.Errorf("'git ls-remote %s %s' failed with '%s'", repo, parts[1], err)
	}
	sha := strings.Split(strings.Split(string(out), "\n")[0], "\t")[0]
	if len(sha) == 0 {
		return nil, fmt.Errorf("ref '%s' does not point to any commit in '%s'", parts[1], parts[0])
	}
	log.Printf("Resolved %s to commit %s", ref, sha)
	return &api.JobSpec{Type: api.PeriodicJob, Job: "dev", Refs: api.Refs{Org: prefix[0], Repo: prefix[1], BaseRef: parts[1], BaseSHA: sha}}, nil
}

func nodeNames(nodes []*api.StepNode) []string {
	var names []string
	for _, node := range nodes {
		name := node.Step.Name()
		if len(name) == 0 {
			name = fmt.Sprintf("<%T>", node.Step)
		}
		names = append(names, name)
	}
	return names
}

func linkNames(links []api.StepLink) []string {
	var names []string
	for _, link := range links {
		name := fmt.Sprintf("<%#v>", link)
		names = append(names, name)
	}
	return names
}

func topologicalSort(nodes []*api.StepNode) ([]*api.StepNode, error) {
	var sortedNodes []*api.StepNode
	var satisfied []api.StepLink
	seen := make(map[api.Step]struct{})
	for len(nodes) > 0 {
		var changed bool
		var waiting []*api.StepNode
		for _, node := range nodes {
			for _, child := range node.Children {
				if _, ok := seen[child.Step]; !ok {
					waiting = append(waiting, child)
				}
			}
			if _, ok := seen[node.Step]; ok {
				continue
			}
			if !api.HasAllLinks(node.Step.Requires(), satisfied) {
				waiting = append(waiting, node)
				continue
			}
			satisfied = append(satisfied, node.Step.Creates()...)
			sortedNodes = append(sortedNodes, node)
			seen[node.Step] = struct{}{}
			changed = true
		}
		if !changed && len(waiting) > 0 {
			for _, node := range waiting {
				var missing []api.StepLink
				for _, link := range node.Step.Requires() {
					if !api.HasAllLinks([]api.StepLink{link}, satisfied) {
						missing = append(missing, link)
					}
					log.Printf("step <%T> is missing dependencies: %s", node.Step, strings.Join(linkNames(missing), ", "))
				}
			}
			return nil, errors.New("steps are missing dependencies")
		}
		nodes = waiting
	}
	return sortedNodes, nil
}

func printExecutionOrder(nodes []*api.StepNode) error {
	ordered, err := topologicalSort(nodes)
	if err != nil {
		return fmt.Errorf("could not sort nodes: %v", err)
	}
	log.Printf("Running %s", strings.Join(nodeNames(ordered), ", "))
	return nil
}

var shaRegex = regexp.MustCompile(`^[0-9a-fA-F]+$`)

// shorten takes a string, and if it looks like a hexadecimal Git SHA it truncates it to
// l characters. The values provided to job spec are not required to be SHAs but could also be
// tags or other git refs.
func shorten(value string, l int) string {
	if len(value) > l && shaRegex.MatchString(value) {
		return value[:l]
	}
	return value
}

func eventRecorder(kubeClient *coreclientset.CoreV1Client, namespace string) record.EventRecorder {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartRecordingToSink(&coreclientset.EventSinkImpl{
		Interface: coreclientset.New(kubeClient.RESTClient()).Events("")})
	return eventBroadcaster.NewRecorder(
		templatescheme.Scheme, coreapi.EventSource{Component: namespace})
}
