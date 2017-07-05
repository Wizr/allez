<template>
    <div>
        <h1 class="title has-text-centered">注册码</h1>
        <tabs :activeName="activeName">
            <tab class="nav-item"
                 :class="{'is-active': tab == activeName}"
                 v-for="tab in tabs"
                 :key="tab"
                 :name="tab"
                 v-on:onTabClick="onTabClick">
                <a>{{ tab }}</a>
            </tab>
        </tabs>
        <router-view @created="handleSubCompCreated"></router-view>
    </div>
</template>

<script>
    import Tabs from './libs/Tabs.vue'
    import Tab from './libs/Tab.vue'
    import router from './router'

    export default {
        components: {
            Tabs,
            Tab,
        },
        data() {
            return {
                activeName: '',
                tabs: ['Editplus', 'Charles'],
            }
        },
        methods: {
            onTabClick(tab) {
                router.push(`${tab.name.toLowerCase()}`)
            },
            handleSubCompCreated(tabName) {
                this.activeName = tabName
            },
        },
    }
</script>

<style scoped>
    a.nav-item.is-active {
        color: blue
    }
</style>
