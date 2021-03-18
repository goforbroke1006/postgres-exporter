# postgres-exporter

Custom implementation for dead tuples statistics.

Grafana dashboard in [dashboard.json](dashboard.json)

### Usage

```bash
docker run --name=pg-exporter -it --rm --network=host \
    goforbroke1006/postgres-exporter \
    --addr=0.0.0.0:54380 --target=postgresql://db_user:db_pass@localhost:5432/db_name?sslmode=disable
```

open http://localhost:54380/metrics

### Local development

```bash
docker-compose down --volumes
docker-compose up -d
```
