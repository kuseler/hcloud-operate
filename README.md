# SSHKeyOperator
SSH Key Deletion/Creation Operator for Hetzner API

## Setup:
Set your `API_TOKEN` environment variable to your Hetzner API Token
## flags:
1. you can use the `-d` flag to delete keys with the parameter name, e.g. `./SSHKeyOperator -d exampleKeyName`
2. you can use the `-c` flag to create keys, the parameters are name and PublicKey, e.g.: `./SSHKeyOperator -c exampleKeyName "$(< /home/user/.ssh/yourpublickey)"`
