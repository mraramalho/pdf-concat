![alt text](image.png)

# 🚀 GoPDF Merger Service

Uma aplicação robusta desenvolvida em **Go** para concatenação de arquivos PDF, projetada para rodar em alta disponibilidade utilizando **Docker Swarm**. A infraestrutura conta com gerenciamento automático de certificados SSL via **Traefik**, controle de taxa (**Rate Limit**) e um ecossistema completo de monitoramento com **Prometheus** e **Grafana**.

## ✨ Funcionalidades

* **Merge de PDFs:** Endpoint otimizado para combinar múltiplos documentos.
* **SSL Automático:** Renovação e gerenciamento de certificados via Let's Encrypt.
* **Segurança:** Middleware de Rate Limit para prevenir abusos (DoS).
* **Monitoramento:** Painéis em tempo real para acompanhar performance e saúde do sistema.
* **Escalabilidade:** Pronto para escala horizontal via Docker Swarm.

---

## 🛠️ Stack Tecnológica

| Componente | Tecnologia | Função |
| --- | --- | --- |
| **Linguagem** | Go | Backend de alta performance |
| **Orquestrador** | Docker Swarm | Gerenciamento de containers e réplicas |
| **Proxy / Ingress** | Traefik v2 | Roteamento, SSL e Rate Limit |
| **Monitoramento** | Prometheus | Coleta de métricas da aplicação |
| **Visualização** | Grafana | Dashboards analíticos |
| **Rede** | Overlay Network | Comunicação isolada entre serviços |

---

## 🏗️ Arquitetura de Rede

A aplicação utiliza duas redes distintas:

1. **`web` (Externa):** Para comunicação entre o Traefik e a internet (tráfego público).
2. **`pdf_internal` (Interna):** Para comunicação isolada entre o PDF Merger, Prometheus e Grafana.

---

## 🚀 Como Implantar (Deploy)

### 1. Pré-requisitos

* Docker e Swarm Mode ativo.
* Rede externa `web` criada:
```bash
docker network create --driver overlay web

```


* Serviço do Traefik rodando com o resolver `letsencryptresolver`.

### 2. Configuração do Monitoramento

Certifique-se de que o arquivo `prometheus.yml` está presente no diretório raiz para coletar as métricas corretamente.

### 3. Deploy da Stack

No diretório do projeto, execute:

```bash
docker stack deploy -c docker-compose.yml pdf

```

### 4. Acesso

* **Aplicação:** `https://pdf.andreramalho.tech`
* **Dashboards:** `https://monitoramento.andreramalho.tech`

---

## 🛡️ Configurações de Segurança (Rate Limit)

A aplicação está protegida por um middleware de limite de tráfego configurado via labels no `docker-compose.yml`:

* **Média:** 5 requisições por segundo.
* **Burst (Pico):** Até 10 requisições simultâneas.

Se o limite for excedido, o Traefik retornará automaticamente o status `429 Too Many Requests`.

---

## 📊 Monitoramento e Logs

Para visualizar os logs da aplicação e diagnosticar problemas:

```bash
# Verificar status das réplicas
docker stack ps pdf

# Visualizar logs em tempo real
docker service logs -f pdf_pdf-merger

```

O **Grafana** está pré-configurado com as credenciais administrativas definidas nas variáveis de ambiente do serviço.

---

## 📝 Manutenção

Para atualizar a imagem da aplicação sem downtime:

```bash
docker service update --image pdf-merger:latest pdf_pdf-merger

```

---

**Desenvolvido por André Ramalho**