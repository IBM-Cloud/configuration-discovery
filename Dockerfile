FROM ibmterraform/terraform-provider-ibm-docker:latest

# We can keep these. Will be useful 
ARG gitSHA
ARG travisBuildNo
ARG buildDate

ARG API_REPO=/go/src/github.com/terraform-provider-ibm-api

COPY . $API_REPO

ENV GOPATH="/go"
ENV PATH=$PATH:/usr/local/go/bin:/go/bin
ENV TERRAFORMER_VERSION="0.8.10"

RUN cd ${API_REPO} && . ./configure_image.sh

EXPOSE 9080

USER appuser

WORKDIR $API_REPO

# Run the application as root process
# To pass HTTP addr and port other than localhost:8080, Docker run with flags
CMD ["terraform-provider-ibm-api"]
