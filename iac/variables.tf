variable "project_id" {
  type    = string
  default = "yummyos-prod"
}

variable "region" {
  type    = string
  default = "asia-southeast2"
}

variable "ssl_certificates" {
  type    = list(string)
  default = ["unilever-yummycorp-com-cert"]
}

variable "ip_address" {
  type    = string
  default = "https://www.googleapis.com/compute/v1/projects/yummyos-prod/global/addresses/pos"
}