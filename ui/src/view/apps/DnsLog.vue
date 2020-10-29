<template>
  <div class="container-app-dns-query">
    <div class="container-image">
      <img src="http://hyuga.co/hyuga.png" alt="image" />
    </div>

    <hr class="container-divider" />

    <div class="container-user-area">
      <AButton type="primary" @click="handleGetSubDomain">
        Get SubDomain
      </AButton>
      <div style="min-height: 32px">
        <div v-if="tokenInformation.token" class="container-token-information">
          <div>Domain: {{ tokenInformation.identity }}</div>
          <div>Token: {{ tokenInformation.token }}</div>
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
          :columns="dnsQueryResult.fields"
          :data-source="dnsQueryResult.dataList"
        ></ATable>
      </div>
      <div v-else class="container-result-http">
        <ATable
          :columns="httpQueryResult.fields"
          :data-source="httpQueryResult.dataList"
        ></ATable>
      </div>
    </div>

    <AModal
      title="设置DNS-Rebinding"
      :visible="showDnsRebindingDialog"
      @ok="handleSetDialogVisibility(false)"
      @cancel="handleSetDialogVisibility(false)"
    >
      <div class="container-settings-rebinding">
        <AInput
          placeholder="请输入 IP"
          style="margin-right: 4px"
          @keyup.enter="handleSetTag"
          v-model="ipInputField"
        >
        </AInput>
        <AButton type="primary" @click="handleSetTag">确认</AButton>
      </div>
      <div
        v-if="ipDataList.length > 0"
        class="container-settings-rebinding-list"
      >
        <ATag
          v-for="(ipAddress, index) in ipDataList"
          :key="index"
          closable="t"
          @close="handleRemoveIp(index)"
          class="ip-item"
        >
          {{ ipAddress }}
        </ATag>
      </div>
      <!-- <div class="container-settings-rebinding-clear">
        <AButton type="primary" @click="handleDropIpDataList">
          Clear Data
        </AButton>
      </div> -->
    </AModal>
  </div>
</template>

<script>
export default {
  name: "DnsLog",
  data: () => ({
    tokenInformation: {},
    dnsQueryResult: {
      fields: [
        { dataIndex: "name", title: "DNS Query Record" },
        {
          dataIndex: "remoteAddr",
          title: "Remote Address"
        },
        {
          dataIndex: "ts",
          title: "Created Time"
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
          title: "Created Time"
        }
      ],
      dataList: []
    },
    queryModes: ["dns", "http"],
    queryMode: "dns",
    inputData: "",
    showDnsRebindingDialog: false,
    ipInputField: "",
    ipDataList: []
  }),
  methods: {
    handleGetSubDomain() {
      fetch("http://api.hyuga.co/v1/users", { mode: "cors", method: "POST" })
        .then(res => res.json())
        .then(res => {
          const { code, data, message } = res;
          if (code === 200) {
            this.tokenInformation = data;
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
      const url = `http://api.hyuga.co/v1/records?type=${this.queryMode}&token=${this.tokenInformation.token}&filter=${this.inputData}`;
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
    handleSetDialogVisibility(visible) {
      this.showDnsRebindingDialog = visible;
    },
    handleSaveIpDataList() {
      sessionStorage.setItem(
        "dns_rebinding_ip_data_list",
        JSON.stringify(this.ipDataList)
      );
    },
    handleDropIpDataList() {
      this.ipDataList = [];
      sessionStorage.removeItem("dns_rebinding_ip_data_list");
      alert("清除成功");
    },
    handleSetTag() {
      if (
        RegExp(
          /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/
        ).test(this.ipInputField) &&
        this.ipDataList.indexOf(this.ipInputField) === -1
      ) {
        this.ipDataList.push(this.ipInputField);
        this.ipInputField = "";
        this.handleSaveIpDataList();
      }
    },
    handleRemoveIp(index) {
      this.ipDataList.splice(index, 1);
      this.handleSaveIpDataList();
    }
  },
  created() {
    const localIpDataList = sessionStorage.getItem("dns_rebinding_ip_data_list");
    this.ipDataList =
      localIpDataList === null ? [] : JSON.parse(localIpDataList);
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

.container-settings-rebinding-clear {
  margin-top: 32px;
  display: flex;
  flex-direction: row;
  justify-content: flex-end;
}
</style>
