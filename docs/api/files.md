# File API

## Upload

**POST** `/api/files`

```bash
curl -X POST http://localhost:8090/api/files \
  -H "Authorization: Bearer TOKEN" \
  -F "file=@image.png" \
  -F "path=/uploads/"
```

## Download

**GET** `/api/files/{collection}/{id}/{filename}`

```bash
curl http://localhost:8090/api/files/posts/usr_123/image.png \
  -H "Authorization: Bearer TOKEN"
```

## Storage Endpoint

**GET** `/api/storage`

```bash
curl http://localhost:8090/api/storage?path=/uploads \
  -H "Authorization: Bearer TOKEN"
```
