<template>
  <div class="app-container">
    <el-form :model="queryParams" ref="queryFormRef" :inline="true" v-show="showSearch" label-width="68px">
      <el-form-item label="岗位编码" prop="postCode">
        <el-input v-model="queryParams.postCode" placeholder="请输入岗位编码" clearable @keyup.enter="handleQuery" />
      </el-form-item>
      <el-form-item label="岗位名称" prop="postName">
        <el-input v-model="queryParams.postName" placeholder="请输入岗位名称" clearable @keyup.enter="handleQuery" />
      </el-form-item>
      <el-form-item label="状态" prop="status">
        <el-select v-model="queryParams.status" placeholder="岗位状态" clearable>
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
        <el-button type="primary" plain @click="handleAdd" v-hasPerms="['system:post:add']"><ep:plus /> 新增</el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button type="danger" plain :disabled="multiple" @click="handleDelete" v-hasPerms="['system:post:remove']"><ep:delete /> 删除</el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button type="warning" plain @click="handleExport" v-hasPerms="['system:post:export']"><ep:download /> 导出</el-button>
      </el-col>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList"></right-toolbar>
    </el-row>

    <el-table border v-loading="loading" :data="postList" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" align="center" />
      <el-table-column label="岗位编号" align="center" prop="postId" />
      <el-table-column label="岗位编码" align="center" prop="postCode" />
      <el-table-column label="岗位名称" align="center" prop="postName" />
      <el-table-column label="岗位排序" align="center" prop="postSort" />
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
          <el-button link type="primary" @click="handleUpdate(scope.row)" v-hasPerms="['system:post:edit']"><ep:edit /> 修改</el-button>
          <el-button link type="danger" @click="handleDelete(scope.row)" v-hasPerms="['system:post:remove']"><ep:delete /> 删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total > 0" :total="total" v-model:page="queryParams.pageNum" v-model:limit="queryParams.pageSize" @pagination="getList" />

    <!-- 添加或修改岗位对话框 -->
    <el-dialog :title="dialog.title" v-model="dialog.visible" width="500px" @close="closeDialog" append-to-body>
      <el-form ref="postFormRef" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="岗位名称" prop="postName">
          <el-input v-model="form.postName" placeholder="请输入岗位名称" />
        </el-form-item>
        <el-form-item label="岗位编码" prop="postCode">
          <el-input v-model="form.postCode" placeholder="请输入编码名称" />
        </el-form-item>
        <el-form-item label="岗位顺序" prop="postSort">
          <el-input-number v-model="form.postSort" controls-position="right" :min="0" />
        </el-form-item>
        <el-form-item label="岗位状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio v-for="dict in sys_normal_disable" :key="dict.value" :label="dict.value">{{ dict.label }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="form.remark" type="textarea" placeholder="请输入内容" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button type="primary" @click="submitForm">确 定</el-button>
          <el-button @click="closeDialog">取 消</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup name="Post">
import { listPost, addPost, delPost, getPost, updatePost } from '@/api/system/post'
import useCurrentInstance from '@/hooks/useCurrentInstance'
import { PostFormData } from '@/types/api/post'
import { Dialog } from '@/types/common'

const { proxy } = useCurrentInstance()
const { sys_normal_disable } = proxy.useDict('sys_normal_disable')

const postList = ref([])
const loading = ref(true)
const showSearch = ref(true)
const ids = ref([])
const single = ref(true)
const multiple = ref(true)
const total = ref(0)

const queryFormRef = ref<ElForm>(null)
const postFormRef = ref<ElForm>(null)

const data = reactive({
  dialog: {
    visible: false,
    title: ''
  } as Dialog,
  form: {
    postSort: 0,
    status: '0'
  } as PostFormData,
  queryParams: {
    pageNum: 1,
    pageSize: 10,
    postCode: undefined,
    postName: undefined,
    status: undefined
  },
  rules: {
    postName: [{ required: true, message: '岗位名称不能为空', trigger: 'blur' }],
    postCode: [{ required: true, message: '岗位编码不能为空', trigger: 'blur' }],
    postSort: [{ required: true, message: '岗位顺序不能为空', trigger: 'blur' }]
  }
})

const { dialog, queryParams, form, rules } = toRefs(data)

/** 查询岗位列表 */
function getList() {
  loading.value = true
  listPost(queryParams.value).then(response => {
    postList.value = response.data.rows
    total.value = response.data.total
    loading.value = false
  })
}
/** 关闭弹窗 */
function closeDialog() {
  dialog.value.visible = false
  postFormRef.value?.resetFields()
  postFormRef.value?.clearValidate()
  form.value.postId = undefined
}
/** 搜索按钮操作 */
function handleQuery() {
  queryParams.value.pageNum = 1
  getList()
}
/** 重置按钮操作 */
function resetQuery() {
  queryFormRef.value.resetFields()
  handleQuery()
}
/** 多选框选中数据 */
function handleSelectionChange(selection: any) {
  ids.value = selection.map((item: any) => item.postId)
  single.value = selection.length != 1
  multiple.value = !selection.length
}
/** 新增按钮操作 */
function handleAdd() {
  dialog.value.visible = true
  dialog.value.title = '添加岗位'
}
/** 修改按钮操作 */
function handleUpdate(row: { [key: string]: any }) {
  const postId = row.postId || ids.value
  getPost(postId).then(response => {
    dialog.value.visible = true
    dialog.value.title = '修改岗位'
    nextTick(() => {
      form.value = response.data
    })
  })
}
/** 提交按钮 */
function submitForm() {
  postFormRef.value.validate((valid: any) => {
    if (valid) {
      if (form.value.postId != undefined) {
        updatePost(form.value).then(() => {
          proxy.$modal.msgSuccess('修改成功')
          closeDialog()
          getList()
        })
      } else {
        addPost(form.value).then(() => {
          proxy.$modal.msgSuccess('新增成功')
          closeDialog()
          getList()
        })
      }
    }
  })
}
/** 删除按钮操作 */
function handleDelete(row: { [key: string]: any }) {
  let postIds = []
  if (row.configId !== undefined) {
    postIds.push(row.postId)
  } else {
    postIds = ids.value
  }
  proxy.$modal
    .confirm('是否确认删除岗位编号为"' + postIds + '"的数据项？')
    .then(function () {
      return delPost({ ids: postIds })
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
    'system/post/export',
    {
      ...queryParams.value
    },
    `post_${new Date().getTime()}.xlsx`
  )
}

getList()
</script>
