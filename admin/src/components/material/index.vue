<template>
    <!-- v-loading="pager.loading" -->
    <div class="material" v-loading="pager.loading">
        <div class="material__left">
            <div class="flex-1 min-h-0">
                <el-scrollbar>
                    <div class="material-left__content pt-4 p-b-4">
                        <el-tree
                            ref="treeRef"
                            node-key="id"
                            :data="cateLists"
                            empty-text=""
                            :highlight-current="true"
                            :expand-on-click-node="false"
                            :current-node-key="cateId"
                            @node-click="handleCatSelect"
                        >
                            <template v-slot="{ data }">
                                <div class="flex flex-1 items-center min-w-0 pr-4">
                                    <img
                                        class="w-[20px] h-[16px] mr-3"
                                        src="@/assets/images/icon_folder.png"
                                    />
                                    <span class="flex-1 truncate mr-2">
                                        <overflow-tooltip :content="data.name" />
                                    </span>
                                    <el-dropdown
                                        v-perms="[
                                            'admin:common:album:cateRename',
                                            'admin:common:album:cateDel'
                                        ]"
                                        v-if="data.id > 0"
                                        :hide-on-click="false"
                                    >
                                        <span class="muted m-r-10">···</span>
                                        <template #dropdown>
                                            <el-dropdown-menu>
                                                <popover-input
                                                    v-perms="['admin:common:album:cateRename']"
                                                    @confirm="handleEditCate($event, data.id)"
                                                    size="default"
                                                    :value="data.name"
                                                    width="400px"
                                                    :limit="20"
                                                    show-limit
                                                    teleported
                                                >
                                                    <div>
                                                        <el-dropdown-item>
                                                            命名分组
                                                        </el-dropdown-item>
                                                    </div>
                                                </popover-input>
                                                <div
                                                    v-perms="['admin:common:album:cateDel']"
                                                    @click="handleDeleteCate(data.id)"
                                                >
                                                    <el-dropdown-item>删除分组</el-dropdown-item>
                                                </div>
                                            </el-dropdown-menu>
                                        </template>
                                    </el-dropdown>
                                </div>
                            </template>
                        </el-tree>
                    </div>
                </el-scrollbar>
            </div>

            <div class="flex justify-center p-2 border-t border-br">
                <popover-input
                    v-perms="['admin:common:album:cateAdd']"
                    @confirm="handleAddCate"
                    size="default"
                    width="400px"
                    :limit="20"
                    show-limit
                    teleported
                >
                    <el-button> 添加分组 </el-button>
                </popover-input>
            </div>
        </div>
        <div class="material__center flex flex-col">
            <div class="operate-btn flex">
                <div class="flex-1 flex">
                    <upload
                        v-if="type == 'image'"
                        v-perms="['admin:common:upload:image']"
                        class="mr-3"
                        :data="{ cid: cateId }"
                        :type="type"
                        :show-progress="true"
                        @change="refresh"
                    >
                        <el-button type="primary">本地上传</el-button>
                    </upload>
                    <upload
                        v-if="type == 'video'"
                        v-perms="['admin:common:upload:video']"
                        class="mr-3"
                        :data="{ cid: cateId }"
                        :type="type"
                        :show-progress="true"
                        @change="refresh"
                    >
                        <el-button type="primary">本地上传</el-button>
                    </upload>
                </div>
                <el-input
                    class="w-60"
                    placeholder="请输入名称"
                    v-model="fileParams.name"
                    @keyup.enter="refresh"
                >
                    <template #append>
                        <el-button @click="refresh">
                            <template #icon>
                                <icon name="el-icon-Search" />
                            </template>
                        </el-button>
                    </template>
                </el-input>
            </div>

            <div class="material-center__content flex flex-col flex-1 mb-1 min-h-0">
                <el-table
                    ref="tableRef"
                    class="mt-4"
                    :data="pager.lists"
                    width="100%"
                    height="100%"
                    size="large"
                    @selection-change="handleSelectionChange"
                >
                    <el-table-column type="selection" width="55" />
                    <el-table-column label="图片" width="100">
                        <template #default="{ row }">
                            <file-item
                                :uri="row.uri"
                                file-size="50px"
                                :type="type"
                                @click.stop="handlePreview(row.uri)"
                            ></file-item>
                        </template>
                    </el-table-column>
                    <el-table-column label="名称" min-width="100" show-overflow-tooltip>
                        <template #default="{ row }">
                            <el-link @click.stop="handlePreview(row.uri)" :underline="false">
                                {{ row.name }}
                            </el-link>
                        </template>
                    </el-table-column>
                    <el-table-column prop="createTime" label="上传时间" min-width="100" />
                    <el-table-column label="操作" width="150" fixed="right">
                        <template #default="{ row }">
                            <div class="inline-block" v-perms="['admin:common:album:albumRename']">
                                <popover-input
                                    @confirm="handleFileRename($event, row.id)"
                                    size="default"
                                    :value="row.name"
                                    width="400px"
                                    :limit="50"
                                    show-limit
                                    teleported
                                >
                                    <el-button type="primary" link> 重命名 </el-button>
                                </popover-input>
                            </div>
                            <div class="inline-block">
                                <el-button type="primary" link @click.stop="handlePreview(row.uri)">
                                    查看
                                </el-button>
                            </div>
                            <div class="inline-block" v-perms="['admin:common:album:albumDel']">
                                <el-button
                                    type="primary"
                                    link
                                    @click.stop="batchFileDelete([row.id])"
                                >
                                    删除
                                </el-button>
                            </div>
                        </template>
                    </el-table-column>
                </el-table>

                <div
                    class="flex flex-1 justify-center items-center"
                    v-if="!pager.loading && !pager.lists.length"
                >
                    暂无数据~
                </div>
            </div>
            <div class="material-center__footer flex justify-between items-center mt-2">
                <div class="flex">
                    <template v-if="mode == 'page'">
                        <!-- <span class="mr-3">
                            <el-checkbox
                                :disabled="!pager.lists.length"
                                v-model="isCheckAll"
                                @change="selectAll"
                                :indeterminate="isIndeterminate"
                            >
                                当页全选
                            </el-checkbox>
                        </span> -->
                        <el-button
                            v-perms="['admin:common:album:albumDel']"
                            :disabled="!select.length"
                            @click="batchFileDelete()"
                        >
                            删除
                        </el-button>
                        <popup
                            v-perms="['admin:common:album:albumMove']"
                            class="ml-3 inline"
                            @confirm="batchFileMove"
                            :disabled="!select.length"
                            title="移动文件"
                        >
                            <template #trigger>
                                <el-button :disabled="!select.length">移动</el-button>
                            </template>

                            <el-select v-model="moveId" placeholder="请选择">
                                <template v-for="item in cateLists" :key="item.id">
                                    <el-option
                                        v-if="item.id !== ''"
                                        :label="item.name"
                                        :value="item.id"
                                    ></el-option>
                                </template>
                            </el-select>
                        </popup>
                    </template>
                </div>
                <pagination
                    v-model="pager"
                    @change="getFileList"
                    layout="total, prev, pager, next, jumper"
                />
            </div>
        </div>
        <div class="material__right" v-if="mode == 'picker'">
            <div class="flex justify-between p-2 flex-wrap">
                <div class="sm flex items-center">
                    已选择 {{ select.length }}
                    <span v-if="limit">/{{ limit }}</span>
                </div>
                <el-button type="primary" link @click="clearSelect">清空</el-button>
            </div>
            <div class="flex-1 min-h-0">
                <el-scrollbar class="ls-scrollbar">
                    <ul class="select-lists flex flex-col p-t-3">
                        <li class="mb-4" v-for="item in select" :key="item.id">
                            <div class="select-item">
                                <del-wrap @close="cancelSelete(item.id)">
                                    <file-item
                                        :uri="item.uri"
                                        file-size="100px"
                                        :type="type"
                                    ></file-item>
                                </del-wrap>
                            </div>
                        </li>
                    </ul>
                </el-scrollbar>
            </div>
        </div>
        <preview v-model="showPreview" :url="previewUrl" :type="type" />
    </div>
</template>

<script lang="ts" setup>
import { useCate, useFile } from './hook'
import FileItem from './file.vue'
import Preview from './preview.vue'
const props = defineProps({
    fileSize: {
        type: String,
        default: '100px'
    },
    limit: {
        type: Number,
        default: 1
    },
    type: {
        type: String,
        default: 'image'
    },
    mode: {
        type: String,
        default: 'picker'
    },
    pageSize: {
        type: Number,
        default: 15
    }
})
const emit = defineEmits(['change'])
const { limit } = toRefs(props)
const typeValue = computed<number>(() => {
    switch (props.type) {
        case 'image':
            return 10
        case 'video':
            return 20
        case 'file':
            return 30
        default:
            return 0
    }
})
// const visible: Ref<boolean> = inject('visible')
const previewUrl = ref('')
const showPreview = ref(false)
const {
    treeRef,
    cateId,
    cateLists,
    handleAddCate,
    handleEditCate,
    handleDeleteCate,
    getCateLists,
    handleCatSelect
} = useCate(typeValue.value)

const {
    tableRef,
    moveId,
    pager,
    fileParams,
    select,
    isCheckAll,
    isIndeterminate,
    getFileList,
    refresh,
    batchFileDelete,
    batchFileMove,

    clearSelect,
    cancelSelete,

    selectItems,
    handleFileRename
} = useFile(cateId, typeValue, limit, props.pageSize)
function handleSelectionChange(val: any[]) {
    console.log('handleSelectionChange', val)
    selectItems(val)
    // multipleSelection.value = val
}
const getData = async () => {
    await getCateLists()
    treeRef.value?.setCurrentKey(cateId.value)
    getFileList()
}
// getData()
const handlePreview = (url: string) => {
    previewUrl.value = url
    showPreview.value = true
}
// watch(
//     () => visible,
//     (val) => {
//         if (val) {
//             getData()
//         }
//     },
//     {
//         immediate: true
//     }
// )
watch(cateId, () => {
    fileParams.name = ''
    refresh()
})

watch(
    select,
    (val: any[]) => {
        emit('change', val)
        if (val.length == pager.lists.length && val.length !== 0) {
            isIndeterminate.value = false
            isCheckAll.value = true
            return
        }
        if (val.length > 0) {
            isIndeterminate.value = true
        } else {
            isCheckAll.value = false
            isIndeterminate.value = false
        }
    },
    {
        deep: true
    }
)

onMounted(() => {
    props.mode == 'page' && getData()
})

defineExpose({
    clearSelect
})
</script>

<style scoped lang="scss">
.material {
    height: 100%;
    min-height: 0px;
    display: flex;
    flex: 1 1 0%;
    &__left {
        border-right-width: 1px;
        border-color: var(--el-border-color);
        display: flex;
        flex-direction: column;
        width: 200px;
        :deep(.el-tree-node__content) {
            height: 36px;
        }
    }
    &__center {
        flex: 1;
        min-width: 0;
        min-height: 0;
        padding: 16px 16px 0;
        .list-icon {
            border-radius: 3px;
            display: flex;
            padding: 5px;
            cursor: pointer;
            &.select {
                color: var(--el-color-primary);
                background-color: var(--el-color-primary-light-8);
            }
        }
        .file-list {
            .file-item-wrap {
                margin-right: 16px;
                line-height: 1.3;
                cursor: pointer;
                .item-selected {
                    display: flex;
                    align-items: center;
                    justify-content: center;
                    position: absolute;
                    top: 0;
                    left: 0;
                    width: 100%;
                    height: 100%;
                    border-radius: 4px;
                    background-color: rgba(0, 0, 0, 0.5);
                    box-sizing: border-box;
                }
                .operation-btns {
                    height: 28px;
                    visibility: hidden;
                }
                &:hover .operation-btns {
                    visibility: visible;
                }
            }
        }
    }
    &__right {
        border-left-width: 1px;
        border-color: var(--el-border-color);
        display: flex;
        flex-direction: column;
        width: 130px;
        .select-lists {
            padding: 10px;

            .select-item {
                width: 100px;
                height: 100px;
            }
        }
    }
}
</style>
