# Validation Contracts

Every user-facing form must have a documented validation contract before implementation. Frontend Zod validation and backend DTO validation must both match this contract.

This file is the single source of truth for all user-facing form validation.

Rules:

- Check this file before implementing any new form.
- If a validation contract does not exist, add it before implementation.
- Frontend validation must use Zod and is for user experience.
- Backend request DTO validation must enforce the same rules and is the final authority.
- Never implement validation on only one side.
- Any validation change must update:
  1. `docs/validation-contracts.md`
  2. Frontend Zod schema
  3. Backend request DTO validation
- Backend validation errors should return field-level messages using the standard validation error response from `docs/api.md`.

---

## Auth Forms

### Register

Fields:

#### name

- required
- trim
- string
- min 2
- max 80

#### username

- optional
- trim
- lowercase
- min 3
- max 30
- allowed characters: lowercase letters, numbers, underscore
- regex: `^[a-z0-9_]+$`

#### email

- required
- trim
- lowercase
- valid email
- max 255

#### targetBand

- required
- number
- allowed values: `5.0`, `5.5`, `6.0`, `6.5`, `7.0`, `7.5`, `8.0`, `8.5`

#### password

- required
- min 8
- max 72
- must contain at least one letter
- must contain at least one number

---

### Login

Fields:

#### email

- required
- trim
- lowercase
- valid email
- max 255

#### password

- required
- min 8
- max 72
