package http_extractors

/*

# Example 1:
POST http://example.com/api
Content-Type: application/json
Authorization: Bearer token123

{
  "name": "John Doe",
  "email": "john@example.com"
}

%{assert StatusCode == 200}%

# Example 2:

# Example 3:

*/
