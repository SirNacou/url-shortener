resource "aws_ssm_parameter" "frontend_bucket" {
  name        = "/url-shortener/frontend-bucket-name"
  description = "S3 Bucket Name for URL Shortener UI"
  type        = "String"
  value       = aws_s3_bucket.frontend.id
}

resource "aws_ssm_parameter" "cloudfront_distribution_id" {
  name        = "/url-shortener/cloudfront-distribution-id"
  description = "CloudFront Distribution ID for URL Shortener UI"
  type        = "String"
  value       = aws_cloudfront_distribution.cdn.id
}