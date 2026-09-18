# Chega Junto Concurseiro

## Instalar na VPS

```bash
curl -fsSL https://raw.githubusercontent.com/reinaldorocha/cjc/main/install.sh | sudo bash
```

O instalador cria a aplicação em `/opt/chega_junto_concurseiro`, usa o MySQL existente do Getfy e publica o site na porta `8082`.

## Atualizar

```bash
sudo /opt/chega_junto_concurseiro/install.sh --update
```

O comando baixa a versão atual, recompila as imagens e reinicia os containers.
