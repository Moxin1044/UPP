<template>
  <div>
    <t-card :title="monitor?.name || ''">
      <template #actions>
        <t-button @click="$router.push('/monitors')">{{ $t('nav.monitors') }}</t-button>
      </template>
      <t-row :gutter="[16, 16]">
        <t-col :span="3">
          <div class="stat-card">
            <div class="stat-icon" :class="latestResult?.status === 'up' ? 'stat-icon--up' : 'stat-icon--down'">
              <t-icon :name="latestResult?.status === 'up' ? 'check-circle' : 'close-circle'" size="24px" />
            </div>
            <div class="stat-info">
              <div class="stat-title">{{ $t('monitor.status') }}</div>
              <div class="stat-value" :class="latestResult?.status === 'up' ? 'stat-value--up' : 'stat-value--down'">
                {{ latestResult?.status === 'up' ? $t('monitor.up') : $t('monitor.down') }}
              </div>
            </div>
          </div>
        </t-col>
        <t-col :span="3">
          <div class="stat-card">
            <div class="stat-icon stat-icon--latency"><t-icon name="time" size="24px" /></div>
            <div class="stat-info">
              <div class="stat-title">{{ $t('monitor.latency') }}</div>
              <div class="stat-value">{{ latestResult?.latency?.toFixed(1) ?? '0' }}<span class="stat-unit">ms</span></div>
            </div>
          </div>
        </t-col>
        <t-col :span="3">
          <div class="stat-card">
            <div class="stat-icon stat-icon--uptime"><t-icon name="chart" size="24px" /></div>
            <div class="stat-info">
              <div class="stat-title">{{ $t('monitor.uptime') }}</div>
              <div class="stat-value">{{ stats?.uptime?.toFixed(1) ?? '0' }}<span class="stat-unit">%</span></div>
            </div>
          </div>
        </t-col>
        <t-col :span="3">
          <div class="stat-card">
            <div class="stat-icon stat-icon--avg"><t-icon name="swap" size="24px" /></div>
            <div class="stat-info">
              <div class="stat-title">{{ $t('monitor.avgLatency') }}</div>
              <div class="stat-value">{{ stats?.avg_latency?.toFixed(1) ?? '0' }}<span class="stat-unit">ms</span></div>
            </div>
          </div>
        </t-col>
      </t-row>
    </t-card>

    <t-card style="margin-top: 16px" :title="$t('dashboard.latencyTrend')">
      <div ref="chartRef" style="height: 350px"></div>
    </t-card>

    <t-card style="margin-top: 16px" :title="$t('monitor.history')">
      <t-table :data="results" :columns="resultColumns" row-key="id" :loading="loading" />
    </t-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { getMonitor, getMonitorResults, getMonitorStats, type MonitorResult, type MonitorStats } from '../api/monitor'
import { useI18n } from 'vue-i18n'
import * as echarts from 'echarts'

const route = useRoute()
const { t } = useI18n()

const darkAxis = { axisLine: { lineStyle: { color: 'rgba(255,255,255,0.1)' } }, axisLabel: { color: 'rgba(255,255,255,0.4)' }, splitLine: { lineStyle: { color: 'rgba(255,255,255,0.04)' } } }

const monitor = ref<any>(null)
const results = ref<MonitorResult[]>([])
const stats = ref<MonitorStats | null>(null)
const latestResult = ref<MonitorResult | null>(null)
const loading = ref(false)
const chartRef = ref<HTMLElement>()
let chart: echarts.ECharts | null = null

const resultColumns = [
  { colKey: 'id', title: 'ID', width: 60 },
  { colKey: 'status', title: t('monitor.status') },
  { colKey: 'latency', title: t('monitor.latency') },
  { colKey: 'message', title: 'Message' },
  { colKey: 'created_at', title: t('monitor.createdAt'), cell: (_h: any, { row }: any) => formatTime(row.created_at) },
]

function formatTime(val: string) {
  if (!val) return ''
  const d = new Date(val)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}:${pad(d.getMonth() + 1)}:${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}

onMounted(async () => {
  const id = Number(route.params.id)
  loading.value = true
  try {
    monitor.value = await getMonitor(id)
    results.value = await getMonitorResults(id, 100)
    stats.value = await getMonitorStats(id)
    latestResult.value = results.value.length > 0 ? results.value[0] : null

    if (chartRef.value) {
      chart = echarts.init(chartRef.value)
      const reversed = [...results.value].reverse()
      chart.setOption({
        tooltip: { trigger: 'axis', backgroundColor: '#1a2540', borderColor: 'rgba(255,255,255,0.06)', textStyle: { color: '#fff' } },
        xAxis: { type: 'category', data: reversed.map((r) => r.created_at), ...darkAxis },
        yAxis: { type: 'value', name: 'ms', nameTextStyle: { color: 'rgba(255,255,255,0.4)' }, ...darkAxis },
        series: [
          {
            type: 'line',
            data: reversed.map((r) => r.latency),
            smooth: true,
            areaStyle: { color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(0, 212, 255, 0.3)' },
              { offset: 1, color: 'rgba(0, 212, 255, 0.02)' },
            ]) },
            itemStyle: { color: '#00d4ff' },
            lineStyle: { width: 2 },
            symbol: 'none',
          },
        ],
        grid: { left: 50, right: 20, top: 20, bottom: 30 },
      })
    }
  } finally {
    loading.value = false
  }

  window.addEventListener('resize', onResize)
})

function onResize() {
  chart?.resize()
}

onUnmounted(() => {
  chart?.dispose()
  window.removeEventListener('resize', onResize)
})
</script>

<style scoped>
.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px 24px;
  background: var(--upp-bg-card);
  border: 1px solid var(--upp-border);
  border-radius: 12px;
  backdrop-filter: blur(12px);
  transition: all 0.3s;
}
.stat-card:hover {
  border-color: var(--upp-border-hover);
  box-shadow: 0 0 20px rgba(0, 180, 255, 0.06);
}
.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.stat-icon--up { background: rgba(82, 196, 26, 0.1); color: #73d13d; }
.stat-icon--down { background: rgba(255, 77, 79, 0.1); color: #ff7875; }
.stat-icon--latency { background: rgba(0, 212, 255, 0.1); color: #00d4ff; }
.stat-icon--uptime { background: rgba(82, 196, 26, 0.1); color: #73d13d; }
.stat-icon--avg { background: rgba(123, 97, 255, 0.1); color: #b37feb; }

.stat-info { flex: 1; min-width: 0; }
.stat-title {
  font-size: 12px;
  color: var(--upp-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 1px;
  margin-bottom: 4px;
}
.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: var(--upp-text-primary);
  line-height: 1;
}
.stat-value--up { color: #73d13d; }
.stat-value--down { color: #ff7875; }
.stat-unit {
  font-size: 14px;
  font-weight: 400;
  color: var(--upp-text-tertiary);
  margin-left: 2px;
}
</style>
