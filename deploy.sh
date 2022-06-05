#!/bin/sh

gradle test && gradle buildLambdaZip && pulumi up --stack dev --yes
