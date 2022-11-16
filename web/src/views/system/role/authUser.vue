<template>
  <div class="app-container">
    <el-form :model="queryParams" ref="queryFormRef" v-show="showSearch" :inline="true">
      <el-form-item label="用户名称" prop="username">
        <el-input v-model="queryParams.username" placeholder="请输入用户名称" clearable style="width: 240px" @keyup.enter="handleQuery" />
      </el-form-item>
      <el-form-item label="手机号码" prop="phoneNumber">
        <el-input v-model="queryParams.phoneNumber" placeholder="请输入手机号码" clearable style="width: 240px" @keyup.enter="handleQuery" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleQuery"><ep:search /> 搜索</el-button>
        <el-button @click="resetQuery"><ep:refresh /> 重置</el-button>
      </el-form-item>
    </el-form>

    <el-row :gutter="10" class="mb8">
      <el-col :span="1.5">
        <el-button type="primary" plain @click="openSelectUser" v-hasPerms="['system:role:add']">
          <ep:plus />添加用户
        </el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button type="danger" plain :disabled="multiple" @click="cancelAuthUserAll" v-hasPerms="['system:role:remove']">
          <ep:circleClose /> 批量取消授权
        </el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button type="warning" plain @click="handleClose">
          <ep:close /> 关闭
        </el-button>
      </el-col>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList"></right-toolbar>
    </el-row>

    <el-table v-loading="loading" :data="userList" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" align="center" />
      <el-table-column label="用户名称" prop="username" :show-overflow-tooltip="true" />
      <el-table-column label="用户昵称" prop="nickname" :show-overflow-tooltip="true" />
      <el-table-column label="邮箱" prop="email" :show-overflow-tooltip="true" />
      <el-table-column label="手机" prop="phoneNumber" :show-overflow-tooltip="true" />
      <el-table-column label="状态" align="center" prop="status">
        <template #default="scope">
          <dict-tag :options="sys_normal_disable" :value="scope.row.status" />
        </template>
      </el-table-column>
      <el-table-column label="创建时间" align="center" prop="createTime" width="180">
        <template #default="scope">
          <span>{{ parseTime(scope.row.createTime) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" class-name="small-padding fixed-width">
        <template #default="scope">
          <el-button link @click="cancelAuthUser(scope.row)" v-hasPerms="['system:role:remove']">
            <ep:circleClose /> 取消授权
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total > 0" :total="total" v-model:page="queryParams.pageNum" v-model:limit="queryParams.pageSize" @pagination="getList" />
    <select-user ref="selectRef" :roleId="queryParams.roleId" @ok="handleQuery" />
  </div>
</template>

<script lang="ts" setup name="AuthUser">
import selectUser from './selectUser.vue'
import { allocatedUserList, authUserCancel, authUserCancelAll } from '@/api/system/role'
import useCurrentInstance from '@/hooks/useCurrentInstance'

const route = useRoute()
const { proxy } = useCurrentInstance()
const { sys_normal_disable } = proxy.useDict('sys_normal_disable')

const userList = ref([])
const loading = ref(true)
const showSearch = ref(true)
const multiple = ref(true)
const total = ref(0)
const userIds = ref([])

const selectRef = ref<any>(null)
const queryFormRef = ref<ElForm>(null)

const queryParams = reactive({
  pageNum: 1,
  pageSize: 10,
  roleId: route.params.roleId,
  username: undefined,
  phoneNumber: undefined
})

/** 查询授权用户列表 */
function getList() {
  loading.value = true
  allocatedUserList(queryParams).then(response => {
    userList.value = response.data.rows
    total.value = response.data.total
    loading.value = false
  })
}
// 返回按钮
function handleClose() {
  const obj = { path: '/system/role' }
  proxy.$tab.closeOpenPage(obj)
}
/** 搜索按钮操作 */
function handleQuery() {
  queryParams.pageNum = 1
  getList()
}
/** 重置按钮操作 */
function resetQuery() {
  queryFormRef.value?.resetFields()
  handleQuery()
}
// 多选框选中数据
function handleSelectionChange(selection:any) {
  userIds.value = selection.map((item:any) => item.userId)
  multiple.value = !selection.length
}
/** 打开授权用户表弹窗 */
function openSelectUser() {
  selectRef.value.show()
}
/** 取消授权按钮操作 */
function cancelAuthUser(row: { [key: string]: any }) {
  proxy.$modal
    .confirm('确认要取消该用户"' + row.username + '"角色吗？')
    .then(function () {
      return authUserCancel({ userId: Number(row.userId), roleId: Number(queryParams.roleId) })
    })
    .then(() => {
      getList()
      proxy.$modal.msgSuccess('取消授权成功')
    })
    .catch(() => {})
}
/** 批量取消授权按钮操作 */
function cancelAuthUserAll() {
  // const roleId = queryParams.roleId
  // const uIds = userIds.value.join(',')
  proxy.$modal
    .confirm('是否取消选中用户授权数据项?')
    .then(function () {
      return authUserCancelAll({ roleId: Number(queryParams.roleId), userIds: userIds.value })
    })
    .then(() => {
      getList()
      proxy.$modal.msgSuccess('取消授权成功')
    })
    .catch(() => {})
}

onMounted(() => {
  getList()
})
</script>
