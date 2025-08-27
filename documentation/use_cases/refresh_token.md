### Title: Refresh Access Token

### Description:

Handles the token refresh flow for authenticated users, allowing them to obtain a new access token using a valid refresh token without requiring re-authentication.

### Primary Actor: Client Application (Frontend)

### Main Flow:

1. The client sends a `POST /auth/refresh-token` request with the refresh token in the request body.
2. The API validates the provided refresh token:
   - Checks if the token is properly formatted (JWT)
   - Verifies the token signature using the secret key
   - Confirms the token has not expired
   - Extracts the user ID from the token claims
3. If the refresh token is valid, the API generates new tokens:
   - Creates a new access token with fresh expiration time
   - Creates a new refresh token with extended expiration time
4. The API responds with `200 OK`, including both new tokens.

### Alternate Flows:

1. **Invalid Refresh Token:**
   - If the refresh token is malformed, tampered with, or has an invalid signature, the API responds with `401 Unauthorized`.
2. **Expired Refresh Token:**
   - If the refresh token has expired, the API responds with `401 Unauthorized`.
3. **Missing Refresh Token:**
   - If the request body doesn't contain a refresh token, the API responds with `400 Bad Request`.
4. **Token Generation Error:**
   - If there's an error generating new tokens, the API responds with `500 Internal Server Error`.

### Request Format:

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Success Response Format:

```json
{
  "success": true,
  "message": "auth.refreshToken.success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### Error Responses:

- **400 Bad Request:** If the request body is malformed or missing the refresh token.
- **401 Unauthorized:** If the refresh token is invalid, expired, or tampered with.
- **500 Internal Server Error:** If there's an error during token generation.

### Security Considerations:

- The refresh token should be stored securely by the client (e.g., httpOnly cookies or secure storage).
- Each refresh operation should invalidate the old refresh token and issue a new one (token rotation).
- Refresh tokens have longer expiration times than access tokens but should still expire eventually.

### Token Lifetimes:

- **Access Token:** 1 hour
- **Refresh Token:** 7 days

### Notes:

- This endpoint allows users to maintain their session without re-authenticating frequently.
- If the refresh token is expired or invalid, the user must go through the full OAuth2 login flow again.
- The new tokens should replace the old ones in the client's storage.
