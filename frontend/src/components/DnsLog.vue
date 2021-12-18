<template>
  <div class="container-app-dns-query">
    <div class="container-user-area">
      <a-card hoverable style="width: 500px">
        <div>
          <a-icon slot="prefix" type="user" />:
          <a-tag v-if="User.ID !== ''">{{ User.ID }}</a-tag>
          <a-icon
            v-if="User.ID !== ''"
            type="copy"
            @click="
              () => {
                copy(User.ID);
              }
            "
          />
        </div>
        <div>
          <a-icon slot="prefix" type="key" />:
          <a-tag v-if="User.Token !== ''">{{ User.Token }}</a-tag>
          <a-icon
            v-if="User.Token !== ''"
            type="copy"
            @click="
              () => {
                copy(User.Token);
              }
            "
          />
        </div>

        <template slot="actions" class="ant-card-actions">
          <a-popconfirm
            title="Delete the current domain name and get a new one?"
            ok-text="Yes"
            cancel-text="No"
            @confirm="handleCreateUser"
          >
            <a-button size="small" type="primary" shape="circle" icon="redo" />
          </a-popconfirm>
          <a-icon
            size="small"
            type="setting"
            @click="
              () => {
                Shows.DnsRebindSetting = true;
                User.OldIPs = User.IPs;
              }
            "
          />
          <a-popconfirm
            title="Wipe all records data?"
            ok-text="Yes"
            cancel-text="No"
            @confirm="handleWipeData"
          >
            <a-button size="small" type="danger" shape="circle" icon="close" />
          </a-popconfirm>
        </template>
      </a-card>

      <div class="container-options-input">
        <a-select
          v-model="CurrentQueryMode"
          style="width: 120px; margin-right: 4px"
          @change="handleSelectQueryModeChange"
        >
          <a-select-option
            v-for="(mode, index) in QueryModes"
            :key="index"
            :value="mode"
          >
            {{ mode }}
          </a-select-option>
        </a-select>

        <a-input
          v-model="Inputs.Filter"
          placeholder="filter characters"
          style="margin-right: 4px"
        ></a-input>
        <a-button
          type="primary"
          style="margin-right: 4px"
          @click="getLogRecords"
        >
          Refresh Record
        </a-button>
      </div>
    </div>

    <div>
      <a-table
        :columns="
          CurrentQueryMode === 'dns' ? DNSResult.fields : HttpResult.fields
        "
        :data-source="
          CurrentQueryMode === 'dns' ? DNSResult.data : HttpResult.data
        "
      >
        <a slot="url" slot-scope="text">{{ text }}</a>
        <span slot="httpraw" slot-scope="rawText">
          <a-icon
            type="copy"
            @click="
              () => {
                copy(rawText);
              }
            "
          />
          <pre>{{ formatHttpRaw(rawText) }}</pre>
        </span>
      </a-table>
    </div>

    <a-modal
      title="Set DNS Rebinding"
      :visible="Shows.DnsRebindSetting"
      @ok="handleUpdateUserDnsRebindingHosts"
      @cancel="
        () => {
          Shows.DnsRebindSetting = false;
        }
      "
    >
      <div v-if="User.ID">
        Host: <a-tag>r.{{ User.ID }}</a-tag
        ><a-icon
          v-if="User.ID !== ''"
          type="copy"
          @click="
            () => {
              copy('r.' + User.ID);
            }
          "
        />
      </div>
      <br />
      <div>
        <template v-for="dns in User.IPs">
          <a-tag
            :key="dns"
            :closable="true"
            @close="
              () => {
                handleCloseDns(dns);
              }
            "
          >
            {{ dns }}
          </a-tag>
        </template>

        <a-input
          v-if="Shows.AddDns"
          ref="input"
          type="text"
          size="small"
          :style="{ width: '78px' }"
          :value="Inputs.DnsRebinding"
          @change="handleTagInputChange"
          @blur="handleTagInputBlur"
          @keyup.enter="handleTagInputConfirm"
        />
        <a-tag
          v-else
          @click="handleShowAddDns"
          style="background: #fff; borderstyle: dashed"
        >
          <a-icon type="plus" /> New DNS
        </a-tag>
      </div>
    </a-modal>
  </div>
</template>


<script>
import dayjs from "dayjs";
import {
  CreateUser,
  DeleteUser,
  GetUserDnsRebindingHosts,
  UpdateUserDnsRebindingHosts,
  GetLogRecords,
  WipeRecodsData,
} from "../utils/apis";
import { equar, validateIPaddress } from "../utils/util";

const formatTimestamp = (created) => {
  const parsed = parseInt(created, 10);
  if (isNaN(parsed)) {
    return "Wrong time";
  }
  return dayjs(parsed * 1000).format("YYYY-MM-DD HH:mm:ss");
};

export default {
  name: "DnsLog",
  data: () => ({
    User: {
      ID: "",
      Token: "",
      IPs: [],
      OldIPs: [],
    },
    QueryModes: ["dns", "http"],
    CurrentQueryMode: "dns",

    DNSResult: {
      fields: [
        { dataIndex: "name", title: "DNS Query Record", key: "name" },
        {
          dataIndex: "remote_addr",
          title: "Remote Address",
          key: "remote_addr",
        },
        {
          key: "created",
          dataIndex: "created",
          title: "Created Time",
          customRender: formatTimestamp,
        },
      ],
      data: [],
    },
    HttpResult: {
      fields: [
        {
          dataIndex: "url",
          title: "URL",
          key: "url",
          scopedSlots: { customRender: "url" },
        },
        { dataIndex: "method", title: "Method", key: "method" },
        {
          dataIndex: "remote_addr",
          title: "Remote Address",
          key: "remote_addr",
        },
        {
          key: "raw",
          dataIndex: "raw",
          title: "HTTP Raw",
          scopedSlots: { customRender: "httpraw" },
        },
        {
          key: "created",
          dataIndex: "created",
          title: "Created Time",
          customRender: formatTimestamp,
        },
      ],
      data: [],
    },
    Inputs: { Filter: "", DnsRebinding: "" },
    Shows: { DnsRebindSetting: false, AddDns: false },
  }),
  methods: {
    handleCloseDns(removedDns) {
      const ips = this.User.IPs.filter((tag) => tag !== removedDns);
      this.User.IPs = ips;
    },
    handleShowAddDns() {
      this.Shows.AddDns = true;
      this.$nextTick(function () {
        this.$refs.input.focus();
      });
    },
    handleTagInputChange(e) {
      this.Inputs.DnsRebinding = e.target.value;
    },
    handleTagInputBlur() {
      this.Shows.AddDns = false;
      this.Inputs.DnsRebinding = "";
    },
    handleTagInputConfirm() {
      if (!validateIPaddress(this.Inputs.DnsRebinding)) {
        this.$message.error("You have entered an invalid IP address!");
      } else {
        const dns = this.Inputs.DnsRebinding;
        let ips = this.User.IPs;
        if (dns && ips.indexOf(dns) === -1) {
          ips = [...ips, dns];
        }
        this.User.IPs = ips;
      }

      this.Shows.AddDns = false;
      this.Inputs.DnsRebinding = "";
    },
    handleUpdateUserDnsRebindingHosts() {
      this.Shows.DnsRebindSetting = false;
      if (!equar(this.User.OldIPs, this.User.IPs)) {
        this.updateUserDnsRebindingHosts();
      }
    },
    handleCreateUser() {
      if (this.User.ID !== "" && this.User.Token !== "") {
        this.deleteUser(true);
      } else {
        this.createUser();
      }
    },
    handleWipeData() {
      this.wipeRecodsData();
    },
    handleSelectQueryModeChange() {
      this.getLogRecords();
    },
    formatHttpRaw(raw) {
      return raw.length > 80 ? `${raw.substring(0, 80)}...` : raw;
    },
    copy: function (msg) {
      this.$copyText(msg).then(
        () => {
          this.$message.success("copied");
        },
        () => {
          this.$message.warning("can not copy");
        }
      );
    },
    fail(msg) {
      this.$message.error(msg);
    },
    createUser() {
      const succ = (data) => {
        this.User.ID = data.id;
        this.User.Token = data.token;
        this.$cookies.set("identity", data.id, -1, "/");
        this.$cookies.set("token", data.token, -1, "/");
        this.$message.success("user created successfully");
        this.getUserDnsRebindingHosts();
      };
      CreateUser(succ, this.fail);
    },
    deleteUser(needCreate) {
      const succ = () => {
        this.User = { ID: "", Token: "", IPs: [], OldIPs: [] };
        this.$cookies.remove("identity");
        this.$cookies.remove("token");
        this.$message.success("user deleted successfully");
        needCreate ? this.createUser() : {};
      };
      DeleteUser(succ, this.fail);
    },
    getUserDnsRebindingHosts() {
      const succ = (data) => {
        this.User.IPs = data;
      };
      const fail = (msg, code) => {
        this.fail(msg);
        if (code === 200) {
          this.User = { ID: "", Token: "", IPs: [], OldIPs: [] };
          this.$cookies.remove("identity");
          this.$cookies.remove("token");
        }
      };
      GetUserDnsRebindingHosts(succ, fail);
    },
    updateUserDnsRebindingHosts() {
      const succ = () => {
        this.$message.success("hosts updated successfully");
      };
      const fail = (msg) => {
        this.fail(msg);
        this.User.IPs = this.User.OldIPs;
      };
      UpdateUserDnsRebindingHosts(this.User.IPs, succ, fail);
    },
    getLogRecords() {
      const succ = (data) => {
        switch (this.CurrentQueryMode) {
          case "dns":
            this.DNSResult.data = data;
            break;
          case "http":
            this.HttpResult.data = data;
            break;
        }
      };
      GetLogRecords(this.CurrentQueryMode, this.Inputs.Filter, succ, this.fail);
    },
    wipeRecodsData() {
      const succ = () => {
        this.DNSResult.data = [];
        this.HttpResult.data = [];
        this.$message.success("records wiped successfully");
      };
      WipeRecodsData(succ, this.fail);
    },
    initUser() {
      this.User.ID = this.$cookies.get("identity");
      this.User.Token = this.$cookies.get("token");
      if (this.User.ID !== "" && this.User.Token !== "") {
        this.getUserDnsRebindingHosts();
        this.getLogRecords();
      }
    },
  },
  created() {
    this.initUser();
  },
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
  .container-user-area {
    margin-top: 20px;
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
}
</style>

