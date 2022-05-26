# tereus-api

Tereus API

## Configuration

Add Bucket lifecycle rule:

```sh
mc ilm import local/tereus <<EOF
{
    "Rules": [
        {
            "Expiration": {
                "Days": 1
            },
            "ID": "ClearedSubmissions",
            "Filter": {
                "Tags": [
                    {
                      "Key": "to-delete",
                      "Value": "true"
                    }
                ]
            },
            "Status": "Enabled"
        }
    ]
}
EOF
```
