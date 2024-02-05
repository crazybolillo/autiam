# autiam
Automate IAM Role token generation on EC2 instances. Simply run `eval $(autiam <role>)`.

## How it works
This simply automates the [process described by AWS](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/iam-roles-for-amazon-ec2.html#instance-metadata-security-credentials)
to obtain access tokens that grant a certain IAM Role permission on EC2 instances.

The program generates the following output:

```shell
admin@ec2:~$ autiam role
export AWS_ACCESS_KEY_ID=...
export AWS_SECRET_ACCESS_KEY=...
export AWS_SESSION_TOKEN=...
```

Where 'role' would be replaced by the actual role's name you want to assume, and the actual tokens
would replace the ellipsis. By running the output through `eval` you automatically set them as environment
variables.

## Why
I run [restic](https://restic.net/) to back up data to S3 buckets on EC2 instances and use temporary IAM Role tokens,
this automates the process of obtaining them so I can focus on running the restic commands.
