output "dataset_id" {
  description = "BigQuery dataset ID"
  value       = google_bigquery_dataset.dataset.dataset_id
}

output "raw_table_id" {
  description = "Raw logs table ID"
  value       = google_bigquery_table.raw_logs.table_id
}

output "service_account_email" {
  description = "Service account email"
  value       = google_service_account.ingestion_sa.email
}

output "api_service_url" {
  description = "API service URL"
  value       = google_cloud_run_v2_service.api.uri
}

output "ingester_service_url" {
  description = "Ingester service URL"
  value       = google_cloud_run_v2_service.ingester.uri
}

output "dbt_job_name" {
  description = "DBT Cloud Run job name"
  value       = google_cloud_run_v2_job.dbt_transform.name
}
