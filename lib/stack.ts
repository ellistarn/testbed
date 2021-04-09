import { InstanceClass, InstanceSize, InstanceType, Vpc } from '@aws-cdk/aws-ec2'
import { Cluster, KubernetesVersion } from '@aws-cdk/aws-eks'
import { HostedZone, PublicHostedZone, ZoneDelegationRecord } from "@aws-cdk/aws-route53"
import { Stack, StackProps } from "@aws-cdk/core"
import { Construct } from "constructs"

export interface TestbedOptions extends StackProps {
  readonly name: string
  readonly domain: string
}

export class TestbedStack extends Stack {
  constructor(scope: Construct, id: string, options: TestbedOptions) {
    super(scope, id, options)

    const vpc = new Vpc(this, 'VPC', {
      cidr: '10.0.0.0/16',
    })

    const parentZone = HostedZone.fromLookup(this, 'ParentZone', { domainName: options.domain, })

    const zone = new PublicHostedZone(this, 'Zone', {
      zoneName: `${options.name}.${parentZone.zoneName}`,
    })

    new ZoneDelegationRecord(this, 'ZoneDelegation', {
      zone: parentZone,
      recordName: zone.zoneName,
      nameServers: zone.hostedZoneNameServers!,
    })

    const cluster = new Cluster(this, 'Cluster', {
      clusterName: options.name,
      vpc: vpc,
      defaultCapacityInstance: InstanceType.of(InstanceClass.T3, InstanceSize.XLARGE2),
      version: KubernetesVersion.V1_18,
    })

    cluster.addHelmChart("karpenter", {
      chart: "karpenter",
      release: "karpenter/karpenter",
      repository: "https://awslabs.github.io/karpenter/charts",
      createNamespace: true,
      version: "0.2.0"
    })

    cluster.addHelmChart("cert-manager", {
      release: "cert-manager",
      chart: "jetstack/cert-manager",
      repository: "https://charts.jetstack.io",
      version: "v1.2.0",
    })

    // Install flux, which will bootstrap all other components from /config
    // Array.of(
    //   '../config/flux-system/components.yaml',
    //   '../config/flux-system/sync.yaml',
    // ).forEach((file, i) => {
    //   cluster.addManifest(`Manifest-${i}`, ...loadAll(readFileSync(join(__dirname, file), 'utf8')))
    // })
  }
}
