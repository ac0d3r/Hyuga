import { apihost } from "./conf";
import { getCookie } from "./cookie";

function SetTokenHeader(obj) {
    obj.Authorization = "Bearer " + getCookie("token")
    return obj
}

function CreateUser(succ, fail) {
    fetch(`${apihost}/api/user/create`, { method: "POST" })
        .then((res) => res.json())
        .then((res) => {
            const { code, msg, data } = res;
            if (code === 0) {
                succ(data);
            } else {
                fail(msg);
            }
        })
        .catch((err) => {
            fail(err.message);
        });
}

function DeleteUser(succ, fail) {
    fetch(`${apihost}/api/user/delete`, { method: "POST", headers: SetTokenHeader({}) })
        .then((res) => res.json())
        .then((res) => {
            const { code, msg, data } = res;
            if (code === 0) {
                succ(data);
            } else {
                fail(msg);
            }
        })
        .catch((err) => {
            fail(err.message);
        });
}

function GetUserDnsRebindingHosts(succ, fail) {
    fetch(`${apihost}/api/user/dns-rebinding`, { headers: SetTokenHeader({}) })
        .then((res) => res.json())
        .then((res) => {
            const { code, msg, data } = res;
            if (code === 0) {
                succ(data);
            } else {
                fail(msg, code);
            }
        })
        .catch((err) => {
            fail(err.message);
        });
}

function UpdateUserDnsRebindingHosts(ips, succ, fail) {
    const opts = {
        method: "POST",
        headers: SetTokenHeader({ 'Content-Type': 'application/json' }),
        body: JSON.stringify({ ip: ips })
    }
    fetch(`${apihost}/api/user/dns-rebinding`, opts)
        .then((res) => res.json())
        .then((res) => {
            const { code, msg, data } = res;
            if (code === 0) {
                succ(data);
            } else {
                fail(msg);
            }
        })
        .catch((err) => {
            fail(err.message);
        });
}

function GetLogRecords(type, filter, succ, fail) {
    const opts = {
        headers: SetTokenHeader({}),
    }
    fetch(`${apihost}/api/record/list?type=${type}&filter=${filter}`, opts)
        .then((res) => res.json())
        .then((res) => {
            const { code, msg, data } = res;
            if (code === 0) {
                succ(data);
            } else {
                fail(msg);
            }
        })
        .catch((err) => {
            fail(err.message);
        });
}

function WipeRecodsData(succ, fail) {
    const opts = {
        method: "POST",
        headers: SetTokenHeader({}),
    }
    fetch(`${apihost}/api/record/clean`, opts)
        .then((res) => res.json())
        .then((res) => {
            const { code, msg, data } = res;
            if (code === 0) {
                succ(data);
            } else {
                fail(msg);
            }
        })
        .catch((err) => {
            fail(err.message);
        });
}

export {
    CreateUser,
    DeleteUser,
    GetUserDnsRebindingHosts,
    UpdateUserDnsRebindingHosts,
    GetLogRecords,
    WipeRecodsData
};