FROM python:3.10-slim AS python-ml

# Устанавливаем зависимости
#RUN pip install --no-cache-dir scikit-learn

WORKDIR /app

# Копируем только нужные файлы для тренировки
COPY countries.db .
COPY ml/requirements.txt .
COPY ml/train_model.py ml/ml_suggest.py ml/

RUN pip install -r requirements.txt

# Обучаем модель
RUN python ml/train_model.py

# -------------------------------------

FROM golang:1.23-alpine AS go-builder

WORKDIR /app

# Устанавливаем зависимости
RUN apk add --no-cache git gcc musl-dev python3 py3-pip

# Копируем исходники
COPY . .

# Собираем Go-приложение
RUN go build -o app

# -------------------------------------

FROM alpine:latest

WORKDIR /app

# Устанавливаем Python для ml_suggest.py
RUN apk add --no-cache python3 py3-setuptools py3-pip
RUN apk add py3-scikit-learn

COPY --from=go-builder /app/app .
COPY --from=python-ml /app/ml /app/ml
COPY --from=python-ml /app/countries.db /app/countries.db

EXPOSE 3030

CMD ["./app"]
