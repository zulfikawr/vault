# Real-time API

Server-Sent Events (SSE) for real-time updates.

## Connect

**GET** `/api/realtime`

```javascript
const eventSource = new EventSource(
  'http://localhost:8090/api/realtime?collection=posts',
  {
    headers: {'Authorization': 'Bearer TOKEN'}
  }
);

eventSource.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Update:', data);
};
```

## Message Format

```json
{
  "action": "create",
  "collection": "posts",
  "record": {"id": "post_123", "title": "New Post"}
}
```

## Actions

- `create` - New record
- `update` - Record updated
- `delete` - Record deleted

See Also: [Realtime Hub](../internal/realtime/hub.go)
