import {ref} from 'vue'

/**
 * 获取 CSS 变量的值
 */
function getCssVariable(variableName: string) {
    return getComputedStyle(document.documentElement)
        .getPropertyValue(variableName)
        .trim()
}

export function useCssVariables() {
    const majorBgColor = ref(getCssVariable('--billadm-color-major-background'))
    const minorBgColor = ref(getCssVariable('--billadm-color-minor-background'))
    const siderWidthSize = ref(getCssVariable('--billadm-size-sider-width'))
    const positiveColor = ref(getCssVariable('--billadm-color-positive'))
    const negativeColor = ref(getCssVariable('--billadm-color-negative'))

    const hoverBgColor = ref(getCssVariable('--billadm-color-icon-hover-bg-color'))
    const uiSizeMenuWidth = ref(getCssVariable('--billadm-ui-size-menu-width'))

    return {
        majorBgColor,
        minorBgColor,
        siderWidthSize,
        positiveColor,
        negativeColor,

        hoverBgColor,
        uiSizeMenuWidth,
    }
}