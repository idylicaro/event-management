### Title: Delete an Event

### Description:
Allows a user to delete an event they created.

### Primary Actor: Authenticated User (Event Owner)

### Main Flow:
1. The user sends a `DELETE /events/{id}` request.
2. The API checks if the user owns the event.
3. The API deletes the event from the database.
4. The API responds with 204 No Content.

### Alternate Flows:
1. Not Authorized: If the user is not the owner of the event, the API responds with 403 Forbidden.
2. Event Not Found: If the event ID does not exist, the API responds with 404 Not Found.