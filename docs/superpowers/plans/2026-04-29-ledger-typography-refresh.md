# Ledger Page Typography Refresh — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Elevate LedgerView card typography and visual details to match the "Warm Editorial Precision" design system.

**Architecture:** Single-file CSS + template change to `LedgerView.vue`. No new files, no logic changes, no test changes (visual-only).

**Tech Stack:** Vue 3 SFC, SCSS (scoped), existing design tokens

---

### Task 1: Typography hierarchy — name, description, date label

**Files:**
- Modify: `app/src/components/ledger_view/LedgerView.vue` (CSS only)

- [ ] **Step 1: Update `.ledger-name` to Playfair Display**

Replace the current `.ledger-name` block:

```css
.ledger-name {
  font-family: var(--billadm-font-display);
  font-size: var(--billadm-size-text-title);
  font-weight: 600;
  color: var(--billadm-color-text-major);
  margin: 0 0 var(--billadm-space-xs) 0;
  line-height: var(--billadm-height-tight);
  letter-spacing: -0.02em;
}
```

- [ ] **Step 2: Update `.ledger-desc` to Source Serif 4 italic**

Replace the current `.ledger-desc` block:

```css
.ledger-desc {
  font-family: var(--billadm-font-body);
  font-size: var(--billadm-size-text-body);
  color: var(--billadm-color-text-secondary);
  margin: 0;
  line-height: var(--billadm-height-snug);
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  font-style: italic;
}

.ledger-desc.is-empty {
  color: var(--billadm-color-text-disabled);
  font-style: italic;
}
```

- [ ] **Step 3: Update `.ledger-meta-item` to small-caps style**

Replace the current `.ledger-meta-item` block:

```css
.ledger-meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: var(--billadm-size-text-small);
  color: var(--billadm-color-text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}
```

- [ ] **Step 4: Update `.empty-state-title` to Playfair Display**

Add `font-family: var(--billadm-font-display);` to `.empty-state-title`:

```css
.empty-state-title {
  font-family: var(--billadm-font-display);
  font-size: var(--billadm-size-text-display-sm);
  font-weight: var(--billadm-weight-bold);
  color: var(--billadm-color-text-major);
  margin: 0 0 var(--billadm-space-sm) 0;
  letter-spacing: -0.01em;
}
```

- [ ] **Step 5: Commit typography changes**

```bash
git add app/src/components/ledger_view/LedgerView.vue
git commit -m "refactor(ui): apply editorial typography hierarchy to ledger cards"
```

---

### Task 2: Card visual treatment — background, accent bar, double-line ornament

**Files:**
- Modify: `app/src/components/ledger_view/LedgerView.vue` (CSS only)

- [ ] **Step 1: Update `.ledger-card` — warm background, left accent bar**

Replace the current `.ledger-card` block (remove `border-top`, add `border-left`, change background):

```css
.ledger-card {
  border-radius: var(--billadm-radius-lg);
  background-color: var(--billadm-color-major-warm);
  border: 1px solid var(--billadm-color-window-border);
  border-left: 4px solid var(--ledger-accent, var(--billadm-ledger-forest));
  transition: box-shadow var(--billadm-transition-normal), transform var(--billadm-transition-fast);
  position: relative;
  overflow: hidden;
}
```

- [ ] **Step 2: Update `.ledger-card:hover` — keep hover on left accent**

Since `border-top-color` no longer exists, update hover to change left border color:

```css
.ledger-card:hover {
  box-shadow: var(--billadm-shadow-lg);
  transform: translateY(-2px);
  border-left-color: var(--billadm-color-primary);
}
```

- [ ] **Step 3: Add featured card double-line ornament**

Add new CSS rule for the double-line decoration on the first card:

```css
/* Featured card double-line ornament (book spine detail) */
.ledger-card.is-featured::before {
  content: '';
  position: absolute;
  top: var(--billadm-space-md);
  left: var(--billadm-space-md);
  right: var(--billadm-space-md);
  height: 4px;
  border-top: 1px solid var(--billadm-color-divider);
  border-bottom: 1px solid var(--billadm-color-divider);
  pointer-events: none;
}
```

- [ ] **Step 4: Remove featured card enlarged icon rules**

Delete the following CSS block (no longer needed — hierarchy via typography, not icon size):

```css
.ledger-card.is-featured .ledger-icon {
  width: 48px;
  height: 48px;
}

.ledger-card.is-featured .ledger-icon svg {
  width: 24px;
  height: 24px;
}
```

- [ ] **Step 5: Commit card treatment changes**

```bash
git add app/src/components/ledger_view/LedgerView.vue
git commit -m "refactor(ui): add editorial card treatment with left accent bar and double-line ornament"
```

---

### Task 3: Icon and empty state refinements

**Files:**
- Modify: `app/src/components/ledger_view/LedgerView.vue` (template + CSS)

- [ ] **Step 1: Reduce icon background opacity from 15% to 10%**

In the template, change the inline `backgroundColor` on `.ledger-icon`:

```html
<!-- Before -->
:style="{
  backgroundColor: `color-mix(in srgb, var(${ledgerColorVars[index % ledgerColorVars.length]}) 15%, transparent)`,
  ...
}"

<!-- After -->
:style="{
  backgroundColor: `color-mix(in srgb, var(${ledgerColorVars[index % ledgerColorVars.length]}) 10%, transparent)`,
  ...
}"
```

- [ ] **Step 2: Increase SVG stroke-width from 1.5 to 1.8**

Change all `stroke-width="1.5"` → `stroke-width="1.8"` and `stroke-width="1.2"` → `stroke-width="1.5"` in the ledger card's inline SVGs:

**Ledger icon SVG:**
```html
<svg viewBox="0 0 24 24" fill="none">
  <path d="M4 19.5A2.5 2.5 0 016.5 17H20" stroke="currentColor" stroke-width="1.8"
    stroke-linecap="round" stroke-linejoin="round" />
  <path d="M6.5 2H20v20H6.5A2.5 2.5 0 014 19.5v-15A2.5 2.5 0 016.5 2z" stroke="currentColor"
    stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" />
</svg>
```

**Date meta icon SVG:**
```html
<svg class="meta-icon" viewBox="0 0 16 16" fill="none">
  <rect x="2" y="3" width="12" height="11" rx="1.5" stroke="currentColor" stroke-width="1.5" />
  <path d="M5 1.5v3M11 1.5v3M2 7h12" stroke="currentColor" stroke-width="1.5"
    stroke-linecap="round" />
</svg>
```

**Edit action button SVG:**
```html
<svg viewBox="0 0 16 16" fill="none">
  <path d="M11.5 2.5l2 2-8 8H3.5v-2l8-8z" stroke="currentColor" stroke-width="1.5"
    stroke-linecap="round" stroke-linejoin="round" />
</svg>
```

**Delete action button SVG:**
```html
<svg viewBox="0 0 16 16" fill="none">
  <path
    d="M3 5h10M6 5V4a1 1 0 011-1h2a1 1 0 011 1v1M11 5v7a1.5 1.5 0 01-1.5 1.5H6.5A1.5 1.5 0 015 12V5"
    stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
</svg>
```

- [ ] **Step 3: Reduce empty state icon opacity and hover**

Change `.empty-state-visual` opacity from `0.25` → `0.15` and hover from `0.35` → `0.25`:

```css
.empty-state-visual {
  width: 96px;
  height: 96px;
  margin-bottom: var(--billadm-space-xl);
  color: var(--billadm-color-primary);
  opacity: 0.15;
  transform: scale(1);
  transition: transform var(--billadm-transition-slow);
}

.empty-state:hover .empty-state-visual {
  transform: scale(1.05);
  opacity: 0.25;
}
```

- [ ] **Step 4: Commit icon and empty state refinements**

```bash
git add app/src/components/ledger_view/LedgerView.vue
git commit -m "refactor(ui): refine ledger card icons and empty state visuals"
```

---

### Task 4: Build verification

**Files:** None (verification only)

- [ ] **Step 1: Verify the frontend builds**

```bash
cd app && npm run build 2>&1 | tail -20
```

Expected: Build succeeds with no errors.

- [ ] **Step 2: Visual verification checklist**

Manually verify in the running app:
- [ ] Ledger names render in Playfair Display (serif editorial font)
- [ ] Descriptions render in italic Source Serif 4
- [ ] Date labels appear in small uppercase
- [ ] Card background is warm off-white (distinct from page)
- [ ] Left accent bar uses the correct ledger color
- [ ] First card (when ≥3 ledgers) shows double-line ornament
- [ ] Icons are more restrained (10% background opacity)
- [ ] Empty state title uses Playfair Display
- [ ] No layout regressions — cards still in grid, responsive breakpoints work
