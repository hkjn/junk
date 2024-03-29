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
bash remote_config.sh
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

## Run Dockerized Terraform

Using the `hkjn/terraform` image, you can load the alias from
`docker-tf.sh` and run Dockerized Terraform commands:

```
$ source docker-tf.sh
$ tf plan -var "k8s_token=1b2256.e1c08a1b230a0b04"
$ tf apply -var "k8s_token=1b2256.e1c08a12320a0b04"
```

## Generate Google Cloud Platform service user

The `var.google_credentials` value refers to a `.json`-formatted
credentials file for a service account user which can update Google Cloud DNS.

It can be generated by going to the API manager at
https://console.cloud.google.com, then creating a new Service
Account. The Compute Engine default service account will work.

## Produce graphs of your infrastructure

Terraform can produce directed graphs in `.dot` format:

```
$ terraform graph > graph.dot
$ dot -Tpng graph.dot -o graph.png && xdg-open graph.png
```

## Debug TF outputs

Some of the data in state files are large, like the `cloud-config.yml`
file used for `user-data`. With the `jq` tool, these can be output in readable format:

```
$ cc=$(cat .terraform/terraform.tfstate | jq '.modules[0].resources."data.template_file.master_init".primary.attributes.rendered'); python -c "print($cc)"
```

## TF bug:

Similar to #5199, the following API response indicates that a resource doesn't exist and that it should be removed from state:

Error refreshing state: 1 error(s) occurred:

```
* aws_route.tf_public_gateway_route: Error while checking if route exists: InvalidRouteTableID.NotFound: The routeTable ID 'rtb-2ddae949' does not exist
        status code: 400, request id: 52c3e323-44ef-4445-835b-f1f738c72bac
```

Workaround was to edit local `.terraform/terraform.tfstate` to remove
the resource, _and_ incrementing `serial` (otherwise a local vs remote
state conflict occurs).
