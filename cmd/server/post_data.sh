curl -X POST http://localhost:8080/api/v2/activity_log/sao \
  -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJjaWQiOiI3ODkiLCJjcGlkIjoiMTIzIiwiZXQiOjE3NDkxOTUyNDMsIml0IjoxNzQ5MTA4ODQzLCJwcGlkIjoic2FvIiwidWlkIjoiNDU2In0.RzKIlW3evp7c6cp_QlR7wlyNj1J0as69XPi981VTd5lmmXJkRzIBm-tk7nC88yZ0QE0BBqHFL7Gj5dcBouW3o_YDqC2o6p6FJIaQyJjzkOcQvI32nDr_PkN9QYXOmzux1dsATF7uWmZckjS2fjk-eoVw_G5s6cNksU776H9VeLrdKAQyA2pw_OX-DM7f_-2ZB2wc7Vq-u3gZYRuoYy6BFpIIyMKFw75bz2OkoKWKQsH7A82dG6uzYZec-tQLruBVQxT_WTHyar_YLORDlJM-vxFMTsNwhmD1drPYlPIz5Kji87qtWFCjNYVQQGe0a6Wa4AmayAz6wEdKGqi0UNXaYA" \
  -H "Content-Encoding: gzip" \
  -H "Content-Type: application/json" \
  -d '{
        "event": "user_login",
        "timestamp": "2025-06-05T12:00:00Z"
      }'