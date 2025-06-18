resource "google_bigquery_dataset" "dataset" {
  dataset_id                  = var.dataset_id
  friendly_name              = "Data Ingestion Dataset"
  description                = "Dataset for raw and processed log data"
  location                   = "US"
  default_table_expiration_ms = 31536000000

  labels = {
    env = "production"
  }
}

resource "google_bigquery_table" "raw_logs" {
  dataset_id = google_bigquery_dataset.dataset.dataset_id
  table_id   = "raw_logs"

  schema = jsonencode([
    {
      name = "payload"
      type = "STRING"
      mode = "REQUIRED"
    },
    {
      name = "ingested_at"
      type = "TIMESTAMP"
      mode = "REQUIRED"
    },
    {
      name = "source"
      type = "STRING"
      mode = "REQUIRED"
    },
    {
      name = "batch_id"
      type = "STRING"
      mode = "REQUIRED"
    }
  ])
}

resource "google_bigquery_table" "processed_logs" {
  dataset_id = google_bigquery_dataset.dataset.dataset_id
  table_id   = "processed_logs"

  schema = jsonencode([
    {
      name = "log_id"
      type = "INTEGER"
      mode = "REQUIRED"
    },
    {
      name = "user_id"
      type = "INTEGER"
      mode = "REQUIRED"
    },
    {
      name = "title"
      type = "STRING"
      mode = "NULLABLE"
    },
    {
      name = "body"
      type = "STRING"
      mode = "NULLABLE"
    },
    {
      name = "word_count"
      type = "INTEGER"
      mode = "NULLABLE"
    },
    {
      name = "ingested_at"
      type = "TIMESTAMP"
      mode = "REQUIRED"
    },
    {
      name = "source"
      type = "STRING"
      mode = "REQUIRED"
    },
    {
      name = "process_date"
      type = "DATE"
      mode = "NULLABLE"
    },
    {
      name = "created_at"
      type = "TIMESTAMP"
      mode = "REQUIRED"
    }
  ])
}