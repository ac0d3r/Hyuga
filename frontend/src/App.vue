<script lang="ts" setup>
import { Position } from '@element-plus/icons-vue';
import Login from './Login.vue';
import Home from './Home.vue';
</script>

<template>
    <el-container>
        <el-header>
            <el-menu mode="horizontal" :ellipsis="false">
                <el-menu-item>
                    <img src="./assets/logo.png" alt="logo" style="width: 100px;">
                </el-menu-item>
                <div class="flex-grow" />
                <Login />
            </el-menu>
        </el-header>

        <el-main>
            <Home />
        </el-main>

        <el-footer class="footer">
            <el-space direction="vertical">
                <el-text size="small">
                    <el-link href="https://github.com/ac0d3r/Hyuga" type="primary" :underline="false">
                        <img src="./assets/logo.png" alt="logo" style="width: 50px;padding: 0px;">
                    </el-link>
                </el-text>
                <el-text size="small">
                    © 巴斯.zznQ 2019-2023
                    <el-link :icon="Position" href="https://blog.macoder.tech" :underline="false"></el-link>
                </el-text>
            </el-space>
        </el-footer>
    </el-container>
</template>

<style>
.flex-grow {
    flex-grow: 1;
}

.footer {
    text-align: center;
}
</style>


