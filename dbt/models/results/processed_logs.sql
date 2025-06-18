SELECT
    CAST(log_id AS INT64) AS log_id,
    CAST(user_id AS INT64) AS user_id,
    title,
    body,
    ARRAY_LENGTH(SPLIT(CONCAT(IFNULL(title, ''), ' ', IFNULL(body, '')), ' ')) AS word_count,
    ingested_at,
    source,
    DATE(ingested_at) AS process_date,
    CURRENT_TIMESTAMP() AS created_at
FROM {{ ref('stg_raw_logs') }}
WHERE log_id IS NOT NULL
    AND user_id IS NOT NULL