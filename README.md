# Countries API

Local API for Countries DB from `Restcountries.com`

## Endpoint

```shell
http://127.0.0.1:3030/api/countries
```

## ML Testing

```shell
cd ml
python3 -m venv venv
source venv/bin/activate

pip install -r requirements.txt
```
Build model:
```shell
make ml-train
```

Test model:
```shell
make ml-test 
```
Output:
```shell
🔍 Test ML-suggestion:
source ml/venv/bin/activate && python3 ml/ml_suggest.py "georiga"
["Georgia", "South Georgia", "Nigeria", "Tonga", "Germany"]
```

## Docker

Build docker image:

```shell
docker build -t country-api .
```

Run image:

```shell
docker run -p 3030:3030 country-api
```