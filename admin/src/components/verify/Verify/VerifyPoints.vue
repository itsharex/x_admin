<template>
    <div style="position: relative">
        <div class="verify-img-out">
            <div
                class="verify-img-panel"
                :style="{
                    width: imgSize.width + 'px',
                    height: imgSize.height + 'px',

                    'margin-bottom': vSpace + 'px'
                }"
            >
                <div
                    class="verify-refresh"
                    style="z-index: 3"
                    @click="refresh"
                    v-show="showRefresh"
                >
                    <i class="iconfont icon-refresh"></i>
                </div>
                <img
                    v-if="pointBackImgBase"
                    :src="pointBackImgBase"
                    ref="canvas"
                    alt=""
                    style="width: 100%; height: 100%; display: block"
                    @click="bindingClick ? canvasClick($event) : undefined"
                />

                <div
                    v-for="(tempPoint, index) in tempPoints"
                    :key="index"
                    class="point-area"
                    :style="{
                        'background-color': '#1abd6c',
                        color: '#fff',
                        'z-index': 9999,
                        width: '20px',
                        height: '20px',
                        'text-align': 'center',
                        'line-height': '20px',
                        'border-radius': '50%',
                        position: 'absolute',
                        top: tempPoint.y - 10 + 'px',
                        left: tempPoint.x - 10 + 'px'
                    }"
                >
                    {{ index + 1 }}
                </div>
            </div>
        </div>

        <div
            class="verify-bar-area"
            :style="{
                width: imgSize.width + 'px',
                color: barAreaColor,
                'border-color': barAreaBorderColor,
                'line-height': '40px'
            }"
        >
            <span class="verify-msg">{{ text }}</span>
        </div>
    </div>
</template>
<script>
/**
 * VerifyPoints
 * @description 点选
 * */
// import { resetSize } from '../utils/util'
import { aesEncrypt } from '../utils/ase'
import { reqGet, reqCheck } from '../api/index'
import { onMounted, reactive, ref, toRefs, getCurrentInstance } from 'vue'
export default {
    name: 'VerifyPoints',
    emits: ['success', 'error'],
    props: {
        //弹出式pop，固定fixed
        mode: {
            type: String,
            default: 'fixed'
        },
        captchaType: {
            type: String
        },
        //间隔
        vSpace: {
            type: Number,
            default: 5
        },
        imgSize: {
            type: Object,
            default() {
                return {
                    width: 310,
                    height: 155
                }
            }
        }
    },
    setup(props, { emit }) {
        const { mode, captchaType } = toRefs(props)
        const { proxy } = getCurrentInstance()
        const secretKey = ref(''), //后端返回的ase加密秘钥
            checkNum = ref(3), //默认需要点击的字数
            fontPos = reactive([]), //选中的坐标信息
            checkPosArr = reactive([]), //用户点击的坐标
            num = ref(1), //点击的记数
            pointBackImgBase = ref(''), //后端获取到的背景图片
            pointTextList = ref([]), //后端返回的点击字体顺序
            backToken = ref(''), //后端返回的token值
            tempPoints = reactive([]),
            text = ref(''),
            barAreaColor = ref(undefined),
            barAreaBorderColor = ref(undefined),
            showRefresh = ref(true),
            bindingClick = ref(true)

        const init = () => {
            //加载页面
            fontPos.splice(0, fontPos.length)
            checkPosArr.splice(0, checkPosArr.length)
            num.value = 1
            getPicture()
        }
        onMounted(() => {
            // 禁止拖拽
            init()
            proxy.$el.onselectstart = function () {
                return false
            }
        })
        const canvas = ref(null)
        const canvasClick = (e) => {
            checkPosArr.push(getMousePos(canvas, e))
            if (num.value == checkNum.value) {
                num.value = createPoint(getMousePos(canvas, e))
                //按比例转换坐标值
                const arr = pointTransform(checkPosArr)
                checkPosArr.length = 0
                checkPosArr.push(...arr)
                //等创建坐标执行完
                setTimeout(() => {
                    const data = {
                        captchaType: captchaType.value,
                        pointJson: secretKey.value
                            ? aesEncrypt(JSON.stringify(checkPosArr), secretKey.value)
                            : JSON.stringify(checkPosArr),
                        token: backToken.value
                    }
                    reqCheck(data).then((res) => {
                        if (res.repCode == '0000') {
                            barAreaColor.value = '#4cae4c'
                            barAreaBorderColor.value = '#5cb85c'
                            text.value = '验证成功'
                            bindingClick.value = false
                            if (mode.value == 'pop') {
                                setTimeout(() => {
                                    refresh()
                                }, 1500)
                            }
                            emit('success', { ...data })
                        } else {
                            emit('error')
                            barAreaColor.value = '#d9534f'
                            barAreaBorderColor.value = '#d9534f'
                            text.value = '验证失败'
                            setTimeout(() => {
                                refresh()
                            }, 700)
                        }
                    })
                }, 400)
            }
            if (num.value < checkNum.value) {
                num.value = createPoint(getMousePos(canvas, e))
            }
        }
        //获取坐标
        const getMousePos = function (obj, e) {
            const x = e.offsetX
            const y = e.offsetY
            return { x, y }
        }
        //创建坐标点
        const createPoint = function (pos) {
            tempPoints.push(Object.assign({}, pos))
            return num.value + 1
        }
        const refresh = function () {
            tempPoints.splice(0, tempPoints.length)
            barAreaColor.value = '#000'
            barAreaBorderColor.value = '#ddd'
            bindingClick.value = true
            fontPos.splice(0, fontPos.length)
            checkPosArr.splice(0, checkPosArr.length)
            num.value = 1
            getPicture()
            text.value = '获取中...'
            showRefresh.value = true
        }

        // 请求背景图片和验证图片
        function getPicture() {
            const data = {
                captchaType: captchaType.value
            }
            reqGet(data).then((res) => {
                if (res.repCode == '0000') {
                    pointBackImgBase.value =
                        'data:image/png;base64,' + res.repData.originalImageBase64
                    backToken.value = res.repData.token
                    secretKey.value = res.repData.secretKey
                    pointTextList.value = res.repData.wordList
                    text.value = '请依次点击【' + pointTextList.value.join(',') + '】'
                } else {
                    text.value = res.repMsg
                }
            })
        }
        //坐标转换函数
        const pointTransform = function (pointArr) {
            const newPointArr = pointArr.map((p) => {
                const x = Math.round((310 * p.x) / parseInt(props.imgSize.width))
                const y = Math.round((155 * p.y) / parseInt(props.imgSize.height))
                return { x, y }
            })
            return newPointArr
        }
        return {
            secretKey,
            checkNum,
            fontPos,
            checkPosArr,
            num,
            pointBackImgBase,
            pointTextList,
            backToken,

            tempPoints,
            text,
            barAreaColor,
            barAreaBorderColor,
            showRefresh,
            bindingClick,
            init,
            canvas,
            canvasClick,
            getMousePos,
            createPoint,
            refresh,
            getPicture,
            pointTransform
        }
    }
}
</script>
