#
# MailHog Dockerfile
#

FROM golang:alpine as builder

# Install MailHog:
RUN apk --no-cache add --virtual build-dependencies \
    git \
    && mkdir -p /root/gocode \
    && export GOPATH=/root/gocode \
    && go install github.com/sskorol/mailhog@latest

FROM alpine:3
# Add mailhog user/group with uid/gid 1000.
# This is a workaround for boot2docker issue #581, see
# https://github.com/boot2docker/boot2docker/issues/581
RUN adduser -D -u 1000 mailhog

COPY --from=builder /root/gocode/bin/mailhog /usr/local/bin/

USER mailhog

WORKDIR /home/mailhog

ENTRYPOINT ["mailhog"]

# Expose the SMTP and HTTP ports:
EXPOSE 1025 8025
