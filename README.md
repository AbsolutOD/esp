# esp

A small Go CLI for managing AWS SSM Parameter Store entries as environment variables. Per-project configuration via `.espFile` lets you address parameters by short name and switch between environments with `--env`.

## Install

```sh
go install github.com/AbsolutOD/esp@latest
```

Or build from source:

```sh
git clone https://github.com/AbsolutOD/esp.git
cd esp
go build -o esp .
```

## Configuration

`esp` requires two AWS environment variables to be set before running any subcommand (other than `--help`):

- `AWS_DEFAULT_REGION` — the AWS region to operate in.
- `AWS_PROFILE` — must be set, but its value is only used by the AWS SDK's standard credential resolution chain. If you authenticate via `AWS_ACCESS_KEY_ID` / `AWS_SECRET_ACCESS_KEY`, you can set `AWS_PROFILE` to anything — `esp` checks only that the variable is non-empty.

### Per-project config: `.espFile`

Run `esp init` in a project directory to generate a `.espFile`. The file maps short parameter names to a path scheme of `/<OrgName>/<Env>/<AppName>/<param>`:

```yaml
backend: ssm
orgname: acme
orgprefix: ACME
appname: api-service
envs:
  - dev
  - staging
  - prod
```

When run inside a directory containing a `.espFile`, `esp` requires the `--env` flag and resolves short names through the path scheme above.

## Commands

| Command | Aliases | Description |
|---|---|---|
| `esp init` | | Initialize a `.espFile` in the current directory. |
| `esp list` | `ls` | List parameters under a path. |
| `esp get <name>` | | Read a single parameter. |
| `esp put --name <name> --value <value>` | `add`, `create` | Write a parameter. Add `--secure` to store as `SecureString`. |
| `esp copy <src> <dst>` | `cp` | Copy a parameter to a new path. |
| `esp move <src> <dst>` | `mv` | Move (copy then delete) a parameter. |
| `esp delete <name>` | `rm` | Delete a parameter. |
| `esp version` | | Print the build version. |

## Common workflows

One-off lookup against a literal SSM path:

```sh
AWS_DEFAULT_REGION=us-east-1 AWS_PROFILE=my-profile esp get /acme/prod/api-service/DB_URL
```

Project-relative lookup (inside a `.espFile` directory):

```sh
esp --env=prod get DB_URL
```

List everything for one environment:

```sh
esp --env=staging list
```

Store a secret as a SecureString parameter:

```sh
esp --env=prod put --name DB_PASSWORD --value "$(pbpaste)" --secure
```

## Logging

`esp` logs to stderr via Go's `slog`. The default level is WARN. Pass `--verbose` to lower the threshold to INFO so you can see what `esp` is doing under the hood:

```sh
esp --verbose --env=dev list
```

---

Released under the Apache License 2.0 — see [LICENSE](LICENSE).
