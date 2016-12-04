run_tf() {
  docker run --rm -it -v $(pwd):/home/tfuser \
         --env-file .aws/credentials.env \
         hkjn/terraform $@
}

alias tf='run_tf $@'
