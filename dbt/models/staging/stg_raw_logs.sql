SELECT
    JSON_EXTRACT_SCALAR(payload, '$.id') AS log_id,
    JSON_EXTRACT_SCALAR(payload, '$.userId') AS user_id,
    JSON_EXTRACT_SCALAR(payload, '$.title') AS title,
    JSON_EXTRACT_SCALAR(payload, '$.body') AS body,
    ingested_at,
    source,
    batch_id
FROM {{ source('data_ingestion', 'raw_logs') }}