FROM golang AS builder
WORKDIR $GOPATH/src/github.com/stone-payments/rumpel
COPY . ./
RUN make
RUN mv $GOPATH/bin/rumpel /rumpel

FROM alpine
COPY --from=builder /rumpel /usr/bin/rumpel
ENV RUMPEL_RULES /.rumpel
VOLUME $RUMPEL_RULES
WORKDIR $RUMPEL_RULES
ENTRYPOINT rumpel -rules=$RUMPEL_RULES