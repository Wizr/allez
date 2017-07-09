<template>
    <div class="columns justify-center">
        <div class="column is-narrow" v-if="value && value['result']">
            <div class="box" v-for="(info, name) in value['result']">
                <span>{{ name }}</span>
                <span class="account-field">{{ info['AppleID'] }}</span>
                <span class="account-field">{{ info['Password'] }}</span>
            </div>
        </div>
        <span class="column is-narrow loading" v-else="">正在召唤，请稍后再试...</span>
    </div>
</template>

<script>
    export default {
        name: 'AppStore',
        asyncComputed: {
            value() {
                return fetch('/api/keygen/appstore', {
                    method: 'GET',
                }).then(response => response.json()).catch(() => 'Internal Error')
            },
        },
        created() {
            this.$emit('created', 'AppStore')
        },
    }
</script>

<style scoped>
    .loading {
        margin-top: 40px;
    }

    .box:not(:last-child) {
        margin-top: 20px;
    }

    .account-field {
        margin-left: 30px;
    }
</style>
