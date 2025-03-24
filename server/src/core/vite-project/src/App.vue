
<template>
  <div class="bg">
    <!-- <img src="" alt=""> -->
    <div class="info">
      <div>用户id</div>
      <div>{{ player.UserId }}</div>
      <div>头像</div>
      <div>{{ player.Icon }}</div>
      <div>卡片</div>
      <div>{{ player.Card }}</div>
      <div>名字</div>
      <div>{{ player.Name }}</div>
      <div>签名</div>
      <div>{{ player.Sign }}</div>
      <div>玩家等级</div>
      <div>{{ player.PlayerLevel }}</div>
      <div>玩家经验</div>
      <div>{{ player.PlayerExp }}</div>
      <div>大世界等级</div>
      <div>{{ player.WorldLevel }}</div>
      <div>大世界现等级</div>
      <div>{{ player.WorldLevelNow }}</div>
      <div>大世界等级冷却</div>
      <div>{{ player.WorldLevelCool }}</div>
      <div>生日</div>
      <div>{{ player.Birth }}</div>
      <div>队伍</div>
      <div>{{ player.ShowTeam }}</div>
      <div>是否隐藏队伍</div>
      <div>{{ player.HideShowTeam }}</div>
      <div>展示卡片</div>
      <div>{{ player.ShowCard }}</div>
    </div>
  </div>
</template>
<script setup lang="ts">
import player from "../save/20221019/player.json";

</script>

<style scoped>
.bg {
  position: relative;
  width: 900px;
  height: 1144px;
  background: url("./assets/character1.png");

}

.info {
  position: absolute;
  left: 40px;
  top: 135px;
  background-color: rgb(22, 20, 16);
  color: white;
  width: 282px;
  display: grid;
  gap: 6px;
  grid-template-columns: auto auto;
  font-size: 23px;
  /* grid-template-rows: 1fr 1fr; */
  justify-content: space-between;
}
</style>
