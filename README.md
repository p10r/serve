# serve

## Requirements

For deployment, an installation of node is needed.

## Deployment

To set up your AWS Lambda deployment,
follow [this](https://www.http4k.org/guide/tutorials/serverless_http4k_with_aws_lambda/) article.

- [Setting up env variables](https://docs.aws.amazon.com/lambda/latest/dg/configuration-envvars.html#configuration-envvars-config)
- [Invoking Lambda Functions](https://docs.aws.amazon.com/AmazonCloudWatch/latest/events/RunLambdaSchedule.html#schedule-create-rule)

## Configuring Env Variables
```
aws lambda update-function-configuration --function-name serve-0ced699 \
    --environment "Variables={ \
    DISCORD_URI=http://uri,\
    FLASH_SCORE_URI=https://uri,\
    FLASH_SCORE_API_KEY=<YOU_SECRET_HERE>\
    }"

```

## Package

```
./gradlew buildLambdaZip && pulumi up --stack dev --yes
```

