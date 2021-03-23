# postgres-exporter

### Local development

```bash
docker-compose down --volumes
docker-compose up -d

go build
go run main.go --addr=0.0.0.0:54380 --target=postgresql://db_user:db_pass@localhost:5432/db_name?sslmode=disable
```
