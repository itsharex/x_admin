<template>
  <view class="page-content">


    <uv-form labelPosition="left" labelWidth="80" :model="form">
      <uv-form-item label="项目uuid" prop="projectKey" borderBottom>
        {{ form.projectKey }}
      </uv-form-item>
      <uv-form-item label="项目名称" prop="projectName" borderBottom>
        {{ form.projectName }}
      </uv-form-item>
      <uv-form-item label="项目类型" prop="projectType" borderBottom>
        {{ form.projectType }}
      </uv-form-item>
      <uv-form-item label="创建时间" prop="createTime" borderBottom>
        {{ form.createTime }}
      </uv-form-item>
      <uv-form-item label="更新时间" prop="updateTime" borderBottom>
        {{ form.updateTime }}
      </uv-form-item>
    </uv-form>
    <uv-button
      v-if="$perms('admin:monitor_project:edit')"
      type="primary"
      text="编辑"
      customStyle="margin: 40rpx 0"
      @click="edit"
    ></uv-button>
  </view>
</template>

<script setup lang="ts">
import { reactive, ref, computed } from "vue";
import { onLoad, onShow, onPullDownRefresh } from "@dcloudio/uni-app";

import { useDictData } from "@/hooks/useDictOptions";
import { monitor_project_detail } from "@/api/monitor_project";
import type { type_monitor_project } from "@/api/monitor_project";

import { toast, alert, toPath } from "@/utils/utils";

let form = ref<type_monitor_project>({
  id: 0,
  projectKey: "",
  projectName: "",
  projectType: "",
  createTime: "",
  updateTime: "",
});
onLoad((e) => {
  console.log("onLoad", e);
  if (e.id) {
    getDetails(e.id);
  }
});
onShow((e) => {
  if (form.value?.id) {
    getDetails(form.value.id);
  }
});
onPullDownRefresh(() => {
  getDetails(form.value.id);
});
function getDetails(id) {
  monitor_project_detail(id)
    .then((res) => {
      uni.stopPullDownRefresh();
      if (res.code == 200) {
        if (res.data) {
          form.value = res?.data;
        }
      } else {
        toast(res.message);
      }
    })
    .catch((err) => {
      uni.stopPullDownRefresh();
      toast(err.message || "网络错误");
    });
}

function edit() {
  toPath("/pages/monitor_project/edit", { id: form.value.id });
}
</script>

<style lang="scss" scoped>
.page-content {
  padding: 10rpx 20rpx 300rpx;
}
</style>
