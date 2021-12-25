# FortiClient VPN no Linux com token Aladdin eToken Pro

Clone este repositório e construa a imagem:

```bash
docker build -t openfortivpn:latest .
```

Crie o `alias` de execução adicionando o seguinte trecho ao arquivo de
inicialização do seu shell (`~/.bashrc` se você usa Bash; `~/.zshrc`, se ZSH):

```bash
# `sudo` é opcional se seu usuário pertencer ao grupo `docker`
alias vpn="sudo docker run --rm -ti --network=host --device=/dev/bus/usb --device=/dev/ppp --cap-add=NET_ADMIN -v ~/.config/openfortivpn:/vpn -v /etc/resolv.conf:/etc/resolv.conf openfortivpn"
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
criar uma interface de rede `ppp` e ao `/dev/usb` para ler os certificados
do token USB.

Idealmente, passaríamos apenas o _device_ do token USB (`--device=/dev/bus/usb/$BUS/$DEVICE`),
mas precisaríamos de algum script para determinar os valores `$BUS` e `$DEVICE`
que formam o caminho do dispositivo.

Já a _capability_ `NET_ADMIN` é um [requisito do driver `ppp`](https://git.io/Jys2R)
(é por esse motivo que o openfortivpn exige o `sudo` pra rodar fora do container).

Para simplificar, essas flags poderiam ser substituídas por `--privileged` e teríamos
o equivalente a rodar `sudo openfortivpn` diretamente no host. Porém, passar amplas
permissões ocultaria o exato nível de acesso do container.
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
