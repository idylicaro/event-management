### Title: Edit an Event

### Description:
Allows a user to update the details of an event they created.

### Primary Actor: Authenticated User (Event Owner)

### Main Flow:
1. The user sends a `PUT /events/{id}` request with updated event details.
2. The API checks if the user owns the event.
3. The API validates the new data.
4. The API updates the event in the database.
5. The API responds with 200 OK and the updated event details.

### Alternate Flows:
1. Not Authorized: If the user is not the owner of the event, the API responds with 403 Forbidden.
2. Event Not Found: If the event ID does not exist, the API responds with 404 Not Found.
