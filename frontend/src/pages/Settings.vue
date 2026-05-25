<template>
  <div>
    <t-card :title="$t('setting.title')">
      <t-form :data="formData" label-width="120px" @submit="onSave">
        <t-form-item :label="$t('setting.platformName')">
          <t-input v-model="formData.platformName" />
        </t-form-item>
        <t-form-item :label="$t('setting.defaultLang')">
          <t-select v-model="formData.defaultLang">
            <t-option value="zh" label="中文" />
            <t-option value="en" label="English" />
          </t-select>
        </t-form-item>
        <t-form-item>
          <t-button theme="primary" type="submit" :loading="saving">{{ $t('setting.save') }}</t-button>
        </t-form-item>
      </t-form>
    </t-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getSettings, setSetting } from '../api/dashboard'
import { useI18n } from 'vue-i18n'
import { MessagePlugin } from 'tdesign-vue-next'

const { t } = useI18n()
const saving = ref(false)
const formData = reactive({
  platformName: '',
  defaultLang: 'zh',
})

onMounted(async () => {
  try {
    const settings = await getSettings()
    formData.platformName = settings.platformName || 'UUP'
    formData.defaultLang = settings.defaultLang || 'zh'
  } catch {
    // ignore
  }
})

async function onSave() {
  saving.value = true
  try {
    await setSetting('platformName', formData.platformName)
    await setSetting('defaultLang', formData.defaultLang)
    MessagePlugin.success(t('common.success'))
  } catch (e: any) {
    MessagePlugin.error(e.message)
  } finally {
    saving.value = false
  }
}
</script>
