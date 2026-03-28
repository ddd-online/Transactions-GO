import type {TransactionRecord} from '@/types/billadm'
import {type EChartsOption, type LegendComponentOption, type PieSeriesOption,} from 'echarts'
import {TransactionTypeToColor, TransactionTypeToLabel,} from '@/backend/constant'
import dayjs from 'dayjs'

/**
 * 构造折线图配置项（按年或月分组）
 */
export function buildLineChart(
    trList: TransactionRecord[],
    options: {
        granularity: 'year' | 'month'
        lineDisplayTypes: string[]
        includeOutlier: boolean
    }
): EChartsOption {
    const {granularity, lineDisplayTypes, includeOutlier} = options

    // 过滤非异常且类型匹配的数据
    let filteredData = trList
        .filter((item) => lineDisplayTypes.includes(item.transactionType))

    if (!includeOutlier) {
        filteredData = filteredData.filter((item) => !item.outlier)
    }

    if (filteredData.length === 0) {
        return {
            tooltip: {trigger: 'axis'},
            legend: {
                data: lineDisplayTypes.map((type) => TransactionTypeToLabel.get(type)),
            } as LegendComponentOption,
            xAxis: {
                type: 'category',
                data: [],
                name: granularity === 'year' ? '年份' : '月份',
            },
            yAxis: {
                type: 'value',
                name: '金额',
            },
            series: lineDisplayTypes.map((type) => ({
                name: TransactionTypeToLabel.get(type),
                type: 'line',
                data: [],
            })),
        }
    }

    // 找出时间范围
    let minYear = Infinity,
        minMonth = Infinity
    let maxYear = -Infinity,
        maxMonth = -Infinity

    filteredData.forEach((item) => {
        const date = dayjs(item.transactionAt * 1000)
        const year = date.year()
        const month = date.month() + 1

        if (year < minYear || (year === minYear && month < minMonth)) {
            minYear = year
            minMonth = month
        }
        if (year > maxYear || (year === maxYear && month > maxMonth)) {
            maxYear = year
            maxMonth = month
        }
    })

    // 生成完整的时间轴（年 or 月）
    const timeLabels: string[] = []
    if (granularity === 'year') {
        for (let y = minYear; y <= maxYear; y++) {
            timeLabels.push(String(y))
        }
    } else {
        let currentYear = minYear
        let currentMonth = minMonth
        while (currentYear < maxYear || (currentYear === maxYear && currentMonth <= maxMonth)) {
            timeLabels.push(`${currentYear}-${String(currentMonth).padStart(2, '0')}`)
            currentMonth++
            if (currentMonth > 12) {
                currentMonth = 1
                currentYear++
            }
        }
    }

    // 初始化每月/每年数据结构
    const initData = () =>
        lineDisplayTypes.reduce(
            (acc, type) => {
                acc[type] = 0
                return acc
            },
            {} as Record<string, number>
        )

    const timeDataMap = new Map<string, Record<string, number>>()
    timeLabels.forEach((label) => {
        timeDataMap.set(label, initData())
    })

    // 聚合数据
    filteredData.forEach((item) => {
        const date = dayjs(item.transactionAt * 1000)
        const label =
            granularity === 'year'
                ? String(date.year())
                : `${date.year()}-${String(date.month() + 1).padStart(2, '0')}`

        const type = item.transactionType
        const amount = item.price

        if (timeDataMap.has(label) && lineDisplayTypes.includes(type)) {
            const record = timeDataMap.get(label)!
            record[type]! += amount
        }
    })

    return {
        tooltip: {trigger: 'axis'},
        legend: {
            data: lineDisplayTypes.map((type) => TransactionTypeToLabel.get(type)),
        } as LegendComponentOption,
        xAxis: {
            type: 'category',
            data: timeLabels,
            name: granularity === 'year' ? '年份' : '月份',
        },
        yAxis: {
            type: 'value',
            name: '金额（元）',
            axisLabel: {
                formatter: '{value}',
            },
        },
        series: lineDisplayTypes.map((type) => {
            const data = timeLabels.map((label) => {
                const value = timeDataMap.get(label)?.[type] ?? 0
                return (value / 100).toFixed(2)
            })
            return {
                name: TransactionTypeToLabel.get(type),
                type: 'line',
                data,
                emphasis: {focus: 'series'},
                color: TransactionTypeToColor.get(type),
            }
        }),
    }
}

/**
 * 构造饼图配置项（按 category 聚合）
 */
export function buildPieChart(
    trList: TransactionRecord[],
    pieOptions: {
        transactionType: string
    }
): EChartsOption {
    const {transactionType} = pieOptions;
    const filteredData = trList.filter(
        (item) => item.transactionType === transactionType
    )

    const categoryMap = new Map<string, number>()
    filteredData.forEach((item) => {
        const {category, price} = item
        categoryMap.set(category, (categoryMap.get(category) || 0) + price)
    })

    const seriesData = Array.from(categoryMap, ([name, value]) => ({
        name,
        value: value / 100,
    }))

    if (seriesData.length === 0) {
        return {
            title: {
                text: '暂无数据',
                left: 'center',
                top: 'center',
                textStyle: {color: '#999'},
            },
            tooltip: {show: false},
            series: [],
        }
    }

    const series: PieSeriesOption[] = [
        {
            name: TransactionTypeToLabel.get(transactionType) || transactionType,
            type: 'pie',
            radius: ['40%', '70%'],
            center: ['40%', '50%'],
            avoidLabelOverlap: false,
            label: {
                show: true,
                formatter: '{b}\n{c}\n({d}%)',
                fontSize: 12,
            },
            emphasis: {
                label: {
                    show: true,
                    fontSize: 14,
                    fontWeight: 'bold',
                },
            },
            labelLine: {
                show: true,
                length: 20,
                length2: 50,
            },
            data: seriesData,
        },
    ]

    return {
        title: {text: ''},
        tooltip: {
            trigger: 'item',
            formatter: '{b}: {c} ({d}%)',
        },
        legend: {
            type: 'scroll',
            orient: 'vertical',
            right: 10,
            top: 20,
            bottom: 20,
        },
        series,
    }
}