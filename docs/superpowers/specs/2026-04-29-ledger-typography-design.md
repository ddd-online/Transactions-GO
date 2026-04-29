# Ledger Page Typography Refresh

**Date**: 2026-04-29
**Scope**: `app/src/components/ledger_view/LedgerView.vue`

## Objective

The ledger page's grid of cards lacks typographic personality and fails to express the app's "Warm Editorial Precision" design language. This is a focused typography and visual-detail refresh — layout structure, component choices, and interaction patterns stay as-is.

## Design Direction

Editorial card treatment — using the design system's existing typefaces (Playfair Display, Source Serif 4) to create typographic hierarchy within each ledger card, plus subtle decorative details inspired by fine bookbinding.

---

## Typography Hierarchy

| Element | Before | After |
|---------|--------|-------|
| Ledger name | 16px body font, semibold | 20px Playfair Display, 600 weight, -0.02em letter-spacing |
| Ledger description | 13px body font | 14px Source Serif 4, italic (when present), disabled color + italic (when empty) |
| Date label | 12px caption, body font | 11px body font, uppercase, 0.04em letter-spacing (small-caps style) |

Font family assignments use existing CSS variables (`--billadm-font-display`, `--billadm-font-body`). No new font imports needed.

## Card Visual Treatment

**Before**: White background, 1px gray border, 3px colored top stripe.

**After**:
- Background: `--billadm-color-major-warm` (micro-warm off-white), creating subtle surface distinction from the page
- Left accent bar: 4px wide colored strip on the left edge (replaces top stripe)
- Top double-line ornament: two 1px horizontal lines at top of featured (first) card only, inspired by book spine tooling — uses `--billadm-color-divider`
- Border: keep 1px `--billadm-color-window-border`, unchanged
- Hover: keep existing lift + shadow, unchanged

## Icon Treatment

- Icon background opacity: reduce from `color-mix(15%)` to `color-mix(10%)` — more restrained
- SVG stroke-width: increase from 1.5 to 1.8 — more grounded

## Featured Card

The first card (when ≥3 ledgers) keeps `grid-column: span 2`. The enlarged icon behavior is removed — hierarchy comes from typography and the double-line ornament, not from scaling elements.

## Empty State

- Title: switch to Playfair Display (currently inherits body font)
- Icon opacity: reduce from 0.25 to 0.15

## Modal

No changes needed — the global Ant Design Modal overrides already apply Playfair Display to modal titles.

## Floating Button

No changes — already consistent with the design system.

## Non-Goals

- Grid layout changes (column count, responsive breakpoints, span rules stay)
- New features or data additions
- Animation/motion changes beyond what exists
- Dark mode changes (existing dark mode tokens will apply automatically)
