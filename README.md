# openai chat bot


## API

- Chat

```shell
curl --request POST \
  --url http://127.0.0.1:3000/v1/chat \
  --header 'Content-Type: application/json' \
  --header 'X-Username: saltfishpr' \
  --data '{
	"message": "js 实现冒泡排序"
}'
```

- Delete Conversation

```shell
curl --request DELETE \
  --url http://127.0.0.1:3000/v1/chat \
  --header 'X-Username: saltfishpr'
```

- Get Records

```shell
curl --request GET \
  --url http://127.0.0.1:3000/v1/chat \
  --header 'X-Username: saltfishpr'
```

## TODO

- custom prompt param
- i18n support
