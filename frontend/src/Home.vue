<script lang="ts" setup>
import { Ref, ref } from "vue";
import { User, UserFilled, Key, Link, Tools, CopyDocument, View, Hide } from '@element-plus/icons-vue';
import useClipboard from 'vue-clipboard3';
import { useStore } from "./lib/store";

const store = useStore();
// user info
const token = ref('***************');
const ticon = ref(View);
let hides = true;

const viewAction = () => {
    if (hides) {
        token.value = store.state.user.token;
        ticon.value = Hide;
        hides = false;
    } else {
        token.value = '***************';
        ticon.value = View;
        hides = true;
    }
}
// copy
const { toClipboard } = useClipboard();
const copy = async (data: string) => {
    try {
        await toClipboard(data);
    } catch (e) {
    }
}

const records: Ref<any[]> = ref([]);
</script>

<template>
    <div style="text-align: center; ">
        <el-space direction="vertical">
            <el-card shadow="always" style="width: 500px;">
                <template #header>
                    <div class="card-header">
                        <el-text>
                            <el-icon>
                                <UserFilled />
                            </el-icon>
                        </el-text>
                        <el-button :icon="Tools" text />
                    </div>
                </template>
                <el-descriptions v-if="store.state.logged" :column="1" size="small" direction="horizontal">
                    <el-descriptions-item>
                        <el-text>
                            <el-icon>
                                <User />
                            </el-icon>
                            <el-tag style="margin-left: 5px;margin-right: 8px;">
                                {{ store.state.user.data.subdomain }}
                            </el-tag>
                            <el-button text size="small" :icon="CopyDocument" circle
                                @click="copy(store.state.user.data.subdomain)" />
                        </el-text>
                    </el-descriptions-item>
                    <el-descriptions-item>
                        <el-text>
                            <el-icon>
                                <Key />
                            </el-icon>
                            <el-tag style="margin-left: 5px;margin-right: 8px;">
                                {{ token }}
                            </el-tag>
                            <el-button text size="small" :icon="ticon" circle @click="viewAction" />
                            <el-button text size="small" :icon="CopyDocument" circle
                                @click="copy(store.state.user.token)" />
                        </el-text>
                    </el-descriptions-item>
                    <el-descriptions-item>
                        <el-text>
                            <el-icon>
                                <Link />
                            </el-icon>
                            <el-tag style="margin-left: 5px;margin-right: 8px;">
                                {{ store.state.user.data.ldap }}
                            </el-tag>
                            <el-button text size="small" :icon="CopyDocument" circle
                                @click="copy(store.state.user.data.ldap)" />
                        </el-text>
                    </el-descriptions-item>
                    <el-descriptions-item>
                        <el-text>
                            <el-icon>
                                <Link />
                            </el-icon>
                            <el-tag style="margin-left: 5px;margin-right: 8px;">
                                {{ store.state.user.data.rmi }}
                            </el-tag>
                            <el-button text size="small" :icon="CopyDocument" circle
                                @click="copy(store.state.user.data.rmi)" />
                        </el-text>
                    </el-descriptions-item>
                </el-descriptions>
            </el-card>

            <el-empty v-if="records.length == 0" description="No Data">
            </el-empty>
            <el-table v-else :data="records" height="550" stripe table-layout="fixed" max-height="550"
                style="width: 1000px;">
                <el-table-column prop="type" label="Type" />
                <el-table-column prop="name" label="Name" />
                <el-table-column prop="remote_addr" label="RemoteAddr" />
                <el-table-column prop="created_at" label="CreatedAt" />
                <el-table-column type="expand">
                    <template #default="props">
                        <div m="4">
                            <h3>Detail</h3>
                            <p m="t-0 b-2">{{ props.row.detail }}</p>
                        </div>
                    </template>
                </el-table-column>
            </el-table>
        </el-space>
    </div>
</template>

<style>
.card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.cell-item {
    display: flex;
    align-items: center;
}
</style>
