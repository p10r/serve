#!/bin/sh

gradle clean && gradle test && gradle buildLambdaZip && pulumi up --stack dev --yes
