## Quotes API README

This README provides an overview of the Quotes API, including available endpoints and example PowerShell commands using `Invoke-RestMethod` as well as `curl.exe`. Replace placeholder values where necessary.

---

### Base URL

```
http://localhost:8080/quotes
```

---

## Endpoints

| HTTP Method | Endpoint                      | Description                          |
| ----------- | ----------------------------- | ------------------------------------ |
| GET         | `/quotes`                     | Retrieve all quotes                  |
| GET         | `/quotes/random`              | Retrieve a random quote              |
| GET         | `/quotes?author={authorName}` | Retrieve quotes by a specific author |
| POST        | `/quotes`                     | Create a new quote                   |
| DELETE      | `/quotes/{id}`                | Delete a quote by its ID             |

---

## Examples

### 1. List all quotes

**PowerShell**

Invoke-RestMethod -Uri 'http://localhost:8080/quotes' -Method GET
```

**curl.exe**

```powershell
curl.exe http://localhost:8080/quotes
```

---

### 2. Get a random quote

**PowerShell**

```powershell
Invoke-RestMethod \
  -Uri 'http://localhost:8080/quotes/random' \
  -Method GET
```

**curl.exe**

```powershell
curl.exe http://localhost:8080/quotes/random
```

---

### 3. Filter quotes by author

Retrieve all quotes by Confucius:

**PowerShell**

```powershell
Invoke-RestMethod \
  -Uri 'http://localhost:8080/quotes?author=Confucius' \
  -Method GET
```

**curl.exe**

```powershell
curl.exe "http://localhost:8080/quotes?author=Confucius"
```

---

### 4. Create a new quote

Replace `author` and `quote` with desired values.

**PowerShell**

```powershell
$body = @{
  author = 'Albert Einstein'
  quote  = 'Life is like riding a bicycle. To keep your balance you must keep moving.'
}

Invoke-RestMethod \
  -Uri 'http://localhost:8080/quotes' \
  -Method POST \
  -ContentType 'application/json' \
  -Body (ConvertTo-Json $body)
```

**curl.exe**

```powershell
$json = @'
{"author":"Albert Einstein","quote":"Life is like riding a bicycle. To keep your balance you must keep moving."}
'@

curl.exe -X POST http://localhost:8080/quotes \
  -H 'Content-Type: application/json' \
  -d $json
```

---

### 5. Delete a quote by ID

Replace `$id` with the quote's identifier.

**PowerShell**

```powershell
$id = 1

Invoke-RestMethod \
  -Uri "http://localhost:8080/quotes/$id" \
  -Method DELETE
```

**curl.exe**

```powershell
$id=1
curl.exe -X DELETE http://localhost:8080/quotes/$id
```

---

## Notes

* Ensure the API server is running on `localhost:8080` before executing commands.
* PowerShell examples use backticks (\`\`\`) for line continuation.
* `Invoke-RestMethod` automatically parses JSON responses into PowerShell objects.
* For `curl.exe` on Windows, avoid PowerShell's built-in `curl` alias by specifying `curl.exe`.

---

Happy quoting!
