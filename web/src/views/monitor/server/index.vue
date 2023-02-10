<template>
  <div class="app-container">
    <el-row>
      <el-col :span="12" class="card-box">
        <el-card>
          <template #header><span>CPU</span></template>
          <div class="prop-wrapper">
            <div class="prop-item-container">
              <div>核心数</div>
              <div v-if="server.cpu">{{ server.cpu.cpuNum }}</div>
            </div>
            <div class="horizontal-divider"></div>
            <div class="prop-item-container">
              <div>用户使用率</div>
              <div v-if="server.cpu">{{ server.cpu.used }}%</div>
            </div>
            <div class="horizontal-divider"></div>
            <div class="prop-item-container">
              <div>系统使用率</div>
              <div v-if="server.cpu">{{ server.cpu.sys }}%</div>
            </div>
            <div class="horizontal-divider"></div>
            <div class="prop-item-container">
              <div>当前空闲率</div>
              <div v-if="server.cpu">{{ server.cpu.free }}%</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="12" class="card-box">
        <el-card>
          <template #header><span>内存</span></template>
          <div class="prop-wrapper">
            <div class="prop-item-container">
              <div>总内存</div>
              <div v-if="server.mem">{{ server.mem.total }}G</div>
            </div>
            <div class="horizontal-divider"></div>
            <div class="prop-item-container">
              <div>已用内存</div>
              <div v-if="server.mem">{{ server.mem.used }}G</div>
            </div>
            <div class="horizontal-divider"></div>
            <div class="prop-item-container">
              <div>剩余内存</div>
              <div v-if="server.mem">{{ server.mem.free }}G</div>
            </div>
            <div class="horizontal-divider"></div>
            <div class="prop-item-container">
              <div>使用率</div>
              <div v-if="server.mem">{{ server.mem.usedPercent }}%</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="24" class="card-box">
        <el-card>
          <template #header><span>服务器信息</span></template>
          <div class="prop-wrapper">
            <div class="prop-item-container">
              <div class="prop-item-container">
                <div>服务器名称&nbsp;</div>
                <div v-if="server.sys">{{ server.sys.computerName }}</div>
              </div>
              <div class="prop-item-container">
                <div>操作系统&nbsp;</div>
                <div v-if="server.sys">{{ server.sys.osName }}</div>
              </div>
            </div>
            <div class="horizontal-divider"></div>
            <div class="prop-item-container">
              <div class="prop-item-container">
                <div>服务器IP&nbsp;</div>
                <div v-if="server.sys">{{ server.sys.computerIp }}</div>
              </div>
              <div class="prop-item-container">
                <div>系统架构&nbsp;</div>
                <div v-if="server.sys">{{ server.sys.osArch }}</div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="24" class="card-box">
        <el-card>
          <template #header><span>磁盘状态</span></template>
          <el-table :data="server.disk">
            <el-table-column label="盘符路径" key="dirName" prop="dirName" />
            <el-table-column label="文件系统" key="fsType" prop="fsType" />
            <!-- <el-table-column label="盘符类型" key="userId" prop="userId" /> -->
            <el-table-column label="总大小" key="total" prop="total" />
            <el-table-column label="可用大小" key="free" prop="free" />
            <el-table-column label="已用大小" key="used" prop="used" />
            <el-table-column label="已用百分比" key="usedPercent" prop="usedPercent" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script lang="ts" setup>
import { getServer } from '@/api/monitor/server'
import useCurrentInstance from '@/hooks/useCurrentInstance'
import { ServerFormData } from '@/types/api/server'

const data = reactive({
  server: {} as ServerFormData
})

const { server } = toRefs(data)

const { proxy } = useCurrentInstance()

function getList() {
  proxy.$modal.loading('正在加载服务监控数据，请稍候！')
  getServer().then(response => {
    server.value = response.data
    proxy.$modal.closeLoading()
  })
}

getList()
</script>

<style>
.prop-wrapper {
  display: flex;
  flex-direction: column;
}
.prop-item-container {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  font-size: 14px;
  color: #606266;
}
</style>
