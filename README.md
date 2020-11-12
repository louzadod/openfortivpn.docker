# FortiClient VPN no Linux com token Aladdin eToken Pro

Crie o `alias` de execução e coloque no inicializador do seu shell (`~/.zshrc`, `~/.bashrc`, ...):

```bash
# `sudo` é opcional se seu usário pertencer ao grupo `docker`
alias vpn="sudo docker run --rm -ti --network=host --privileged -v ~/.config/openfortivpn:/vpn -v /etc/resolv.conf:/etc/resolv.conf registry.senado.leg.br/fparente/openfortivpn:latest"
```

Crie o arquivo `~/.config/openfortivpn/config.cfg` com o conteúdo:

```ini
host = IP do gateway
port = Porta do gateway
# user-cert e trusted-cert serão deduzidos automaticamente
```

Inicie a VPN:

```bash
vpn
```

E só! 🤓
