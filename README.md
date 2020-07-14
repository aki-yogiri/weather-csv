# weather-csv

weather-csvは[weather-store](https://github.com/aki-yogiri/weather-store)に蓄積された気象データを
CSVファイルで取得するAPIサービスです。

# API

## 指定した地域の気象データを取得する

```
GET /weather?location=Tokyo,jp
```

## 指定した地域の特定の期間の気象データを取得する

2020/07/20以前の期間

```
GET /weather?location=Tokyo,jp&dtend=2020-07-20T23:59:59Z
```

2020/06/20以降の期間

```
GET /weather?location=Tokyo,jp&dtstart=2020-06-20T00:00:00Z
```


2020/06/20 - 2020/07/20までの期間

```
GET /weather?location=Tokyo,jp&dtstart=2020-06-20T00:00:00Z&dtend=2020-07-20T23:59:59Z
```

# Build Image

Docker でのビルドを想定しています。

```
$ git clone https://github.com/aki-yogiri/weather-csv.git
$ cd weather-csv
$ sudo docker build -t weather-csv:v1.0.1 .
```

# Deploy on Kubernetes

```
$ kubectl apply -f <path>/<to>/<weather-csv>/kubernetes/weather-csv.yaml
```


# Configuration

weather-crawlerは以下の環境変数を利用します。

| variable | default | |
|----------|---------|-|
| API_HOST | 0.0.0.0 | weather-csvサービスのホスト名 |
| API_PORT | 8080 | weather-csvサービスのポート番号 |
| STORE_HOST | weather-store | weather-storeサービスのホスト名 |
| STORE_PORT | 80 | weather-storeサービスのポート番号 |
