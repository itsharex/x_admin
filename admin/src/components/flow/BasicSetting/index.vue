<template>
    <div class="setting-container">
        <el-form
            ref="elFormRef"
            :model="formData"
            :rules="rules"
            label-width="100px"
            label-position="top"
        >
            <el-form-item label="审批名称" prop="flowName">
                <el-input
                    v-model="formData.flowName"
                    placeholder="请输入审批名称"
                    clearable
                    :style="{ width: '100%' }"
                >
                </el-input>
            </el-form-item>
            <el-form-item label="选择分组" prop="flowGroup">
                <el-select
                    v-model="formData.flowGroup"
                    placeholder="请选择选择分组"
                    clearable
                    :style="{ width: '100%' }"
                >
                    <el-option
                        v-for="item in dictData.flow_group"
                        :key="item.id"
                        :label="item.name"
                        :value="parseInt(item.value)"
                    ></el-option>
                </el-select>
            </el-form-item>
            <!-- <el-form-item label="谁可以发起审批" prop="approver">
        <span style="font-size: 12px; color: #aaa">默认所有人</span>
      </el-form-item> -->

            <el-form-item label="流程描述" prop="flowRemark">
                <el-input
                    v-model="formData.flowRemark"
                    type="textarea"
                    placeholder="请输入流程描述"
                    :maxlength="100"
                    show-word-limit
                    :autosize="{ minRows: 4, maxRows: 4 }"
                    :style="{ width: '100%' }"
                ></el-input>
            </el-form-item>
        </el-form>
    </div>
</template>

<script lang="ts" setup>
import { onMounted, ref, shallowRef } from 'vue'
import type { FormInstance } from 'element-plus'
import { useDictData } from '@/hooks/useDictOptions'
defineOptions({
    name: 'BasicSetting'
})
const { dictData } = useDictData(['flow_group'])
const props = defineProps(['tabName', 'conf'])
const formData = ref({
    flowName: '',
    flowImg: '',
    flowGroup: undefined,
    flowRemark: undefined
})
const rules = {
    flowName: [
        {
            required: true,
            message: '请输入审批名称',
            trigger: 'blur'
        }
    ],
    flowGroup: [
        {
            required: true,
            message: '请选择选择分组',
            trigger: 'change'
        }
    ]
}
onMounted(() => {
    if (typeof props.conf === 'object' && props.conf !== null) {
        Object.assign(formData.value, props.conf)
    }
})
const elFormRef = shallowRef<FormInstance>()
function getData() {
    return new Promise((resolve, reject) => {
        elFormRef.value.validate((valid) => {
            if (!valid) {
                reject({ target: this.tabName })
                return
            }
            // this.formData.flowImg = this.activeIcon
            resolve({ formData: { ...formData.value }, target: props.tabName }) // TODO 提交表单
        })
    })
}
defineExpose({
    getData
})
</script>

<style lang="scss" scoped>
.icon-item {
    width: 40px;
    height: 40px;
    margin: 6px;
    position: relative;
    cursor: pointer;
    opacity: 0.5;
    &.active {
        opacity: 1;
    }
    &:hover {
        opacity: 1;
    }
}
.setting-container {
    width: 600px;
    margin: auto;
    background: #fff;
    padding: 16px;
}
</style>
