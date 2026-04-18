### E2E Tests: Admin Settings Case-Study Workflow

**Suite ID:** `ADMIN-SETTINGS-E2E`
**Feature:** Settings hub navigation and admin case-study workflow orchestration.

---

## Test Case: `ADMIN-SETTINGS-E2E-001` - Complete canonical-first happy path

**Priority:** `high`

**Tags:**
- type → @e2e
- feature → @admin-settings

**Description/Objective:** Verify Settings stays the entry point, the workflow runs on a dedicated subpage, and admins can confirm/start publish and import in the MVP canonical-first flow.

### Flow Steps:
1. Open `/admin/settings` as an authenticated admin.
2. Navigate to the case-study workflow subpage.
3. Start a run from an existing canonical source.
4. Confirm and execute publish, then confirm and execute import.

### Expected Result:
- Branding settings remain on `/admin/settings` only.
- The workflow page shows the MVP generation-unavailable message.
- Published canonical URL and imported project ID are visible with operator logs.

---

## Test Case: `ADMIN-SETTINGS-E2E-002` - Resume after failed import

**Priority:** `high`

**Tags:**
- type → @e2e
- feature → @admin-settings

**Description/Objective:** Verify a failed import can resume from saved workflow state without losing previously published canonical evidence.

### Flow Steps:
1. Open a saved failed run on `/admin/settings/case-studies?run=...`.
2. Inspect the persisted failure and published canonical URL.
3. Resume from the latest checkpoint.

### Expected Result:
- The run initially shows failed state with preserved canonical URL.
- Resume completes import successfully.
- Logs show resume progress and the project ID becomes visible.
