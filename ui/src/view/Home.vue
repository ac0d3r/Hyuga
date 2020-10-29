<template>
  <ALayout id="app-viewer" class="container-page">
    <ALayoutHeader class="container-page-header">
      <div class="container-page-header-logo">
        <img
          src="http://hyuga.co/hyuga.png"
          alt="logo"
        />
      </div>
      <div class="container-page-header-menu">
        <AMenu mode="horizontal" theme="light">
          <AMenuItem key="home"><a-icon type="home" />Home</AMenuItem>
          <AMenuItem key="home1"><a-icon type="home" />Home1</AMenuItem>
          <AMenuItem key="home2"><a-icon type="home" />Home2</AMenuItem>
          <AMenuItem key="home3"><a-icon type="home" /><a href="https://antdv.com" target="_blank" rel="noopener noreferrer">11</a></AMenuItem>
        </AMenu>
      </div>
      <!-- <div>
        <AButton>
          <AIcon type="user" />
          登录
        </AButton>
      </div> -->
    </ALayoutHeader>
    <ALayoutContent class="container-content">
      <div class="section-content">
        <div class="section-content-header">
          <div class="section-content-header-left">
            <h1 class="title">这是标题</h1>
            <p class="description">这是一段副标题描述</p>
            <div class="action-buttons">
              <AButton
                type="primary"
                shape="round"
                size="large"
                style="width:110px;"
              >
                登录
              </AButton>
              <AButton
                shape="round"
                size="large"
                style="width:110px;margin: 0px 16px;"
              >
                注册账户
              </AButton>
            </div>
          </div>
          <div class="section-content-header-right">
            <a-carousel autoplay>
              <img
                src="https://static.ctfhub.com/website/cover/3.png"
                alt="image"
              />
              <img
                src="https://static.ctfhub.com/website/cover/1.png"
                alt="image"
              />
            </a-carousel>
          </div>
        </div>
        <div class="section-content-title-join">
          <h2>Hyuga<span> 工具列表</span></h2>
        </div>
        <div class="section-content-detail">
          <ACol
            v-for="(item, index) in detailList"
            :key="index"
            :span="8"
            class="detail-list-item"
          >
            <ACard hoverable>
              <template slot="actions" class="ant-card-actions">
                <a-icon key="setting" type="setting" />
                <a-icon key="edit" type="edit" />
                <a-icon key="ellipsis" type="ellipsis" />
              </template>
              <ACardMeta :title="item.title" :description="item.description">
                <AAvatar slot="avatar" :src="item.avatar" />
              </ACardMeta>
            </ACard>
          </ACol>
        </div>
      </div>
    </ALayoutContent>
    <ALayoutFooter class="container-footer">
      Copyright © 2020 Buzz2d0 All Rights Reserved.
      <ADivider type="vertical" />
      苏ICP备2020051327号
    </ALayoutFooter>
  </ALayout>
</template>

<script>
export default {
  name: "Home",
  data: () => ({
    detailList: [
      {
        title: "卡片标题",
        description: "卡片描述卡片描述卡片描述",
        avatar: "https://static.ctfhub.com/website/index/2-skill.svg"
      },
      {
        title: "卡片标题",
        description: "卡片描述卡片描述卡片描述",
        avatar: "https://static.ctfhub.com/website/index/2-skill.svg"
      },
      {
        title: "卡片标题",
        description: "卡片描述卡片描述卡片描述",
        avatar: "https://static.ctfhub.com/website/index/2-skill.svg"
      },
      {
        title: "卡片标题",
        description: "卡片描述卡片描述卡片描述",
        avatar: "https://static.ctfhub.com/website/index/2-skill.svg"
      },
      {
        title: "卡片标题",
        description: "卡片描述卡片描述卡片描述",
        avatar: "https://static.ctfhub.com/website/index/2-skill.svg"
      },
      {
        title: "卡片标题",
        description: "卡片描述卡片描述卡片描述",
        avatar: "https://static.ctfhub.com/website/index/2-skill.svg"
      }
    ]
  })
};
</script>

<style lang="scss" scoped>
.container-page {
  width: calc(100% - 4px);
  max-width: 100vw !important;
  min-height: 100vh;
  .container-page-header {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
    box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
    background-color: #ffffff;
    &-logo {
      img {
        max-width: 32px;
      }
    }
    &-menu {
    }
  }
  .container-content {
    min-height: 100%;
    .section-content {
      margin: auto;
      max-width: 1280px;
      padding-top: 24px;
      display: flex;
      flex-direction: column;
      justify-content: flex-start;
      align-items: flex-start;
      &-header {
        display: flex;
        flex-direction: row;
        justify-content: space-around;
        margin: 16px 0;
        align-items: center;
        height: 512px;
        width: 100%;
        &-left {
          width: 40%;
          display: flex;
          flex-direction: column;
          justify-content: space-between;
          align-items: flex-start;
          max-width: 480px;
          height: 245px;
          max-height: 245px !important;
          .title {
            font-size: 54px;
            color: #000000d9;
          }
          .description {
            font-size: 16px;
            color: rgb(105, 123, 140);
          }
          .action-buttons {
            display: flex;
            flex-direction: row;
            align-items: center;
            justify-content: flex-start;
          }
        }
        &-right {
          width: 40%;
          max-width: 598px;
          height: 324px;
          img {
            height: 324px !important;
            width: 598px !important;
          }
        }
      }
      &-title-join {
        text-align: center;
        font-weight: 300;
        font-size: 28px;
        color: #314659;
        letter-spacing: 0.6px;
        width: 100%;
        span {
          font-weight: 600;
        }
      }
      &-detail {
        width: 100%;
        margin-top: 16px;
        display: flex;
        flex-direction: row;
        justify-content: flex-start;
        flex-wrap: wrap;
        .detail-list-item {
          padding: 24px;
        }
      }
    }
  }
  .container-footer {
    width: inherit;
    display: flex;
    flex-direction: row;
    justify-content: center;
  }
}
</style>
