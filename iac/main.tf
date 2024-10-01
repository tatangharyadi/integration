terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.51.0"
    }
  }
}

provider "google" {
  project = var.project_id
}

resource "google_compute_region_network_endpoint_group" "neg" {
  provider              = google-beta
  project               = var.project_id
  region                = var.region
  name                  = "pos-api-neg"
  network_endpoint_type = "SERVERLESS"
  serverless_deployment {
    platform = "apigateway.googleapis.com"
    resource = google_api_gateway_gateway.api_gw.gateway_id
  }
}

resource "google_compute_backend_service" "default" {
  name                  = "pos-api-backend-service"
  load_balancing_scheme = "EXTERNAL"

  backend {
    group = google_compute_region_network_endpoint_group.neg.id
  }
}

resource "google_compute_url_map" "urlmap" {
  name            = "pos-url-map"
  default_service = google_compute_backend_service.default.name
}

resource "google_compute_target_https_proxy" "default" {
  name             = "pos-https-proxy"
  url_map          = google_compute_url_map.urlmap.name
  ssl_certificates = var.ssl_certificates
}

resource "google_compute_global_forwarding_rule" "default" {
  name       = "pos-https-rule"
  target     = google_compute_target_https_proxy.default.self_link
  port_range = "443"
  ip_address = var.ip_address
}