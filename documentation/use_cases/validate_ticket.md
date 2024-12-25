### Title: Validate a Ticket

### Description:
Checks the validity of a ticket for a specific event.

### Primary Actor:
Event Organizer or System

### Main Flow:
1. The organizer or system sends a `POST /events/{id}/tickets/validate` request with the ticket code in the body.
2. The API checks the database for the ticket.
3. If valid, the API responds with 200 OK and a validation success message.
4. If invalid, the API responds with 400 Bad Request and an error message.
### Alternate Flows:
1. Ticket Not Found: If the ticket does not exist, the API responds with 404 Not Found.
2. Expired Ticket: If the ticket is expired, the API responds with 410 Gone.