terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 4.0"
    }
  }
  required_version = ">= 1.0"
}

provider "google" {
  project = var.project_id
  region  = var.region
}

resource "google_project_service" "required_apis" {
  for_each = toset([
    "bigquery.googleapis.com",
    "run.googleapis.com",
    "cloudscheduler.googleapis.com",
    "iam.googleapis.com"
  ])
  
  service = each.value
  disable_on_destroy = false
}

resource "google_cloud_run_v2_service" "api" {
  name     = "data-ingestion-api"
  location = var.region

  template {
    containers {
      image = "gcr.io/${var.project_id}/data-ingestion:latest"
      command = ["./api"]
      
      env {
        name  = "GCP_PROJECT_ID"
        value = var.project_id
      }
      env {
        name  = "DATASET"
        value = var.dataset_id
      }
      env {
        name  = "PROD_TABLE"
        value = "processed_logs"
      }
      env {
        name  = "PORT"
        value = "8080"
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
      }
    }

    service_account = google_service_account.ingestion_sa.email
  }

  traffic {
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
    percent = 100
  }
}

resource "google_cloud_run_v2_service" "ingester" {
  name     = "data-ingestion-ingester"
  location = var.region

  template {
    containers {
      image = "gcr.io/${var.project_id}/data-ingestion:latest"
      command = ["./ingester"]
      
      env {
        name  = "API_ENDPOINT"
        value = var.api_endpoint
      }
      env {
        name  = "GCP_PROJECT_ID"
        value = var.project_id
      }
      env {
        name  = "DATASET"
        value = var.dataset_id
      }
      env {
        name  = "RAW_TABLE"
        value = "raw_logs"
      }
      env {
        name  = "SOURCE"
        value = "placeholder_api"
      }
      env {
        name  = "FETCH_INTERVAL"
        value = "5m"
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
      }
    }

    service_account = google_service_account.ingestion_sa.email
  }

  traffic {
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
    percent = 100
  }
}

resource "google_cloud_run_v2_job" "dbt_transform" {
  name     = "dbt-transform"
  location = var.region

  template {
    template {
      containers {
        image = "gcr.io/${var.project_id}/dbt-transform:latest"
        
        env {
          name  = "GCP_PROJECT_ID"
          value = var.project_id
        }
        env {
          name  = "DATASET"
          value = var.dataset_id
        }

        resources {
          limits = {
            cpu    = "1"
            memory = "2Gi"
          }
        }
      }

      service_account = google_service_account.ingestion_sa.email
      max_retries     = 3
    }
  }
}

resource "google_cloud_run_service_iam_member" "api_public" {
  service  = google_cloud_run_v2_service.api.name
  location = google_cloud_run_v2_service.api.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}

resource "google_cloud_scheduler_job" "dbt_transform" {
  name        = "dbt-transform-daily"
  description = "Daily DBT transformation job"
  schedule    = "0 2 * * *"
  time_zone   = "UTC"

  http_target {
    uri         = "https://${var.region}-run.googleapis.com/apis/run.googleapis.com/v1/namespaces/${var.project_id}/jobs/${google_cloud_run_v2_job.dbt_transform.name}:run"
    http_method = "POST"

    oidc_token {
      service_account_email = google_service_account.ingestion_sa.email
    }
  }

  retry_config {
    retry_count = 3
  }
}