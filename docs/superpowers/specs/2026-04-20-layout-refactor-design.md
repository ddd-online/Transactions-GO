# Layout Refactor Design

## Context

The current `Layout.vue` has a structural issue: the bottom bar (`app-footer`) is a sibling to `app-main`, meaning it exists in the DOM regardless of which page is displayed. While `AppBottomBar` internally uses route-based visibility, the footer wrapper always renders.

Goal: Move bottom bar inside the content area so pages that don't need it don't have it in their DOM at all.

## Structure

**Before:**
```
app-shell-body (column)
├── app-header
├── app-main (row)
│   ├── app-sidebar
│   └── app-content (router-view)
└── app-footer (always rendered)
```

**After:**
```
app-shell-body (column)
├── app-header
└── app-main (row)
    ├── app-sidebar
    └── app-content (column)
        ├── router-view (flex: 1)
        └── app-footer (only mounted when route needs it)
```

## Changes

### Layout.vue

1. Move `<footer class="app-footer">` inside `<main class="app-content">` as a sibling to `<router-view />`
2. Add `display: flex; flex-direction: column;` to `.app-content`
3. Add `flex: 1` to router-view wrapper (or ensure content fills space)

### CSS Adjustments

`.app-content`:
```css
display: flex;
flex-direction: column;
```

Router-view container (implicit via parent):
```css
flex: 1;
overflow: auto;
```

## No Logic Changes

`AppBottomBar.vue` already handles route-based visibility internally via:
```js
const showStatistics = computed(() => {
  return route.path === '/tr_view' || route.path === '/da_view'
})
```
This remains unchanged.

## Files Affected

- `app/src/components/Layout.vue` — structural change only
