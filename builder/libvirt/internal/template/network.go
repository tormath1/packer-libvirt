package template

const NetworkXML = `
<network>
  <name>{{ .NetworkName }}</name>
  <forward mode='{{ .NetworkMode }}'>
    <nat>
      <port start='1024' end='65535'/>
    </nat>
  </forward>
  <bridge name='{{ .NetworkBridgeName }}' stp='on' delay='0'/>
  <ip address='192.168.122.1' netmask='255.255.255.0'>
    <dhcp>
      <range start='192.168.122.2' end='192.168.122.254'/>
    </dhcp>
  </ip>
</network>
`
