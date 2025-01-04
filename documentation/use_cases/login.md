### Title: OAuth2 Login with Google (First Time User Included)

### Description:
Handles the OAuth2 login flow for authenticating users via Google, including creating a user record if it's their first time logging in.

### Primary Actor: Client Application (Frontend)

### Main Flow:
1. The client sends a `GET /auth/google/url` request to initiate the OAuth2 flow.
2. The API returns the URL to Google's OAuth2 server with the appropriate parameters.
3. The client redirects the user to the obtained url.
4. The user logs in and consents on Google's page.
5. Google redirects the user back to the API's callback endpoint with an authorization code.
6. The client sends a `POST /auth/google/callback` request with the authorization code in the body.
7. The API exchanges the code for an access token with Google.
8. The API validates the token and retrieves the user's profile (e.g., email, name, profile picture).
9. **First-Time Login Check:**
   - If the user does not exist in the database, the API creates a new user record.
   - Assigns default roles or settings as needed.
10. The API responds with `200 OK`, including a JWT token and user details.

### Alternate Flows:
1. **Invalid Authorization Code:**
   - If the authorization code is invalid or expired, the API responds with `400 Bad Request`.
2. **Google Service Unavailability:**
   - If Google’s OAuth2 service is unreachable, the API responds with `503 Service Unavailable`.
3. **User Denies Consent:**
   - If the user denies consent on Google's page, the API redirects back with an error message.

### Error Responses:
- **400 Bad Request:** If the authorization code is invalid or expired.
- **401 Unauthorized:** If the returned access token is invalid or tampered with.
- **503 Service Unavailable:** If Google’s service is unavailable.
- **409 Conflict:** If there’s an issue creating the user (e.g., duplicate email).

### Notes:
- First-time users may also trigger additional flows, such as sending a welcome email or redirecting to a setup page.
