package main

import (
	"fmt"
	"os"

	"github.com/jonasvinther/nomad-gitops-operator/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// func Dain() {
// 	fmt.Println("Nomad GitOps Operator")

// 	// https://github.com/hashicorp/nomad-openapi/
// 	// https://docs.google.com/presentation/d/1h4OOjPFOHbDJsbtuQZRYDjotyBH1YZs7V8L7qmEjRXc/edit#slide=id.gd36c5fdcb4_1_200
// 	client, err := v1.NewClient()
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}

// 	opts := v1.DefaultQueryOpts()

// 	jobHCL, _ := ReadFromFile("./jobs/traefik.nomad")

// 	// fmt.Println(string(jobHCL[:]))
// 	j, err := client.Jobs().Parse(context.Background(), string(jobHCL[:]), false, false)
// 	jo, _, err := client.Jobs().Plan(context.Background(), j, false)

// 	fmt.Printf("%+v\n", *jo)
// 	client.Jobs().Post(context.Background(), j)

// 	jobs, meta, err := client.Jobs().GetJobs(opts.Ctx())
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

// 	for index, job := range *jobs {
// 		fmt.Println(index)
// 		fmt.Println(*job.Name)
// 	}

// 	// fmt.Println(*job.ID)
// 	fmt.Printf("%v", &meta)

// }

// func ReadFromFile(file string) (data []byte, err error) {
// 	byteValue, err := ioutil.ReadFile(file)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return byteValue, nil
// }
