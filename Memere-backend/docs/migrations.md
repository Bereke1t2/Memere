# Migration Strategy — Expand/Contract (Zero-Downtime)

## Rule

Every database migration **must be backward-compatible** with the previous
application version. This allows a rolling update — some pods still run the
old code while others run the new code — without downtime or data loss.

**Required pattern: Expand → Deploy → Contract**

| Phase | What you do | Both versions work? |
|---|---|---|
| **Expand** | Add new nullable column, new table, or new index | ✅ Old code ignores it |
| **Deploy** | Roll out the new application version | ✅ Reads new column when present |
| **Contract** | Add NOT NULL constraint, drop old column, etc. | ✅ Old version is gone |

Contract migrations land in a **separate, later release** — never in the
same migration as the expand.

## Why Migrations Run Before the App

In the CI/CD pipeline (and as a Kubernetes init container in Skill 5), the
migration runner executes `/migrate -direction up` against the target database
**before** the new application pods start. This means:

1. Schema changes are applied while the old pods are still serving traffic.
2. Because the migration is expand-only, the old code continues to work.
3. The new pods start; traffic shifts to them.
4. The old pods terminate. Only then is the schema fully owned by the new code.

```
Timeline:
  [old pods running]
       │
       ├── migrate up (expand-only changes)
       │
       ├── new pods start → rolling update
       │
       └── old pods terminate
```

## Verification in CI

`ci.yml` runs `make migrate-up` against the PostgreSQL service container on
every PR. If a migration fails to apply, the CI job fails before tests run.

## Checked-In Migrations

All migrations in `migrations/` follow the expand/contract rule:

| Migration | Type | Notes |
|---|---|---|
| 0001–0007 | Initial schema | Greenfield — no compatibility concern |
| 0008 | Expand | Adds `deleted_at` to sections/lessons |
| 0009–0018 | Expand | New tables, nullable columns, indexes |
| 0019 | Expand | Performance indexes only — `CREATE INDEX IF NOT EXISTS` |

No migration has ever dropped a column without a preceding expand phase.

## Checklist for New Migrations

Before merging a migration PR:

- [ ] New columns are `NULL` or have a `DEFAULT` value (old app ignores them)
- [ ] `NOT NULL` constraints only tighten an existing nullable column in a
      **separate, later** migration after the new app is fully deployed
- [ ] Dropped columns are already unused by all running application versions
- [ ] Index creation uses `CREATE INDEX IF NOT EXISTS` (idempotent)
- [ ] Migration passes `make migrate-up && make migrate-down` locally

## Secret: the `cmd/migrate` Runner

The runner is the binary embedded in the API image (`/migrate`). It uses
`golang-migrate` with an `iofs` source over the embedded `migrations/` FS,
so no external `migrate` binary is required and migrations are always in sync
with the binary that runs them.
