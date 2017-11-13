# summon-cerberus

Cerberus provider for [Summon](https://cyberark.github.io/summon).

Provides access to secrets stored in [Cerberus](http://engineering.nike.com/cerberus/).

## Usage

[Set summon-cerberus as your Summon provider](https://github.com/cyberark/summon#flags).

Make sure to set `CERBERUS_API` via environment variable.  
Give summon a path to an object in Cerberus regardless of Safety Deposit Box path (`/app/` or `/shared/`) and it will fetch it for you and
print the value to stdout.

### Example 1
```bash
$ export CERBERUS_API='https://mycerberus_endpoint.com'
$ cat > /tmp/my_secrets.yml <<-EOF
	DB_USER: product_name
	DB_PASSWORD: !var app/product/$ENVTAG/dbpassword
	DATADOG_API_TOKEN: !var shared/datadog/$ENVTAG/datadog_api_token
EOF
$ summon --provider summon-cerberus \
         -f /tmp/my_secrets.yml \
         -D ENVTAG=myenv \
         cat @SUMMONENVFILE

DB_USER=product_name
DB_PASSWORD=Wylb6owWawtenJab
DATADOG_API_TOKEN=6d4f1e2992a11a332550aa555e630f0dc
```

### Example 2
```bash
$ export CERBERUS_API='https://mycerberus_endpoint.com'
$ summon --provider summon-cerberus \
         -D ENVTAG=myenv
         --yaml 'DATADOG_API_TOKEN: !var app/product/$ENVTAG/datadog_api_token' \
         printenv | grep DATADOG_API_TOKEN

DATADOG_API_TOKEN=6d4f1e2992a11a332550aa555e630f0dc
```

### Example 3
```bash
$ export CERBERUS_API='https://mycerberus_endpoint.com'
$ DATADOG_API_KEY=$(summon-cerberus shared/product/myenv/datadog_api_token)
$ echo $DATADOG_API_KEY
6d4f1e2992a11a332550aa555e630f0dc
```

## Configuration

summon-cerberus uses the [official AWS Go SDK](https://github.com/aws/aws-sdk-go).
It will use the credentials file or environment variables [as they explain](https://github.com/aws/aws-sdk-go#configuring-credentials).

Additionally, see [Summon Usage](https://github.com/cyberark/summon#usage) documentation.

## Limitations

summon-cerberus provider assumes the usage of IAM profiles and currently does not support usage of AWS API key/secret. As such, it is unusable anywhere but EC2 instances.  
Improvements are required (PRs welcome) to make it support AWS key/secret the way AWS CLI tool does.

## Authors
credit goes to [@burdzz](https://github.com/burdzz) and [@anapsix](https://github.com/anapsix)
