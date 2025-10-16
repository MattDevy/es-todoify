es-setup:
	curl -fsSL https://elastic.co/start-local | sh -s -- --edot

es-start:
	./elastic-start-local/start.sh

es-stop:
	./elastic-start-local/stop.sh

es-uninstall:
	./elastic-start-local/uninstall.sh

es-restart: es-stop es-start

es-creds:
	@source ./elastic-start-local/.env && \
	echo "Elasticsearch username: elastic" && \
	echo "Elasticsearch password: $${ES_LOCAL_PASSWORD}" && \
	echo "Elasticsearch API key: $${ES_LOCAL_API_KEY}"

PHONY: es-setup es-start es-stop es-uninstall es-restart es-creds
