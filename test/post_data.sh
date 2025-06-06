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

# sds real jwt token
JWT='eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJjaWQiOiAiNjg5NjBjOTQtOWJlNi00MzQzLWE0Y2EtNjQwOGRlN2FhMzMxIiwgImNwaWQiOiAiY2xyIiwgInBwaWQiOiAic2RzIiwgIml0IjogMTc1MDE0MDczNywgInVpZCI6ICIiLCAicGwiOiAiIiwgImV0IjogMTc1MjczMjczN30.2VbSxMD6PX9ClO14Ac6u6E_QBxKXppjP6abjGN6HnH9J97q-VviKwmyPO4jHb3HmS2Uo82GW6aoZYnWNgNngrTWNItW8jeHyFvScgvUxLSvuFcEK2VZaRqWUgchNOGQvlP3B3Fqq_VCFdd5vPyPCDzdPpi5a7WvtM1RVisENd7krij84F34nHHvCfpUkAyqrJ0owMv0aMizbn5Q93gdVxFqNXhzIGFzd7_GX_niQUNVLpd8ANoauLQwfsKNwhPhs6Rw5x8sfReAS0kosuNHn7m5wKiWDIRvuunHuct9OcQBWQ1RbysHYLBstbVofBvpKNDWh44c7avjmCmQTxPz9atp4nIoAbl_ynk4QApdbVkODg695zv_B0MSGR4lq4i_uswwZlSUr8l-X1coUqLZIhetluS-1xde91HTwmrKTDyfS_w_94jf4PTKAeyKjFDTl2quyTd3KEHIMQUSR-KibFCsQGGguWg6edkiYsXLdxCIlb9iZnh6snMlm8g-mAj9gtxBl52PCBrtfsCn6GjeBXzmq1mqpdQYCHezrUN6nZPloiBrNdp6wcdPEjHXw52duPqFUU24Eg5xOB5xWXERMtBt5h3e6NOZyyXYK7IPN97W_zPB2uQZFrEGtGSak9t4kDk5B_tIyQ5YUwdF_fZHlpH2KBM84uX4IwjOOCttrE1g'
curl -X POST http://localhost:8080/api/v2/activity_log/sds \
  -H "Authorization: Bearer $JWT" \
  -H "Content-Encoding: gzip" \
  -H "Content-Type: application/json" \
  -d '{
        "event": "user_login",
        "timestamp": "2025-06-05T12:00:00Z"
      }'