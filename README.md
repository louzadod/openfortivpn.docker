# FortiClient VPN no Linux com tokens Aladdin eToken Pro e Safenet 5110

Crie o `alias` de execução adicionando o seguinte trecho ao arquivo de
inicialização do seu shell (`~/.bashrc` se você usa Bash; `~/.zshrc`, se ZSH):

```bash
alias vpn="sudo docker run --rm -ti --network=host --device=/dev/bus/usb --device=/dev/ppp --cap-add=NET_ADMIN -v ~/.config/openfortivpn:/vpn -v /etc/resolv.conf:/etc/resolv.conf ghcr.io/fabianonunes/openfortivpn.docker:1.9.2"
```

> **Atenção!** A criação do `alias` não afeta os terminais que já estavam
abertos. Portanto, após ajustar o arquivo de inicialização do shell, abra
um outro terminal ou recarregue-o com `source ~/.bashrc` ou `source ~/.zshrc`.

Inicie a VPN:

```bash
vpn
```

E só! 🤓

## Atalhos

* `vpn reconfigure`: abre formulário de configuração da VPN
* `vpn edit`: permite edição manual do arquivo de configuração

## FAQ

<!-- markdownlint-disable no-inline-html -->
<details>
<summary>Por que utilizar --network=host?</summary>

Para a VPN funcionar, o `openfortivpn` cria uma interface `ppp` e adiciona
rotas IP estáticas à tabela de roteamento do kernel. Por exemplo, ele pode
rotear todas as conexões com destino a 172.16.0.0/12 para a interface `ppp0`.

Se não utilizássemos `--network=host`, essas rotas só funcionariam dentro do
próprio container.
</details>

<details>
<summary>
Por que o container precisa de permissões em determinados <em>devices</em>
e da capacidade <code>NET_ADMIN</code>?
</summary>

O openfortivpn precisa de permissões de acesso ao `/dev/ppp` do host para
criar uma interface de rede `ppp`. Já o acesso ao `/dev/bus/usb` permite
a leitura dos certificados do token USB.

Idealmente, passaríamos apenas o device exato do token USB (`--device=/dev/bus/usb/$BUS/$DEVICE`),
mas precisaríamos de algum script para determinar os valores `$BUS` e `$DEVICE`
que formam o caminho do dispositivo, uma vez que eles não são determinísticos.

Já a _capability_ `NET_ADMIN` é um [requisito do driver `ppp`](https://git.io/Jys2R)
(é por esse motivo que o openfortivpn exige o `sudo` pra rodar fora do container).

Para simplificar, essas flags poderiam ser substituídas por `--privileged` e teríamos
o equivalente a rodar `sudo openfortivpn` diretamente no host. Porém, passar amplas
permissões ocultaria o nível exato de acesso do container.
</details>

<details>
<summary>Por que preciso montar o /etc/resolv.conf dentro do container?</summary>

Além de criar uma interface `ppp` e adicionar rotas IP, o `openfortivpn`
também precisa configurar o DNS para que o cliente possa acessar os domínios
da rede sob a VPN.
</details>

<details>
<summary>Qual a função do utilitário `vpnconfig`?</summary>

Nada mais do que um formulário que permite criar um arquivo de configuração do
`openfortivpn` sem passar por toda aquela cerimônia de identificação de
certificados.

Ele detecta automaticamente os certificados elegíveis do token bem como o
hash do certificado do servidor e os guarda nos respectivos atributos do arquivo
de configuração. Caso haja mais de um certificado elegível no token, o usuário
pode escolher qual usar.
</details>

<details>
<summary>Como utilizar minha própria imagem?</summary>

Clone este repositório e construa a imagem:

```bash
docker build -t localhost/openfortivpn:latest .
```

No alias de inicialização, substitua a imagem `ghcr.io/fabianonunes/openfortivpn.docker`
por `localhost/openfortivpn:latest`.

</details>

<details>
<summary>Como tratar o erro [PKCS11 ENGINE_load_private_key: error:26096080:engine routines:ENGINE_load_private_key:failed loading private key]?</summary>

Esse erro informa que o PIN do token foi inserido incorretamente. Se você
optou por guardá-lo no arquivo de configuração, corrija-o com `vpn reconfigure`.

🔴 **CUIDADO**: dependendo das configurações do token, um determinado número de tentativas
inválidas pode bloquear o PIN. Nesse caso, você precisará do PUK para desbloqueá-lo no aplicativo
SafeNet Authentication Client no Windows.
</details>

<details>
<summary>
Como condicionar a inicialização/desligamento da VPN à inserção/remoção do token?
</summary>

Para inicializar a VPN automaticamente, o PIN do token deve estar salvo no
arquivo de configuração. Caso não esteja, execute `vpn reconfigure` para
reconfigurar os atributos da VPN e repassar o PIN.

Crie um arquivo de regras do `udev` e uma unidade service do `systemd`:

**/etc/udev/rules.d/99-eToken.rules**:

_Ajuste o `idVendor` e o `idProduct` de acordo com o token. O Aladdin eToken
é 0529:0600; o Safenet 5110, 0529:0620. Para conferir o identificador do seu token
utilize o comando `lsusb`_

```cfg
ACTION=="add", SUBSYSTEM=="usb" , ATTRS{idVendor}=="0529", ATTRS{idProduct}=="0600", TAG+="systemd", ENV{SYSTEMD_ALIAS}="/dev/meutoken"
ACTION=="remove", SUBSYSTEM=="usb", ENV{PRODUCT}=="529/600/*", TAG+="systemd"
```

**/etc/systemd/system/openfortivpn.service**:

_Substitua `$USUARIO$` pelo nome correto do usuário._

```ini
[Unit]
Wants=pcscd.service
BindsTo=dev-meutoken.device

After=network-online.target
Wants=network-online.target

[Service]
Restart=always
RestartSec=1
StartLimitBurst=3
ExecStartPre=-/usr/bin/docker rm %n
ExecStartPre=/bin/sleep 2
ExecStart=/usr/bin/docker run --rm --name %n --network=host --device=/dev/bus/usb --device=/dev/ppp --cap-add=NET_ADMIN -v /home/$USUARIO$/.config/openfortivpn:/vpn -v /etc/resolv.conf:/etc/resolv.conf localhost/openfortivpn

[Install]
WantedBy=dev-meutoken.device
```

Reincie o `systemd` e o `udev` para que as regras sejam aplicadas:

```bash
sudo udevadm control --reload
sudo systemctl daemon-reload
```

Habilite o servico:

```bash
sudo systemctl enable openfortivpn
```

Para acompanhar os logs do serviço, execute:

```bash
sudo journalctl -fu openfortivpn.service
```

</details>
