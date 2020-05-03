package template

const VolumeXML = `
<volume>
  <!-- Providing a name for the volume which is unique to the pool -->
  <name>{{ .VolumeName }}</name>
  <!-- Providing the total storage allocation for the volume -->
  <allocation>{{ .VolumeAllocation }}</allocation>
  <!-- Providing the logical capacity for the volume -->
  <capacity unit="{{. VolumeCapacityUnit }}">{{ .VolumeCapacity }}</capacity>
  <!-- Provides information about the representation of the volume on the local host -->
  <target>
    <!-- Provides the location at which the volume can be accessed on the local filesystem, as an absolute path -->
    <path>{{ .VolumeTargetPath }}</path>
  </target>
</volume>
`
