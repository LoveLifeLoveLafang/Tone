<template>
    <div class="tab">
        <router-link class="tab-item" tag="div" to="/recommend">
            <span class="tab-link">推荐</span>
        </router-link>
        <router-link class="tab-item" tag="div" to="/singer">
            <span class="tab-link">歌手</span>
        </router-link>
        <router-link class="tab-item" tag="div" to="/rank">
            <span class="tab-link">排行</span>
        </router-link>
        <router-link class="tab-item" tag="div" to="/search">
            <span class="tab-link">搜索</span>
        </router-link>
    </div>
</template>

<script type="text/javascript">
export default {}
</script>

<style type="text/css">
.tab {
    display: flex;
    height: 44px;
    line-height: 44px;
    font-size: 16px;
}

.tab-item {
    flex:1;
    text-align: center;
}

.tab-link {
    padding-bottom: 5px;
    color: #A794C1;
}
.router-link-active .tab-link{
    color: #7A6090;
    border-bottom: 2px solid #7A6090;
}
</style>
