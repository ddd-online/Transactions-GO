import {defineStore} from 'pinia'
import {computed, ref} from 'vue'
import {createLedger, deleteLedgerById, modifyLedger, queryAllLedgers} from "@/backend/api/ledger.ts"
import NotificationUtil from "@/backend/notification"
import type {Ledger} from "@/types/billadm"


export const useLedgerStore = defineStore('ledger', () => {
    const ledgers = ref([] as Ledger[])
    const currentLedger = ref({} as Ledger | null)

    const currentLedgerId = computed(() => {
        return currentLedger.value ? currentLedger.value.id : ''
    })

    const currentLedgerName = computed(() => {
        return currentLedger.value ? currentLedger.value.name : ''
    })

    const init = async () => {
        await refreshLedgers();
    }

    // 访问后端更新账本
    const refreshLedgers = async () => {
        try {
            ledgers.value = []
            const ledgersFromServer = await queryAllLedgers()
            if (ledgersFromServer && Array.isArray(ledgersFromServer)) {
                ledgersFromServer.sort((a, b) => a.createdAt - b.createdAt)
                ledgersFromServer.forEach(ledger => {
                    ledgers.value.push(ledger)
                })
            }
            if (currentLedger.value !== null) {
                setCurrentLedger(currentLedger.value.id)
            }
            if (currentLedger.value === null && ledgers.value.length > 0) {
                const firstLedger = ledgers.value[0]
                if (firstLedger) {
                    setCurrentLedger(firstLedger.id)
                }
            }
        } catch (error) {
            NotificationUtil.error('查询账本失败', `${error}`)
        }
    }

    // 添加账本
    const createLedgerAction = async (name: string, description: string = '') => {
        try {
            await createLedger(name, description)
            await refreshLedgers()
            NotificationUtil.success(`创建账本 ${name} 成功`)
        } catch (error) {
            NotificationUtil.error(`创建账本 ${name} 失败`, `${error}`)
        }
    }

    // 删除账本
    const deleteLedger = async (id: string) => {
        try {
            await deleteLedgerById(id);
            await refreshLedgers();
            NotificationUtil.success(`删除账本成功`);
        } catch (error) {
            NotificationUtil.error(`删除账本失败`);
        }
    }

    // 修改账本
    const modifyLedgerAction = async (id: string, name: string, description: string = '') => {
        try {
            await modifyLedger(id, name, description);
            await refreshLedgers();
            NotificationUtil.success(`修改账本成功`);
        } catch (error) {
            NotificationUtil.error(`修改账本失败`, `${error}`);
        }
    }

    // 设置当前账本
    const setCurrentLedger = (id: string) => {
        if (id === null) {
            currentLedger.value = null
            return
        }
        const ledger: Ledger | undefined = ledgers.value.find(l => l.id === id)
        if (ledger) {
            currentLedger.value = JSON.parse(JSON.stringify(ledger)) // 创建副本，避免直接引用
        } else {
            currentLedger.value = null
        }
    }

    return {
        ledgers,
        currentLedger,
        currentLedgerId,
        currentLedgerName,
        init,
        refreshLedgers,
        createLedger: createLedgerAction,
        deleteLedger,
        modifyLedger: modifyLedgerAction,
        setCurrentLedger,
    }
})
