package nomad

import (
	"fmt"
	"io/ioutil"
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

func Apply() (string, error) {
	// https://github.com/hashicorp/nomad-openapi/
	// https://docs.google.com/presentation/d/1h4OOjPFOHbDJsbtuQZRYDjotyBH1YZs7V8L7qmEjRXc/edit#slide=id.gd36c5fdcb4_1_200

	applyJob("./test/data/jobs/traefik.nomad")
	applyJob("./test/data/jobs/hello_world.nomad")

	return "Succes", nil
}

func ReadFromFile(file string) (data []byte, err error) {
	byteValue, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return byteValue, nil
}

func applyJob(filePath string) (string, error) {
	client, err := v1.NewClient()
	if err != nil {
		fmt.Println(err.Error())
	}

	opts := v1.DefaultQueryOpts()

	jobHCL, _ := ReadFromFile(filePath)

	parsedJob, err := client.Jobs().Parse(opts.Ctx(), string(jobHCL[:]), false, false)
	if err != nil {
		return "", err
	}
	_, _, err = client.Jobs().Plan(opts.Ctx(), parsedJob, false)
	if err != nil {
		return "", err
	}

	client.Jobs().Post(opts.Ctx(), parsedJob)

	return "Succes", nil
}
