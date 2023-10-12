FROM ubuntu:22.04

RUN apt update

RUN apt install -y curl git bash unzip

RUN curl https://rtx.pub/rtx-latest-linux-x64 >/usr/local/bin/rtx && chmod +x /usr/local/bin/rtx

ENV RTX_DEBUG=1
# RUN rtx install -y go
RUN rtx use -g -y hyperfine

RUN rtx use -g -y go

RUN rtx which hyperfine

ENV PATH="${PATH}:/root/.local/share/rtx/installs/go/latest/go/bin:/root/.local/share/rtx/installs/hyperfine/latest/bin"

RUN echo $PATH

RUN go version

RUN hyperfine --version

ENV GIT_PAT=""

ENV GIT_USER="jinyus"
ENV GIT_EMAIL="jinyus@users.noreply.github.com"

# the repo that will be clone. Most likely your fork
ENV GIT_REPO="https://github.com/jinyus/related_post_gen.git"

# incase you use a different name for your fork
ENV GIT_REPO_NAME="related_post_gen"

ENV TEST_NAME="all"

ENV BRANCH="main"

ENV DEVICE="Workflow-VM-2vCPU-7GBram"

ENV RUN2="20000"

ENV RUN3="60000"

WORKDIR /app

COPY docker_start.sh /docker_start.sh

CMD ["/docker_start.sh"]
