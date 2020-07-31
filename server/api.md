# REST API for Tempest Server

This file documents the REST API for Tempest.

All failure states (ie. when `status` = `false`) will contain an additional
`message` field.

## Accessible Without Authorization

### Route: `/get-nonce`

Gets a temporary server nonce for your IP address (required for login).

#### Method: `GET`

#### Response

| Field | Type |
| ----- | ---- |
| status | bool |
| nonce | int |

### Route: `/login`

Logs into the server.

#### Method: `POST`

#### Args

| Field | Type | Value |
| ----- | ---- | ----- |
| cnonce | int | client-generated nonce |
| auth-hash | string | `sha256(pwd + snonce + cnonce)` |

#### Response

| Field | Type | Value |
| ----- | ---- | ----- |
| status | bool | Whether or not sign in was successful |
| token | string | JWT Token given if sign in succeeds |

## Authorization Required

The JWT token provided during the login process should be placed in the HTTP
header in a field called `Authorization`.

### Route: `/req-upload`

Requests an upload ID to upload a file to a specific location.

#### Method: `GET`

#### Args

| Field | Type | Value |
| ----- | ---- | ----- |
| path | string | Path to the file location |
| isfile | bool | ... |

#### Response 

| Field | Type | Value |
| ----- | ---- | ----- |
| status | bool | ... |
| upload-id | int | The upload id |

### Route: `/upload?uid=upload_id`

Uploads a file to the server.

#### Method: `POST`

#### CONTAINS FILE (Args)

#### Response

| Field | Type | Value |
| ----- | ---- | ----- |
| status | bool | ... |
| message | string | ... |

### Route: `/req-download`

Requests a download ID to download a file.

#### Args

| Field | Type | Value |
| ----- | ---- | ----- |
| path | string | ... |
| isfile | bool | ... |

#### Response

| Field | Type | Value |
| ----- | ---- | ----- |
| status | bool | ... |
| download-id | int | The download id |

### Route: `/download?did=download_id`

Downloads a file from the server.

#### CONTAINS FILE (Response)
