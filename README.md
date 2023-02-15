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

## Build

```shell
docker buildx build --platform=linux/amd64,linux/arm64 --pull --tag saltfishpr/openai-chat:latest --tag saltfishpr/openai-chat:$(date '+%Y%m%d%H%M%S') . --push
```


## Deploy

```shell
docker run -p 3000:3000 -e OPENAI_API_KEY='token' --name openai-chat saltfishpr/openai-chat:latest
```
