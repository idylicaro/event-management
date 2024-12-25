### Title: Generate Paid/Free Ticket

### Description:
Generates a unique ticket for a user after successful payment (for paid events) or registration (for free events).

### Primary Actor: Authenticated User

### Main Flow:
1. The user sends a `POST /events/{id}/tickets` request.
2. For paid events:
3. The API integrates with Stripe to process the payment.
4. Upon successful payment, the API generates a ticket.

### For free events:
1. The API generates the ticket without payment.
2. The API stores the ticket in the database.
3. The API responds with 201 Created and the ticket details.

### Alternate Flows:
1. Payment Failure: If payment fails, the API responds with 402 Payment Required.
2. Event Not Found: If the event ID does not exist, the API responds with 404 Not Found.