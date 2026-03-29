import { ref } from 'vue'

/**
 * 获取 CSS 变量的值
 */
function getCssVariable(variableName: string): string {
    return getComputedStyle(document.documentElement)
        .getPropertyValue(variableName)
        .trim()
}

// 全局 CSS 变量 ref
const majorBgColor = ref(getCssVariable('--billadm-color-major-background'))
const minorBgColor = ref(getCssVariable('--billadm-color-minor-background'))
const siderWidthSize = ref(getCssVariable('--billadm-size-sider-width'))
const positiveColor = ref(getCssVariable('--billadm-color-positive'))
const negativeColor = ref(getCssVariable('--billadm-color-negative'))
const hoverBgColor = ref(getCssVariable('--billadm-color-icon-hover-bg'))
const uiSizeMenuWidth = ref(getCssVariable('--billadm-ui-size-menu-width'))
const headerHeight = ref(getCssVariable('--billadm-size-header-height'))
const textMajor = ref(getCssVariable('--billadm-color-text-major'))
const textMinor = ref(getCssVariable('--billadm-color-text-minor'))

/**
 * 刷新所有 CSS 变量（用于动态主题切换后）
 */
function refreshCssVariables() {
    majorBgColor.value = getCssVariable('--billadm-color-major-background')
    minorBgColor.value = getCssVariable('--billadm-color-minor-background')
    siderWidthSize.value = getCssVariable('--billadm-size-sider-width')
    positiveColor.value = getCssVariable('--billadm-color-positive')
    negativeColor.value = getCssVariable('--billadm-color-negative')
    hoverBgColor.value = getCssVariable('--billadm-color-icon-hover-bg')
    uiSizeMenuWidth.value = getCssVariable('--billadm-ui-size-menu-width')
    headerHeight.value = getCssVariable('--billadm-size-header-height')
    textMajor.value = getCssVariable('--billadm-color-text-major')
    textMinor.value = getCssVariable('--billadm-color-text-minor')
}

export function useCssVariables() {
    return {
        // 背景色
        majorBgColor,
        minorBgColor,
        siderWidthSize,
        // 颜色
        positiveColor,
        negativeColor,
        hoverBgColor,
        uiSizeMenuWidth,
        // 布局
        headerHeight,
        // 文字
        textMajor,
        textMinor,
        // 方法
        refreshCssVariables,
    }
}

/**
 * 切换主题
 */
export function useTheme() {
    const setTheme = (theme: 'light' | 'dark') => {
        document.documentElement.setAttribute('data-theme', theme)
        // 主题切换后刷新 CSS 变量
        setTimeout(refreshCssVariables, 0)
    }

    return {
        setTheme,
    }
}
