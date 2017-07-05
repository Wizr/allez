<template>
    <nav class="nav">
        <div class="nav-left">
            <div class="tab-bar" :style="barStyle"></div>
            <slot></slot>
        </div>
    </nav>
</template>

<script>
    export default {
        name: 'Tabs',
        props: {
            activeName: '',
        },
        data() {
            return {
                tabs: [],
            }
        },
        computed: {
            barStyle: {
                cache: false,
                get() {
                    let style = {}
                    let offset = 0
                    let tabWidth = 0

                    this.tabs.every(tab => {
                        if (tab.active) {
                            tabWidth = tab.$el.clientWidth
                            return false
                        } else {
                            offset += tab.$el.clientWidth
                            return true
                        }
                    })

                    style.width = tabWidth + 'px'
                    style.transform = `translateX(${offset}px)`
                    return style
                },
            },
        },
        methods: {
            _addTab(tab) {
                this.tabs.push(tab)
            },
        },
    }
</script>

<style scoped>
    .tab-bar {
        width: 80px;
        height: 3px;
        position: absolute;
        bottom: 0;
        background-color: #20a0ff;
        transition: transform .3s cubic-bezier(.645, .045, .355, 1);
    }

    .nav {
        border-bottom: 1px solid #ccc;
    }
</style>
