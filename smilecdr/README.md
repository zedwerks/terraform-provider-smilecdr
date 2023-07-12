# Smile CDR JSON Admin API library

Smile CDR includes an administration API based on RESTful JSON Web Services.  The API uses simple JSON-based REST calls to check status and configure the CDR.

The default base URL for this API will be found on port 9000 and would be accessible on a localhost installation at http://localhost:9000

## For More details

See [Smile CDR Admin API](https://smilecdr.com/docs/json_admin_endpoints/json_admin_api.html)

## Swagger

[JSON Admin API](http://localhost:9000/swagger-ui.html)

## Authentication

The default auth for the Admin APIs is Basic Digest.. i.e. username and password Base 64 encoded in the Authorization Header.

In future, this library will also support OAuth 2.0 Client Credentials Grant.
