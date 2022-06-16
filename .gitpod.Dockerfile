FROM gitpod/workspace-full:2022-06-15-18-10-29

RUN curl -1sLf 'https://packages.vectorized.io/nzc4ZYQK3WRGd9sy/redpanda/cfg/setup/bash.deb.sh' | sudo -E bash

RUN sudo apt install -y \
    redpanda \
    curl \
    kafkacat \
    htop \
    redis-server \
    && sudo rm -rf /var/lib/apt/lists/*