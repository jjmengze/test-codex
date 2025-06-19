#!/bin/sh
#{"error":"Token is expired"}
curl -X POST http://localhost:8080/api/v2/activity_log/sao \
  -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJjaWQiOiI3ODkiLCJjcGlkIjoiMTIzIiwiZXQiOjE3NDk2MDY0MzUsIml0IjoxNzQ5NTIwMDM1LCJwcGlkIjoic2FvIiwidWlkIjoiNDU2In0.TGklEBE7SdwHWwqtabdEv0a0-N19yyuFGghbe7BYbshPeatwxezOK4TlkzEjpoJBbjA7ThNe8L3UOFpMNlmZJSuG5POY6HncDrG1iTqvZcB1QxnFUvaVc7OQSVtDGJ8ZMe4EB4pqrhRuE5zF3PLRD33ltsgCFipKCkilDIiYZ8G9P9vawV8selu1LMkiM0O2kfn-F312-uo8jtkYKG6VBvblwLYiG4-G8nH5eLt1dd-LOlw5ClzNQBfDppejZqH4tIUnIBHUhCH7kMxfSSVKqsbL09U0R6GHMCQ6Y-6T_soPD4DNqwIAq9PaxJFobbmbcRfH7JLtj7NV0d7G-OGqgQ" \
  -H "Content-Encoding: gzip" \
  -H "Content-Type: application/json" \
  -d '{
        "event": "user_login",
        "timestamp": "2025-06-05T12:00:00Z"
      }'

JWT=$(go run ./jwt/jwt_token_getter.go | tr -d '\n')

# mock jwt cloud pass
curl -X POST http://localhost:8080/api/v2/activity_log/sao \
  -H "Authorization: Bearer $JWT" \
  -H "Content-Encoding: gzip" \
  -H "Content-Type: application/json" \
  -d '{
        "event": "user_login",
        "timestamp": "2025-06-05T12:00:00Z"
      }'

# sds real jwt token put your real token
JWT=''
curl -X POST http://localhost:8080/api/v2/activity_log/sds \
  -H "Authorization: Bearer $JWT" \
  -H "Content-Encoding: gzip" \
  -H "Content-Type: application/json" \
  -d '{
        "event": "user_login",
        "timestamp": "2025-06-05T12:00:00Z"
      }'