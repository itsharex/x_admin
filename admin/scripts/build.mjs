import { execaCommand } from 'execa'
import path from 'path'

import fsExtra from 'fs-extra'
const { existsSync, remove, copy } = fsExtra
const cwd = process.cwd()
//打包发布路径，谨慎改动
const releaseRelativePath = '../frontend'
const distPath = path.resolve(cwd, 'dist')
const releasePath = path.resolve(cwd, releaseRelativePath)

const goProjectPath = path.resolve(cwd, '../server')
async function build_go() {
    await execaCommand('swag fmt', {
        stdio: 'inherit',
        encoding: 'utf8',
        cwd: goProjectPath
    })
    await execaCommand('swag init', {
        stdio: 'inherit',
        encoding: 'utf8',
        cwd: goProjectPath
    })
    await execaCommand('goreleaser release --snapshot --clean', {
        stdio: 'inherit',
        encoding: 'utf8',
        cwd: goProjectPath
    })
}
async function build_frontend() {
    await execaCommand('npm run prod', { stdio: 'inherit', encoding: 'utf8', cwd })
    if (existsSync(releasePath)) {
        await remove(releasePath)
    }
    console.log(`文件正在复制 ==> ${releaseRelativePath}`)
    try {
        await copyFile(distPath, releasePath)
    } catch (error) {
        console.log(error)
    }
    console.log(`文件已复制 ==> ${releaseRelativePath}`)
}
async function build() {
    try {
        console.log('开始打包')
        await Promise.allSettled([build_go(), build_frontend()])
    } catch (error) {
        console.log(error)
    }
}
function copyFile(sourceDir, targetDir) {
    return new Promise((resolve, reject) => {
        copy(sourceDir, targetDir, (err) => {
            if (err) {
                reject(err)
            } else {
                resolve()
            }
        })
    })
}

build()
