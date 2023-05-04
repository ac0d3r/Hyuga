import { useCookies } from "vue3-cookies";
import { ElMessage } from 'element-plus';
import { store } from "./store";

const { cookies } = useCookies();

// CUSTOM: 自定义部署，需修改 github client_id
const client_id = "";
const githubOauth = `https://github.com/login/oauth/authorize?scope=user:email&client_id=${client_id}`;


// message
const succ = (msg: string) => {
    ElMessage({
        message: msg,
        type: 'success',
    });
};

const warn = (msg: string) => {
    ElMessage({
        message: msg,
        type: 'warning',
    });
};

const fail = (msg: string) => {
    ElMessage({
        message: msg,
        type: 'error',
    });
};

const apihost: string = function () {
    return `${window.location.protocol}//${window.location.host}`;
}();



function isLogin(): boolean {
    return cookies.get("sid") !== null && cookies.get("token") !== null;
}

function logout() {
    fetch(`${apihost}/api/v2/user/logout`, { method: "POST" })
        .then((res) => res.json())
        .then((res) => {
            const { code, msg, _ } = res;
            if (code === 0) {
                // 清除 cookie
                cookies.remove("sid");
                cookies.remove("token");
                store.state.logged = false;
                succ('logout success');
            } else {
                warn(msg);
            }
        }).catch((err) => {
            fail(err.message);
        });
}

function getUserInfo(succcallback: Function) {
    fetch(`${apihost}/api/v2/user/info`, { method: "GET" })
        .then((res) => res.json())
        .then((res) => {
            const { code, msg, data } = res;
            if (code !== 0) {
                warn(msg);
            } else {
                // update user info
                store.state.user = data;
                succcallback();
            }
        }).catch((err) => {
            fail(err.message);
        });
}

function setNotify(notify: any) {
    fetch(`${apihost}/api/v2/user/notify`, {
        method: "POST",
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(notify),
    })
        .then((res) => res.json())
        .then((res) => {
            const { code, msg, _ } = res;
            if (code === 0) {
                getUserInfo(() => { });
                succ('set notify success');
            } else {
                warn(msg);
            }
        }).catch((err) => {
            fail(err.message);
        });
}

function resetToken(succcallback: Function) {
    fetch(`${apihost}/api/v2/user/reset`, { method: "POST" })
        .then((res) => res.json())
        .then((res) => {
            const { code, msg, data } = res;
            if (code === 0) {
                getUserInfo(succcallback);
                succ('reset token success');
            } else {
                warn(msg);
            }
        }).catch((err) => {
            fail(err.message);
        });
}


export {
    isLogin,
    githubOauth,
    logout,
    getUserInfo,
    setNotify,
    resetToken
};
