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

### Bootstrapping using a local path

> Get help with `./nomoporator bootstrap fs -h`

```
Bootstrap Nomad using a local path

Usage:
  nomoperator bootstrap fs [path] [flags]

Flags:
      --base-dir string   Path to the base directory (default "./")
      --delete            Enable delete missing jobs
  -h, --help              help for fs
      --path string       glob pattern relative to the base-dir (default "**/*.nomad")
      --var-path string   var glob pattern relative to the base-dir (default "**/*.vars.yml")
      --watch             Enable watch mode

Global Flags:
  -a, --address string   Address of the Nomad server
```

Use it like this:
```
./nomoperator bootstrap fs --base-dir /path/to/base/dir --path jobs/*.nomad
```

### Bootstrapping using a git repository

> Get help with `./nomoporator bootstrap git -h`

```
Bootstrap Nomad using a git repository

Usage:
  nomoperator bootstrap git [git repo] [flags]

Flags:
      --branch string                  git branch (default "main")
      --delete                         Enable delete missing jobs (default true)
  -h, --help                           help for git
      --password string                SSH private key password
      --path string                    glob pattern relative to the repository root (default "**/*.nomad")
      --ssh-insecure-ignore-host-key   Ignore insecure SSH host key
      --ssh-key string                 SSH private key
      --url string                     git repository URL
      --username string                SSH username (default "git")
      --var-path string                var glob pattern relative to the repository root (default "**/*.vars.yml")
      --watch                          Enable watch mode (default true)

Global Flags:
  -a, --address string   Address of the Nomad server
```

Use it like this:
```
./nomoperator bootstrap git --url https://github.com/jonasvinther/nomad-state.git --path jobs/*.nomad --branch main
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
        args    = ["bootstrap", "git", "--url", "https://github.com/jonasvinther/nomad-state.git", "--branch", "main", "--path", "jobs/*.nomad"]
      }
      artifact {
        source      = "https://github.com/jonasvinther/nomad-gitops-operator/releases/download/v0.0.2/nomad-gitops-operator_0.0.2_linux_amd64.tar.gz"
        destination = "local"
        mode        = "any"
      }
    }
  }
}
```

## SSH

You can use SSH keys to connect to a private git repository.

* Generate a public and private key

```bash
ssh-keygen -t ed25519 -C "nomoperator" -f "nomoperatordeploykey" -N ""
```

If you would like to set password remove `-N ""` and enter the password. Make sure to set `--username sshusername ` and `--pasword sshpassword` when running nomoperator.

* Configure the server git repository with public key

* Generate `known_hosts` for the git server in `/path_to/known_hosts` which is accessible via nomoperator.

```bash
ssh-keyscan -t ed25519 github.com
```

If your git server uses non started port use the `-p` flag.

```bash
ssh-keyscan -t ed25519 -p 2222 mygitserver.com
```

If you would like to avoid using hosts files you can set `--ssh-insecure-ignore-host-key=true`. This is highly discouraged due to security risks.

* Run as nomad job

```yaml
job "nomoperator" {
  datacenters = ["dc1"]
  group "nomoperator" {
    count = 1
    task "nomoperator" {
      driver = "exec"
      env {
        SSH_KNOWN_HOSTS = "/path_to/known_hosts"
        SSH_KEY = <<EOF
-----BEGIN OPENSSH PRIVATE KEY-----
......
-----END OPENSSH PRIVATE KEY-----
EOF

      }
      config {
        command = "nomoperator"
        args    = ["bootstrap", "git", "--url", "git@github.com:jonasvinther/nomad-state.git", "--branch", "main", "--path", "/prod-env", "--username", "git", "--password", "", "--ssh-key", "$SSH_KEY"]
      }
      artifact {
        source      = "https://github.com/jonasvinther/nomad-gitops-operator/releases/download/v0.0.2/nomad-gitops-operator_0.0.2_linux_amd64.tar.gz"
        destination = "local"
        mode        = "any"
      }
    }
  }
}
```

## Variables

Variables are yml files. All keys and values in items should be of type string.

```yaml
path: nomad/jobs/jobname
items:
  key1: "value1"
  key2: "value2"
```
