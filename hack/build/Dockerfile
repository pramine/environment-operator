FROM alpine:3.5
RUN apk update && apk add --no-cache libgit2 openssl openssh curl bash git openssh libcurl
ADD _output/bin/environment-operator /environment-operator
CMD ["/environment-operator"]
