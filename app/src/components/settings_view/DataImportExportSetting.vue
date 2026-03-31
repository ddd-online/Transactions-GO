<template>
  <div class="data-import-export">
    <div class="section">
      <div class="section-title">数据导出</div>
      <div class="section-desc">将一段时间内的消费记录导出为 JSON 文件</div>

      <div class="export-form">
        <div class="form-row">
          <div class="form-item-half">
            <div class="form-label">开始时间</div>
            <a-date-picker v-model:value="exportStartTime" style="width: 100%" />
          </div>
          <div class="form-item-half">
            <div class="form-label">结束时间</div>
            <a-date-picker v-model:value="exportEndTime" style="width: 100%" />
          </div>
        </div>
        <billadm-button type="primary" @click="handleExport">
          <template #icon>
            <DownloadOutlined />
          </template>
          导出消费记录
        </billadm-button>
      </div>
    </div>

    <a-divider />

    <div class="section">
      <div class="section-title">数据导入</div>
      <div class="section-desc">从 JSON 文件导入消费记录到当前账本</div>

      <div class="import-form">
        <billadm-button type="primary" @click="handleImportSelect">
          <template #icon>
            <UploadOutlined />
          </template>
          选择文件导入
        </billadm-button>
        <input
          type="file"
          ref="fileInputRef"
          accept=".json"
          style="display: none"
          @change="handleFileChange"
        />
      </div>

      <!-- 导入预览 -->
      <div v-if="importPreview.length > 0" class="import-preview">
        <div class="preview-header">
          <span class="preview-title">导入预览 (共 {{ importPreview.length }} 条记录)</span>
          <billadm-button type="text" size="small" @click="clearImportPreview">清除</billadm-button>
        </div>
        <a-table
          :columns="previewColumns"
          :data-source="importPreview"
          :pagination="{ pageSize: 5 }"
          size="small"
          :scroll="{ y: 200 }"
        />
        <div class="preview-footer">
          <billadm-button @click="clearImportPreview">取消</billadm-button>
          <billadm-button type="primary" @click="confirmImport">确认导入</billadm-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { DownloadOutlined, UploadOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import type { TransactionRecord } from '@/types/billadm'
import { useLedgerStore } from '@/stores/ledgerStore'
import { queryTrOnCondition, batchCreateTrForLedger } from '@/backend/api/tr'
import dayjs from 'dayjs'

const ledgerStore = useLedgerStore()

// 导出相关
const exportStartTime = ref<dayjs.Dayjs>(dayjs().startOf('month'))
const exportEndTime = ref<dayjs.Dayjs>(dayjs().endOf('month'))

// 导入相关
const fileInputRef = ref<HTMLInputElement | null>(null)
const importPreview = ref<ExportRecord[]>([])

interface ExportRecord {
  transactionAt: number
  transactionType: string
  category: string
  tags: string[]
  flags: string[]
  description: string
  price: number
}

const previewColumns = [
  { title: '时间', dataIndex: 'transactionAt', width: 120 },
  { title: '类型', dataIndex: 'transactionType', width: 80 },
  { title: '分类', dataIndex: 'category', width: 100 },
  { title: '标签', dataIndex: 'tags', width: 150 },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '价格', dataIndex: 'price', width: 100 },
]

const handleExport = async () => {
  if (!ledgerStore.currentLedgerId) {
    message.warning('请先选择一个账本')
    return
  }

  if (!exportStartTime.value || !exportEndTime.value) {
    message.warning('请选择导出时间范围')
    return
  }

  try {
    const startTs = exportStartTime.value.startOf('day').unix()
    const endTs = exportEndTime.value.endOf('day').unix()

    const result = await queryTrOnCondition({
      ledgerId: ledgerStore.currentLedgerId,
      tsRange: [startTs, endTs],
      offset: 0,
      limit: 10000,
    })

    if (result.items.length === 0) {
      message.info('该时间范围内没有消费记录')
      return
    }

    // 转换为导出格式
    const exportData: ExportRecord[] = result.items.map((tr: TransactionRecord) => ({
      transactionAt: tr.transactionAt,
      transactionType: tr.transactionType,
      category: tr.category,
      tags: tr.tags || [],
      flags: tr.outlier ? ['outlier'] : [],
      description: tr.description,
      price: tr.price,
    }))

    const jsonContent = JSON.stringify(exportData, null, 2)

    // 调用保存对话框
    const saveResult = await window.electronAPI.saveDialog({
      title: '保存导出文件',
      defaultPath: `transactions_export_${dayjs().format('YYYYMMDD_HHmmss')}.json`,
      filters: [{ name: 'JSON Files', extensions: ['json'] }],
    })

    if (saveResult.canceled || !saveResult.filePath) {
      return
    }

    // 写入文件
    const writeResult = await window.electronAPI.writeFile(saveResult.filePath, jsonContent)
    if (writeResult.success) {
      message.success(`已导出 ${exportData.length} 条消费记录到 ${saveResult.filePath}`)
    } else {
      message.error(`写入文件失败: ${writeResult.error}`)
    }
  } catch (error) {
    message.error(`导出失败: ${error}`)
  }
}

const handleImportSelect = () => {
  fileInputRef.value?.click()
}

const handleFileChange = async (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return

  try {
    const text = await file.text()
    const data = JSON.parse(text)

    if (!Array.isArray(data)) {
      message.error('文件格式错误：期望一个数组')
      return
    }

    // 验证数据格式
    const validData: ExportRecord[] = []
    for (let i = 0; i < data.length; i++) {
      const item = data[i]
      if (typeof item.transactionAt !== 'number') {
        message.error(`第 ${i + 1} 条记录缺少有效的 transactionAt 字段`)
        continue
      }
      if (typeof item.transactionType !== 'string') {
        message.error(`第 ${i + 1} 条记录缺少有效的 transactionType 字段`)
        continue
      }
      validData.push({
        transactionAt: item.transactionAt,
        transactionType: item.transactionType,
        category: item.category || '',
        tags: Array.isArray(item.tags) ? item.tags : [],
        flags: Array.isArray(item.flags) ? item.flags : [],
        description: item.description || '',
        price: typeof item.price === 'number' ? item.price : 0,
      })
    }

    if (validData.length === 0) {
      message.error('没有有效的消费记录')
      return
    }

    importPreview.value = validData
    message.info(`已加载 ${validData.length} 条记录，请确认导入`)
  } catch (error) {
    message.error(`读取文件失败: ${error}`)
  }

  // 清空 input 值以便重复选择同一文件
  target.value = ''
}

const clearImportPreview = () => {
  importPreview.value = []
}

const confirmImport = async () => {
  if (!ledgerStore.currentLedgerId) {
    message.warning('请先选择一个账本')
    return
  }

  if (importPreview.value.length === 0) {
    message.warning('没有可导入的记录')
    return
  }

  try {
    const trList: TransactionRecord[] = importPreview.value.map(record => ({
      ledgerId: ledgerStore.currentLedgerId,
      transactionId: '', // 空表示新建
      transactionAt: record.transactionAt,
      transactionType: record.transactionType,
      category: record.category,
      tags: record.tags,
      description: record.description,
      price: record.price,
      outlier: record.flags.includes('outlier'),
    }))

    const count = await batchCreateTrForLedger(trList)
    message.success(`成功导入 ${count} 条消费记录`)
    clearImportPreview()
  } catch (error) {
    message.error(`导入失败: ${error}`)
  }
}
</script>

<style scoped>
.data-import-export {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.section {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.section-title {
  font-size: 16px;
  font-weight: 500;
  color: var(--billadm-color-text-major);
}

.section-desc {
  font-size: 13px;
  color: var(--billadm-color-text-minor);
}

.export-form,
.import-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 8px;
}

.form-row {
  display: flex;
  gap: 16px;
}

.form-item-half {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-label {
  font-size: 13px;
  color: var(--billadm-color-text-minor);
}

.import-preview {
  margin-top: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.preview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.preview-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--billadm-color-text-major);
}

.preview-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}
</style>
