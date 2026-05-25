<template>
  <div>
    <t-card :title="$t('user.title')">
      <template #actions>
        <t-button theme="primary" @click="showCreate = true">{{ $t('user.create') }}</t-button>
      </template>
      <t-table :data="users" :columns="columns" :loading="loading" row-key="id">
        <template #role="{ row }">
          <t-tag :theme="row.role === 'admin' ? 'warning' : 'default'">
            {{ row.role === 'admin' ? $t('user.admin') : $t('user.user') }}
          </t-tag>
        </template>
        <template #operation="{ row }">
          <t-button variant="text" @click="changePassword(row)">{{ $t('user.changePassword') }}</t-button>
          <t-popconfirm :content="$t('common.confirm')" @confirm="onDelete(row.id)">
            <t-button variant="text" theme="danger">{{ $t('common.delete') }}</t-button>
          </t-popconfirm>
        </template>
      </t-table>
    </t-card>

    <t-dialog v-model:visible="showCreate" :header="$t('user.create')" @confirm="onCreateSubmit" width="400px">
      <t-form :data="createData" label-width="100px">
        <t-form-item :label="$t('user.username')">
          <t-input v-model="createData.username" />
        </t-form-item>
        <t-form-item :label="$t('user.password')">
          <t-input v-model="createData.password" type="password" />
        </t-form-item>
        <t-form-item :label="$t('user.role')">
          <t-select v-model="createData.role">
            <t-option value="admin" :label="$t('user.admin')" />
            <t-option value="user" :label="$t('user.user')" />
          </t-select>
        </t-form-item>
      </t-form>
    </t-dialog>

    <t-dialog v-model:visible="showPassword" :header="$t('user.changePassword')" @confirm="onPasswordSubmit" width="400px">
      <t-form label-width="100px">
        <t-form-item :label="$t('user.password')">
          <t-input v-model="newPassword" type="password" />
        </t-form-item>
      </t-form>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getUsers, createUser, updateUserPassword, deleteUser, type User } from '../api/user'
import { useI18n } from 'vue-i18n'
import { MessagePlugin } from 'tdesign-vue-next'

const { t } = useI18n()

const users = ref<User[]>([])
const loading = ref(false)
const showCreate = ref(false)
const showPassword = ref(false)
const passwordUserId = ref(0)
const newPassword = ref('')
const createData = reactive({ username: '', password: '', role: 'user' })

const columns = [
  { colKey: 'id', title: 'ID', width: 60 },
  { colKey: 'username', title: t('user.username') },
  { colKey: 'role', title: t('user.role') },
  { colKey: 'created_at', title: 'Created' },
  { colKey: 'operation', title: t('common.operation'), width: 200 },
]

async function loadData() {
  loading.value = true
  try {
    users.value = await getUsers()
  } finally {
    loading.value = false
  }
}

function changePassword(row: User) {
  passwordUserId.value = row.id
  newPassword.value = ''
  showPassword.value = true
}

async function onCreateSubmit() {
  await createUser(createData)
  showCreate.value = false
  createData.username = ''
  createData.password = ''
  createData.role = 'user'
  MessagePlugin.success(t('common.success'))
  loadData()
}

async function onPasswordSubmit() {
  await updateUserPassword(passwordUserId.value, newPassword.value)
  showPassword.value = false
  MessagePlugin.success(t('common.success'))
}

async function onDelete(id: number) {
  await deleteUser(id)
  MessagePlugin.success(t('common.success'))
  loadData()
}

onMounted(loadData)
</script>
