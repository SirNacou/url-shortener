data "archive_file" "lambda_zip" {
  type        = "zip"
  output_path = "${path.module}/bootstrap.zip"

  source {
    content  = "#!/bin/sh\necho 'Initial provisioning placeholder'"
    filename = "bootstrap"
  }
}

resource "aws_lambda_function" "api" {
  function_name = "${var.project_name}-api"
  runtime       = "provided.al2023"
  architectures = ["arm64"]
  handler       = "bootstrap"
  role          = aws_iam_role.lambda_exec.arn

  filename = data.archive_file.lambda_zip.output_path

  environment {
    variables = {
      TABLE_NAME = aws_dynamodb_table.urls.name
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