alias tf='docker run -it --rm --env-file .aws/credentials.env -v $(pwd):/home/tfuser hkjn/terraform $@'
