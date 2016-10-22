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

A new token can be generated with:
```
python -c 'import random; print("{0:x}.{1:x}".format(random.SystemRandom().getrandbits(3*8), random.SystemRandom().getrandbits(8*8)))'
```

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

## Remote access

To access the kubernetes cluster remotely, the easiest way is to copy
out `admin.conf` from the k8s master, as described under the `kubeadm`
"Limitations" section:

- http://kubernetes.io/docs/getting-started-guides/kubeadm/

You'll need to change the `server` entry to specify the public IP for
the master, instead of the internal IP.

After that, just adding `--kubeconfig` should work:

```
kubectl --kubeconfig ./admin.conf get nodes
```

## Producing graphs of your infrastructure

Terraform can produce directed graphs in `.dot` format:

```
$ terraform graph > graph.dot
$ dot -Tpng graph.dot -o graph.png && xdg-open graph.png
```