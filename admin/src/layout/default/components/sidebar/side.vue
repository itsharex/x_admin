<template>
    <div class="side" :style="sideStyle">
        <side-logo v-if="settingStore.showLogo" :show-title="!isCollapsed" :theme="sideTheme" />
        <side-menu
            :routes="routes"
            :isCollapsed="isCollapsed"
            :width="settingStore.sideWidth"
            :unique-opened="settingStore.isUniqueOpened"
            :config="menuProp"
            :theme="sideTheme"
            @select="handleSelect"
        />
    </div>
</template>

<script setup lang="ts">
import type { RouteRecordRaw } from 'vue-router'
import useAppStore from '@/stores/modules/app'
import useSettingStore from '@/stores/modules/setting'
import useUserStore from '@/stores/modules/user'
import SideLogo from './logo.vue'
import SideMenu from './menu.vue'
defineOptions({
    name: 'Side'
})
const appStore = useAppStore()
const isCollapsed = computed(() => {
    if (appStore.isMobile) {
        return false
    } else {
        return appStore.isCollapsed
    }
})

const settingStore = useSettingStore()
const sideTheme = computed(() => settingStore.sideTheme)
const userStore = useUserStore()

const routes = computed(() => userStore.routes as RouteRecordRaw[])

const sideStyle = computed(() => {
    return sideTheme.value == 'dark'
        ? {
              '--side-dark-color': settingStore.sideDarkColor
          }
        : ''
})
const menuProp = computed(() => {
    return {
        backgroundColor: sideTheme.value == 'dark' ? settingStore.sideDarkColor : '',
        textColor: sideTheme.value == 'dark' ? 'var(--el-color-white)' : '',
        activeTextColor: sideTheme.value == 'dark' ? 'var(--el-color-white)' : ''
    }
})
const handleSelect = () => {
    if (appStore.isMobile) {
        appStore.toggleCollapsed(true)
    }
}
</script>

<style lang="scss" scoped>
.side {
    position: relative;
    z-index: 999;
    border-right-width: 1px;
    border-color: var(--el-border-color-light);
    height: 100%;
    display: flex;
    flex-direction: column;

    background-color: var(--side-dark-color, var(--el-bg-color));
}
</style>
