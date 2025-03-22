# TCP Server

## Install Postgres

### Add Secret for Postgres

```shell
> kubectl apply -f ./k8s/PostgresSecret.yaml
```

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: postgres-secret
  labels:
    app: postgres
data:
  POSTGRES_DB: "pigeonproject_db"
  POSTGRES_USER: "pigeonproject"
  POSTGRES_PASSWORD: "password"
```

### Add Postgres

```shell
> kubectl apply -f ./k8s/Postgres.yaml
```

## Add Secret for JWT

```shell
> kubectl apply -f ./k8s/MatchmakerSecret.yaml
```

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: matchmaker-config
  labels:
    app: matchmaker
data:
  JWT_SECRET: "ADD_SECRET_HERE"
```