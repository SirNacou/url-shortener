# Infrastructure

OpenTofu root module that provisions the URL shortener's AWS resources.
Managed as a single flat module — one file per AWS service — so that
CloudFront, S3, Lambda, and DynamoDB can reference each other directly
without module plumbing.

## Layout

One file per service. Files are read together as a single module; OpenTofu
resolves cross-file references automatically.

| File                | Contains                                                                   |
| ------------------- | -------------------------------------------------------------------------- |
| `versions.tofu`     | OpenTofu/provider version constraints, S3 backend, AWS provider config     |
| `variables.tofu`    | Input variables for the root module                                        |
| `outputs.tofu`      | Outputs consumed by deploy tooling and CI                                  |
| `s3.tofu`           | Frontend static-hosting bucket, public-access block, CloudFront OAC policy |
| `cloudfront.tofu`   | CDN distribution (S3 frontend + Lambda API origins) and OAC                |
| `lambda.tofu`       | API Lambda (Go + Lambda Web Adapter) and its function URL                  |
| `iam.tofu`          | IAM role and inline DynamoDB policy for the Lambda                         |
| `logging.tofu`      | CloudWatch log group and log-permission attachment for the Lambda          |
| `dynamodb.tofu`     | DynamoDB table storing shortened URL records                               |
| `ssm.tofu`          | SSM parameters exposing resource IDs to deploy tooling                     |
| `Taskfile.yml`      | Task runner wrappers around OpenTofu commands                              |

## Dependency flow

```
DynamoDB ───────────────┐
IAM ────────────────────┤
Lambda ─────────────────┼──► CloudFront ◄── S3
CloudWatch (logging) ───┘        │  │          │
                                 │  └──► SSM ──┘
                                 └──► outputs
```

- CloudFront serves the S3 frontend for the site and proxies `/api/*` and
  `/s/*` to the Lambda function URL.
- The Lambda reads/writes DynamoDB and logs to CloudWatch.
- SSM parameters and outputs expose bucket/distribution IDs for deployment
  (frontend sync, cache invalidation, backend updates).

## Commands

```sh
task init       # OpenTofu init (providers, modules, backend)
task fmt        # Format configuration files
task validate   # Validate configuration
task plan       # Show the execution plan
task plan:out   # Write a saved plan to tfplan
task apply      # Apply the current configuration
task apply:plan # Apply a previously saved plan
task destroy    # Destroy all managed resources
task outputs    # Show root module outputs
```

## State

State is stored remotely in the S3 backend declared in `versions.tofu`
(`opentofu-state-419212279550`, key `url-shortener/terraform.tfstate`).
Local state files, `.terraform/`, and generated artifacts (`bootstrap.zip`,
`tfplan`) are gitignored.
