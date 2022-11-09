<template>
  <div :class="classObj" class="app-wrapper" :style="{ '--current-color': theme }">
    <div v-if="device === 'mobile' && sidebarObj.opened" class="drawer-bg" @click="handleClickOutside" />
    <sidebar v-if="!sidebarObj.hide" class="sidebar-container" />
    <div :class="{ hasTagsView: needTagsView, sidebarHide: sidebarObj.hide }" class="main-container">
      <div :class="{ 'fixed-header': fixedHeader }">
        <navbar @setLayout="setLayout" />
        <tags-view v-if="needTagsView" />
      </div>
      <app-main />
      <settings ref="settingRef" />
    </div>
  </div>
</template>

<script lang="ts" setup>
import useStore from '@/store'
const { app, setting } = useStore()

const theme = computed(() => setting.theme)
const sideTheme = computed(() => setting.sideTheme)
const sidebarObj = computed(() => app.sidebar)
const device = computed(() => app.device)
const needTagsView = computed(() => setting.tagsView)
const fixedHeader = computed(() => setting.fixedHeader)

const classObj = computed(() => ({
  hideSidebar: !sidebarObj.value.opened,
  openSidebar: sidebarObj.value.opened,
  withoutAnimation: sidebarObj.value.withoutAnimation,
  mobile: device.value === 'mobile'
}))

const { width, height } = useWindowSize()
const WIDTH = 992 // refer to Bootstrap's responsive design

watchEffect(() => {
  if (device.value === 'mobile' && sidebarObj.value.opened) {
    app.closeSideBar({ withoutAnimation: false })
  }
  if (width.value - 1 < WIDTH) {
    app.toggleDevice('mobile')
    app.closeSideBar({ withoutAnimation: true })
  } else {
    app.toggleDevice('desktop')
  }
})

function handleClickOutside() {
  app.closeSideBar({ withoutAnimation: false })
}

const settingRef = ref(null)
function setLayout() {
  settingRef.value.openSetting()
}
</script>

<style lang="scss" scoped>
@import '@/assets/styles/mixin.scss';
@import '@/assets/styles/variables.module.scss';

.app-wrapper {
  @include clearfix;
  position: relative;
  height: 100%;
  width: 100%;

  &.mobile.openSidebar {
    position: fixed;
    top: 0;
  }
}

.drawer-bg {
  background: #000;
  opacity: 0.3;
  width: 100%;
  top: 0;
  height: 100%;
  position: absolute;
  z-index: 999;
}

.fixed-header {
  position: fixed;
  top: 0;
  right: 0;
  z-index: 9;
  width: calc(100% - #{$base-sidebar-width});
  transition: width 0.28s;
}

.hideSidebar .fixed-header {
  width: calc(100% - 54px);
}

.sidebarHide .fixed-header {
  width: 100%;
}

.mobile .fixed-header {
  width: 100%;
}
</style>
