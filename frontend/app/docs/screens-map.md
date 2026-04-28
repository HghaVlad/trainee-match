# Screens Map (v2 — approved)

> Approved by product owner. Source for all routing and navigation work.

## 1. Architectural decisions

- **Multi-company**: a user may belong to several companies via `recruiter` / `admin` roles. Active company stored in session; URL path always carries `companyId` for unambiguous deep-linking.
- **`/auth/me`** is assumed available; bootstrap reads it.
- **`/companies/me`** returns paginated `[{id, name, openVacanciesCount, logoKey, createdAt, role}]` (role added per agreement; if backend ships without role we'll fetch members for active company).
- **Application API** is cookie-authenticated (httpOnly) like the rest. Uses snapshots, immutable after apply.
- **Resume status** will be `draft|published` enum once backend fixes contract; until then a thin in/out adapter normalizes.
- **Default resume**: stored in localStorage keyed by userId. Apply modal preselects it; user may override in dropdown.
- Polnotext `q` search, PDF/DOCX upload, logo upload, admin section — **deferred / hidden** in UI.

## 2. Routes

### Public / Auth
| Route | Purpose | Main API |
|---|---|---|
| `/` | Landing | (opt) `GET /vacancies?limit=6` |
| `/login` | Login by `username + password` | `POST /auth/login` → `GET /auth/me` (+ `GET /companies/me` if Company) |
| `/register` | Register `Candidate` or `Company` | `POST /auth/register` |
| `/403` | Forbidden | — |
| `*` | NotFound | — |

### Vacancy catalog (public)
| Route | Purpose | Main API |
|---|---|---|
| `/vacancies` | List + filters + sort + cursor | `GET /vacancies` |
| `/vacancies/:vacancyId` | Details + Apply CTA | `GET /vacancies/{vacancy-id}` |

### Company catalog (public)
| Route | Purpose | Main API |
|---|---|---|
| `/companies` | List | `GET /companies` |
| `/companies/:companyId` | Profile + their vacancies | `GET /companies/{id}`, `GET /vacancies?company_id=` |

### Candidate area
| Route | Purpose | Main API |
|---|---|---|
| `/me/profile` | Contact profile | `GET /candidate/me`, `POST/PATCH /candidate/` |
| `/me/resumes` | List + default flag | `GET /resume`, `PATCH /resume/{id}` |
| `/me/resumes/new` | Create | `POST /resume/` |
| `/me/resumes/:resumeId` | Editor + publish | `GET/PATCH /resume/{id}`, `GET /skill/list` |
| `/me/applications` | My applications | `GET /api/v1/applications` |
| `/me/applications/:applicationId` | Details + history + withdraw | `GET /api/v1/applications/{id}` + `/history` + `/withdraw` |

### Company area (multi-company)
Prefix: `/company/:companyId/...`. Active companyId synced with `CompanySwitcher`.
| Route | Purpose | API | Access |
|---|---|---|---|
| `/company/new` | Create company | `POST /companies` | Any Company role |
| `/company` | Redirect to active company dashboard or `/company/new` if none | `GET /companies/me` | Company |
| `/company/:companyId/dashboard` | Analytics: summary + funnel + dynamics | `/api/v1/hr/companies/{id}/analytics/...` | member |
| `/company/:companyId/profile` | View / edit company | `GET/PATCH/DELETE /companies/{id}` | edit/delete = admin |
| `/company/:companyId/members` | Members list + add/edit/delete | `GET (new) / POST / PATCH / DELETE /companies/{id}/members` | view: any; manage: admin |
| `/company/:companyId/vacancies` | Vacancies list with status filter | `GET /companies/{id}/vacancies` | member |
| `/company/:companyId/vacancies/new` | Create | `POST /companies/{id}/vacancies` | member |
| `/company/:companyId/vacancies/:vacancyId` | Edit + publish/archive/delete | `GET/PATCH/DELETE` + `/publish` + `/archive` | member; delete = admin |
| `/company/:companyId/vacancies/:vacancyId/applications` | Vacancy applications + filters | `GET /api/v1/hr/vacancies/{vacancyId}/applications` | member |
| `/company/:companyId/vacancies/:vacancyId/analytics` | Vacancy analytics | `/api/v1/hr/vacancies/{id}/analytics/...` | member |
| `/company/:companyId/applications` | All company applications | `GET /api/v1/hr/companies/{companyId}/applications` | member |
| `/company/:companyId/applications/:applicationId` | Application details + history + change status | `GET /api/v1/hr/applications/{id}` + `/history` + `/status` | member |

## 3. Modals (business-critical)

1. **Apply Vacancy** — preselects default published resume; allows override; optional `coverLetter` (max 2000).
2. **Withdraw Application** — confirm + optional comment.
3. **Change Application Status (HR)** — uses `allowedActions[]` from API (no client state machine).
4. **Confirm Publish Resume** — final validation message.
5. **Re-edit Published Resume** — warns about implicit `draft` transition.
6. **Confirm Archive / Delete Vacancy**.
7. **Add Member to Company** (admin only).
8. **Confirm Delete Company** (admin only).
9. **CompanySwitcher** — Header dropdown, not a modal.

## 4. Status semantics

`ApplicationStatus`: `submitted | seen | interview | rejected | offer | withdrawn`.

- Active = `submitted | seen | interview` (1 active per vacancy per candidate; reapply allowed once previous is `withdrawn|rejected`).
- Terminal for HR action buttons (UI disabled): `rejected | offer | withdrawn`.
- Allowed actions are server-driven via `allowedActions[]` to avoid duplicating state machine on the client.

## 5. Default resume rule

- Stored in `localStorage` under key `tm.defaultResumeId.<userId>`.
- Cleared on resume delete or status change to `draft` if it was the default published one.
- UI: star icon on `/me/resumes` row marks default; only one at a time.

## 6. Hidden / deferred

- Full-text search input on `/vacancies` (no `q` in API).
- PDF/DOCX resume upload.
- Company logo upload (URL only; logoKey shown if backend produces a serving URL).
- Admin section.
- HR vacancy moderation status.

## 7. Open follow-ups (low risk, parking)

- Backend to fix resume `status: integer` to enum `draft|published`.
- Backend to add company-member list, `/companies/me`, `status` field on company-vacancy list item.
