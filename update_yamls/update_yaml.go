// cmd/main.go
package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"strings"
)

// QuotedString type to preserve quoted strings
type QuotedString string

// UnmarshalYAML implements the yaml.Unmarshaler interface
func (qs *QuotedString) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}

	// Preserve quotes by adding them explicitly
	if strings.HasPrefix(str, "\"") && strings.HasSuffix(str, "\"") {
		*qs = QuotedString(str)
	} else {
		*qs = QuotedString("\"" + str + "\"")
	}
	return nil
}

var (
	copiedOverJobs    []string
	notCopiedOverJobs []string
	allJobs           []string
)

//type Presubmit struct {
//	Name              string            `yaml:"name" json:"name"`
//	Cluster           string            `yaml:"cluster,omitempty" json:"cluster,omitempty"`
//	PathAlias         string            `yaml:"path_alias,omitempty" json:"path_alias"`
//	Optional          bool              `yaml:"optional" json:"optional"`
//	Decorate          bool              `yaml:"decorate,omitempty" json:"decorate"`
//	DecorationConfig  DecorationConfig  `yaml:"decoration_config,omitempty" json:"decoration_config,omitempty"`
//	SkipIfOnlyChanged string            `yaml:"skip_if_only_changed,omitempty" json:"skip_if_only_changed,omitempty"`
//	RunIfChanged      string            `yaml:"run_if_changed,omitempty" json:"run_if_changed,omitempty"`
//	AlwaysRun         bool              `yaml:"always_run" json:"always_run"`
//	MaxConcurrency    int               `yaml:"max_concurrency,omitempty" json:"max_concurrency,omitempty"`
//	Labels            map[string]string `yaml:"labels,omitempty" json:"labels,omitempty"`
//	Branches          []string          `yaml:"branches,omitempty" json:"branches"`
//	ExtraRefs         []ExtraRef        `yaml:"extra_refs,omitempty" json:"extra_refs"`
//	Spec              Spec              `yaml:"spec,omitempty" json:"spec"`
//	Annotations       TestAnnotation    `yaml:"annotations,omitempty" json:"annotations"`
//}

type Periodic struct {
	Name              string            `yaml:"name" json:"name"`
	Cluster           string            `yaml:"cluster,omitempty" json:"cluster,omitempty"`
	PathAlias         string            `yaml:"path_alias,omitempty" json:"path_alias"`
	Optional          bool              `yaml:"optional,omitempty" json:"optional"`
	Decorate          bool              `yaml:"decorate,omitempty" json:"decorate"`
	DecorationConfig  DecorationConfig  `yaml:"decoration_config,omitempty" json:"decoration_config,omitempty"`
	Interval          string            `yaml:"interval,omitempty" json:"interval,omitempty"`
	SkipIfOnlyChanged string            `yaml:"skip_if_only_changed,omitempty" json:"skip_if_only_changed,omitempty"`
	RunIfChanged      string            `yaml:"run_if_changed,omitempty" json:"run_if_changed,omitempty"`
	AlwaysRun         bool              `yaml:"always_run,omitempty" json:"always_run"`
	MaxConcurrency    int               `yaml:"max_concurrency,omitempty" json:"max_concurrency,omitempty"`
	Labels            map[string]string `yaml:"labels,omitempty" json:"labels,omitempty"`
	Branches          []string          `yaml:"branches,omitempty" json:"branches"`
	ExtraRefs         []ExtraRef        `yaml:"extra_refs,omitempty" json:"extra_refs"`
	Spec              Spec              `yaml:"spec,omitempty" json:"spec"`
	Annotations       TestAnnotation    `yaml:"annotations,omitempty" json:"annotations"`
}

type ExtraRef struct {
	Org       string `yaml:"org,omitempty" json:"org,omitempty"`
	Repo      string `yaml:"repo,omitempty" json:"repo,omitempty"`
	BaseRef   string `yaml:"base_ref,omitempty" json:"base_ref,omitempty"`
	PathAlias string `yaml:"path_alias,omitempty" json:"path_alias,omitempty"`
}

type Spec struct {
	ServiceAccountName string      `yaml:"serviceAccountName,omitempty" json:"serviceAccountName"`
	Containers         []Container `yaml:"containers,omitempty" json:"containers"`
}

type Container struct {
	Image           string          `yaml:"image,omitempty" json:"image"`
	Command         []string        `yaml:"command,omitempty" json:"command"`
	Args            []string        `yaml:"args,omitempty" json:"args,omitempty"`
	Env             []EnvVar        `yaml:"env,omitempty" json:"env,omitempty"`
	SecurityContext SecurityContext `yaml:"securityContext,omitempty" json:"securityContext,omitempty"`
	Resources       Resources       `yaml:"resources,omitempty" json:"resources,omitempty"`
}

type EnvVar struct {
	Name  string `yaml:"name" json:"name"`
	Value string `yaml:"value" json:"value"`
}

type SecurityContext struct {
	Privileged bool `yaml:"privileged,omitempty" json:"privileged"`
}

type Resources struct {
	Limits   ResourceDetails `yaml:"limits,omitempty" json:"limits,omitempty"`
	Requests ResourceDetails `yaml:"requests,omitempty" json:"requests,omitempty"`
}

type ResourceDetails struct {
	CPU    string `yaml:"cpu,omitempty" json:"cpu,omitempty"`
	Memory string `yaml:"memory,omitempty" json:"memory,omitempty"`
}

type DecorationConfig struct {
	Timeout string `yaml:"timeout,omitempty" json:"timeout,omitempty"`
}

//type Config struct {
//	Presubmits map[string][]Presubmit `yaml:"presubmits" json:"presubmits"`
//}

type Config struct {
	Periodics []Periodic `yaml:"periodics" json:"periodics"`
}

type TestAnnotation struct {
	TestgridDashboards string `yaml:"testgrid-dashboards,omitempty" json:"testgrid_dashboards,omitempty"`
	TestgridTabName    string `yaml:"testgrid-tab-name,omitempty" json:"testgrid_tab_name,omitempty"`
	TestgridAlertEmail string `yaml:"testgrid-alert-email,omitempty" json:"testgrid_alert_email,omitempty"`
	Description        string `yaml:"description,omitempty" json:"description,omitempty"`
}

// LoadYAML loads the YAML file
func LoadYAML(filename string) (Config, error) {
	yamlFile := os.Args[1]
	data, err := os.ReadFile(yamlFile)
	if err != nil {
		fmt.Printf("Error reading YAML file: %v\n", err)
		os.Exit(1)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %v\n", err)
		os.Exit(1)
	}

	decoder := yaml.NewDecoder(strings.NewReader(string(data)))
	decoder.SetStrict(false)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Printf("Error decoding YAML file with strict mode: %v\n", err)
		os.Exit(1)
	}

	return config, nil
}

// SaveYAML saves the YAML file
func SaveYAML(filename string, newConfig Config) error {
	data, err := yaml.Marshal(&newConfig)
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
func PrintAllJobs() {
	fmt.Println("All jobs:")
	for _, job := range allJobs {
		fmt.Println(job)
	}
	fmt.Println("")

	fmt.Println("Jobs that were copied over:")
	for _, job := range copiedOverJobs {
		fmt.Printf("  - [ ] " + job + "\n")
	}
	fmt.Println("")

	fmt.Println("Jobs that were not copied over:")
	for _, job := range notCopiedOverJobs {
		fmt.Printf("  - " + job + "\n")
	}
}

// MigrateSpecToWi migrates the spec to -wi
func MigrateSpecToWiPeriodics(config Config) Config {
	newPresubmits := []Periodic{}
	for _, oldJob := range config.Periodics {

		// store all the jobs
		allJobs = append(allJobs, oldJob.Name)

		// Create a new spec only if it has the specified presets
		if oldJob.Labels["preset-azure-cred-only"] == "true" || oldJob.Labels["preset-azure-capz-sa-cred"] == "true" {

			// exit if the job has more than 1 container
			numContainers := len(oldJob.Spec.Containers)
			if numContainers > 1 {
				log.Fatalf("Job %s has more than 1 container", oldJob.Name)
			}

			isE2E := false
			for _, env := range oldJob.Spec.Containers[0].Args {
				if strings.Contains(env, "ci-e2e.sh") {
					isE2E = true
				}
			}

			// migrate only e2e jobs
			if isE2E {
				newJob := oldJob

				// remove old labels
				newJob.Labels = make(map[string]string)
				for k, v := range oldJob.Labels {
					if _, ok := newJob.Labels["preset-azure-cred-only"]; ok {
						continue
					}
					if _, ok := newJob.Labels["preset-azure-capz-sa-cred"]; ok {
						continue
					}
					newJob.Labels[k] = v
				}

				// add new label preset-azure-cred-wi
				newJob.Labels["preset-azure-cred-wi"] = "true"

				// add new serviceAccountName: prowjob-default-sa to spec
				newJob.Spec.ServiceAccountName = "prowjob-default-sa"

				// save to newPresubmits
				newPresubmits = append(newPresubmits, newJob)
			} else {
				// do not migrate non-e2e jobs, ie. conformance jobs for now
				// save to newPresubmits
				newPresubmits = append(newPresubmits, oldJob)
			}
			// save the list of jobs being copied over even if they are not migrated right now.
			// TODO: wait for conformance jobs to be green
			copiedOverJobs = append(copiedOverJobs, oldJob.Name)
		} else {
			notCopiedOverJobs = append(notCopiedOverJobs, oldJob.Name)
			newPresubmits = append(newPresubmits, oldJob)
		}
	}

	return Config{newPresubmits}
}

// CheckAndCreateSpec checks for the specified presets and creates a new spec if they are found
//func CheckAndCreateSpec(config Config) Config {
//	newPresubmits := make(map[string][]Presubmit)
//	for repo, jobs := range config.Presubmits {
//		for _, oldJob := range jobs {
//			newPresubmits[repo] = append(newPresubmits[repo], oldJob)
//			allJobs = append(allJobs, oldJob.Name)
//
//			// work on creating a new spec
//			if oldJob.Labels["preset-azure-cred-only"] == "true" || oldJob.Labels["preset-azure-capz-sa-cred"] == "true" {
//				newJob := oldJob
//
//				// add a -wi to its name
//				newJob.Name = oldJob.Name + "-wi"
//
//				// update labels
//				newJob.Labels = make(map[string]string)
//				for k, v := range oldJob.Labels {
//					if _, ok := newJob.Labels["preset-azure-cred-only"]; ok {
//						continue
//					}
//					if _, ok := newJob.Labels["preset-azure-capz-sa-cred"]; ok {
//						continue
//					}
//					newJob.Labels[k] = v
//				}
//
//				// add new preset-azure-cred-wi label
//				newJob.Labels["preset-azure-cred-wi"] = "true"
//
//				// add new serviceAccountName: prowjob-default-sa to spec
//				newJob.Spec.ServiceAccountName = "prowjob-default-sa"
//
//				// update testgrid tab name
//				newJob.Annotations = TestAnnotation{
//					TestgridDashboards: oldJob.Annotations.TestgridDashboards,
//					TestgridTabName:    oldJob.Annotations.TestgridTabName + "-wi",
//					TestgridAlertEmail: oldJob.Annotations.TestgridAlertEmail,
//					Description:        oldJob.Annotations.Description,
//				}
//
//				// save to newPresubmits
//				newPresubmits[repo] = append(newPresubmits[repo], newJob)
//				copiedOverJobs = append(copiedOverJobs, oldJob.Name)
//			} else {
//				notCopiedOverJobs = append(notCopiedOverJobs, oldJob.Name)
//			}
//		}
//	}
//	return Config{newPresubmits}
//}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <input-yaml-file>", os.Args[0])
	}
	filename := os.Args[1]

	log.Println("Starting update_yaml")
	log.Println("+--------------------------------+")

	presubmits, err := LoadYAML(filename)
	if err != nil {
		log.Fatalf("Failed to load YAML file: %v", err)
	}

	// TO copy over jobs
	//log.Println("Following jobs will be updated:")
	//newConfig := CheckAndCreateSpec(presubmits)

	// To Migrate to wi
	log.Println("Following jobs will be updated:")
	newConfig := MigrateSpecToWiPeriodics(presubmits)

	err = SaveYAML(filename, newConfig)
	if err != nil {
		log.Fatalf("Failed to save YAML file: %v", err)
	}
	fmt.Println("YAML file updated successfully")

	PrintAllJobs()
	log.Println("+--------------------------------+")
}
