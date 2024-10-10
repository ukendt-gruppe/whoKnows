# Project Structure

This diagram represents the dependency structure of our project:

```mermaid
graph LR
    subgraph External Dependencies
        AK[github.com/mattn/go-sqlite3]
        AL[gorilla/sessions]
        AM[golang.org/x/crypto/bcrypt]
    end

    subgraph Utils
        AH[utils/password.go]
        AI[bcrypt.GenerateFromPassword]
        AJ[bcrypt.CompareHashAndPassword]
        Y[utils.JSONResponse]
    end

    subgraph Database
        AA[db/db.go]
        AB[sql.Open]
        AC[os.ReadFile]
        AD[DB.Exec]
        AE[DB.Query]
        AF[DB.QueryRow]
        AG[utils.HashPassword]
        B[db.InitDB]
        K[db.DB]
        M[db.QueryDB]
        O[db.GetUser]
        Q[db.CreateUser]
    end

    subgraph API Handlers
        S[handlers/api.go]
        T[Search]
        U[Weather]
        V[Register]
        W[Login]
        X[Logout]
    end

    subgraph Web Handlers
        C[handlers.SearchHandler]
        D[handlers.LoginHandler]
        E[handlers.RegisterHandler]
        F[handlers.LogoutHandler]
        G[handlers.AboutHandler]
        H[handlers.WeatherHandler]
    end

    subgraph Middleware
        I[middleware.LoggingMiddleware]
        J[middleware.SessionMiddleware]
    end

    subgraph Main
        A[main.go]
    end

    AH --> AM
    AI --> AH
    AJ --> AH

    AA --> AK
    AA --> AB
    AA --> AC
    AA --> AD
    AA --> AE
    AA --> AF
    AA --> AG
    B --> K
    B --> L[os.ReadFile]
    M --> AA
    O --> AA
    Q --> AA

    S --> T
    S --> U
    S --> V
    S --> W
    S --> X
    T --> M
    T --> Y
    U --> Y
    V --> Z[r.ParseForm]
    V --> Y
    W --> Z
    W --> Y
    X --> Y

    C --> M
    C --> N[templates.ExecuteTemplate]
    D --> O
    D --> P[utils.CheckPasswordHash]
    D --> N
    E --> O
    E --> Q
    E --> N
    F --> N
    G --> N
    H --> N

    J --> R[sessions.Store]
    J --> O
    J --> AL

    A --> B
    A --> C
    A --> D
    A --> E
    A --> F
    A --> G
    A --> H
    A --> I
    A --> J

    classDef default fill:#f9f,stroke:#333,stroke-width:2px,color:black;
    classDef external fill:#ff9,stroke:#333,stroke-width:2px,color:black;
    class AK,AL,AM external;