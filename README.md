# serve

Serverless function running on AWS Lambda to publish today's games for a given list of volleyball leagues to a discord
server.

## Requirements

For deployment, an installation of node is needed.

## Deployment

To set up your AWS Lambda deployment,
follow [this](https://www.http4k.org/guide/tutorials/serverless_http4k_with_aws_lambda/) article.

- [Setting up env variables](https://docs.aws.amazon.com/lambda/latest/dg/configuration-envvars.html#configuration-envvars-config)
- [Invoking Lambda Functions](https://docs.aws.amazon.com/AmazonCloudWatch/latest/events/RunLambdaSchedule.html#schedule-create-rule)

## Configuring Env Variables

Create a file called `env.ts` in the root of the project with the following content:

```typescript
export const envVariables = {
    variables: {
        DISCORD_URI: "<YOUR-DISCORD-URL>",
        FLASH_SCORE_URI: "<YOUR-FLASHSCORE-URI>",
        FLASH_SCORE_API_KEY: "<YOUR-FLASHSCORE-API-KEY>",
    },
}
```

## Deployment

The script expects that Gradle is [installed](https://sdkman.io/install). 

```
./deploy.sh
```

## TODO

Create scheduled CloudWatch event via pulumi

