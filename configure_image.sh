# source this file
# expected current dir is ${GIT_SCHEMATICS_JOB_PATH}

# in order to avoid conflicts with other scripts and shells w/o support for local variables
# - prefix 'tfinst' is used for all vars and functions,
# - variables within functions also contain a function-specific, short prefix in their name

# Assumption 
# Copy your Terraform provider's plugin(s) to folder ~/.terraform.d/plugins/{darwin,linux}_amd64/, as appropriate.
# Assuming above step is already done in the base ibm tf provider image 

tfinst_download_file () {
  tfinst_df_cmd=$1
  tfinst_df_file=$2
  echo "download_file: cmd: ${tfinst_df_cmd}"
  eval "${tfinst_df_cmd}"
  tfinst_df_rc=$?
  if [ ${tfinst_df_rc} -ne 0 ] ; then
    echo "download_file: download error: rc=${tfinst_df_rc}" >&2
    exit 1
  fi
  echo "download_file: current dir:"
  pwd
  echo "download_file: saved file:"
  ls -l ${tfinst_df_file}
}

# write permission for root only
chmod 755 /go
chmod 755 /go/bin


# Compiles and installs the packages named by the import paths,
# along with their dependencies.
# -o apiserver to change the name .. change in Makefile as well then. 
echo "\n### Installing terraform-provider-ibm-api as /go/bin/terraform-provider-ibm-api"
cd $API_REPO
GOMOD=/go/bin go install -v -ldflags "-X main.commit=${gitSHA} -X main.travisBuildNumber=${travisBuildNo} -X main.buildDate=${buildDate}"
if [ $? -ne 0 ] ; then
  echo "ERROR: failed while making go install command"
  exit 1
fi
mv /go/bin/terraform-provider-ibm-api /go/bin
chmod 700 /go/bin/terraform-provider-ibm-api

# from now on use /tmp as current dir
cd /tmp


# installs terraformer from release. Method 1
# Future use - Uncomment and use after ibmcloud is available as release. Comment Method 2
# tfinst_bin="terraformer"
# echo "\n### Installing Terraformer as /go/bin/${tfinst_bin}"
# #tfinst_url="https://api.github.com/repos/GoogleCloudPlatform/terraformer/releases/assets/${TERRAFORMER_ASSETID}"
# tfinst_cmd="curl -LO https://github.com/GoogleCloudPlatform/terraformer/releases/download/${TERRAFORMER_VERSION}/terraformer-ibmcloud-linux-amd64"
# tfinst_download_file "${tfinst_cmd}" ${TERRAFORMER_NAME}
# chmod +x terraformer-ibmcloud-linux-amd64
# mv terraformer-ibmcloud-linux-amd64 /go/bin/terraformer


# builds and installs terraformer. Method 2
# Current use - Remove after ibmcloud is available as release.
tfinst_bin="terraformer"
tfinst_bin_ibm="terraformer-ibm"
echo "\n### Cloning, Building and Installing Terraformer as /go/bin/${tfinst_bin}"
tfclone_url="https://github.com/GoogleCloudPlatform/terraformer.git"
# clone - Not cloning under gopath. Since it's go mods, it should be fine
git clone -v ${tfclone_url}
echo "\n### Cloned terraformer"
cd terraformer
# download dependencies
echo "\n### Downloading go mod dependencies"
go mod vendor
# build and install
#GOBIN=/go/bin go install -v
echo "\n### Building and installing"
go env
echo "\n Here is the go path $GOPATH"
#go build -v -o $GOPATH/bin/${tfinst_bin}
# GO111MODULE=auto GOBIN=/go/bin GOFLAGS=-mod=vendor 
echo "\n### Building and installing"
go build -v
if [ $? -ne 0 ] ; then
  echo "ERROR: failed while making go build command for terraform"
  exit 1
fi
echo "\n### Installation done, setting permissions"
# Executable permission
#mv /go/bin/${tfinst_bin} /go/bin
mv ${tfinst_bin} /go/bin/${tfinst_bin}  # todo: @srikar - check this once
echo "Giving permissions"
chmod +x /go/bin/${tfinst_bin}

# clean up /tmp
echo "\n### Cleaning up"
#rm -rf /tmp/*

#appuser
chmod -R 775 /go
addgroup -g 1001 -S appuser && adduser -u 1001 -S appuser -G appuser
chown -R appuser:appuser /go