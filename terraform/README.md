# terraform

Some demo Terraform plans.


## Build k8s cluster

Provide a k8s token when doing `plan` / `apply`:
```
terraform plan -var 'k8s_token=123456.43b453f7de265cba'
```

```
terraform apply -var 'k8s_token=123456.43b453f7de265cba'
```

*TODO(hkjn)*: Give python snippet for generating token.

## Remote state

The remote state via S3 bucket is configured with:
```
terraform remote config -backend=s3 -backend-config="bucket=terraform-test-state" -backend-config="key=network/terr
aform.tfstate" -backend-config="region=eu-west-1"
```

This seems to require that `AWS_ACCESS_KEY_ID` /
`AWS_SECRET_ACCESS_KEY` / `AWS_DEFAULT_REGION` is set. (I.e. it
doesn't use the `vars.creds_file`, for some reason).

After that, the `.tfstate` can be pulled/pushed with:
```
terraform remote push
terraform remote pull
```