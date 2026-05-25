<template>
  <div>
    <t-card :title="$t('project.title')">
      <template #actions>
        <t-button theme="primary" @click="showCreate = true">{{ $t('project.create') }}</t-button>
      </template>
      <t-table :data="projects" :columns="columns" :loading="loading" row-key="id">
        <template #operation="{ row }">
          <t-button variant="text" @click="editProject(row)">{{ $t('common.edit') }}</t-button>
          <t-button variant="text" @click="viewProject(row)">{{ $t('nav.monitors') }}</t-button>
          <t-popconfirm :content="$t('common.confirm')" @confirm="onDelete(row.id)">
            <t-button variant="text" theme="danger">{{ $t('common.delete') }}</t-button>
          </t-popconfirm>
        </template>
      </t-table>
    </t-card>

    <t-dialog v-model:visible="showCreate" :header="editing ? $t('project.edit') : $t('project.create')" @confirm="onSubmit">
      <t-form :data="formData" label-width="100px">
        <t-form-item :label="$t('project.name')">
          <t-input v-model="formData.name" />
        </t-form-item>
        <t-form-item :label="$t('project.description')">
          <t-textarea v-model="formData.description" />
        </t-form-item>
      </t-form>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getProjects, createProject, updateProject, deleteProject, type Project } from '../api/project'
import { useI18n } from 'vue-i18n'
import { MessagePlugin } from 'tdesign-vue-next'

const router = useRouter()
const { t } = useI18n()

const projects = ref<Project[]>([])
const loading = ref(false)
const showCreate = ref(false)
const editing = ref<Project | null>(null)
const formData = reactive({ name: '', description: '' })

const columns = [
  { colKey: 'id', title: 'ID', width: 60 },
  { colKey: 'name', title: t('project.name') },
  { colKey: 'description', title: t('project.description') },
  { colKey: 'created_at', title: 'Created' },
  { colKey: 'operation', title: t('common.operation'), width: 200 },
]

async function loadData() {
  loading.value = true
  try {
    projects.value = await getProjects()
  } finally {
    loading.value = false
  }
}

function editProject(row: Project) {
  editing.value = row
  formData.name = row.name
  formData.description = row.description
  showCreate.value = true
}

function viewProject(row: Project) {
  router.push(`/projects/${row.id}`)
}

async function onDelete(id: number) {
  await deleteProject(id)
  MessagePlugin.success(t('common.success'))
  loadData()
}

async function onSubmit() {
  if (editing.value) {
    await updateProject(editing.value.id, formData)
  } else {
    await createProject(formData)
  }
  showCreate.value = false
  editing.value = null
  formData.name = ''
  formData.description = ''
  MessagePlugin.success(t('common.success'))
  loadData()
}

onMounted(loadData)
</script>
