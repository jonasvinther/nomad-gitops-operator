package nomad

import (
	"fmt"
	"os"

	v1 "github.com/hashicorp/nomad-openapi/v1"
)

func List() {
	client, err := v1.NewClient()
	if err != nil {
		fmt.Println(err.Error())
	}

	opts := v1.DefaultQueryOpts()

	jobs, _, err := client.Jobs().GetJobs(opts.Ctx())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for index, job := range *jobs {
		fmt.Println(index)
		fmt.Println(*job.Name)
	}
}

// https://github.com/hashicorp/nomad-openapi/
// https://docs.google.com/presentation/d/1h4OOjPFOHbDJsbtuQZRYDjotyBH1YZs7V8L7qmEjRXc/edit#slide=id.gd36c5fdcb4_1_200
func ApplyJob(job string) (string, error) {
	client, err := v1.NewClient()
	if err != nil {
		fmt.Println(err.Error())
	}

	opts := v1.DefaultQueryOpts()

	parsedJob, err := client.Jobs().Parse(opts.Ctx(), job, false, false)
	if err != nil {
		return "", err
	}
	_, _, err = client.Jobs().Plan(opts.Ctx(), parsedJob, false)
	if err != nil {
		return "", err
	}

	res, _, err := client.Jobs().Post(opts.Ctx(), parsedJob)

	return *res.EvalID, nil
}
