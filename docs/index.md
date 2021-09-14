# Yandex Cloud Provider

## Install

```shell
cloudquery init yandex-cloud/yandex
```

## Authentication

There are four ways to authenticate cloudquery with Yandex Cloud account:

- IAM-token.
- OAuth-token. Get it [here](https://oauth.yandex.ru/authorize?response_type=token&client_id=1a6990aa636648e9b2ef855fa7bec2fb).
- Service account key.
- Authentication from service account on instance.

To authenticate with IAM and OAuth token specify `YC_TOKEN`. Path to file or service account 
key itself should be passed with `YC_SERVICE_ACCOUNT_KEY_FILE`. Authentication from service account on instance is used by default.

`YC_SERVICE_ACCOUNT_KEY_FILE` has higher priority then `YC_TOKEN`.

## Other options

You can specify `organizations_ids`, `clouds_ids` and `folders_ids` to point provider which resources you want to fetch.
Provider will extract all clouds according to specified organizations and all folders according to specified and extracted clouds.

There is an opportunity to specify YC API endpoint via `YC_API_ENDPOINT` or field `api_endpoint` in config. If specified both environment variable will be used. 