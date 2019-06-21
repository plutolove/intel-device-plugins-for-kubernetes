FROM golang:1.11-alpine as builder
ARG DIR=/home/master/workdir/go-workspace/src/github.com/intel/intel-device-plugins-for-kubernetes

ENV GOPATH /home/master/workdir/go-workspace

ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

WORKDIR $DIR
COPY . .
RUN cd cmd/gpu_plugin; go install
RUN chmod a+x /home/master/workdir/go-workspace/bin/gpu_plugin


FROM alpine
COPY --from=builder /home/master/workdir/go-workspace/bin/gpu_plugin /usr/bin/intel_gpu_device_plugin
CMD ["/usr/bin/intel_gpu_device_plugin"]
