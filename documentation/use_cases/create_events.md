### Title: Create an Event

### Description:
Allows a user to create a new event by providing details such as title, description, location, time, and ticket information.

### Primary Actor: Authenticated User

### Main Flow:
1. The user sends a `POST /events` request with event details in the request body.
2. The API validates the input data.
3. The API saves the event in the database.
4. The API responds with 201 Created and the created event details.

### Alternate Flows:
1. Invalid Input: If the input data is invalid, the API responds with 400 Bad Request and an error message.