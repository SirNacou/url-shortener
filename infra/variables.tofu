variable "project_name" {
  type    = string
  default = "url-shortener"
}
variable "aws_region" {
  type    = string
  default = "ap-southeast-1"
}
variable "environment" {
  type        = string
  default     = "prod"
  description = "Deployment environment name (e.g. dev, prod)"
}