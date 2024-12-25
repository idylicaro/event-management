### Title: List Events

### Description:
Allows a user to retrieve a list of events, optionally filtered or paginated.

### Primary Actor: Authenticated or Anonymous User

### Main Flow:
1. The user sends a `GET /events` request.
2. The API retrieves a list of events from the database.
3. The API responds with 200 OK and the list of events in JSON format.

### Alternate Flows:
1. No Events: If no events are available, the API responds with an empty list.