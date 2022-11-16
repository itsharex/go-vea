<template>
  <div class="app-container">
    <el-row :gutter="20">
      <!--部门数据-->
      <el-col :span="4" :xs="24">
        <div class="head-container">
          <el-input v-model="deptName" placeholder="请输入部门名称" clearable prefix-icon="Search" style="margin-bottom: 20px" />
        </div>
        <div class="head-container">
          <el-tree
            :data="deptOptions"
            :props="{ label: 'deptName', children: 'children' }"
            :expand-on-click-node="false"
            :filter-node-method="filterNode"
            ref="deptTreeRef"
            node-key="id"
            highlight-current
            default-expand-all
            @node-click="handleNodeClick"
          />
        </div>
      </el-col>
      <!--用户数据-->
      <el-col :span="20" :xs="24">
        <el-form :model="queryParams" ref="queryFormRef" :inline="true" v-show="showSearch" label-width="68px">
          <el-form-item label="用户名称" prop="username">
            <el-input v-model="queryParams.username" placeholder="请输入用户名称" clearable style="width: 240px" @keyup.enter="handleQuery" />
          </el-form-item>
          <el-form-item label="手机号码" prop="phoneNumber">
            <el-input v-model="queryParams.phoneNumber" placeholder="请输入手机号码" clearable style="width: 240px" @keyup.enter="handleQuery" />
          </el-form-item>
          <el-form-item label="状态" prop="status">
            <el-select v-model="queryParams.status" placeholder="用户状态" clearable style="width: 240px">
              <el-option v-for="dict in sys_normal_disable" :key="dict.value" :label="dict.label" :value="dict.value" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleQuery"> <ep:search /> 搜索 </el-button>
            <el-button @click="resetQuery"> <ep:refresh /> 重置 </el-button>
          </el-form-item>
        </el-form>

        <el-row :gutter="10" class="mb8">
          <el-col :span="1.5">
            <el-button type="primary" plain @click="handleAdd" v-hasPerms="['system:user:add']"><ep:plus /> 新增</el-button>
          </el-col>
          <el-col :span="1.5">
            <el-button type="danger" plain :disabled="multiple" @click="handleDelete" v-hasPerms="['system:user:remove']"><ep:delete /> 删除</el-button>
          </el-col>
          <el-col :span="1.5">
            <el-button type="info" plain @click="handleImport" v-hasPerms="['system:user:import']"><ep:download />导入</el-button>
          </el-col>
          <el-col :span="1.5">
            <el-button type="warning" plain @click="handleExport" v-hasPerms="['system:user:export']"><ep:top />导出</el-button>
          </el-col>
          <right-toolbar v-model:showSearch="showSearch" @queryTable="getList" :columns="columns"></right-toolbar>
        </el-row>

        <el-table border v-loading="loading" :data="userList" @selection-change="handleSelectionChange">
          <el-table-column type="selection" width="50" align="center" />
          <el-table-column label="用户编号" align="center" key="userId" prop="userId" v-if="columns[0].visible" />
          <el-table-column label="用户名称" align="center" key="username" prop="username" v-if="columns[1].visible" :show-overflow-tooltip="true" />
          <el-table-column label="用户昵称" align="center" key="nickname" prop="nickname" v-if="columns[2].visible" :show-overflow-tooltip="true" />
          <el-table-column label="部门" align="center" key="deptName" prop="deptName" v-if="columns[3].visible" :show-overflow-tooltip="true" />
          <el-table-column label="手机号码" align="center" key="phoneNumber" prop="phoneNumber" v-if="columns[4].visible" width="120" />
          <el-table-column label="状态" align="center" key="status" v-if="columns[5].visible">
            <template #default="scope">
              <el-switch v-model="scope.row.status" active-value="0" inactive-value="1" @change="handleStatusChange(scope.row)"></el-switch>
            </template>
          </el-table-column>
          <el-table-column label="创建时间" align="center" prop="createTime" v-if="columns[6].visible" width="160">
            <template #default="scope">
              <span>{{ parseTime(scope.row.createTime) }}</span>
            </template>
          </el-table-column>
          <el-table-column label="操作" align="center" width="180" class-name="small-padding fixed-width">
            <template #default="scope">
              <el-tooltip content="修改" placement="top" v-if="scope.row.userId !== 1">
                <el-button link type="primary" @click="handleUpdate(scope.row)" v-hasPerms="['system:user:edit']"><ep:edit /></el-button>
              </el-tooltip>
              <el-tooltip content="重置密码" placement="top" v-if="scope.row.userId !== 1">
                <el-button link type="primary" @click="handleResetPwd(scope.row)" v-hasPerms="['system:user:resetPwd']"><ep:key /></el-button>
              </el-tooltip>
              <el-tooltip content="分配角色" placement="top" v-if="scope.row.userId !== 1">
                <el-button link type="primary" @click="handleAuthRole(scope.row)" v-hasPerms="['system:user:edit']"><ep:circle-check /></el-button>
              </el-tooltip>
              <el-tooltip content="删除" placement="top" v-if="scope.row.userId !== 1">
                <el-button link type="danger" @click="handleDelete(scope.row)" v-hasPerms="['system:user:remove']"><ep:delete /></el-button>
              </el-tooltip>
            </template>
          </el-table-column>
        </el-table>
        <pagination v-show="total > 0" :total="total" v-model:page="queryParams.pageNum" v-model:limit="queryParams.pageSize" @pagination="getList" />
      </el-col>
    </el-row>

    <!-- 添加或修改用户配置对话框 -->
    <el-dialog :title="dialog.title" v-model="dialog.visible" width="600px" @close="closeDialog" append-to-body>
      <el-form :model="form" :rules="rules" ref="userFormRef" label-width="80px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="用户昵称" prop="nickname">
              <el-input v-model="form.nickname" placeholder="请输入用户昵称" maxlength="30" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="归属部门" prop="deptId">
              <el-tree-select
                v-model="form.deptId"
                :data="deptOptions"
                :props="{ value: 'deptId', label: 'deptName', children: 'children' }"
                value-key="id"
                placeholder="请选择归属部门"
                check-strictly
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="手机号码" prop="phoneNumber">
              <el-input v-model="form.phoneNumber" placeholder="请输入手机号码" maxlength="11" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="form.email" placeholder="请输入邮箱" maxlength="50" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item v-if="form.userId == undefined" label="用户名称" prop="username">
              <el-input v-model="form.username" placeholder="请输入用户名称" maxlength="30" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item v-if="form.userId == undefined" label="用户密码" prop="password">
              <el-input v-model="form.password" placeholder="请输入用户密码" type="password" maxlength="20" show-password />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="用户性别" prop="gender">
              <el-select v-model="form.gender" placeholder="请选择">
                <el-option v-for="dict in sys_user_sex" :key="dict.value" :label="dict.label" :value="dict.value"></el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-radio-group v-model="form.status">
                <el-radio v-for="dict in sys_normal_disable" :key="dict.value" :label="dict.value">{{ dict.label }}</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="12">
            <el-form-item label="岗位" prop="postIds">
              <el-select v-model="form.postIds" multiple placeholder="请选择">
                <el-option v-for="item in postOptions" :key="item.postId" :label="item.postName" :value="item.postId" :disabled="item.status == 1"></el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="角色" prop="roleIds">
              <el-select v-model="form.roleIds" multiple placeholder="请选择">
                <el-option v-for="item in roleOptions" :key="item.roleId" :label="item.roleName" :value="item.roleId" :disabled="item.status == 1"></el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="24">
            <el-form-item label="备注" prop="remark">
              <el-input v-model="form.remark" type="textarea" placeholder="请输入内容"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button type="primary" @click="submitForm">确 定</el-button>
          <el-button @click="closeDialog">取 消</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 用户导入对话框 -->
    <el-dialog :title="upload.title" v-model="upload.open" width="400px" append-to-body>
      <el-upload
        ref="uploadRef"
        :limit="1"
        accept=".xlsx, .xls"
        :headers="upload.headers"
        :action="upload.url + '?updateSupport=' + upload.updateSupport"
        :disabled="upload.isUploading"
        :on-progress="handleFileUploadProgress"
        :on-success="handleFileSuccess"
        :auto-upload="false"
        drag
      >
        <el-icon class="el-icon--upload"><ep:upload /></el-icon>
        <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
        <template #tip>
          <div class="el-upload__tip text-center">
            <div class="el-upload__tip"><el-checkbox v-model="upload.updateSupport" />是否更新已经存在的用户数据</div>
            <span>仅允许导入xls、xlsx格式文件。</span>
            <el-link type="primary" :underline="false" style="font-size: 12px; vertical-align: baseline" @click="importTemplate">下载模板</el-link>
          </div>
        </template>
      </el-upload>
      <template #footer>
        <div class="dialog-footer">
          <el-button type="primary" @click="submitFileForm">确 定</el-button>
          <el-button @click="upload.open = false">取 消</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup name="User">
import { getToken } from '@/utils/auth'
import { changeUserStatus, listUser, resetUserPwd, delUser, getUser, updateUser, addUser, deptTreeSelect } from '@/api/system/user'
import useCurrentInstance from '@/hooks/useCurrentInstance'
import { UserFormData } from '@/types/api/user'
import { PostFormData } from '@/types/api/post'
import { RoleFormData } from '@/types/api/role'
import { Dialog } from '@/types/common'

const router = useRouter()
const { proxy } = useCurrentInstance()
const { sys_normal_disable, sys_user_sex } = proxy.useDict('sys_normal_disable', 'sys_user_sex')

const userList = ref([])
const loading = ref(true)
const showSearch = ref(true)
const ids = ref([])
const single = ref(true)
const multiple = ref(true)
const total = ref(0)
const dateRange = ref([])
const deptName = ref('')
const deptOptions = ref(undefined)
const initPassword = ref(undefined)
const postOptions = ref<PostFormData[]>([])
const roleOptions = ref<RoleFormData[]>([])
/*** 用户导入参数 */
const upload = reactive({
  // 是否显示弹出层（用户导入）
  open: false,
  // 弹出层标题（用户导入）
  title: '',
  // 是否禁用上传
  isUploading: false,
  // 是否更新已经存在的用户数据
  updateSupport: 0,
  // 设置上传的请求头部
  headers: { Authorization: 'Bearer ' + getToken() },
  // 上传的地址
  url: import.meta.env.VITE_APP_BASE_API + '/system/user/importData'
})
// 列显隐信息
const columns = ref([
  { key: 0, label: `用户编号`, visible: true },
  { key: 1, label: `用户名称`, visible: true },
  { key: 2, label: `用户昵称`, visible: true },
  { key: 3, label: `部门`, visible: true },
  { key: 4, label: `手机号码`, visible: true },
  { key: 5, label: `状态`, visible: true },
  { key: 6, label: `创建时间`, visible: true }
])

const data = reactive({
  dialog: {
    visible: false,
    title: ''
  } as Dialog,
  form: {
    status: '0'
  } as UserFormData,
  queryParams: {
    pageNum: 1,
    pageSize: 10,
    username: undefined,
    phoneNumber: undefined,
    status: undefined,
    deptId: undefined
  },
  rules: {
    username: [
      { required: true, message: '用户名称不能为空', trigger: 'blur' },
      { min: 2, max: 20, message: '用户名称长度必须介于 2 和 20 之间', trigger: 'blur' }
    ],
    nickname: [{ required: true, message: '用户昵称不能为空', trigger: 'blur' }],
    password: [
      { required: true, message: '用户密码不能为空', trigger: 'blur' },
      { min: 5, max: 20, message: '用户密码长度必须介于 5 和 20 之间', trigger: 'blur' }
    ],
    email: [{ type: 'email', message: '请输入正确的邮箱地址', trigger: ['blur', 'change'] }],
    phoneNumber: [{ pattern: /^1[3|4|5|6|7|8|9][0-9]\d{8}$/, message: '请输入正确的手机号码', trigger: 'blur' }]
  }
})

const { dialog, queryParams, form, rules } = toRefs(data)

const deptTreeRef = ref<ElTree>(null) // 部门树
const queryFormRef = ref<ElForm>(null) // 查询表单
const userFormRef = ref<ElForm>(null) // 用户表单

/** 通过条件过滤节点  */
const filterNode = (value: string, data:any) => {
  if (!value) return true
  return data.deptName.indexOf(value) !== -1
}
/** 根据名称筛选部门树 */
watch(deptName, val => {
  deptTreeRef.value.filter(val)
})
/** 查询部门下拉树结构 */
function getDeptTree() {
  deptTreeSelect().then(response => {
    deptOptions.value = response.data
  })
}
/** 查询用户列表 */
function getList() {
  loading.value = true
  listUser(queryParams.value).then(res => {
    loading.value = false
    userList.value = res.data.rows
    total.value = res.data.total
  })
}
/** 节点单击事件 */
function handleNodeClick(data: { [key: string]: any }) {
  queryParams.value.deptId = data.deptId
  handleQuery()
}
/** 搜索按钮操作 */
function handleQuery() {
  queryParams.value.pageNum = 1
  getList()
}
/** 重置按钮操作 */
function resetQuery() {
  dateRange.value = []
  queryFormRef.value?.resetFields()
  queryParams.value.deptId = undefined
  deptTreeRef.value.setCurrentKey(null)
  handleQuery()
}
/** 删除按钮操作 */
function handleDelete(row: { [key: string]: any }) {
  let userIds: string[] = []
  if (row.userId !== undefined) {
    userIds.push(row.userId)
  } else {
    userIds = ids.value
  }
  proxy.$modal
    .confirm('是否确认删除用户编号为"' + userIds + '"的数据项？')
    .then(function () {
      return delUser({ids: userIds})
    })
    .then(() => {
      getList()
      proxy.$modal.msgSuccess('删除成功')
    })
    .catch(() => {})
}
/** 导出按钮操作 */
function handleExport() {
  proxy.download(
    'system/user/export',
    {
      ...queryParams.value
    },
    `user_${new Date().getTime()}.xlsx`
  )
}
/** 用户状态修改  */
function handleStatusChange(row: { [key: string]: any }) {
  let text = row.status === '0' ? '启用' : '停用'
  proxy.$modal
    .confirm('确认要"' + text + '""' + row.username + '"用户吗?')
    .then(function () {
      return changeUserStatus(row.userId, row.status)
    })
    .then(() => {
      proxy.$modal.msgSuccess(text + '成功')
    })
    .catch(function () {
      row.status = row.status === '0' ? '1' : '0'
    })
}
/** 更多操作 */
function handleCommand(command: string, row: { [key: string]: any }) {
  switch (command) {
    case 'handleResetPwd':
      handleResetPwd(row)
      break
    case 'handleAuthRole':
      handleAuthRole(row)
      break
    default:
      break
  }
}
/** 跳转角色分配 */
function handleAuthRole(row: { [key: string]: any }) {
  const userId = row.userId
  router.push('/system/user-auth/role/' + userId)
}
/** 重置密码按钮操作 */
function handleResetPwd(row: { [key: string]: any }) {
  ElMessageBox.prompt('请输入"' + row.username + '"的新密码', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      closeOnClickModal: false,
      inputPattern: /^.{5,20}$/,
      inputErrorMessage: '用户密码长度必须介于 5 和 20 之间'
    })
    .then(({ value }) => {
      resetUserPwd({userId: row.userId, password: value}).then(() => {
        proxy.$modal.msgSuccess('修改成功，新密码是：' + value)
      })
    })
    .catch(() => {})
}
/** 选择条数  */
function handleSelectionChange(selection: any) {
  ids.value = selection.map((item: any) => item.userId)
  single.value = selection.length != 1
  multiple.value = !selection.length
}
/** 导入按钮操作 */
function handleImport() {
  upload.title = '用户导入'
  upload.open = true
}
/** 下载模板操作 */
function importTemplate() {
  proxy.download('system/user/importTemplate', {}, `user_template_${new Date().getTime()}.xlsx`)
}
/**文件上传中处理 */
const handleFileUploadProgress = () => {
  upload.isUploading = true
}
/** 文件上传成功处理 */
const handleFileSuccess = (response: any, file: UploadFile, fileList: UploadFiles) => {
  upload.open = false
  upload.isUploading = false
  proxy.$refs['uploadRef'].handleRemove(file)
  proxy.$alert("<div style='overflow: auto;overflow-x: hidden;max-height: 70vh;padding: 10px 20px 0;'>" + response.msg + '</div>', '导入结果', {
    dangerouslyUseHTMLString: true
  })
  getList()
}
/** 提交上传文件 */
function submitFileForm() {
  proxy.$refs['uploadRef'].submit()
}
/** 关闭用户弹窗 */
function closeDialog() {
  dialog.value.visible =  false
  userFormRef.value?.resetFields()
  userFormRef.value?.clearValidate()
}
/** 新增按钮操作 */
function handleAdd() {
  getUser("-1").then(response => {
    postOptions.value = response.data.posts
    roleOptions.value = response.data.roles
    dialog.value.visible = true
    dialog.value.title = '添加用户'
    form.value.password = initPassword.value
  })
}
/** 修改按钮操作 */
function handleUpdate(row: { [key: string]: any }) {
  const userId = row.userId || ids.value
  getUser(userId).then(response => {
    form.value = response.data.user
    postOptions.value = response.data.posts
    roleOptions.value = response.data.roles
    form.value.postIds = response.data.postIds
    form.value.roleIds = response.data.roleIds
    dialog.value.visible = true
    dialog.value.title = '修改用户'
    form.value.password = ''
  })
}
/** 提交按钮 */
function submitForm() {
  userFormRef.value.validate((valid: any) => {
    if (valid) {
      const data = {
        sysUser: form.value,
        postIds: form.value.postIds,
        roleIds: form.value.roleIds
      }
      if (form.value.userId != undefined) {
        updateUser(data).then(() => {
          proxy.$modal.msgSuccess('修改成功')
          closeDialog()
          getList()
        })
      } else {
        addUser(data).then(() => {
          proxy.$modal.msgSuccess('新增成功')
          closeDialog()
          getList()
        })
      }
    }
  })
}

onMounted(() => {
  // 初始化部门树
  getDeptTree()
  // 初始化用户列表
  getList()
})
</script>
