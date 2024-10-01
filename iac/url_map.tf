resource "google_compute_url_map" "urlmap" {
  name            = "pos-url-map"
  default_service = google_compute_backend_service.default.name

  host_rule {
    hosts        = ["unilever.yummycorp.com"]
    path_matcher = "pos-matcher"
  }

  path_matcher {
    name            = "pos-matcher"
    default_service = google_compute_backend_service.default.name

    route_rules {
      priority = 1000
      match_rules {
        prefix_match = "/"
      }
      route_action {
        weighted_backend_services {
          backend_service = google_compute_backend_service.default.name
          weight          = 100
        }
      }
    }
  }
}
