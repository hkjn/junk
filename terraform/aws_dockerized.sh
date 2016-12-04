run_aws() {
  docker run --rm -it \
         --env-file .aws/credentials.env \
         hkjn/aws $@
}

alias awscli='run_aws $@'
