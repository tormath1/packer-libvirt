package template

const PoolXML = `
<pool type="{{ .PoolType }}">
  <!-- Providing a name for the pool which is unique to the host -->
  <name>{{ .PoolName }}</name>
  <target>
    <!-- Provides the location at which the pool will be mapped into the local filesystem namespace, as an absolute path -->
    <path>{{ .PoolTargetPath }}</path>
  </target>
</pool>
`
