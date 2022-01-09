FROM golang as builder

WORKDIR /app
ADD vpnconfig .
RUN go build -ldflags "-w"

FROM debian:bullseye

RUN set -ex;                                    \
    apt-get update;                             \
    apt-get install -y --no-install-recommends  \
      libengine-pkcs11-openssl                  \
      nano                                      \
      openfortivpn                              \
      openssl                                   \
      pcscd                                     \
      unzip                                     \
      wget;                                     \
    rm -rf /var/lib/apt/lists/*;                \
    ln -s libcrypto.so.1.1 /usr/lib/x86_64-linux-gnu/libcrypto.so

COPY --from=builder /app/vpnconfig /usr/bin/vpnconfig
COPY SHA256SUMS entrypoint.sh /

ARG DRIVER_URL="http://repositorio.serpro.gov.br/drivers/safenet/SafeNetAuthenticationClient-9.1_Linux_Ubuntu-RedHat(32-64bits).zip"
RUN set -ex;                                                                  \
    wget "$DRIVER_URL" -O /tmp/safenet.zip;                                   \
    sha256sum -c SHA256SUMS;                                                  \
    unzip /tmp/safenet.zip -d /tmp/;                                          \
    dpkg -x /tmp/SafenetAuthenticationClient-BR-10.0.37-0_amd64.deb /;        \
    mv /usr/share/eToken/drivers/aks-ifdh.bundle /usr/lib/pcsc/drivers;       \
    ln -s libAksIfdh.so.10.0                                                  \
      /usr/lib/pcsc/drivers/aks-ifdh.bundle/Contents/Linux/libAksIfdh.so;     \
    mkdir -p /etc/pkcs11/modules;                                             \
    echo "module: /usr/lib/libeToken.so" > /etc/pkcs11/modules/safenet.conf;  \
    echo "enable-in:" > /etc/pkcs11/modules/p11-kit-trust.module;

ENTRYPOINT ["/entrypoint.sh"]
CMD ["start"]
