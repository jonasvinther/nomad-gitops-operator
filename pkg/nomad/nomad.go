package nomad

import (
	"fmt"

	nc "github.com/hashicorp/nomad-openapi/clients/go/v1"
	v1 "github.com/hashicorp/nomad-openapi/v1"
)

type Client struct {
	nc *v1.Client
}

func NewClient() (*Client, error) {
	nc, err := v1.NewClient()
	if err != nil {
		return nil, err
	}

	client := &Client{}
	client.nc = nc

	return client, nil
}

func (client *Client) GetJob(name string) (*nc.Job, error) {
	opts := v1.DefaultQueryOpts()
	job, _, err := client.nc.Jobs().GetJob(opts.Ctx(), name)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (client *Client) ListJobs() (map[string]*nc.Job, error) {
	opts := v1.DefaultQueryOpts()

	joblist, _, err := client.nc.Jobs().GetJobs(opts.Ctx())
	if err != nil {
		return nil, err
	}

	jobs := make(map[string]*nc.Job)

	for _, job := range *joblist {
		// fmt.Println(*job.Name)
		j, _ := client.GetJob(*job.Name)
		jobs[*job.Name] = j
		// fmt.Println(j.Meta)
	}

	return jobs, nil
}

func (client *Client) ParseJob(job string) (*nc.Job, error) {
	opts := v1.DefaultQueryOpts()

	parsedJob, err := client.nc.Jobs().Parse(opts.Ctx(), job, false, false)
	if err != nil {
		return nil, err
	}

	return parsedJob, nil
}

// https://github.com/hashicorp/nomad-openapi/
// https://docs.google.com/presentation/d/1h4OOjPFOHbDJsbtuQZRYDjotyBH1YZs7V8L7qmEjRXc/edit#slide=id.gd36c5fdcb4_1_200
func (client *Client) ApplyJob(job *nc.Job) (string, error) {
	opts := v1.DefaultQueryOpts()

	// Adding metadata to identify the jobs managed by the Nomoporator
	metadata := make(map[string]string)
	metadata["nomoporater"] = "true"
	metadata["uid"] = "nomoporator"
	job.SetMeta(metadata)
	// fmt.Printf("JobName: %s \n", job.GetName())

	_, _, err := client.nc.Jobs().Plan(opts.Ctx(), job, false)
	if err != nil {
		return "", fmt.Errorf("error while running nomad plan: %s", err)
	}

	res, _, err := client.nc.Jobs().Post(opts.Ctx(), job)

	if err != nil {
		return "", fmt.Errorf("error while running nomad post: %s", err)
	}

	return *res.EvalID, nil
}

func (client *Client) DeleteJob(job *nc.Job) error {
	opts := v1.DefaultQueryOpts()

	_, _, err := client.nc.Jobs().Delete(opts.Ctx(), job.GetName(), true, true)
	if err != nil {
		return err
	}

	return nil
}
