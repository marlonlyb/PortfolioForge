### E2E Tests: Storefront smoke review

**Suite ID:** `STORE-SMOKE`
**Feature:** Basic local visual review for the public storefront.

---

## Test Case: `STORE-SMOKE-001` - Home page fixture renders

**Priority:** `critical`

**Tags:**
- type → @e2e
- feature → @storefront

**Description/Objective:** Verify the local Vite app can boot with mocked public data and show the storefront catalog.

**Preconditions:**
- Vite dev server runs through Playwright webServer.
- Public API routes are intercepted with fixture responses.

### Flow Steps:
1. Open `/`.
2. Wait for the catalog fixture card.
3. Confirm the project link is visible.

### Expected Result:
- The page loads without live backend dependencies.
- The catalog shows the smoke fixture project.

## Test Case: `STORE-SMOKE-002` - Detail page fixture opens

**Priority:** `critical`

**Tags:**
- type → @e2e
- feature → @storefront

**Description/Objective:** Verify the smoke flow can move from the catalog to a project detail page.

**Preconditions:**
- Same mocked public API fixture as the home smoke test.

### Flow Steps:
1. Open `/`.
2. Click the fixture project card.
3. Verify the detail heading and assistant toggle.

### Expected Result:
- The project detail route resolves correctly.
- The public detail view renders the expected fixture content.

## Test Case: `STORE-SMOKE-003` - Landing search stays usable across breakpoints

**Priority:** `critical`

**Tags:**
- type → @e2e
- feature → @storefront

**Description/Objective:** Verify the ChatMLB landing search keeps the prompt → suggestion → result flow usable on representative mobile, tablet, and desktop viewports.

**Preconditions:**
- Same mocked public API fixture as the other storefront smoke tests.
- Locale is forced to English for stable prompt and control labels.

### Flow Steps:
1. Open `/` in mobile, tablet, and desktop viewports.
2. Confirm the landing hero, search input, and Printer 05 quick prompt remain visible.
3. Click the quick prompt and verify the landing result card appears.
4. Clear the input, type `Print`, and select the preserved landing suggestion.

### Expected Result:
- The landing search UI remains reachable and clickable at each breakpoint.
- The quick prompt still feeds the existing catalog/result flow.
- Suggestion selection still updates the search field and preserves the catalog result path.
