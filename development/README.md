docker-compose.yml here contains setup of the demo app services with jeager and logging to Loki enabled. It assumes you run both Loki and Jaeger on your local machine (probably through grafana/devenv blocks).

You can comment parts of the compose file as you need. Some setup is different depending on whether you run it on linux or mac so comment or uncomment appropriate parts. 