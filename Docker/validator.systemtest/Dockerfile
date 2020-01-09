FROM debian:stretch-slim

RUN apt-get update && \
    apt-get -y install \
        git \
        vim \
        telnet \
        make build-essential libssl-dev zlib1g-dev libbz2-dev libreadline-dev \
        libsqlite3-dev wget curl llvm libncurses5-dev libncursesw5-dev xz-utils tk-dev

RUN curl -L https://raw.githubusercontent.com/yyuu/pyenv-installer/master/bin/pyenv-installer | bash && \
    /root/.pyenv/bin/pyenv install 3.7.2

WORKDIR "/root"

ARG GO_VERSION

RUN wget https://dl.google.com/go/go$GO_VERSION.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go$GO_VERSION.linux-amd64.tar.gz && \
    mkdir ~/.go

ENV GOPATH="/root/.go"
ENV GOROOT="/usr/local/go"
ENV PATH="/usr/local/go/bin:/root/.go/bin:${PATH}"

RUN mkdir -p $GOPATH/src/github.com/unification-com
WORKDIR $GOPATH/src/github.com/unification-com

RUN git clone https://github.com/unification-com/mainchain.git
WORKDIR $GOPATH/src/github.com/unification-com/mainchain
RUN make install

COPY Docker/assets /root/.und_mainchain

ENV PATH="/root/.pyenv/versions/3.7.2/bin:${PATH}"
ENV LC_ALL C.UTF-8
ENV LANG C.UTF-8
ENV PYTHONPATH /ev-provision/src

WORKDIR /
RUN git clone https://github.com/unification-com/ev-provision.git

WORKDIR /ev-provision/src
RUN pip install -r requirements.txt

COPY ev-provision /ev-provision

WORKDIR $GOPATH/src/github.com/unification-com
RUN rm -rf mainchain
COPY . $GOPATH/src/github.com/unification-com/mainchain
WORKDIR $GOPATH/src/github.com/unification-com/mainchain
RUN make install

WORKDIR /ev-provision/src

RUN echo "python -m mainnet.systemtest.upgrade" >> /root/.bash_history && \
    echo "alias ll='ls -la'" >> /root/.bashrc

CMD /bin/sleep 600
