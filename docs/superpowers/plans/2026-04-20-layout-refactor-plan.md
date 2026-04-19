# Layout Refactor Implementation Plan

> **For agentic workers:** Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Move bottom bar inside content area so pages that don't need it don't have it in DOM.

**Architecture:** Simple structural refactor - move `<footer class="app-footer">` from sibling to `app-main` to inside `app-content` as sibling to `router-view`. Add flexbox to content area so router-view fills space.

**Tech Stack:** Vue 3, Scoped CSS

---

## Task 1: Refactor Layout.vue Structure

**Files:**
- Modify: `app/src/components/Layout.vue`

- [ ] **Step 1: Read current Layout.vue**

Verify current structure before making changes.

- [ ] **Step 2: Move footer inside app-content**

Change this:
```html
<div class="app-main">
  <aside class="app-sidebar">
    <app-left-bar />
  </aside>
  <main class="app-content">
    <router-view />
  </main>
</div>
<footer class="app-footer">
  <app-bottom-bar />
</footer>
```

To this:
```html
<div class="app-main">
  <aside class="app-sidebar">
    <app-left-bar />
  </aside>
  <main class="app-content">
    <router-view />
    <footer class="app-footer">
      <app-bottom-bar />
    </footer>
  </main>
</div>
```

- [ ] **Step 3: Update .app-content CSS**

Change `.app-content` from:
```css
.app-content {
  flex: 1;
  background-color: var(--billadm-color-major-warm);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}
```

To:
```css
.app-content {
  display: flex;
  flex-direction: column;
  flex: 1;
  background-color: var(--billadm-color-major-warm);
  overflow: hidden;
}
```

- [ ] **Step 4: Add flex: 1 to router-view**

Add a wrapper class around router-view or use CSS to make router-view fill available space:

```html
<router-view class="app-router-view" />
```

```css
.app-router-view {
  flex: 1;
  overflow: auto;
}
```

- [ ] **Step 5: Commit**

```bash
git add app/src/components/Layout.vue
git commit -m "refactor(ui): move footer inside content area"
```

---

## Verification

After refactor:
1. Run `npm run dev` in app directory
2. Navigate to Ledger page - bottom bar should not appear
3. Navigate to Settings page - bottom bar should not appear
4. Navigate to Transaction Record page - bottom bar should appear
5. Navigate to Data Analysis page - bottom bar should appear
