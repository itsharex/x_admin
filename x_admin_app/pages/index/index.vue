<template>
	<view class="content">
		<uv-swiper :list="list" circular indicator indicatorMode="dot"></uv-swiper>
		<!-- <uv-notice-bar :text="text" direction="column"></uv-notice-bar> -->

		<uv-grid :border="false">
			<uv-grid-item v-for="(item, index) in baseList" :key="index" @click="toList(item)">
				<uv-icon :customStyle="{ padding: '50rpx 0 20rpx' }" :name="item.icon" :size="26"></uv-icon>

				<text class="grid-text">{{ item.title }}</text>
			</uv-grid-item>
		</uv-grid>
	</view>
</template>

<script setup lang="ts">
	import {
		ref
	} from "vue";

	import {
		toPath,
		toast,
		scanCode,
		queryToObj,
	} from "@/utils/utils";

	let list = [
		"https://cdn.uviewui.com/uview/swiper/swiper1.png",
		"https://cdn.uviewui.com/uview/swiper/swiper2.png",
		"https://cdn.uviewui.com/uview/swiper/swiper3.png",
	];
	let text = ["Fota", "轻松管理"];


	let baseList = [{
			icon: "/static/index/equipment.png",
			path: "/pages/monitor_project/index",
			title: "项目监控",
		},
		{
			icon: "/static/index/equipment.png",
			path: "/pages/monitor_client/index",
			title: "项目用户",
		},


		{
			icon: "scan",
			// path: "/pages/equipment/equipment",
			title: "扫码维护",
			fn: function() {
				// console.log('this', this);
				scanCode().then((path) => {
					// toPath(path);
					if (!path) {
						return
					}
					var query = queryToObj(path);
					if (!query.state) {
						toast("请扫描设备二维码");
						return;
					}
					toPath("/pages/equipment/details", {
						number: query.state,
					});
				});
			},
		},


	];

	function toList(item) {
		if (item.fn) {
			item.fn();
		} else {
			toPath(item.path);
		}
	}
</script>

<style>
	.grid-text {
		font-size: 28rpx;
	}
</style>