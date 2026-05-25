<template>
  <div>
    <t-card :title="$t('monitor.title')">
      <template #actions>
        <t-button theme="primary" @click="openCreate">{{ $t('monitor.create') }}</t-button>
      </template>
      <t-table :data="monitors" :columns="columns" :loading="loading" row-key="id">
        <template #status="{ row }">
          <t-tag :theme="row.enabled ? 'success' : 'default'">
            {{ row.enabled ? $t('monitor.enabled') : $t('monitor.disabled') }}
          </t-tag>
        </template>
        <template #operation="{ row }">
          <t-button variant="text" @click="toggleEnabled(row)">
            {{ row.enabled ? $t('monitor.disabled') : $t('monitor.enabled') }}
          </t-button>
          <t-button variant="text" @click="$router.push(`/monitors/${row.id}`)">{{ $t('common.view') }}</t-button>
          <t-button variant="text" @click="openEdit(row)">{{ $t('common.edit') }}</t-button>
          <t-popconfirm :content="$t('common.confirm')" @confirm="onDelete(row.id)">
            <t-button variant="text" theme="danger">{{ $t('common.delete') }}</t-button>
          </t-popconfirm>
        </template>
      </t-table>
    </t-card>

    <t-dialog v-model:visible="showForm" :header="editing ? $t('monitor.edit') : $t('monitor.create')" @confirm="onSubmit" width="500px">
      <t-form :data="formData" label-width="120px">
        <t-form-item :label="$t('monitor.name')">
          <t-input v-model="formData.name" />
        </t-form-item>
        <t-form-item :label="$t('monitor.type')">
          <t-select v-model="formData.type">
            <t-option value="http" label="HTTP" />
            <t-option value="https" label="HTTPS" />
            <t-option value="tcp" label="TCP" />
            <t-option value="ping" label="Ping" />
          </t-select>
        </t-form-item>
        <t-form-item :label="$t('monitor.target')">
          <t-input v-model="formData.target" placeholder="https://example.com" />
        </t-form-item>
        <t-form-item :label="$t('monitor.interval')">
          <t-input-number v-model="formData.interval" :min="1" />
        </t-form-item>
        <t-form-item :label="$t('monitor.timeout')">
          <t-input-number v-model="formData.timeout" :min="1" />
        </t-form-item>
        <t-form-item v-if="!editing" :label="$t('nav.projects')">
          <t-select v-model="formData.project_id">
            <t-option v-for="p in projectList" :key="p.id" :value="p.id" :label="p.name" />
          </t-select>
        </t-form-item>
      </t-form>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getMonitors, createMonitor, updateMonitor, toggleMonitor, deleteMonitor, type Monitor } from '../api/monitor'
import { getProjects, type Project } from '../api/project'
import { useI18n } from 'vue-i18n'
import { MessagePlugin } from 'tdesign-vue-next'

const { t } = useI18n()

const monitors = ref<Monitor[]>([])
const projectList = ref<Project[]>([])
const loading = ref(false)
const showForm = ref(false)
const editing = ref<Monitor | null>(null)
const formData = reactive({
  name: '',
  type: 'http',
  target: '',
  interval: 60,
  timeout: 10,
  project_id: 0,
})

const columns = [
  { colKey: 'id', title: 'ID', width: 60 },
  { colKey: 'name', title: t('monitor.name') },
  { colKey: 'type', title: t('monitor.type') },
  { colKey: 'target', title: t('monitor.target') },
  { colKey: 'interval', title: t('monitor.interval') },
  { colKey: 'status', title: t('monitor.status') },
  { colKey: 'operation', title: t('common.operation'), width: 320 },
]

async function loadData() {
  loading.value = true
  try {
    monitors.value = await getMonitors()
    projectList.value = await getProjects()
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editing.value = null
  formData.name = ''
  formData.type = 'http'
  formData.target = ''
  formData.interval = 60
  formData.timeout = 10
  formData.project_id = projectList.value[0]?.id ?? 0
  showForm.value = true
}

function openEdit(row: Monitor) {
  editing.value = row
  formData.name = row.name
  formData.type = row.type
  formData.target = row.target
  formData.interval = row.interval
  formData.timeout = row.timeout
  formData.project_id = row.project_id
  showForm.value = true
}

async function toggleEnabled(row: Monitor) {
  await toggleMonitor(row.id, !row.enabled)
  loadData()
}

async function onDelete(id: number) {
  await deleteMonitor(id)
  MessagePlugin.success(t('common.success'))
  loadData()
}

async function onSubmit() {
  if (editing.value) {
    await updateMonitor(editing.value.id, formData)
  } else {
    await createMonitor(formData)
  }
  showForm.value = false
  MessagePlugin.success(t('common.success'))
  loadData()
}

onMounted(loadData)
</script>
