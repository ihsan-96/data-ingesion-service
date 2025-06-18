variable "project_id" {
  description = "GCP Project ID"
  type        = string
}

variable "region" {
  description = "GCP Region"
  type        = string
  default     = "us-central1"
}

variable "dataset_id" {
  description = "BigQuery dataset ID"
  type        = string
  default     = "data_ingestion"
}

variable "api_endpoint" {
  description = "External API endpoint URL"
  type        = string
  default     = "https://jsonplaceholder.typicode.com/posts"
}
