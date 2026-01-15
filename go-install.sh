#!/usr/bin/env bash

GO_VERSION="1.22.4"


if ! command -v go &>/dev/null; then
  echo "安装 Go ${GO_VERSION}..."

  GO_TARBALL="go${GO_VERSION}.linux-amd64.tar.gz"
  GO_URL="https://go.dev/dl/${GO_TARBALL}"

  cd /tmp

  if ! curl -fL "$GO_URL" -o "$GO_TARBALL"; then
    echo "下载失败，请检查网络（${GO_URL}）"
    exit 1
  fi

  rm -rf /usr/local/go
  tar -C /usr/local -xzf "$GO_TARBALL"
  ln -sf /usr/local/go/bin/go /usr/bin/go

  echo "Go 安装完成: $(go version)"
else
  echo "Go 已存在: $(go version)"
fi
