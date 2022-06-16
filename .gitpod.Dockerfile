FROM gitpod/workspace-full:2022-06-15-18-10-29

RUN sudo apt install -y \
    curl \
    kafkacat \
    htop \
    redis-server \
    && sudo rm -rf /var/lib/apt/lists/*