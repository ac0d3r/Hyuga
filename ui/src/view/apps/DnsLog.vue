<template>
  <div class="container-app-dns-query">
    <div class="container-image">
      <img :src="LogoImg" alt="hyuga" />
    </div>

    <hr class="container-divider" />

    <div class="container-user-area">
      <AButton type="primary" @click="handleGetSubDomain">
        Get SubDomain
      </AButton>
      <div style="min-height: 32px">
        <div
          v-if="userinfos.identity !== '' && userinfos.token !== ''"
          class="container-token-information"
        >
          <div>Domain: {{ userinfos.identity }}</div>
          <div>Token: {{ userinfos.token }}</div>
        </div>
      </div>
      <div class="container-options-input">
        <ATooltip
          placement="top"
          title="设置DNS-Rebinding"
          arrow-point-at-center
        >
          <AIcon
            type="setting"
            :style="{ fontSize: '20px', marginRight: '4px' }"
            @click="handleSetDialogVisibility(true)"
          ></AIcon>
        </ATooltip>

        <ASelect v-model="queryMode" style="width: 120px;margin-right: 4px">
          <ASelectOption
            v-for="(mode, index) in queryModes"
            :key="index"
            :value="mode"
          >
            {{ mode.toUpperCase() }}
          </ASelectOption>
        </ASelect>
        <AInput
          v-model="inputData"
          placeholder="请输入过滤字符"
          style="margin-right: 4px"
        ></AInput>
        <AButton
          type="primary"
          @click="handleRefreshRecord"
          style="margin-right: 4px"
        >
          Refresh Record
        </AButton>
        <!-- <AButton type="primary" @click="handleDropIpDataList">
          Clear Data
        </AButton> -->
      </div>
    </div>

    <div class="container-result-area">
      <div v-if="queryMode === queryModes[0]" class="container-result-dns">
        <ATable
          :rowKey="data => data.ts"
          :columns="dnsQueryResult.fields"
          :data-source="dnsQueryResult.dataList"
        ></ATable>
      </div>
      <div v-else class="container-result-http">
        <ATable
          :rowKey="data => data.ts"
          :columns="httpQueryResult.fields"
          :data-source="httpQueryResult.dataList"
        >
        </ATable>
      </div>
    </div>

    <AModal
      title="设置DNS-Rebinding"
      :visible="showDnsRebindingDialog"
      @ok="handleSetDialogVisibility(false)"
      @cancel="handleSetDialogVisibility(false)"
    >
      <div v-if="userinfos.identity" class="container-dns-rebinding-host">
        DNS rebinding host: <a-tag>r.{{ userinfos.identity }}</a-tag>
      </div>
      <div class="container-settings-rebinding">
        <AInput
          placeholder="请输入 IP"
          style="margin-right: 4px"
          @keyup.enter="handleSetRebindingHosts"
          v-model="ipInputField"
        >
        </AInput>
        <AButton type="primary" @click="handleSetRebindingHosts">确认</AButton>
      </div>
      <div
        v-if="rebindingHosts.length > 0"
        class="container-settings-rebinding-list"
      >
        <ATag
          v-for="(ipAddress, index) in rebindingHosts"
          :key="index"
          :closable="true"
          @close="handleRemoveIp(index)"
          class="ip-item"
        >
          {{ ipAddress }}
        </ATag>
      </div>
      <!-- <div class="container-settings-rebinding-clear">
        <AButton
          type="danger"
          icon="close"
          shape="circle"
          size="small"
          @click="handleDropIpDataList"
        />
      </div> -->
    </AModal>
  </div>
</template>

<script>
import dayjs from "dayjs";
import { getCookie } from "../../utils/cookie";
import { apihost } from "../../utils/conf";

const formatTimestamp = ts => {
  return dayjs(ts / 1000000).format("YYYY-MM-DD HH:mm:ss");
};

export default {
  name: "DnsLog",
  data: () => ({
    LogoImg: require("../../assets/logo.png"),
    userinfos: {},
    dnsQueryResult: {
      fields: [
        {
          dataIndex: "name",
          title: "DNS Query Record"
        },
        {
          dataIndex: "remoteAddr",
          title: "Remote Address"
        },
        {
          dataIndex: "ts",
          title: "Created Time",
          customRender: formatTimestamp
        }
      ],
      dataList: []
    },
    httpQueryResult: {
      fields: [
        { dataIndex: "url", title: "HTTP Request Record" },
        {
          dataIndex: "method",
          title: "Method"
        },
        {
          dataIndex: "remoteAddr",
          title: "Remote Address"
        },
        {
          dataIndex: "cookies",
          title: "Cookies"
        },
        {
          dataIndex: "ts",
          title: "Created Time",
          customRender: formatTimestamp
        }
      ],
      dataList: []
    },
    queryModes: ["dns", "http"],
    queryMode: "dns",
    inputData: "",
    showDnsRebindingDialog: false,
    ipInputField: "",
    localStoreRebindingHostKey: "dns_rebinding_hosts",
    rebindingHosts: []
  }),
  methods: {
    initUserinfos() {
      const identity = getCookie("identity");
      const token = getCookie("token");
      if (identity !== "" && token !== "") {
        this.userinfos.identity = getCookie("identity");
        this.userinfos.token = getCookie("token");
      } else {
        this.handleGetSubDomain();
      }
    },
    initRebindingHosts() {
      const localIpDataList = sessionStorage.getItem(
        this.localStoreRebindingHostKey
      );
      if (localIpDataList === null) {
        this.handleGetUserRebindingHosts();
      } else {
        this.rebindingHosts = JSON.parse(localIpDataList);
      }
    },
    handleGetSubDomain() {
      fetch(`${apihost}/v1/users`, { mode: "cors", method: "POST" })
        .then(res => res.json())
        .then(res => {
          const { code, data, message } = res;
          if (code === 200) {
            this.userinfos = data;
            document.cookie = `identity=${data.identity}`;
            document.cookie = `token=${data.token}`;
          } else {
            alert(message);
          }
        })
        .catch(err => {
          alert(err.message);
        });
    },
    handleRefreshRecord() {
      const url = `${apihost}/v1/records?type=${this.queryMode}&token=${this.userinfos.token}&filter=${this.inputData}`;
      fetch(url, {
        mode: "cors"
      })
        .then(res => res.json())
        .then(res => {
          const { code, data, message } = res;
          if (code === 200) {
            this.queryMode === this.queryModes[0]
              ? (this.dnsQueryResult.dataList = data)
              : (this.httpQueryResult.dataList = data);
          } else {
            alert(message);
          }
        })
        .catch(err => {
          alert(err.message);
        });
    },
    handleGetUserRebindingHosts() {
      const identity = this.userinfos.identity.split(".")[0];
      const url = `${apihost}/v1/users/${identity}/dns-rebinding?token=${this.userinfos.token}`;
      fetch(url, {
        mode: "cors"
      })
        .then(res => res.json())
        .then(res => {
          const { code, data, message } = res;
          if (code === 200) {
            this.rebindingHosts =
              data.rebinding_hosts === null ? [] : data.rebinding_hosts;
            this.handleSaveRebindingHosts2Local();
          } else {
            alert(message);
          }
        })
        .catch(err => {
          alert(err.message);
        });
    },
    handleSetUserRebindingHosts() {
      const identity = this.userinfos.identity.split(".")[0];
      const url = `${apihost}/v1/users/${identity}/dns-rebinding`;
      const reqOptions = {
        method: "POST",
        mode: "cors",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          token: this.userinfos.token,
          hosts: this.rebindingHosts
        })
      };
      fetch(url, reqOptions)
        .then(res => res.json())
        .then(res => {
          const { code, data, message } = res;
          if (code === 200) {
            console.log(data);
            this.handleSaveRebindingHosts2Local();
          } else {
            alert(message);
          }
        })
        .catch(err => {
          alert(err.message);
        });
    },
    handleSetDialogVisibility(visible) {
      this.showDnsRebindingDialog = visible;
    },
    handleSaveRebindingHosts2Local() {
      sessionStorage.setItem(
        self.localStoreRebindingHostKey,
        JSON.stringify(this.rebindingHosts)
      );
    },
    handleDropRebindingHosts() {
      this.rebindingHosts = [];
      this.handleSetUserRebindingHosts();
      sessionStorage.removeItem("dns_rebinding_ip_data_list");
      alert("清除成功");
    },
    handleSetRebindingHosts() {
      if (
        RegExp(
          /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/
        ).test(this.ipInputField) &&
        this.rebindingHosts.indexOf(this.ipInputField) === -1
      ) {
        this.rebindingHosts.push(this.ipInputField);
        this.ipInputField = "";
        this.handleSetUserRebindingHosts();
      } else {
        console.log(this.ipInputField);
        alert("输入IP格式错误！");
      }
    },
    handleRemoveIp(index) {
      this.rebindingHosts.splice(index, 1);
      this.handleSetUserRebindingHosts();
    }
  },
  created() {
    this.initUserinfos();
    this.initRebindingHosts();
  }
};
</script>

<style lang="scss" scoped>
.container-app-dns-query {
  margin: auto;
  max-width: 1280px;
  .container-image {
    margin: 12px 0 0 0;
    display: flex;
    flex-direction: row;
    justify-content: center;
  }
  .container-divider {
    height: 2px;
    border: none;
    border-top: 2px dashed #419b4d;
    margin: 16px 0;
  }
  .container-user-area {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: flex-start;

    .container-token-information {
      display: flex;
      flex-direction: column;
      align-items: center;
      font-size: 16px;
      font-weight: bold;
      margin: 24px 0;
    }

    .container-options-input {
      display: flex;
      flex-direction: row;
      justify-content: center;
      align-items: center;
      margin: 16px 0;
    }
  }

  .container-result-area {
  }
}
</style>

<style lang="scss">
.container-settings-rebinding {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: flex-start;
  flex-wrap: nowrap;
}
.container-settings-rebinding-list {
  display: flex;
  flex-direction: row;
  justify-content: flex-start;
  align-items: flex-start;
  flex-wrap: wrap;
  max-width: calc(100% - 64px);
  margin-top: 4px;
  .ip-item {
    margin-top: 4px;
  }
}
.container-dns-rebinding-host {
  padding-bottom: 10px;
}
.container-settings-rebinding-clear {
  margin-top: 32px;
  display: flex;
  flex-direction: row;
  justify-content: flex-end;
}
</style>
