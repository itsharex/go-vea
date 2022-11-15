<template>
  <div class="app-container">
    <el-form :model="queryParams" ref="queryFormRef" :inline="true" v-show="showSearch">
      <el-form-item label="菜单名称" prop="menuName">
        <el-input v-model="queryParams.menuName" placeholder="请输入菜单名称" clearable @keyup.enter="handleQuery" />
      </el-form-item>
      <el-form-item label="状态" prop="status">
        <el-select v-model="queryParams.status" placeholder="菜单状态" clearable>
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
        <el-button type="primary" plain @click="handleAdd" v-hasPerms="['system:menu:add']"><ep:plus />新增</el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button type="info" plain @click="toggleExpandAll"><ep:sort />展开/折叠</el-button>
      </el-col>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList"></right-toolbar>
    </el-row>

    <el-table
      border
      v-if="refreshTable"
      v-loading="loading"
      :data="menuList"
      row-key="menuId"
      :default-expand-all="isExpandAll"
      :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
    >
      <el-table-column prop="menuName" label="菜单名称" :show-overflow-tooltip="true" width="160"></el-table-column>
      <el-table-column prop="icon" label="图标" align="center" width="100">
        <template #default="scope">
          <svg-icon :icon-class="scope.row.icon" />
        </template>
      </el-table-column>
      <el-table-column prop="orderNum" label="排序" width="60"></el-table-column>
      <el-table-column prop="perms" label="权限标识" :show-overflow-tooltip="true"></el-table-column>
      <el-table-column prop="component" label="组件路径" :show-overflow-tooltip="true"></el-table-column>
      <el-table-column prop="status" label="状态" width="80">
        <template #default="scope">
          <dict-tag :options="sys_normal_disable" :value="scope.row.status" />
        </template>
      </el-table-column>
      <el-table-column label="创建时间" align="center" prop="createTime">
        <template #default="scope">
          <span>{{ parseTime(scope.row.createTime) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" width="250" class-name="small-padding fixed-width">
        <template #default="scope">
          <el-button link type="primary" @click="handleUpdate(scope.row)" v-hasPerms="['system:menu:edit']"><ep:edit />修改</el-button>
          <el-button link type="primary" @click="handleAdd(scope.row)" v-hasPerms="['system:menu:add']"><ep:plus />新增</el-button>
          <el-button link type="danger" @click="handleDelete(scope.row)" v-hasPerms="['system:menu:remove']"><ep:delete />删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 添加或修改菜单对话框 -->
    <el-dialog :title="dialog.title" v-model="dialog.visible" width="680px" @close="closeDialog" append-to-body>
      <el-form ref="menuFormRef" :model="form" :rules="rules" label-width="100px">
        <el-row>
          <el-col :span="24">
            <el-form-item label="上级菜单" prop="parentId">
              <el-tree-select
                style="width: 100%"
                v-model="form.parentId"
                :data="menuOptions"
                :props="{ value: 'menuId', label: 'menuName', children: 'children' }"
                value-key="menuId"
                placeholder="选择上级菜单"
                check-strictly
              />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="菜单类型" prop="menuType">
              <el-radio-group v-model="form.menuType">
                <el-radio label="M">目录</el-radio>
                <el-radio label="C">菜单</el-radio>
                <el-radio label="F">按钮</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="24" v-if="form.menuType != 'F'">
            <el-form-item label="菜单图标" prop="icon">
              <el-popover placement="bottom-start" :width="540" v-model:visible="showChooseIcon" trigger="click" @show="showSelectIcon">
                <template #reference>
                  <el-input v-model="form.icon" placeholder="点击选择图标" v-click-outside="hideSelectIcon" readonly>
                    <template #prefix>
                      <svg-icon v-if="form.icon" :icon-class="form.icon" class="el-input__icon" style="height: 32px; width: 16px" />
                      <el-icon v-else style="height: 32px; width: 16px"><ep:search /></el-icon>
                    </template>
                  </el-input>
                </template>
                <icon-select ref="iconSelectRef" @selected="selected" />
              </el-popover>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="菜单名称" prop="menuName">
              <el-input v-model="form.menuName" placeholder="请输入菜单名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="显示排序" prop="orderNum">
              <el-input-number v-model="form.orderNum" controls-position="right" :min="0" />
            </el-form-item>
          </el-col>
          <el-col :span="12" v-if="form.menuType != 'F'">
            <el-form-item prop="isFrame">
              <template #label>
                <span>
                  <el-tooltip content="选择是外链则路由地址需要以`http(s)://`开头" placement="top">
                    <el-icon><ep:question-filled /></el-icon> </el-tooltip
                  >是否外链
                </span>
              </template>
              <el-radio-group v-model="form.isFrame">
                <el-radio :label=0>是</el-radio>
                <el-radio :label=1>否</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="12" v-if="form.menuType != 'F'">
            <el-form-item prop="path">
              <template #label>
                <span>
                  <el-tooltip content="访问的路由地址，如：`user`，如外网地址需内链访问则以`http(s)://`开头" placement="top">
                    <el-icon><ep:question-filled /></el-icon>
                  </el-tooltip>
                  路由地址
                </span>
              </template>
              <el-input v-model="form.path" placeholder="请输入路由地址" />
            </el-form-item>
          </el-col>
          <el-col :span="12" v-if="form.menuType == 'C'">
            <el-form-item prop="component">
              <template #label>
                <span>
                  <el-tooltip content="访问的组件路径，如：`system/user/index`，默认在`views`目录下" placement="top">
                    <el-icon><ep:question-filled /></el-icon>
                  </el-tooltip>
                  组件路径
                </span>
              </template>
              <el-input v-model="form.component" placeholder="请输入组件路径" />
            </el-form-item>
          </el-col>
          <el-col :span="12" v-if="form.menuType != 'M'">
            <el-form-item prop="perms">
              <el-input v-model="form.perms" placeholder="请输入权限标识" maxlength="100" />
              <template #label>
                <span>
                  <el-tooltip content="控制器中定义的权限字符，如：@PreAuthorize(`@ss.hasPermi('system:user:list')`)" placement="top">
                    <el-icon><ep:question-filled /></el-icon>
                  </el-tooltip>
                  权限字符
                </span>
              </template>
            </el-form-item>
          </el-col>
          <el-col :span="12" v-if="form.menuType == 'C'">
            <el-form-item prop="query">
              <el-input v-model="form.query" placeholder="请输入路由参数" maxlength="255" />
              <template #label>
                <span>
                  <el-tooltip content="访问路由的默认传递参数，如：`{'id': 1, 'user': 'ry'}`" placement="top">
                    <el-icon><ep:question-filled /></el-icon>
                  </el-tooltip>
                  路由参数
                </span>
              </template>
            </el-form-item>
          </el-col>
          <el-col :span="12" v-if="form.menuType == 'C'">
            <el-form-item prop="isCache">
              <template #label>
                <span>
                  <el-tooltip content="选择是则会被`keep-alive`缓存，需要匹配组件的`name`和地址保持一致" placement="top">
                    <el-icon><ep:question-filled /></el-icon>
                  </el-tooltip>
                  是否缓存
                </span>
              </template>
              <el-radio-group v-model="form.isCache">
                <el-radio :label=0>缓存</el-radio>
                <el-radio :label=1>不缓存</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="12" v-if="form.menuType != 'F'">
            <el-form-item prop="visible">
              <template #label>
                <span>
                  <el-tooltip content="选择隐藏则路由将不会出现在侧边栏，但仍然可以访问" placement="top">
                    <el-icon><ep:question-filled /></el-icon>
                  </el-tooltip>
                  显示状态
                </span>
              </template>
              <el-radio-group v-model="form.visible">
                <el-radio v-for="dict in sys_show_hide" :key="dict.value" :label="dict.value">{{ dict.label }}</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="12" v-if="form.menuType != 'F'">
            <el-form-item prop="status">
              <template #label>
                <span>
                  <el-tooltip content="选择停用则路由将不会出现在侧边栏，也不能被访问" placement="top">
                    <el-icon><ep:question-filled /></el-icon>
                  </el-tooltip>
                  菜单状态
                </span>
              </template>
              <el-radio-group v-model="form.status">
                <el-radio v-for="dict in sys_normal_disable" :key="dict.value" :label="dict.value">{{ dict.label }}</el-radio>
              </el-radio-group>
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
  </div>
</template>

<script lang="ts" setup name="Menu">
import { addMenu, delMenu, getMenu, listMenuTree, updateMenu } from '@/api/system/menu'
import { ClickOutside as vClickOutside } from 'element-plus'
import { MenuFormData,MenuOptions } from '@/types/api/menu'
import useCurrentInstance from '@/hooks/useCurrentInstance'
import { Dialog } from '@/types/common'

const { proxy } = useCurrentInstance()
const { sys_show_hide, sys_normal_disable } = proxy.useDict('sys_show_hide', 'sys_normal_disable')

const menuList = ref([])
const loading = ref(true)
const showSearch = ref(true)
const menuOptions = ref<MenuOptions[]>([])
const isExpandAll = ref(false)
const refreshTable = ref(true)
const showChooseIcon = ref(false)
const iconSelectRef = ref<IconSelect>(null)

const queryFormRef = ref<ElForm>(null)
const menuFormRef = ref<ElForm>(null)

const data = reactive({
  dialog: {
    visible: false,
    title: ''
  } as Dialog,
  form: {
    parentId: 0,
    menuType: 'M',
    isFrame: 1,
    isCache: 0,
    visible: '0',
    status: '0'
  } as MenuFormData,
  queryParams: {
    menuName: undefined,
    visible: undefined,
    status: undefined
  },
  rules: {
    menuName: [{ required: true, message: '菜单名称不能为空', trigger: 'blur' }],
    orderNum: [{ required: true, message: '菜单顺序不能为空', trigger: 'blur' }],
    path: [{ required: true, message: '路由地址不能为空', trigger: 'blur' }]
  }
})

const { dialog, queryParams, form, rules } = toRefs(data)

/** 查询菜单列表 */
function getList() {
  loading.value = true
  listMenuTree(queryParams.value).then(response => {
    menuList.value = response.data
    loading.value = false
  })
}
/** 查询菜单下拉树结构 */
function getTreeSelect() {
  menuOptions.value = []
  listMenuTree(queryParams.value).then(response => {
    const menu = { menuId: 0, menuName: '主类目', children: [] }
    menu.children = response.data
    menuOptions.value.push(menu)
  })
}
function closeDialog() {
  dialog.value.visible =  false
  menuFormRef.value?.resetFields()
  menuFormRef.value?.clearValidate()
}
/** 展示下拉图标 */
function showSelectIcon() {
  iconSelectRef.value?.reset()
  showChooseIcon.value = true
}
/** 选择图标 */
function selected(name:string) {
  form.value.icon = name
  showChooseIcon.value = false
}
/** 图标外层点击隐藏下拉列表 */
function hideSelectIcon(event: any) {
  var elem = event.relatedTarget || event.srcElement || event.target || event.currentTarget
  var className = elem.className
  if (className !== 'el-input__inner') {
    showChooseIcon.value = false
  }
}
/** 搜索按钮操作 */
function handleQuery() {
  getList()
}
/** 重置按钮操作 */
function resetQuery() {
  queryFormRef.value.resetFields()
  handleQuery()
}
/** 新增按钮操作 */
function handleAdd(row: { [key: string]: any }) {
  getTreeSelect()
  if (row != null && row.menuId) {
    form.value.parentId = row.menuId
  } else {
    form.value.parentId = 0
  }
  dialog.value.visible = true
  dialog.value.title = '添加菜单'
}
/** 展开/折叠操作 */
function toggleExpandAll() {
  refreshTable.value = false
  isExpandAll.value = !isExpandAll.value
  nextTick(() => {
    refreshTable.value = true
  })
}
/** 修改按钮操作 */
async function handleUpdate(row: { [key: string]: any }) {
  getTreeSelect()
  getMenu(row.menuId).then(response => {
    form.value = response.data
    dialog.value.visible = true
    dialog.value.title = '修改菜单'
  })
}
/** 提交按钮 */
function submitForm() {
  menuFormRef.value.validate((valid:any) => {
    if (valid) {
      if (form.value.menuId != undefined) {
        updateMenu(form.value).then(() => {
          proxy.$modal.msgSuccess('修改成功')
          closeDialog()
          getList()
        })
      } else {
        addMenu(form.value).then(() => {
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
  proxy.$modal
    .confirm('是否确认删除名称为"' + row.menuName + '"的数据项?')
    .then(function () {
      return delMenu(row.menuId)
    })
    .then(() => {
      getList()
      proxy.$modal.msgSuccess('删除成功')
    })
    .catch(() => {})
}

onMounted(() => {
  getList()
})
</script>
