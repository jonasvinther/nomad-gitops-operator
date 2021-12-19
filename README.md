# nomad-gitops-operator
A GitOps operator for Hashicorp Nomad

## Run as Nomad job
```yaml
job "nomoperator" {
  datacenters = ["dc1"]
  group "nomoperator" {
    count = 1
    task "nomoperator" {
      driver = "exec"
      config {
        command = "nomoperator"
        args    = ["bootstrap", "git", "--url", "https://github.com/jonasvinther/nomad-state.git"]
      }
      artifact {
        source      = "https://github.com/jonasvinther/nomad-gitops-operator/releases/download/v0.0.1-pre/nomad-gitops-operator_0.0.1-pre_linux_amd64.tar.gz"
        destination = "local"
        mode        = "any"
      }
    }
  }
}
```