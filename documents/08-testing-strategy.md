# Stage 8 — Testing Strategy

## The Testing Pyramid

```
        /\
       /  \        E2E Tests
      /    \       skip for v1 — overkill for 50 users
     /──────\
    /        \     Integration Tests
   /          \    a few critical ones
  /────────────\
 /              \  Unit Tests
/________________\ focus here — fast, high value
```

## What to Test

| What | Type | Priority | Why |
|------|------|----------|-----|
| Point calculation logic | Unit | 🔴 Must have | Silent wrongness — wrong scores, trust collapses undetected |
| Invite token sign/verify | Unit | 🔴 Must have | Security — wrong means anyone joins any leaderboard |
| Prediction upsert unique constraint | Integration | 🟡 Should have | DB-level behaviour, can't unit test |
| OAuth callback + JWT issue | Integration | 🟡 Should have | Auth is high stakes |
| Full prediction → leaderboard flow | E2E | 🟢 Nice to have | Defer to v2 |

## Unit Test Examples (Go)

### Point Calculation

```go
// domain/prediction/point_test.go

func TestCalculatePoint(t *testing.T) {
    tests := []struct {
        name       string
        match      Match
        prediction string
        want       int
    }{
        {"home win correct",  Match{HomeScore: 2, AwayScore: 0}, "home_win", 1},
        {"draw correct",      Match{HomeScore: 1, AwayScore: 1}, "draw",     1},
        {"away win correct",  Match{HomeScore: 0, AwayScore: 2}, "away_win", 1},
        {"home win wrong",    Match{HomeScore: 2, AwayScore: 0}, "draw",     0},
        {"draw wrong",        Match{HomeScore: 1, AwayScore: 1}, "home_win", 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := CalculatePoint(tt.match, tt.prediction)
            if got != tt.want {
                t.Errorf("got %d, want %d", got, tt.want)
            }
        })
    }
}
```

### Invite Token

```go
// domain/auth/invite_test.go

func TestInviteToken_SignAndVerify(t *testing.T) {
    leaderboardID := "some-uuid"
    secret := "test-secret"

    token, err := SignInviteToken(leaderboardID, secret, 10*time.Minute)
    if err != nil {
        t.Fatal(err)
    }

    got, err := VerifyInviteToken(token, secret)
    if err != nil {
        t.Fatalf("valid token rejected: %v", err)
    }
    if got != leaderboardID {
        t.Errorf("got leaderboard %s, want %s", got, leaderboardID)
    }
}

func TestInviteToken_Expired(t *testing.T) {
    token, _ := SignInviteToken("uuid", "secret", -1*time.Minute) // already expired

    _, err := VerifyInviteToken(token, "secret")
    if err == nil {
        t.Error("expected error for expired token, got nil")
    }
}
```

## Rules

- Test files sit next to the files they test (`prediction_test.go` next to `prediction.go`)
- Table-driven tests for any function with multiple input cases
- Test names: `TestFunctionName_Scenario_ExpectedResult`
- Never test implementation details — test behaviour
- Unit tests in `domain/` have zero external dependencies
- Integration tests use a real test DB (spin up via Docker in CI)

## SA Principle
>
> Test behaviour, not implementation.
> Tests should break when the system does the wrong thing —
> not when you rename a function.
