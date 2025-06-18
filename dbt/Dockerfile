FROM python:3.11-slim

RUN pip install dbt-bigquery==1.6.0

WORKDIR /app

COPY . /app/

ENV DBT_PROFILES_DIR=/app

CMD ["dbt", "run"]