# FortiClient VPN no Linux com token Aladdin eToken Pro

Clone este repositório e construa a imagem:

```bash
docker build -t openfortivpn:latest .
```

Crie o `alias` de execução adicionando o seguinte trecho ao arquivo de inicialização do seu shell (`~/.bashrc` se você usa Bash; `~/.zshrc`, se ZSH):

```bash
# `sudo` é opcional se seu usário pertencer ao grupo `docker`
alias vpn="sudo docker run --rm -ti --network=host --privileged -v ~/.config/openfortivpn:/vpn -v /etc/resolv.conf:/etc/resolv.conf openfortivpn"
```

> **Atenção!** A criação do `alias` não afeta os terminais que já estavam abertos. Portanto, após ajustar
o arquivo de inicialização do shell, abra um outro terminal ou recarregue-o com `source ~/.bashrc` ou `source ~/.zshrc)`.

Inicie a VPN:

```bash
vpn
```

E só! 🤓

## Atalhos

* `vpn reconfigure`: abre formulário de configuração da VPN
* `vpn edit`: permite edição manual do arquivo de configuração
* `vpn p11tool`: p11tool, programa que permite operar dispositivos #PKCS11
