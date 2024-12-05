# Mandatory I

**Ukendt Gruppe**:  
Ziggy Langesen, Simon Formann, Rolf Kauffmann & Asbjørn Emil Toft

---

## Problems with the Legacy Code

We have chosen three different severity categories for the problems: **high**, **medium**, and **low**.

### High Severity Problems

- **MD5**:  
  In `app.py`, MD5 is used for password hashing. This hash function has been deprecated and should be replaced by alternatives such as `bcrypt` or `SHA-256`.

- **Outdated Dependencies**:  
  The versions of the dependencies in `requirements.txt` are from 2010, except for `setuptools`. These should be upgraded to the newest versions to address known vulnerabilities.

- **SQL Injection Vulnerability**:  
  In `app.py`, there are multiple instances where user input is directly interpolated into SQL queries. Parameterized queries should be used for all database interactions.

- **Hardcoded Sensitive Data**:  
  In `app.py`, `DATABASE_PATH` and `SECRET_KEY` are visible in the code on GitHub. These should be moved to an `.env` file or stored as GitHub secrets.

- **HTTP Protocol**:  
  The project runs on HTTP, making it vulnerable to attacks like man-in-the-middle due to the lack of encryption.

### Medium Severity Problem(s)

- **Python 2**:  
  Throughout `app.py`, the `Makefile`, and the shell script, Python 2 is used. Python 2 was sunsetted in 2020 and is no longer supported, including security vulnerabilities. Upgrading to Python 3 using a tool like `2to3` is the recommended solution.

### Low Severity Problem(s)

- **Code Organization**:  
  There is no separation of functionalities (e.g., database handling, routing, and security) into different files. This lack of structure makes the project harder to scale, maintain, and understand.

---

## Dependency Graph

We’ve created two dependency graphs:  
- **Legacy System**: [Legacy System Structure](https://github.com/ukendt-gruppe/whoKnows/blob/main/docs/PROJECT_STRUCTURE_LEGACY.png)  
- **Rewritten Solution** (Week 40): [Updated Structure](https://github.com/ukendt-gruppe/whoKnows/blob/main/docs/PROJECT_STRUCTURE.md)

From the charts, we observe how we restructured the project during the rewrite process and addressed the identified problems.

### Comparison: Legacy vs. New System

| Aspect                | Legacy System           | New System                                      |
|-----------------------|-------------------------|------------------------------------------------|
| **Framework**         | Flask (Python)         | Custom Go HTTP server with Gorilla             |
| **Database**          | Direct SQLite3 usage   | Go's `database/sql` package with SQLite3 driver|
| **Password Hashing**  | MD5 (insecure)         | bcrypt (secure)                                |
| **Dependency Mgmt**   | `requirements.txt`     | Go modules (`go.mod`)                          |
| **Code Organization** | Single `app.py` file   | Modular mono repo (`db`, `handlers`, etc.)     |

---

## Addressing Vulnerabilities

### High Severity

1. **MD5 for Password Hashing**:  
   **Mitigated**: The new system uses `bcrypt` (`utils/password.go`), a secure hashing algorithm.

2. **Outdated Dependencies**:  
   **Mitigated**: Go modules ensure the use of recent, compatible dependency versions.

3. **SQL Injection Vulnerability**:  
   **Mitigated**: Parameterized queries are implemented (`db/db.go`).

4. **Hardcoded Sensitive Data**:  
   **Partially Mitigated**: Using environment variables or config files. Work in progress with GitHub secrets.

5. **HTTP Protocol**:  
   **Not Mitigated Yet**: Planning to implement HTTPS with `http.ListenAndServeTLS()`, using SSL/TLS certificates.

### Medium Severity

1. **Python 2**:  
   **Mitigated**: Entire system rewritten in Go.

### Low Severity

1. **Code Organization**:  
   **Mitigated**: System is modular, organized into separate packages for better scalability.

---

## Additional Improvements

- **Type Safety**: Go's static typing prevents certain runtime errors.
- **Concurrency**: Goroutines and channels enable better concurrent operations.
- **Performance**: Go is typically faster than interpreted Python.
- **Middleware**: Logging and session middleware improve security and debugging.

---

## OpenAPI Specification

We implemented the **Swaggo** library to generate API documentation adhering to the OpenAPI Specification.

- **Specification File**: [Swagger JSON](https://github.com/ukendt-gruppe/whoKnows/blob/CreateOpenAPI/src/backend/docs/swagger.json)  
- **Swagger Editor**: [Swagger Editor](https://editor-next.swagger.io/)

---

## Branching Strategy

### Implementation and Enforcement

1. **Feature Development**:  
   Developers create short-lived branches for features/bug fixes.

2. **Pull Requests**:  
   - Code reviews ensure quality and consistency.  
   - Automated workflows (testing, building, deployment) via GitHub Actions.

3. **Merging**:  
   Merges occur after approval and passing CI checks. The main branch is protected.

---

### Rationale for Choice

- **Code Quality**: Ensures consistent, high-quality contributions.
- **Collaboration**: Encourages teamwork and knowledge sharing.
- **Isolation**: Feature branches prevent disruptions to the main codebase.
- **Continuous Integration**: Rapid feedback loop for issue detection.
- **Flexibility**: Promotes frequent, small deployments.

---

### Challenges Faced

1. **Work Schedule**: Limited to ~2 days/week, impacting momentum and focus.  
2. **Context Switching**: Gaps between sessions require reacclimation.  
3. **Review Coordination**: Delays due to scheduling conflicts.  
4. **Branching Consistency**: Occasional bypassing of formal processes.

---

### Improvement Suggestions

- **Sprint Planning**: Adjust for distributed schedules.  
- **Enhanced Documentation**: Reduce context-switching difficulties.  
- **Asynchronous Reviews**: Leverage GitHub's review tools between work days.  
- **Pair Programming**: Foster collaboration on complex tasks.  
- **Automated Assignments**: Assign reviewers automatically to streamline reviews.  
- **Documentation-as-Code**: Integrate docs tightly into development.  

---

## Conclusion

Our pull request-based workflow, combined with short-lived feature branches and CI/CD practices, has significantly improved collaboration and code quality. Despite challenges with workflow adaptation, we are optimistic about refining these processes to maximize their benefits.
