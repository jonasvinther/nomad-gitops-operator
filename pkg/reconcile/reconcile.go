package reconcile

import (
	"fmt"
	"io"
	"time"

	"nomad-gitops-operator/pkg/nomad"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/util"
)

type ReconcileOptions struct {
	Path   string
	Watch  bool
	Delete bool
	Fs     func() (billy.Filesystem, error)
}

func Run(opts ReconcileOptions) error {
	// Create Nomad client
	client, err := nomad.NewClient()
	if err != nil {
		fmt.Printf("Error %s\n", err)
	}

	// Reconcile
	for true {
		fs, err := opts.Fs()

		if err != nil {
			return err
		}

		files, err := util.Glob(fs, opts.Path)
		if err != nil {
			return err
		}

		desiredStateJobs := make(map[string]interface{})

		// Parse and apply all jobs from within the git repo
		for _, filePath := range files {
			f, err := fs.Open(filePath)
			if err != nil {
				return err
			}

			b, err := io.ReadAll(f)
			if err != nil {
				return err
			}

			// Parse job
			hcl := string(b)
			job, err := client.ParseJob(hcl)
			if err != nil {
				// If a parse error occurs we skip the job an continue with the next job
				fmt.Printf("Failed to parse file [%s]: %s\n", filePath, err)
				continue
			}
			desiredStateJobs[*job.Name] = job

			// Apply job
			fmt.Printf("Applying job [%s][%s]\n", *job.Name, filePath)
			_, err = client.ApplyJob(job, hcl)
			if err != nil {
				return err
			}
		}

		// List all jobs managed by Monoporator
		currentStateJobs, err := client.ListJobs()
		if err != nil {
			fmt.Printf("Error %s\n", err)
		}

		// Check if job has the required metadata
		// Check if job is one of the parsed jobs
		for _, job := range currentStateJobs {
			meta := job.Meta

			if _, isManaged := meta["nomoporater"]; isManaged {
				// If the job is managed by Nomoporator and is part of the desired state
				if _, inDesiredState := desiredStateJobs[*job.Name]; inDesiredState {

				} else {
					if opts.Delete {
						fmt.Printf("Deleting job [%s]\n", *job.Name)
						err = client.DeleteJob(job)
						if err != nil {
							fmt.Println(err)
						}
					}
				}
			}
		}

		if !opts.Watch {
			return nil
		}

		time.Sleep(30 * time.Second)
	}

	return nil
}
