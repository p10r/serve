import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";
import {RolePolicyAttachment} from "@pulumi/aws/iam";
import {envVariables} from "./env";

const defaultRole = new aws.iam.Role("http4k-default-role", {
    assumeRolePolicy: `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
`
});

new RolePolicyAttachment("http4k-default-role-policy",
    {
        role: defaultRole,
        policyArn: aws.iam.ManagedPolicies.AWSLambdaBasicExecutionRole
    });

const lambdaFunction = new aws.lambda.Function("serve", {
    code: new pulumi.asset.FileArchive("build/distributions/serve.zip"),
    handler: "de.p10r.ServeAppFunction",
    role: defaultRole.arn,
    runtime: "java11",
    timeout: 60,
    memorySize: 512,
    environment: envVariables
});

const logGroupApi = new aws.cloudwatch.LogGroup("api-route", {
    name: "serve",
});

const apiGatewayPermission = new aws.lambda.Permission("serve-gateway-permission", {
    action: "lambda:InvokeFunction",
    "function": lambdaFunction.name,
    principal: "apigateway.amazonaws.com"
});

const api = new aws.apigatewayv2.Api("serve-api", {
    protocolType: "HTTP"
});

const apiDefaultStage = new aws.apigatewayv2.Stage("default", {
    apiId: api.id,
    autoDeploy: true,
    name: "$default",
    accessLogSettings: {
        destinationArn: logGroupApi.arn,
        format: `{"requestId": "$context.requestId", "requestTime": "$context.requestTime", "httpMethod": "$context.httpMethod", "httpPath": "$context.path", "status": "$context.status", "integrationError": "$context.integrationErrorMessage"}`
    }
})

const lambdaIntegration = new aws.apigatewayv2.Integration("serve-api-lambda-integration", {
    apiId: api.id,
    integrationType: "AWS_PROXY",
    integrationUri: lambdaFunction.arn,
    payloadFormatVersion: "1.0"
});

let serverlessHttp4kApiRoute = "serve";
const apiDefaultRole = new aws.apigatewayv2.Route(serverlessHttp4kApiRoute + "-api-route", {
    apiId: api.id,
    routeKey: `$default`,
    target: pulumi.interpolate`integrations/${lambdaIntegration.id}`
});

export const publishedUrl = apiDefaultStage.invokeUrl;
