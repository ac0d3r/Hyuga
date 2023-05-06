<script lang="ts" setup>
import { nextTick, Ref, ref } from "vue";
import { ElInput, ElMessage } from 'element-plus';
import { User, UserFilled, Key, Link, Tools, CopyDocument, View, Hide, Refresh, InfoFilled, SuccessFilled } from '@element-plus/icons-vue';
import useClipboard from 'vue-clipboard3';
import { useStore } from "./lib/store";
import { resetToken, setNotify, setRebinding } from "./lib/user";
import { validateIP } from "./lib/utils";

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
// setting
const opened = ref(false);
const reset = () => {
    resetToken(
        () => {
            if (!hides) {
                token.value = store.state.user.token;
            }
        }
    );
};
const confirm = () => {
    setNotify(store.state.user.notify);
};
// dns rebinding
const inputVisible = ref(false);
const inputValue = ref('');
const InputRef = ref<InstanceType<typeof ElInput>>();
const handleDNSClose = (dns: string) => {
    store.state.user.rebinding.splice(store.state.user.rebinding.indexOf(dns), 1)
    setRebinding();
}
const handleDNSConfirm = () => {
    if (store.state.user.rebinding == null)
        store.state.user.rebinding = [];

    if (inputValue.value) {
        if (validateIP(inputValue.value)) {
            store.state.user.rebinding.push(inputValue.value);
            setRebinding();
        } else {
            ElMessage({ message: 'invalid ip address', type: 'warning', });
        }
    }
    inputVisible.value = false;
    inputValue.value = '';
}
const showInput = () => {
    inputVisible.value = true
    nextTick(() => {
        InputRef.value!.input!.focus()
    })
}
// records
const records: Ref<any[]> = ref([]);
const ws = new WebSocket(`${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${window.location.host}/api/v2/user/record`)
ws.onmessage = (msg: any) => {
    records.value.push(JSON.parse(msg.data));
}
const typeTag = (type: number) => {
    if (type == 0) {
        return '';
    } else if (type == 1) {
        return 'success';
    } else if (type == 2) {
        return 'warning';
    } else if (type == 3) {
        return 'danger';
    }
    return 'info';
}
const strTag = (type: number) => {
    if (type == 0) {
        return 'dns';
    } else if (type == 1) {
        return 'http';
    } else if (type == 2) {
        return 'ldap';
    } else if (type == 3) {
        return 'rmi';
    }
    return 'unknown';
}
const parseDetail = (detail: any): any => {
    if (detail !== null && detail.raw !== null) {
        let lines = [];
        for (let line of detail.raw.split('\r\n')) {
            if (line !== '')
                lines.push(line);
        }
        return lines;
    }
    return null;
}
const parseTime = (time: number): string => {
    const date = new Date(time * 1000);
    return date.toLocaleString();
}
</script>

<template>
    <div style="text-align: center;">
        <el-space direction="vertical">
            <el-card shadow="always" style="width: 500px;">
                <template #header>
                    <div class="card-header">
                        <el-text>
                            <el-icon>
                                <UserFilled />
                            </el-icon>
                            Profile
                        </el-text>
                        <el-button :icon="Tools" text @click="store.state.logged ? opened = true : opened = false" />
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
                <el-table-column prop="type" label="Type">
                    <template #default="scope">
                        <el-tag :type="typeTag(scope.row.type)" disable-transitions>{{ strTag(scope.row.type)
                        }}</el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="name" label="Name" />
                <el-table-column prop="remote_addr" label="RemoteAddr" />
                <el-table-column prop="created_at" label="CreatedAt">
                    <template #default="scope">
                        {{ parseTime(scope.row.created_at) }}
                    </template>
                </el-table-column>
                <el-table-column type="expand">
                    <template #default="props">
                        <div m="4">
                            <h4>Detail</h4>
                            <el-row>
                                <el-col v-for="line in parseDetail(props.row.detail)" :span="16">
                                    <el-text size="small">{{ line }}</el-text>
                                </el-col>
                            </el-row>
                        </div>
                    </template>
                </el-table-column>
            </el-table>
        </el-space>
    </div>

    <!-- settting  -->
    <el-drawer v-model="opened" title="Settings" direction="rtl">
        <el-form size="small">
            <el-divider content-position="left" style="margin-top: 0px;">Token Setting</el-divider>
            <el-form-item label="Token">
                <el-popconfirm width="300" title="Are you sure to reset the API token?" confirm-button-text="Yes"
                    cancel-button-text="No" :icon="InfoFilled" icon-color="#626AEF" @confirm="reset">
                    <template #reference>
                        <el-button type="primary" :icon="Refresh">Reset</el-button>
                    </template>
                </el-popconfirm>
            </el-form-item>
            <el-divider content-position="left">DNS Rebinding</el-divider>
            <el-form-item label="Domain">
                <el-tag>
                    {{ store.state.user.data.rdomain }}
                </el-tag>
                <el-button text size="small" :icon="CopyDocument" circle @click="copy(store.state.user.data.rdomain)" />
            </el-form-item>
            <el-form-item label="DNS">
                <el-tag v-for="d in store.state.user.rebinding" closable @close="handleDNSClose(d)"
                    style="margin-right: 5px;">
                    {{ d }}
                </el-tag>
                <el-input v-if="inputVisible" ref="InputRef" v-model="inputValue" size="small"
                    @keyup.enter="handleDNSConfirm" @blur="handleDNSConfirm" />
                <el-button v-else size="small" @click="showInput">
                    + New DNS
                </el-button>
            </el-form-item>
            <el-divider content-position="left">Notify Setting</el-divider>
            <el-form-item label="Enable">
                <el-switch v-model="store.state.user.notify.enable"
                    style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949" />
            </el-form-item>

            <el-text>Bark</el-text>
            <el-form-item label="Key" style="margin-left:12px;">
                <el-input v-model="store.state.user.notify.bark.key" />
            </el-form-item>
            <el-form-item label="Server" style="margin-left:12px;">
                <el-input v-model="store.state.user.notify.bark.server" />
            </el-form-item>

            <el-text>DingTalk</el-text>
            <el-form-item label="Token" style="margin-left:12px;">
                <el-input v-model="store.state.user.notify.dingtalk.token" />
            </el-form-item>
            <el-form-item label="Secret" style="margin-left:12px;">
                <el-input v-model="store.state.user.notify.dingtalk.secret" />
            </el-form-item>

            <el-text>Lark</el-text>
            <el-form-item label="Token" style="margin-left:12px;">
                <el-input v-model="store.state.user.notify.lark.token" />
            </el-form-item>
            <el-form-item label="Secret" style="margin-left: 12px;">
                <el-input v-model="store.state.user.notify.lark.secret" />
            </el-form-item>

            <el-text>Feishu</el-text>
            <el-form-item label="Token" style="margin-left:12px;">
                <el-input v-model="store.state.user.notify.feishu.token" />
            </el-form-item>
            <el-form-item label="Secret" style="margin-left:12px;">
                <el-input v-model="store.state.user.notify.feishu.secret" />
            </el-form-item>

            <el-text>ServerChan</el-text>
            <el-form-item label="UserID" style="margin-left:12px;">
                <el-input v-model="store.state.user.notify.serverchan.user_id" />
            </el-form-item>
            <el-form-item label="SendKey" style="margin-left:12px;">
                <el-input v-model="store.state.user.notify.serverchan.send_key" />
            </el-form-item>
            <el-form-item style="float: right;">
                <el-button type="primary" :icon="SuccessFilled" @click="confirm">Confrim</el-button>
            </el-form-item>
        </el-form>
    </el-drawer>
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
