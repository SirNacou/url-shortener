# Outputs consumed by deploy tooling and CI (frontend sync, cache invalidation).

output "aws_region" {
  value       = var.aws_region
  description = "Target AWS deployment region"
}

output "s3_bucket_name" {
  value       = aws_s3_bucket.frontend.id
  description = "Target S3 bucket name for frontend sync"
}

output "lambda_function_name" {
  value       = aws_lambda_function.api.function_name
  description = "Lambda function name for backend updates"
}

output "cloudfront_distribution_id" {
  value       = aws_cloudfront_distribution.cdn.id
  description = "CloudFront distribution ID for cache invalidations"
}

output "cloudfront_url" {
  value       = "https://${aws_cloudfront_distribution.cdn.domain_name}"
  description = "The public HTTPS entrypoint for both frontend and backend"
}