#!/usr/bin/env node
import * as cdk from '@aws-cdk/core'
import 'source-map-support/register'
import { TestbedStack } from './lib/stack'

new TestbedStack(new cdk.App(), "TestbedStack", {
  name: "testbed",
  domain: `${process.env.USER}.people.aws.dev`,
  env: {
    account: process.env.CDK_DEFAULT_ACCOUNT,
    region: process.env.CDK_DEFAULT_REGION
  }
})
