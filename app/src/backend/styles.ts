import type { CSSProperties } from 'vue';

/**
 * Shared button styles for edit and delete actions
 */
export const editButtonStyle: CSSProperties = {
    color: 'var(--billadm-color-positive)',
};

export const deleteButtonStyle: CSSProperties = {
    color: 'var(--billadm-color-negative)',
};

/**
 * Shared layout styles
 */
export const contentStyle: CSSProperties = {
    backgroundColor: 'var(--billadm-color-major-background)',
    overflowY: 'auto',
    marginBottom: 'auto',
};

export const headerStyle: CSSProperties = {
    height: 'auto',
    backgroundColor: 'var(--billadm-color-major-background)',
    padding: '0 0 16px 0',
    display: 'flex',
    alignItems: 'start',
    justifyContent: 'center',
};
