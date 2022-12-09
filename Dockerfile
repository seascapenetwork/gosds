FROM huayuanzhen/sdsbase:v1

ARG GOPROXY https://goproxy.cn

ENV GOPROXY="https://goproxy.cn"

WORKDIR /data/code

#COPY . /data/code

#RUN go mod tidy && go mod download github.com/blocklords/gosds
