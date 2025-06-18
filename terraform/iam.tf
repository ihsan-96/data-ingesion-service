resource "google_service_account" "ingestion_sa" {
  account_id   = "data-ingestion-sa"
  display_name = "Data Ingestion Service Account"
  description  = "Service account for data ingestion pipeline"
}

resource "google_service_account_key" "ingestion_sa_key" {
  service_account_id = google_service_account.ingestion_sa.name
}

resource "google_project_iam_member" "bigquery_data_editor" {
  project = var.project_id
  role    = "roles/bigquery.dataEditor"
  member  = "serviceAccount:${google_service_account.ingestion_sa.email}"
}

resource "google_project_iam_member" "bigquery_job_user" {
  project = var.project_id
  role    = "roles/bigquery.jobUser"
  member  = "serviceAccount:${google_service_account.ingestion_sa.email}"
}

resource "google_project_iam_member" "run_invoker" {
  project = var.project_id
  role    = "roles/run.invoker"
  member  = "serviceAccount:${google_service_account.ingestion_sa.email}"
}

resource "google_project_iam_member" "scheduler_job_runner" {
  project = var.project_id
  role    = "roles/cloudscheduler.jobRunner"
  member  = "serviceAccount:${google_service_account.ingestion_sa.email}"
}
