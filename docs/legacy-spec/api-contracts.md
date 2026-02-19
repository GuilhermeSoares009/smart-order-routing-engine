# API Contracts

## POST /api/v1/routes
Request
```json
{
  "orderId": "ord-1",
  "symbol": "XYZ",
  "side": "BUY",
  "quantity": 100
}
```

Response
```json
{
  "routeId": "route-1",
  "destination": "venue-a",
  "reason": "LOW_LATENCY"
}
```

## GET /api/v1/routes/{id}
Response
```json
{
  "routeId": "route-1",
  "status": "SENT",
  "destination": "venue-a"
}
```
