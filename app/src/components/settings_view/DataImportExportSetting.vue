<template>
  <div class="data-import-export">
    <a-row :gutter="16">
      <!-- 导出区域 -->
      <a-col :span="12">
        <a-card title="数据导出" :bordered="false" class="section-card">
          <p class="section-desc">将一段时间内的消费记录导出为 JSON 文件</p>

          <a-form layout="vertical">
            <a-row :gutter="8">
              <a-col :span="12">
                <a-form-item label="开始时间">
                  <a-date-picker v-model:value="exportStartTime" style="width: 100%" />
                </a-form-item>
              </a-col>
              <a-col :span="12">
                <a-form-item label="结束时间">
                  <a-date-picker v-model:value="exportEndTime" style="width: 100%" />
                </a-form-item>
              </a-col>
            </a-row>
          </a-form>

          <div class="card-footer">
            <a-button type="primary" :loading="exportLoading" @click="handleExport">
              <template #icon>
                <DownloadOutlined />
              </template>
              导出消费记录
            </a-button>
          </div>
        </a-card>
      </a-col>

      <!-- 导入区域 -->
      <a-col :span="12">
        <a-card title="数据导入" :bordered="false" class="section-card">
          <p class="section-desc">从 JSON 文件导入消费记录到当前账本</p>

          <input
            type="file"
            ref="fileInputRef"
            accept=".json"
            style="display: none"
            @change="handleFileChange"
          />

          <div class="card-footer">
            <a-button type="primary" @click="handleImportSelect">
              <template #icon>
                <UploadOutlined />
              </template>
              选择文件导入
            </a-button>
          </div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 导入预览弹窗 -->
    <a-modal
      v-model:open="importPreviewVisible"
      title="导入预览"
      width="900px"
      :footer="null"
      centered
    >
      <div class="import-preview">
        <a-alert
          :message="`共 ${importPreview.length} 条记录待导入`"
          type="info"
          show-icon
          style="margin-bottom: 16px"
        />

        <a-table
          :columns="previewColumns"
          :data-source="importPreview"
          :pagination="{ pageSize: 5 }"
          size="small"
          :scroll="{ y: 300 }"
          row-key="transactionAt"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.dataIndex === 'transactionAt'">
              {{ formatTimestamp(record.transactionAt) }}
            </template>
            <template v-else-if="column.dataIndex === 'transactionType'">
              <a-tag :color="getTypeColor(record.transactionType)">
                {{ getTypeLabel(record.transactionType) }}
              </a-tag>
            </template>
            <template v-else-if="column.dataIndex === 'tags'">
              <a-tag v-for="tag in record.tags" :key="tag" color="green">{{ tag }}</a-tag>
            </template>
            <template v-else-if="column.dataIndex === 'flags'">
              <a-tag v-if="record.flags.includes('outlier')" color="orange">离群值</a-tag>
            </template>
            <template v-else-if="column.dataIndex === 'price'">
              {{ formatPrice(record.price) }}
            </template>
          </template>
        </a-table>

        <div class="preview-footer">
          <a-space>
            <a-button @click="importPreviewVisible = false">取消</a-button>
            <a-button type="primary" :loading="importLoading" @click="confirmImport">
              确认导入
            </a-button>
          </a-space>
        </div>
      </div>
    </a-modal>
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
const exportLoading = ref(false)

// 导入相关
const fileInputRef = ref<HTMLInputElement | null>(null)
const importPreview = ref<ExportRecord[]>([])
const importPreviewVisible = ref(false)
const importLoading = ref(false)

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
  { title: '标记', dataIndex: 'flags', width: 80 },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '价格', dataIndex: 'price', width: 100 },
]

const formatTimestamp = (ts: number) => {
  return dayjs(ts * 1000).format('YYYY-MM-DD')
}

const formatPrice = (cents: number) => {
  return (cents / 100).toFixed(2)
}

const getTypeColor = (type: string) => {
  const colorMap: Record<string, string> = {
    income: 'green',
    expense: 'red',
    transfer: 'orange',
  }
  return colorMap[type] || 'blue'
}

const getTypeLabel = (type: string) => {
  const labelMap: Record<string, string> = {
    income: '收入',
    expense: '支出',
    transfer: '转账',
  }
  return labelMap[type] || type
}

const handleExport = async () => {
  if (!ledgerStore.currentLedgerId) {
    message.warning('请先选择一个账本')
    return
  }

  if (!exportStartTime.value || !exportEndTime.value) {
    message.warning('请选择导出时间范围')
    return
  }

  exportLoading.value = true

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
      message.success(`已导出 ${exportData.length} 条消费记录`)
    } else {
      message.error(`写入文件失败: ${writeResult.error}`)
    }
  } catch (error) {
    message.error(`导出失败: ${error}`)
  } finally {
    exportLoading.value = false
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
    importPreviewVisible.value = true
    message.info(`已加载 ${validData.length} 条记录`)
  } catch (error) {
    message.error(`读取文件失败: ${error}`)
  }

  // 清空 input 值以便重复选择同一文件
  target.value = ''
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

  importLoading.value = true

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
    importPreviewVisible.value = false
    importPreview.value = []
  } catch (error) {
    message.error(`导入失败: ${error}`)
  } finally {
    importLoading.value = false
  }
}
</script>

<style scoped>
.data-import-export {
  width: 100%;
}

.section-card {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.section-card :deep(.ant-card-body) {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.section-desc {
  color: var(--billadm-color-text-minor);
  margin-bottom: 16px;
}

.card-footer {
  display: flex;
  justify-content: flex-end;
  margin-top: auto;
  padding-top: 16px;
}

.import-preview {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.preview-footer {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid var(--billadm-color-window-border);
}
</style>
