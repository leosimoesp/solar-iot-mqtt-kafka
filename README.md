<div align="center"><h2>Monitorando geração de energia solar com MQTT e Kafka</h2></div>



![solar_farm_iot](https://github.com/leosimoesp/solar-iot-mqtt-kafka/assets/7965954/32b042d3-d51a-45fa-badf-8ed669f3ba46)

</br>
</br>

# Introdução
Imagine que nossa fazenda solar produz energia comercialmente e temos a necessidade de monitorar e saber o quanto de energia estamos produzindo diariamente, bem como se os sensores de IoT estão funcionando corretamente em cada planta de geração.
O objetivo aqui é criar uma pipeline de dados para coletar os dados produzidos pelos inversores de cada placa solar e conseguir processar esses dados para monitorar a produção de energia diária dessa fazenda por planta de geração.

</br>
</br>

# Solução

Em resumo a solução seria:

![kafka-mqtt-arch](https://github.com/leosimoesp/solar-iot-mqtt-kafka/assets/7965954/1424a398-ba32-4824-82d4-bb540d62f0d6)

### Simulador: 

- Utilizando um script iremos gerar os dados simulados de sensores de plantas de energia solar enviados a cada 15 minutos. Considere que em média no ano o nascer do sol seria em torno das 06h da manhã e 18h da noite o ponte do sol. Os sensores irão enviar dados durante essa janel de tempo.
Este simulador foi desenvolvido com Golang. Essa aplicação possui um <b>"generator"</b> que irá gerar valores para os sensores configurados.
Os dados então serão enviados para um broker MQTT utilizando o cliente do pacote `paho.golang`

### MQTT: 

- MQTT é um protocolo de comunicação de baixa latência, com baixa complexidade e baixo consumo de banda para internet das coisas (IoT).
O transporte das mensagens entre dispositivos remotos utilizam publish/subscribe. O cliente e o servidor(broker) se comunicam de forma assíncrona.
Nessa simulação utlizaremos o `mosquitto` message broker que implementa o `MQTT`. </br> Maiores informações em: https://mosquitto.org

### Kafka:

- Apache Kafka é uma plataforma de streaming de eventos usada para coletar, processar, armazenar e integrar dados em escala. Ela tem vários casos de uso, incluindo streaming distribuído, processamento de fluxo, integração de dados e mensageria pub/sub.
Nessa demonstração utilizaremos um broker `kafka` para receber os dados dos sensores coletados via `MQTT` broker. 
Mais informações em: https://kafka.apache.org

### Kafka Connect:

- Kafka Connect é uma ferramenta de integração de dados que permite a criação de conexões de dados entre Kafka e outros sistemas de armazenamento de dados. Ela permite a criação de conexões customizadas para importar e exportar dados entre Kafka e outros sistemas.
Por meio deste cnector foi possível integrar diretamente os dados do `MQTT` broker e injetá-lo diretamente no tópico kafka.
Mais informações em: https://kafka.apache.org/documentation/#connect


### Telegraf: 
- Telegraf é um agente executado do lado do servidor para ajudar na coleta de métricas de seus sistemas, sensores e aplicações.
Ele é open source e nessa simulação será utilizado para salvar os dados dos sensores coletados pelo tópico kafka e injetar
no nosso banco de dados `InfluxDB`.
Mais informações em: https://www.influxdata.com/time-series-platform/telegraf


### InfluxDB:

- InfluxDB é um banco de dados de série temporal utilizado para armazenar e processar grandes quantidades de dados.
Ele é muito utlizado em aplicações IoT, monitoramento e analytics.
Mais informações em: https://www.influxdata.com


### Grafana:

- Grafana nos permite termos observabilidade de nossa simulação. Podemos acompanhar pelos dashboards a quantidade de energia
solar produzida pelos sensores das fazendas solares. A simulação está programada para gerar dados aleaórios de geração de 
energia solar entre as 06h da manhã até as 18h da tarde. Os valores são aleatórios baseados em faixas de valores de luminosidade
solar e eficiência considerando o hemisfério sul.
Podemos também monitorar a ocorrência de falha nos sensores e nessa simulação, podemos passar alguns parâmetros para validarmos
via grafana o momento de ocorrência até a solução.
Mais informações em: https://grafana.com


## Como executar a simulação?

Requisitos:

- Docker 24.0.7
- Docker Compose version v2.21.0
- Go 1.19+
- Mosquitto-clients (Linux)

1 - Clone o projeto no github executando: 

```
git@github.com:leosimoesp/solar-iot-mqtt-kafka.git
```

2 - Configurando o MQTT Broker

Nessa simulação, utilizamos o Mosquitto como broker MQTT.

No diretório `tools/mosquito/config` criar um arquivo `password.txt`

Na raiz do projeto crie um arquivo `.env` com as seguintes configurações:

```
BROKER_SERVER_URL="localhost:1883"
MQTT_TOPIC="solar-farm-sensors"
MQTT_CLIENTID="client1"
MQTT_PASSWORD="<PWD>"
MQTT_USERNAME="client1"
```

*Observe que na pasta scripts tem o script `init.sh` esse script vai criar 10 users e senhas
para serem utilizados nessa simulação. Essa abordagem é apenas para essa demonstração e não
deve ser utilizada em produção.
Ajuste a senha `MQTT_PASSWORD` conforme o que está no script `init.sh` e gerado durante a execução pelo makefile conforme
a seguir.

Inicialize somente o mosquito utilizando o docker compose desta forma:

```
docker compose up -d mosquitto
```

Crie um user admin via makefile: 

```
make create-user user=admin password=<YOUR_PWD>
```

Crie os users clients a serem utilizados na simulação:

```
make create-mqtt-clients
```
*Não esqueça de ajustar a senha gerada pelo script.sh no arquivo `.env`
Será necessário termos esse client para publicar mensagens mqtt durante a simulação.

@TODO adjust readme.md

Create user

make create-user user=test password=test

mosquitto_sub -h localhost -p 1883 -u admin -P '123456' -t test

mosquitto_pub -h localhost -p 1883 -u admin -P 'Ftrc,W2#' -t test -m "{"hello":2}"

Create a local .env file at this project root with the content

BROKER_SERVER_URL=""
MQTT_TOPIC=""
MQTT_CLIENTID=""
MQTT_PASSWORD=""
MQTT_USERNAME=""

INFLUXDB_USERNAME=
INFLUXDB_PASSWORD=
DOCKER_INFLUXDB_INIT_ORG=
DOCKER_INFLUXDB_INIT_BUCKET=
DOCKER_INFLUXDB_INIT_ADMIN_TOKEN=

INFLUX_TOKEN=""
INFLUX_HOST=""

GF_SECURITY_ADMIN_USER=
GF_SECURITY_ADMIN_PASSWORD=
GF_INSTALL_PLUGINS="grafana-piechart-panel, agenty-flowcharting-panel"

1. make start
2. make create-user user=admin password=Ftrc,W2#
3. make create-mqtt-clients

Download zip file from https://www.confluent.io/hub/confluentinc/kafka-connect-mqtt with jars
Save it into project root

unpack the file at /tmp dir

At /tmp dir create a new folder called jars

Copy all the files from /lib unpacked dir to /tmp/jars

mkdir /tmp/jars && unzip confluentinc-kafka-connect-mqtt-1.7.0.zip -d /tmp/jars

Fill the file connect-mqtt-source.json.example with yor username e password
Remove filename extension .example suffix

At project root make start

docker exec -it kafka /bin/bash

[root@kafka bin]# kafka-topics --bootstrap-server kafka:9092 --create --topic solar-farm-sensors
Created topic solar-farm-sensors.

curl -d @connect-mqtt-source.json -H "Content-Type: application/json" -X POST http://localhost:8083/connectors

{"name":"mqtt-source","config":{"connector.class":"io.confluent.connect.mqtt.MqttSourceConnector","tasks.max":"1","mqtt.username":"admin","mqtt.password":"Ftrc,W2#","mqtt.server.uri":"tcp://mosquitto:1883","mqtt.topics":"solar-farm-sensors","kafka.topic":"solar-farm-sensors","value.converter":"org.apache.kafka.connect.converters.ByteArrayConverter","confluent.topic.bootstrap.servers":"kafka:9092","confluent.topic.replication.factor":"1","name":"mqtt-source"},"tasks":[],"type":"source"}

mosquitto_sub -h 0.0.0.0 -p 1883 -u admin -P 'Ftrc,W2#' -t solar-farm-sensors

mosquitto_pub -h 0.0.0.0 -p 1883 -u admin -P 'Ftrc,W2#' -m hello -t solar-farm-sensors

make build

kafka-console-consumer --bootstrap-server kafka:9092 --topic solar-farm-sensors --from-beginning

#remove topic
kafka-topics --bootstrap-server kafka:9092 --delete --topic solar-farm-sensors

docker exec broker-tutorial kafka-topics \
 --delete \
 --zookeeper zookeeper:2181 \
 --topic blog-dummy

kafka-topics --bootstrap-server http://kafka:9092 --list

https://www.avsystem.com/blog/how-to-solve-iot-device-management-challenges/