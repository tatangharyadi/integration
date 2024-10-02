resource "google_api_gateway_api" "api_gw" {
  provider = google-beta
  project  = var.project_id
  api_id   = "pos-api"
}

resource "google_api_gateway_api_config" "api_gw" {
  provider = google-beta
  project  = var.project_id

  api                  = google_api_gateway_api.api_gw.api_id
  api_config_id_prefix = "pos-config"

  openapi_documents {
    document {
      path     = "spec.yaml"
      contents = filebase64("openapi.yaml")
    }
  }
  lifecycle {
    create_before_destroy = true
  }
}

resource "google_api_gateway_gateway" "api_gw" {
  provider   = google-beta
  project    = var.project_id
  region     = "asia-northeast1"
  api_config = google_api_gateway_api_config.api_gw.id
  gateway_id = "pos-gateway"
}
