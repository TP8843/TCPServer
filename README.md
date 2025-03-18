# TCP Server

## Install Postgres

### Install Longhorn
```shell
> kubectl apply -f https://raw.githubusercontent.com/longhorn/longhorn/v1.6.0/deploy/longhorn.yaml
```

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