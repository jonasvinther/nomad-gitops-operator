package nomad

import (
	"fmt"

	nc "github.com/hashicorp/nomad/api"
)

type Client struct {
	nc *nc.Client
}

func NewClient() (*Client, error) {
	config := nc.DefaultConfig()
	nc, err := nc.NewClient(config)
	if err != nil {
		return nil, err
	}

	client := &Client{nc}

	return client, nil
}

func (client *Client) GetJob(name string) (*nc.Job, error) {
	job, _, err := client.nc.Jobs().Info(name, nil)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (client *Client) ListJobs() (map[string]*nc.Job, error) {
	joblist, _, err := client.nc.Jobs().List(nil)
	if err != nil {
		return nil, err
	}

	jobs := make(map[string]*nc.Job)

	for _, job := range joblist {
		j, _ := client.GetJob(job.Name)
		jobs[job.Name] = j
	}

	return jobs, nil
}

func (client *Client) ParseJob(job string) (*nc.Job, error) {
	parsedJob, err := client.nc.Jobs().ParseHCL(job, false)
	if err != nil {
		return nil, err
	}

	return parsedJob, nil
}

func (client *Client) ApplyJob(job *nc.Job, hcl string) (string, error) {
	// Adding metadata to identify the jobs managed by the Nomoporator
	job.SetMeta("nomoporater", "true")
	job.SetMeta("uid", "nomoporator")

	// fmt.Printf("JobName: %s \n", job.GetName())

	_, _, err := client.nc.Jobs().Plan(job, false, nil)
	if err != nil {
		return "", fmt.Errorf("error while running nomad plan: %s", err)
	}

	res, _, err := client.nc.Jobs().RegisterOpts(job, &nc.RegisterOptions{
		Submission: &nc.JobSubmission{
			Source: hcl,
			Format: "hcl2",
		},
	}, nil)

	if err != nil {
		return "", fmt.Errorf("error while registering nomad job: %s", err)
	}

	return res.EvalID, nil
}

func (client *Client) DeleteJob(job *nc.Job) error {
	_, _, err := client.nc.Jobs().Deregister(*job.Name, true, nil)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) GetVariableItems(name string) (nc.VariableItems, error) {
	variableItems, _, err := client.nc.Variables().GetVariableItems(name, nil)
	if err != nil {
		return nil, err
	}

	return variableItems, nil
}

func (client *Client) UpdateVariable(v *nc.Variable) error {
	_, _, err := client.nc.Variables().Update(v, nil)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) DeleteVariable(path string) error {
	_, err := client.nc.Variables().Delete(path, nil)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) ListVariables() (map[string]*nc.Variable, error) {
	variablelist, _, err := client.nc.Variables().List(nil)
	if err != nil {
		return nil, err
	}

	variables := make(map[string]*nc.Variable)

	for _, variable := range variablelist {
		v, _, err := client.nc.Variables().Read(variable.Path, nil)
		if err != nil {
			return nil, err
		}
		variables[variable.Path] = v
	}

	return variables, nil
}
