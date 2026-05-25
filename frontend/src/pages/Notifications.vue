<template>
  <div>
    <t-card :title="$t('notification.title')">
      <template #actions>
        <t-button theme="primary" @click="openCreate">{{ $t('notification.create') }}</t-button>
      </template>
      <t-table :data="notifications" :columns="columns" :loading="loading" row-key="id">
        <template #enabled="{ row }">
          <t-tag :theme="row.enabled ? 'success' : 'default'">
            {{ row.enabled ? $t('notification.enabled') : 'Disabled' }}
          </t-tag>
        </template>
        <template #operation="{ row }">
          <t-button variant="text" @click="testNotif(row)">{{ $t('notification.test') }}</t-button>
          <t-button variant="text" @click="editNotif(row)">{{ $t('common.edit') }}</t-button>
          <t-popconfirm :content="$t('common.confirm')" @confirm="onDelete(row.id)">
            <t-button variant="text" theme="danger">{{ $t('common.delete') }}</t-button>
          </t-popconfirm>
        </template>
      </t-table>
    </t-card>

    <t-dialog v-model:visible="showForm" :header="editing ? $t('notification.edit') : $t('notification.create')" @confirm="onSubmit" width="500px">
      <t-form :data="formData" label-width="120px">
        <t-form-item :label="$t('notification.type')">
          <t-select v-model="formData.type">
            <t-option value="lark" :label="$t('notification.lark')" />
            <t-option value="dingtalk" :label="$t('notification.dingtalk')" />
            <t-option value="webhook" :label="$t('notification.webhook')" />
            <t-option value="email" :label="$t('notification.email')" />
          </t-select>
        </t-form-item>
        <t-form-item v-if="formData.type === 'lark' || formData.type === 'dingtalk'" :label="$t('notification.webhookUrl')">
          <t-input v-model="formData.config.webhook_url" />
        </t-form-item>
        <t-form-item v-if="formData.type === 'webhook'" :label="$t('notification.webhookUrl')">
          <t-input v-model="formData.config.url" />
        </t-form-item>
        <t-form-item v-if="formData.type === 'email'" :label="$t('notification.smtpHost')">
          <t-input v-model="formData.config.smtp_host" />
        </t-form-item>
        <t-form-item v-if="formData.type === 'email'" :label="$t('notification.smtpPort')">
          <t-input v-model="formData.config.smtp_port" />
        </t-form-item>
        <t-form-item v-if="formData.type === 'email'" :label="$t('notification.from')">
          <t-input v-model="formData.config.from" />
        </t-form-item>
        <t-form-item v-if="formData.type === 'email'" :label="$t('notification.password')">
          <t-input v-model="formData.config.password" type="password" />
        </t-form-item>
        <t-form-item v-if="formData.type === 'email'" :label="$t('notification.to')">
          <t-input v-model="formData.config.to" />
        </t-form-item>
      </t-form>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import {
  getNotifications,
  createNotification,
  updateNotification,
  deleteNotification,
  testNotification,
  type Notification,
} from '../api/notification'
import { useI18n } from 'vue-i18n'
import { MessagePlugin } from 'tdesign-vue-next'

const { t } = useI18n()

const notifications = ref<Notification[]>([])
const loading = ref(false)
const showForm = ref(false)
const editing = ref<Notification | null>(null)
const formData = reactive({
  type: 'lark',
  config: {} as Record<string, any>,
})

const columns = [
  { colKey: 'id', title: 'ID', width: 60 },
  { colKey: 'type', title: t('notification.type') },
  { colKey: 'enabled', title: t('notification.enabled') },
  { colKey: 'operation', title: t('common.operation'), width: 250 },
]

async function loadData() {
  loading.value = true
  try {
    notifications.value = await getNotifications()
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editing.value = null
  formData.type = 'lark'
  formData.config = {}
  showForm.value = true
}

function editNotif(row: Notification) {
  editing.value = row
  formData.type = row.type
  try {
    formData.config = JSON.parse(row.config)
  } catch {
    formData.config = {}
  }
  showForm.value = true
}

async function testNotif(row: Notification) {
  try {
    await testNotification(row.id)
    MessagePlugin.success(t('common.success'))
  } catch (e: any) {
    MessagePlugin.error(e.message)
  }
}

async function onDelete(id: number) {
  await deleteNotification(id)
  MessagePlugin.success(t('common.success'))
  loadData()
}

async function onSubmit() {
  if (editing.value) {
    await updateNotification(editing.value.id, { type: formData.type, config: formData.config })
  } else {
    await createNotification({ type: formData.type, config: formData.config })
  }
  showForm.value = false
  MessagePlugin.success(t('common.success'))
  loadData()
}

onMounted(loadData)
</script>
