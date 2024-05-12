FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN go build -o getnf .

FROM ubuntu:latest
RUN apt update && apt install -y sudo vim curl fontconfig xz-utils
RUN useradd -m tester && echo "tester:pass" | chpasswd && adduser tester sudo && chown -R tester:tester /home/tester
USER tester
ARG HOME="/home/tester"
RUN mkdir -p ${HOME}/.local/bin/
ENV PATH="${HOME}/.local/bin:${PATH}"
WORKDIR ${HOME}

COPY --chown=tester:tester --from=builder /app/getnf ${HOME}/.local/bin

CMD ["/bin/bash"]
