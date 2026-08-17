data "archive_file" "dummy_lambda" {
  type        = "zip"
  output_path = "${path.module}/bootstrap.zip"

  source {
    content  = "#!/bin/sh\necho 'Initial provisioning placeholder'"
    filename = "bootstrap"
  }
}

resource "aws_lambda_function" "api" {
  function_name = "url-shortener-api"
  runtime       = "provided.al2023"
  architectures = ["arm64"]
  handler       = "bootstrap"
  role          = aws_iam_role.lambda_exec.arn

  filename         = data.archive_file.dummy_lambda.output_path
  source_code_hash = data.archive_file.dummy_lambda.output_base64sha256

  layers = [
    "arn:aws:lambda:${var.aws_region}:753240598075:layer:LambdaAdapterLayerArm64:24"
  ]

  environment {
    variables = {
      PORT                         = "8080"
      AWS_LWA_PORT                 = "8080"
      AWS_LWA_READINESS_CHECK_PATH = "/healthz"
      TABLE_NAME                   = aws_dynamodb_table.urls.name
    }
  }

  lifecycle {
    ignore_changes = [
      filename,
      source_code_hash
    ]
  }
}

resource "aws_lambda_function_url" "api_url" {
  function_name      = aws_lambda_function.api.function_name
  authorization_type = "NONE"
}
