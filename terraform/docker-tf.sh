alias tf='docker run --rm --env-file .aws/credentials.env -v $(pwd):/home/tfuser hkjn/terraform $@'
