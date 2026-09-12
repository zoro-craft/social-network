# 🐛 Groups Feature — Audit & Bug Tracker

Tracking known security vulnerabilities, database logic flaws, and UI/UX issues in the **Groups** feature (Go backend + Next.js frontend).

**Legend:** 🔴 Critical · 🟠 High (data integrity) · 🟡 Medium (UX)

## 📑 Summary

| # | Issue | Severity | Area |
|---|-------|----------|------|
| 1 | Membership control bypass | 🔴 | Backend |
| 2 | Private comments leak | 🔴 | Backend |
| 3 | Private posts leak | 🔴 | Backend |
| 4 | User enumeration | 🔴 | Backend |
| 5 | Group info leak / wrong HTTP method | 🔴 | Backend |
| 6 | Rejected requests cause lockout | 🟠 | Database |
| 7 | Premature DB insertion | 🟠 | Database |
| 8 | Re-invitation blocked | 🟠 | Database |
| 9 | Duplicate member crash | 🟠 | Database |
| 10 | Missing DB transactions | 🟠 | Database |
| 11 | Silent UI failures / stale state | 🟡 | Frontend |
| 12 | Broken group detail page | 🟡 | Frontend |

---

## 🔴 Critical — Security & Access Control

- [ ] **1. Membership control bypass**
  Any user can invite themselves to a group and accept the invite. `ValidCreateInvitation` doesn't check who the caller is.
  📍 `backend/internal/groups/services/invitation_validator.go`

- [ ] **2. Private comments leak**
  The comments API doesn't check group membership, so any logged-in user can read private discussion threads.
  📍 `backend/internal/comments/getcommentbygoup.go`

- [ ] **3. Private posts leak**
  Group posts are forced to `public`. Non-members can view posts, comments, and reactions through unrelated public endpoints (e.g. user profiles).
  📍 `backend/internal/posts/GetSinglePost.go`, `backend/internal/profile/profile.go`

- [ ] **4. User enumeration**
  The "get users not in group" endpoint has no authorization check, so anyone can scrape the full user list.
  📍 `backend/internal/users/users.go`

- [ ] **5. Group info leak & wrong HTTP method**
  `GetGroupInfo` exposes group metadata to non-members and uses `POST` for what should be a `GET` (read-only) request.
  📍 `backend/internal/groups/handlers/get_all_groups_handler.go`

---

## 🟠 High — Database Integrity & Backend Logic

- [ ] **6. Rejected requests cause lockout**
  Rejected join requests stay in the database. A user trying again hits the `UNIQUE(group_id, user_id)` constraint, causing a `500` error instead of a clean `400`.
  📍 `backend/internal/groups/repositories/join_repository.go`

- [ ] **7. Premature DB insertion**
  Join requests are written to the database *before* checking that the creator or group actually exists, leaving orphaned rows behind on failure.
  📍 `backend/internal/groups/handlers/join_request_handler.go`

- [ ] **8. Re-invitation blocked**
  `SaveInvitation` treats old rejected invitations as still active, so admins can't re-invite someone who declined before.
  📍 `backend/internal/groups/repositories/invitation_repository.go`

- [ ] **9. Duplicate member crash**
  Accepting an invite for someone who's already a member triggers a SQLite `UNIQUE` constraint error and crashes the request.
  📍 `backend/internal/groups/handlers/invitation_response_handler.go`

- [ ] **10. Missing DB transactions**
  Accept flows (invites/join requests) aren't atomic. If adding the member succeeds but deleting the notification fails, the database is left in a partially-updated state.
  📍 `backend/internal/groups/handlers/invitation_response_handler.go`

---

## 🟡 Medium — Frontend & UX (Next.js)

- [ ] **11. Silent UI failures & stale state**
  The groups list uses stale cached data. `GroupCard` has no loading/success feedback, so users click repeatedly and trigger `500` errors.
  📍 `frontend/src/app/groups/page.js`, `frontend/src/app/groups/GroupCard.jsx`

- [ ] **12. Broken group detail page**
  Shows admin-only buttons (Invite, Create Post) to non-members, which fail with `403` on click. Also shows a blank white page if the group ID doesn't exist.
  📍 `frontend/src/app/groups/[id]/page.js`