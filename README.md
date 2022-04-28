# Nomoporator
Nomoporator is a GitOps operator for Hashicorp Nomad.

## How to use

### Setting up Nomoporator

#### Environment variables
It's also possible configure Nomoporator via environment variables by setting them like this:
```
NOMAD_ADDR - Required to overide the default of http://127.0.0.1:4646.
NOMAD_TOKEN - Required with ACLs enabled.
NOMAD_CACERT - Required with TLS enabled.
NOMAD_CLIENT_CERT - Required with TLS enabled.
NOMAD_CLIENT_KEY - Required with TLS enabled.
```

### Bootstrapping using a git repository

> Get help with `./nomoporator bootstrap git -h`

```
Bootstrap Nomad using a git repository

Usage:
  nomoperator bootstrap git [git repo] [flags]

Flags:
      --branch string   git branch (default "main")
  -h, --help            help for git
      --path string     path relative to the repository root (default "/")
      --url string      git repository URL

Global Flags:
  -a, --address string   Address of the Nomad server
```

Use it like this:
```
./nomoperator bootstrap git --url https://github.com/jonasvinther/nomad-state.git --path /jobs --branch main
```

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
        args    = ["bootstrap", "git", "--url", "https://github.com/jonasvinther/nomad-state.git", "--branch", "main", "--path", "/prod-env"]
      }
      artifact {
        source      = "https://github.com/jonasvinther/nomad-gitops-operator/releases/download/v0.0.1/nomad-gitops-operator_0.0.1_linux_amd64.tar.gz"
        destination = "local"
        mode        = "any"
      }
    }
  }
}
```