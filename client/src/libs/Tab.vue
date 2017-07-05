<template>
    <div @click="_onTabClick">
        <slot></slot>
    </div>
</template>

<script>
    export default {
        name: 'Tab',
        props: {
            name: '',
        },
        computed: {
            active() {
                return this.$parent.activeName === this.name
            },
        },
        methods: {
            _onTabClick() {
                this.$emit('onTabClick', this)
            },
        },
        mounted() {
            this.$parent._addTab(this)
        },
    }
</script>
